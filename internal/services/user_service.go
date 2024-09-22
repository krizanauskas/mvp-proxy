package services

import (
	"time"

	"krizanauskas.github.com/mvp-proxy/internal/storage"
)

type UserBandwidthControllerI interface {
	UpdateBandwidthUsed(username string, used int) error
	GetAvailableBandwidth(username string) int
}

type UserServiceI interface {
	UserBandwidthControllerI
	GetBandwidthUsed(username string) int
	AddToHistory(user, host string, time time.Time) error
	GetHistory(user string) ([]string, error)
}

type UserService struct {
	userHistoryStore   storage.UserHistoryStoreI
	userBandwidthStore storage.UserBandwidthStoreI
}

func NewUserService(userHistoryStore storage.UserHistoryStoreI, userBandwidhtStore storage.UserBandwidthStoreI) UserService {
	return UserService{
		userHistoryStore,
		userBandwidhtStore,
	}
}

func (s UserService) GetHistory(username string) ([]string, error) {
	return []string{"https://google.lt", "https://delfi.lt "}, nil
}

func (s UserService) GetBandwidthUsed(username string) int {
	bandwitdthUSed := s.userBandwidthStore.GetBandwidthUsed(username)

	return bandwitdthUSed
}

func (s UserService) AddToHistory(user string, host string, time time.Time) error {
	s.userHistoryStore.AddHistory(user, host, time)

	return nil
}

func (s UserService) UpdateBandwidthUsed(username string, bytesUsed int) error {
	s.userBandwidthStore.IncreaseBandwidthUsed(username, bytesUsed)

	return nil
}

func (s UserService) GetAvailableBandwidth(username string) int {
	return s.userBandwidthStore.GetAllocatedBandwidth() - s.userBandwidthStore.GetBandwidthUsed(username)
}
