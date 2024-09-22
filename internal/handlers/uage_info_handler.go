package handlers

import (
	"fmt"
	"net/http"

	"krizanauskas.github.com/mvp-proxy/internal/services"
	"krizanauskas.github.com/mvp-proxy/internal/storage"
)

type UsageInfoHandler struct {
	bandwidthUsedGetter storage.UserBandwidthStoreI
}

func NewUsageInfoHandler(bandwidthUsedGetter storage.UserBandwidthStoreI) UsageInfoHandler {
	return UsageInfoHandler{
		bandwidthUsedGetter,
	}
}

func (h UsageInfoHandler) Handle(w http.ResponseWriter, _ *http.Request) {
	bandwidthAllocated := h.bandwidthUsedGetter.GetAllocatedBandwidth()
	bandwidthUsed := h.bandwidthUsedGetter.GetBandwidthUsed(services.AuthUser)

	bandwidthUsedKB := float64(bandwidthUsed) / 1000

	fmt.Fprintf(w, "Used bandwidth: %2.f KB of %d KB", bandwidthUsedKB, (bandwidthAllocated / 1000))
}
