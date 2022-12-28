package route

import (
	"encoding/json"
	"fmt"
	"github.com/lpuig/batec/stockmanagement/src/backend/logger"
	mgr "github.com/lpuig/batec/stockmanagement/src/backend/manager"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/actor"
	"net/http"
)

func GetActors(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetActors").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	w.Header().Set("Content-Type", "application/json")

	err := mgr.GetActors(w)
	if err != nil {
		AddError(w, logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.Response = http.StatusOK
}

func UpdateActors(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	logmsg := logger.TimedEntry("Route").AddRequest("UpdateActors").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	w.Header().Set("Content-Type", "application/json")
	if r.Body == nil {
		AddError(w, logmsg, "request Actors missing", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	updatedActors := []*actor.Actor{}
	err := json.NewDecoder(r.Body).Decode(&updatedActors)
	if err != nil {
		AddError(w, logmsg, fmt.Sprintf("misformatted request body :%v", err.Error()), http.StatusBadRequest)
		return
	}

	err = mgr.UpdateActors(updatedActors)
	if err != nil {
		AddError(w, logmsg, fmt.Sprintf("error updating actors:%v", err.Error()), http.StatusInternalServerError)
		return
	}

	logmsg.Response = http.StatusOK
}
