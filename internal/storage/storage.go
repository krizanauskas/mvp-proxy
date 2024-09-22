package storage

import (
	"fmt"
	"sync"
	"time"
)

type UserHistoryStoreI interface {
	GetHistory(username string) []string
	AddHistory(username, host string, time time.Time)
}

type UserHistoryGetterI interface {
	GetHistory(username string) []string
}

type BandwidthUsedGetterI interface {
	GetBandwidthUsed(username string) int
}

type UserBandwidthStoreI interface {
	BandwidthUsedGetterI
	IncreaseBandwidthUsed(username string, usedBanwith int)
	GetAllocatedBandwidth() int
}

type userHistory struct {
	Host string
	Time time.Time
}

type Storage struct {
	UserHistoryStore   map[string][]userHistory
	UserBandwidthStore map[string]int
	AllocatedBandwidth int
	historyMutex       sync.RWMutex
	bandwidthMutex     sync.RWMutex
}

var (
	globalStorage *Storage
	once          sync.Once
)

func New(allocatedBandwidthMB int) *Storage {
	// store bandwidth in MB
	allocatedBandwidthBytes := allocatedBandwidthMB * 1000 * 1000

	return &Storage{
		UserHistoryStore:   make(map[string][]userHistory),
		UserBandwidthStore: make(map[string]int),
		AllocatedBandwidth: allocatedBandwidthBytes,
	}
}

func Initialize(allocatedBandwidth int) {
	once.Do(func() {
		globalStorage = New(allocatedBandwidth)
	})
}

type UserHistoryStore struct{}

func NewUserHistoryStore() *UserHistoryStore {
	return &UserHistoryStore{}
}

func (s *UserHistoryStore) GetHistory(username string) []string {
	globalStorage.historyMutex.RLock()
	defer globalStorage.historyMutex.RUnlock()

	userHistories, exists := globalStorage.UserHistoryStore[username]
	if !exists {
		return []string{}
	}

	var historyData []string

	for _, userHistory := range userHistories {
		historyData = append(historyData, fmt.Sprintf("URL: %s, accessed at: %s", userHistory.Host, userHistory.Time.Format("2006-01-02 15:04:05")))
	}

	return historyData
}

func (s *UserHistoryStore) AddHistory(username, host string, time time.Time) {
	userHistory := userHistory{
		Host: host,
		Time: time,
	}

	globalStorage.historyMutex.Lock()
	defer globalStorage.historyMutex.Unlock()

	globalStorage.UserHistoryStore[username] = append(globalStorage.UserHistoryStore[username], userHistory)
}

type UserBandwidthStore struct{}

func NewUserBandwidthStore() *UserBandwidthStore {
	return &UserBandwidthStore{}
}

func (s *UserBandwidthStore) GetBandwidthUsed(username string) int {
	globalStorage.bandwidthMutex.RLock()
	defer globalStorage.bandwidthMutex.RUnlock()

	bandwidthData := globalStorage.UserBandwidthStore[username]
	return bandwidthData
}

func (s *UserBandwidthStore) IncreaseBandwidthUsed(username string, usedBandwidth int) {
	globalStorage.bandwidthMutex.Lock()
	defer globalStorage.bandwidthMutex.Unlock()

	globalStorage.UserBandwidthStore[username] += usedBandwidth
}

func (s *UserBandwidthStore) GetAllocatedBandwidth() int {
	return globalStorage.AllocatedBandwidth
}
