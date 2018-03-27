# perfchat

**p2pclient** is implementation of chat simulator, which works as a service on top of ethereum p2p engine

It has ability to spawn any number of workers by using addnodes.sh  script. 

#Installation
Install  dependencies
```bash
# if ubuntu/debian
sudo apt install gcc
go get github.com/ethereum/go-ethereum
go install github.com/ethereum/go-ethereum/cmd/p2psim
```

Download source of p2pclient and install it:
```bash
go get github.com/phob0s-pl/perfchat
go install github.com/phob0s-pl/perfchat/cmd/p2pclient
```

#Running
```bash
# In one window
p2pclient

# In second window we will communicate with p2pclient and spawn nodes there
$GOPATH/src/github.com/phob0s-pl/perfchat/cmd/p2pclient/addnodes.sh 10

```

#Verification
Each 10s all nodes print simple statistics:
```bash
INFO [03-27|21:53:45] Stats: node.id=0f6e7cebda9c65fa msg_received=5964  groups_created=19 groups_exited=3 msg_sent=87277
```

Webinfo with the same statistics is available at
```
http://localhost:8888/
```


