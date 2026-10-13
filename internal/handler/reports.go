package handler

import (
	"fmt"
	"net/http"

	"github.com/pangaunn/harness-engineering-workshop/internal/model"
)

func (h *handler) revenueReport(w http.ResponseWriter, r *http.Request) {
	date, err := model.ParseDate(r.URL.Query().Get("date"))
	if err != nil {
		writeError(w, fmt.Errorf("%w: %v", errBadRequest, err))
		return
	}
	rep, err := h.Reports.Get(r.Context(), date)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, rep)
}
