package main

import (
	"io/ioutil"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/rpc"
)

var (
	roomCreatePeriod = time.Second * 1
	roomExitPeriod   = time.Second * 5
	msgSendPeriod    = time.Millisecond
	roomJoinPeriod   = time.Second * 2
	statsPeriod      = time.Second * 10
)

type Stats struct {
	Msgsend      uint64 `json:"msgsend"`
	Msgreceived  uint64 `json:"msgreceived"`
	Roomscreated uint64 `json:"roomscreated"`
	Roomsexited  uint64 `json:"roomsexited"`
}

// ChatService runs a chat simulation service between nodes
// where each node can send private message, create chat-group or send message to group
type ChatService struct {
	stats          Stats
	chat           *Chat
	id             discover.NodeID
	log            log.Logger
	received       int64
	msgChan        []chan string
	roomJoinChan   []chan string
	roomCreateChan []chan string
	roomExitChan   []chan string
}

func newChatService(id discover.NodeID) *ChatService {
	return &ChatService{
		id:  id,
		log: log.New("node.id", id),
	}
}

func (c *ChatService) Protocols() []p2p.Protocol {
	return []p2p.Protocol{{
		Name:     "chat-service",
		Version:  1,
		Length:   4,
		Run:      c.Run,
		NodeInfo: c.Info,
	}}
}

func (c *ChatService) APIs() []rpc.API {
	return nil
}

func (c *ChatService) Start(server *p2p.Server) error {
	c.log.Info("chat-service starting")
	c.chat = NewChat(c.id.String())

	go c.PrintStats()
	go c.PeriodicMsgSend()
	go c.PeriodicGroupCreate()
	go c.PeriodicGroupExit()
	go c.PeriodicGroupJoin()

	return nil
}

func (c *ChatService) PeriodicMsgSend() {
	for range time.Tick(msgSendPeriod) {
		group := c.chat.GetRandomGroup()
		if group == "" {
			continue
		}
		for i := range c.msgChan {
			atomic.AddUint64(&c.stats.Msgsend, 1)
			c.msgChan[i] <- group
		}
	}
}

func (c *ChatService) PeriodicGroupExit() {
	for range time.Tick(roomExitPeriod) {
		group := c.chat.ExitRandomGroup()
		if group == "" {
			continue
		}
		atomic.AddUint64(&c.stats.Roomsexited, 1)
		for i := range c.roomExitChan {
			c.roomExitChan[i] <- group
		}
	}
}

func (c *ChatService) PeriodicGroupJoin() {
	for range time.Tick(roomJoinPeriod) {
		group := c.chat.JoinRandomGroup()
		if group == "" {
			continue
		}

		for i := range c.roomJoinChan {
			c.roomJoinChan[i] <- group
		}
	}
}

func (c *ChatService) PeriodicGroupCreate() {
	for range time.Tick(roomCreatePeriod) {
		group := c.chat.CreateRandomGroup()
		if group == "" {
			continue
		}
		atomic.AddUint64(&c.stats.Roomscreated, 1)

		for i := range c.roomCreateChan {
			c.roomCreateChan[i] <- group
		}
	}
}

func (c *ChatService) Stop() error {
	c.log.Info("chat-service stopping")
	return nil
}

func (c *ChatService) Info() interface{} {
	return Stats{
		Msgreceived:  atomic.LoadUint64(&c.stats.Msgreceived),
		Roomscreated: atomic.LoadUint64(&c.stats.Roomscreated),
		Roomsexited:  atomic.LoadUint64(&c.stats.Roomsexited),
		Msgsend:      atomic.LoadUint64(&c.stats.Msgsend),
	}
}

const (
	CreateRoomCode = iota
	JoinRoomCode
	SendMsgCode
	// ExitGroup if called by owner, will delete chat group
	// If called by user, will exit chat group
	ExitGroup
)

// Run implements the chat protocol between peers
func (c *ChatService) Run(peer *p2p.Peer, rw p2p.MsgReadWriter) error {
	go c.ReceiveMessages(peer, rw)
	msgs := make(chan string)
	roomC := make(chan string)
	roomE := make(chan string)
	roomJ := make(chan string)

	c.roomCreateChan = append(c.roomCreateChan, roomC)
	c.roomExitChan = append(c.roomExitChan, roomE)
	c.msgChan = append(c.msgChan, msgs)
	c.roomJoinChan = append(c.roomJoinChan, roomJ)

	for {
		select {
		case roomName := <-msgs:
			if c.chat.IsGroupMember(peer.ID().String(), roomName) {
				if err := p2p.Send(rw, SendMsgCode, roomName); err != nil {
					log.Warn("failed to send msg", "err", err)
				}
			}
		case roomName := <-roomC:
			if err := p2p.Send(rw, CreateRoomCode, roomName); err != nil {
				log.Warn("failed to send msg", "err", err)
			}
		case roomName := <-roomE:
			if err := p2p.Send(rw, ExitGroup, roomName); err != nil {
				log.Warn("failed to send msg", "err", err)
			}
		case roomName := <-roomJ:
			if err := p2p.Send(rw, JoinRoomCode, roomName); err != nil {
				log.Warn("failed to send msg", "err", err)
			}
		}
	}
}

func (c *ChatService) PrintStats() {
	for range time.Tick(statsPeriod) {
		c.log.Info("Stats:",
			"msg_received", atomic.LoadUint64(&c.stats.Msgreceived),
			"groups_created", atomic.LoadUint64(&c.stats.Roomscreated),
			"groups_exited", atomic.LoadUint64(&c.stats.Roomsexited),
			"msg_sent", atomic.LoadUint64(&c.stats.Msgsend),
		)
	}

}

func (c *ChatService) ReceiveMessages(peer *p2p.Peer, rw p2p.MsgReadWriter) {
	for {
		msg, err := rw.ReadMsg()
		if err != nil {
			c.log.Warn("failed to read msg", "err", err)
			continue
		}
		payloadRaw, err := ioutil.ReadAll(msg.Payload)
		if err != nil {
			c.log.Warn("failed to get payload", "err", err)
			continue
		}
		payload := string(decodePayload(payloadRaw))

		atomic.AddUint64(&c.stats.Msgreceived, 1)
		if msg.Code == CreateRoomCode {
			c.log.Trace("Notified that room was created", "roomname", payload)
			c.chat.AddKnownGroup(payload, peer.ID().String())
		}
		if msg.Code == JoinRoomCode {
			c.log.Trace("Notified that user joined room", "roomname", payload)
			c.chat.AddUserToGroup(payload, peer.ID().String())
		}
		if msg.Code == SendMsgCode {
			c.log.Trace("Received message")
		}

		if msg.Code == ExitGroup {
			c.log.Trace("Notified that user exited or deleted room", "roomname", payload)
			c.chat.ExitGroup(payload, peer.ID().String())
		}
	}
}
