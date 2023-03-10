package route

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lpuig/batec/stockmanagement/src/backend/logger"
	mgr "github.com/lpuig/batec/stockmanagement/src/backend/manager"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/stock"
	"net/http"
	"strconv"
)

func GetStocks(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetStocks").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	w.Header().Set("Content-Type", "application/json")

	err := mgr.GetStocks(w)
	if err != nil {
		AddError(w, logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.Response = http.StatusOK
}

func GetStockById(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetStocks").AddUser(mgr.CurrentUser.Name)
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
	err = json.NewEncoder(w).Encode(sr.Stock)
	if err != nil {
		AddError(w, logmsg, "could not marshall stock. "+err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.Response = http.StatusOK
}

func UpdateStocks(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	logmsg := logger.TimedEntry("Route").AddRequest("UpdateStocks").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	w.Header().Set("Content-Type", "application/json")
	if r.Body == nil {
		AddError(w, logmsg, "request Stocks missing", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	updatedStocks := []*stock.Stock{}
	err := json.NewDecoder(r.Body).Decode(&updatedStocks)
	if err != nil {
		AddError(w, logmsg, fmt.Sprintf("misformatted request body :%v", err.Error()), http.StatusBadRequest)
		return
	}

	err = mgr.UpdateStocks(updatedStocks)
	if err != nil {
		AddError(w, logmsg, fmt.Sprintf("error updating Stock:%v", err.Error()), http.StatusInternalServerError)
		return
	}

	logmsg.Response = http.StatusOK
}
