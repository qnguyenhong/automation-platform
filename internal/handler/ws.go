package handler

import (
	"net/http"

	"github.com/qnguyenhong/automation-platform/pkg/httputil"
)

func (d *Deps) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	if err := d.Hub.Upgrade(w, r); err != nil {
		httputil.Error(w, http.StatusInternalServerError, "failed to upgrade websocket")
		return
	}
}
