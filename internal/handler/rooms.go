package handler

import (
	"fmt"
	"net/http"

	"github.com/pangaunn/harness-engineering-workshop/internal/service/room"
)

func (h *handler) listRooms(w http.ResponseWriter, r *http.Request) {
	types, err := h.Rooms.Types(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"check_in_from": fmt.Sprintf("%02d:00", room.CheckInHour),
		"room_types":    types,
	})
}
