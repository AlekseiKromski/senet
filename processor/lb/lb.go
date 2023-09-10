package lb

import (
	"encoding/json"
	"fmt"
	"github.com/AlekseiKromski/at-socket-server/core"
	"log"
	"senet/processor/storage"
	"senet/processor/storage/creators/chat"
	"senet/processor/storage/models"
	"sync"
	"time"
)

type LoadBalancer struct {
	db    storage.Storage
	mutex *sync.Mutex

	clients core.Clients
	users   map[string][]models.Chat
}

func NewLoadBalancer(storage storage.Storage, clients core.Clients) *LoadBalancer {
	lb := &LoadBalancer{
		db:      storage,
		mutex:   &sync.Mutex{},
		users:   make(map[string][]models.Chat),
		clients: clients,
	}
	go lb.backgroundCleanup()

	return lb
}

func (lb *LoadBalancer) CreateChat(creator chat.ChatCreator) (models.Chat, error) {
	chat, err := lb.db.CreateChat(creator)

	if err != nil {
		return models.Chat{}, fmt.Errorf("db error: %v", err)
	}

	lb.mutex.Lock()
	for _, userID := range creator.GetMembers() {
		client := lb.clients[userID]
		if client == nil {
			continue
		}

		payload, err := json.Marshal(&chat)
		if err != nil {
			log.Printf("cannot prepare payload for sending to client: %v", err)
		}

		if err := client.Send(string(payload), CREATED_NEW_CHAT); err != nil {
			log.Printf("cannot send data to client: %v", err)
		}
		lb.users[userID] = append(lb.users[userID], chat)
	}
	lb.mutex.Unlock()

	c := chat
	c.Messages = nil

	return c, nil
}

/*
GetChats - If we have any chats is memory, return then.
If not, take it from db, save in memory and return
*/
func (lb *LoadBalancer) GetChats(userID string) ([]models.Chat, error) {
	lb.mutex.Lock()
	memSavedChats := lb.users[userID]

	if len(memSavedChats) != 0 {
		lb.mutex.Unlock()
		return memSavedChats, nil
	}

	lb.mutex.Unlock()

	chats, err := lb.db.GetChats(userID)
	if err != nil {
		return nil, fmt.Errorf("cannot get chats from databse: %v", err)
	}

	lb.mutex.Lock()
	lb.users[userID] = chats
	defer lb.mutex.Unlock()

	return lb.users[userID], nil
}

func (lb *LoadBalancer) backgroundCleanup() {
	for {
		lb.mutex.Lock()
		if len(lb.users) != 0 {
			lb.users = make(map[string][]models.Chat)
			fmt.Printf("lb: users cleaned")
		}
		lb.mutex.Unlock()

		<-time.After(4 * time.Minute) // or time.sleep(time.Second)
	}
}
