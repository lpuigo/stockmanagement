package route

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lpuig/batec/stockmanagement/src/backend/logger"
	mgr "github.com/lpuig/batec/stockmanagement/src/backend/manager"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/movement"
	"net/http"
	"strconv"
)

func GetMovements(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetMovements").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	w.Header().Set("Content-Type", "application/json")

	err := mgr.GetMovements(w)
	if err != nil {
		AddError(w, logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.Response = http.StatusOK
}

func GetMovementsForStock(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetMovements").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	sid, err := strconv.Atoi(vars["sid"])
	if err != nil {
		AddError(w, logmsg, "mis-formatted stock id '"+vars["sid"]+"'", http.StatusBadRequest)
		return
	}

	sr := mgr.Stocks.GetById(sid)
	if sr == nil {
		AddError(w, logmsg, fmt.Sprintf("stock with id %d does not exist", sid), http.StatusNotFound)
		return
	}

	err = mgr.GetMovementsForStockId(w, sid)
	if err != nil {
		AddError(w, logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.Response = http.StatusOK
}

func UpdateMovements(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	logmsg := logger.TimedEntry("Route").AddRequest("UpdateMovements").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	w.Header().Set("Content-Type", "application/json")
	if r.Body == nil {
		AddError(w, logmsg, "request movements missing", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	updatedMovements := []*movement.Movement{}
	err := json.NewDecoder(r.Body).Decode(&updatedMovements)
	if err != nil {
		AddError(w, logmsg, fmt.Sprintf("misformatted request body :%v", err.Error()), http.StatusBadRequest)
		return
	}

	err = mgr.UpdateMovements(updatedMovements)
	if err != nil {
		AddError(w, logmsg, fmt.Sprintf("error updating movement:%v", err.Error()), http.StatusInternalServerError)
		return
	}

	logmsg.Response = http.StatusOK
}
