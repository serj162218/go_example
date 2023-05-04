package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"websocket_with_redis/initializers"

	"github.com/google/uuid"
	"github.com/olahol/melody"
)

// To store the connected socket infomation.
type SocketInfo struct {
	UUID string
	ID   int
}

// The structure of message which is using for communication between websocket.
type message struct {
	Msg    string
	From   int
	To     int
	Status int
}

var mappingSession map[string]*melody.Session

func Initialize() *melody.Melody {
	websocket := melody.New()
	mappingSession = map[string]*melody.Session{}

	websocket.HandleConnect(func(session *melody.Session) {
		//set uuid => session in map
		uuid := uuid.NewString()
		session.Set("info", &SocketInfo{
			UUID: uuid,
		})
		mappingSession[uuid] = session
	})

	websocket.HandleDisconnect(func(session *melody.Session) {
		value, exists := session.Get("info")

		if !exists {
			return
		}

		//Removing this connect's infomation from map and redis.
		info := value.(*SocketInfo)
		delete(mappingSession, info.UUID)
		initializers.Redisdb.SRem(context.TODO(), "Conn"+strconv.Itoa(info.ID), info.UUID)
	})
	websocket.HandleMessage(func(session *melody.Session, msg []byte) {
		//Json parse
		data := message{}
		err := json.Unmarshal(msg, &data)
		if err != nil {
			fmt.Println(err)
		}

		value, exists := session.Get("info")

		if !exists {
			return
		}

		info := value.(*SocketInfo)
		/*
		* status 1 => send message, 2 => connect to websocket.
		 */
		switch data.Status {
		case 1:
			r, _ := json.Marshal(data)

			//In this part, you can store the data into database first, then sending the message.
			memberTo := initializers.Redisdb.SMembers(context.TODO(), fmt.Sprint("Conn", data.To))
			memberToArr, _ := memberTo.Result()
			for _, e := range memberToArr {
				if session, isExist := mappingSession[e]; isExist {
					session.Write(r)
				}
			}

			memberFrom := initializers.Redisdb.SMembers(context.TODO(), fmt.Sprint("Conn", data.From))
			memberFromArr, _ := memberFrom.Result()
			for _, e := range memberFromArr {
				if session, isExist := mappingSession[e]; isExist {
					session.Write(r)
				}
			}
		case 2:
			loginWs(&data, session, info)
		}
	})

	return websocket
}

func loginWs(data *message, session *melody.Session, info *SocketInfo) {
	//Save this id => uuid to redis then change session's SocketInfo Id
	var redisKey string = "Conn"
	saveToRedis(info.UUID, data.From, redisKey)
	info.ID = data.From
}

func saveToRedis(uuid string, id int, key string) {
	err := initializers.Redisdb.SAdd(context.TODO(), key+strconv.Itoa(id), uuid).Err()
	if err != nil {
		panic(err)
	}
}
