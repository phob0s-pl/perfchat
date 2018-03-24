package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	api "github.com/phob0s-pl/perfchat/apiv1"
	"github.com/phob0s-pl/perfchat/chat"
	log "github.com/sirupsen/logrus"
)

const (
	// ConfigPath is default configuration path
	ConfigPath = "client.conf"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
}

func main() {
	configPath := flag.String("conf", ConfigPath, "path to config file")
	flag.Parse()

	config, err := ReadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to read config, err=%s", err)
	}

	if config.Debug {
		log.SetLevel(log.DebugLevel)
	}
	log.Debugf("Read configuration from %q: %+v", *configPath, config)

	admin := &chat.User{
		Role:   chat.AdminRole,
		AuthID: config.AuthID,
		Token:  config.Token,
		Name:   "admin",
	}
	adminClient := api.NewClient(admin, config.Address)

	for i := uint(0); i < config.Workers; i++ {
		workerUser := &chat.User{
			AuthID: RandString(),
			Name:   RandString(),
			Role:   chat.UserRole,
			Token:  RandString(),
		}
		if err := adminClient.AddUser(workerUser); err != nil {
			log.Fatalf("Failed to add user, err=%s", err)
		}
		go worker(config, workerUser)
	}

	select {}
}

func worker(c *Config, user *chat.User) {
	var (
		msgCount uint
		msg      = fmt.Sprintf("%s_msg", user.Name)
		client   = api.NewClient(user, c.Address)
		roomT    = time.NewTicker(time.Duration(c.RoomOp) * time.Millisecond)
		messageT = time.NewTicker(time.Duration(c.MessageToUserChance) * time.Millisecond)
		myRoom   string
	)

	for {
		select {
		case <-messageT.C:
			users, err := client.GetUsers()
			if err != nil {
				log.Debugf("Failed to list users, err=%s", err)
				continue
			}
			if len(users) == 0 {
				continue
			}

			rooms, err := client.GetRooms()
			if err != nil {
				log.Debugf("Failed to list rooms, err=%s", err)
				continue
			}
			randomUser := users[randSrc.Uint32()%uint32(len(users))]
			if randomUser.Name == "admin" {
				continue
			}
			dstRoom := findUserInRoom(randomUser.Name, rooms)

			if err := client.SendMessage(&api.Message{Content: msg, Room: dstRoom}); err != nil {
				log.Debugf("Failed to send msg, err=%s", err)
			}
			log.Debugf("[%d]msg sent", msgCount)
			msgCount++

		case <-roomT.C:
			_ = client.RoomDelete(myRoom)
			myRoom = RandString()
			_ = client.RoomCreate(myRoom)

		}
	}
}

func findUserInRoom(username string, rooms []api.Room) string {
	for _, room := range rooms {
		for _, roomUser := range room.Users {
			if roomUser == username {
				return room.Name
			}
		}
	}
	return "i_want_to_sleep_:<"
}
