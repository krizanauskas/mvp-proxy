package handlers

import (
	"fmt"
	"net/http"

	"krizanauskas.github.com/mvp-proxy/internal/services"
	"krizanauskas.github.com/mvp-proxy/internal/storage"
)

type HistoryHandler struct {
	userHistoryGetter storage.UserHistoryGetterI
}

func NewHistoryHandler(userHistoryGetter storage.UserHistoryGetterI) HistoryHandler {
	return HistoryHandler{
		userHistoryGetter,
	}
}

func (h HistoryHandler) Handle(w http.ResponseWriter, _ *http.Request) {
	historyData := h.userHistoryGetter.GetHistory(services.AuthUser)

	for _, history := range historyData {
		fmt.Fprintf(w, "%s \n", history)
	}
}
