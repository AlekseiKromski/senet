package lb

import (
	"fmt"
	"senet/processor/storage"
	"senet/processor/storage/models"
	"sync"
	"time"
)

type LoadBalancer struct {
	db    storage.Storage
	mutex *sync.Mutex

	users []*models.User
}

func NewLoadBalancer(storage storage.Storage) *LoadBalancer {
	lb := &LoadBalancer{
		db:    storage,
		mutex: &sync.Mutex{},
	}
	go lb.backgroundCleanup()

	return lb
}

func (lb *LoadBalancer) GetUsers() ([]*models.User, error) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	if len(lb.users) != 0 {
		fmt.Printf("from BUF")
		return lb.users, nil
	}

	users, err := lb.db.GetUsers()

	lb.users = users

	return lb.users, err
}

func (lb *LoadBalancer) backgroundCleanup() {
	for {
		lb.mutex.Lock()
		if len(lb.users) != 0 {
			lb.users = []*models.User{}
			fmt.Printf("lb: cleaned")
		}
		lb.mutex.Unlock()

		<-time.After(1 * time.Minute) // or time.sleep(time.Second)
	}
}
