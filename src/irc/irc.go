package irc

import (
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/patrickmn/go-cache"
)

type TwitchUserTracker struct {
	Name        string
	DisplayName string
	Id          string
}

type UserInfo struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Id          string `json:"providerId"`
	Online      bool   `json:"online"`
}

type TwitchUserDatabase struct {
	Client  *twitch.Client
	UserMap map[string]*cache.Cache
}

func CreateTwitchDatabase() *TwitchUserDatabase {
	client := twitch.NewAnonymousClient()
	log.Println("Created Twitch IRC client")
	channelList := strings.Split(os.Getenv("TWITCH_CHANNELS"), ",")
	userMap := make(map[string]*cache.Cache)
	client.Join(channelList...)
	client.OnPrivateMessage(createPrivateMsgCallback(userMap))
	client.OnConnect(func() { log.Println("Connected!") })
	log.Println("Client configuration finished")

	// Connect via goroutine because it's a blocking operation (infinite loop)
	go client.Connect()

	return &TwitchUserDatabase{
		Client:  client,
		UserMap: userMap,
	}
}

func createPrivateMsgCallback(userMap map[string]*cache.Cache) func(twitch.PrivateMessage) {
	// this mutex should be in scope of all callbacks
	mutex := &sync.Mutex{}
	return func(msg twitch.PrivateMessage) {
		go func() {
			channelName := msg.Channel
			mutex.Lock()
			value, exists := userMap[channelName]
			if !exists {
				// channel doesn't exist on the map yet, make it
				userMap[channelName] = cache.New(15*time.Minute, time.Hour)
				value = userMap[channelName]
			}
			mutex.Unlock()

			userName := msg.User.Name
			userTracker := &TwitchUserTracker{
				Name:        userName,
				DisplayName: msg.User.DisplayName,
				Id:          msg.User.ID,
			}

			value.Set(userName, userTracker, cache.DefaultExpiration)
		}()

	}
}

func (db *TwitchUserDatabase) ReadUserInfo(channel, user string) UserInfo {
	falseResponse := &UserInfo{
		Online: false,
	}
	userMap := db.UserMap
	if value, exists := userMap[channel]; !exists {
		return *falseResponse
	} else {
		twitchInfo, exists := value.Get(user)
		if !exists {
			return *falseResponse
		}
		return UserInfo{
			Name:        twitchInfo.(*TwitchUserTracker).Name,
			DisplayName: twitchInfo.(*TwitchUserTracker).DisplayName,
			Id:          twitchInfo.(*TwitchUserTracker).Id,
			Online:      true,
		}
	}
}
