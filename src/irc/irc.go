package irc

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
)

type TwitchUserTracker struct {
	Name      string
	Id        string
	Timestamp time.Time
}

type UserInfo struct {
	Name   string `json:"name"`
	Id     string `json:"id"`
	Online bool   `json:"online"`
}

type TwitchUserDatabase struct {
	Client  *twitch.Client
	UserMap map[string]map[string]TwitchUserTracker
}

func CreateTwitchDatabase() *TwitchUserDatabase {
	client := twitch.NewAnonymousClient()
	log.Println("Created Twitch IRC client")
	channelList := strings.Split(os.Getenv("TWITCH_CHANNELS"), ",")
	userMap := make(map[string]map[string]TwitchUserTracker)
	client.Join(channelList...)
	client.OnPrivateMessage(createPrivateMsgCallback(userMap))
	log.Println("Client configuration finished")

	// err := client.Connect()
	// if err != nil {
	// 	log.Println("an error occurred")
	// 	log.Fatal(err)
	// 	panic(err)
	// }

	log.Println("Client connected")
	return &TwitchUserDatabase{
		Client:  client,
		UserMap: userMap,
	}
}

func createPrivateMsgCallback(userMap map[string]map[string]TwitchUserTracker) func(twitch.PrivateMessage) {
	return func(msg twitch.PrivateMessage) {
		channelName := msg.Channel
		if value, exists := userMap[channelName]; !exists {
			// channel doesn't exist on the map yet, make it
			userMap[channelName] = make(map[string]TwitchUserTracker)
		} else {
			userName := msg.User.Name
			userTracker := TwitchUserTracker{
				Name:      userName,
				Id:        msg.User.ID,
				Timestamp: msg.Time,
			}
			value[userName] = userTracker
		}

	}
}

func (db *TwitchUserDatabase) ReadUserInfo(channel, user string) UserInfo {
	userMap := db.UserMap
	if value, exists := userMap[channel]; !exists {
		return UserInfo{
			Online: false,
		}
	} else {
		twitchInfo := value[user]
		duration, err := time.ParseDuration("-15m")
		if err != nil {
			return UserInfo{
				Online: false,
			}
		}

		return UserInfo{
			Name:   twitchInfo.Name,
			Id:     twitchInfo.Id,
			Online: twitchInfo.Timestamp.After(time.Now().Add(duration)),
		}
	}
}
