package route

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lpuig/batec/stockmanagement/src/backend/logger"
	mgr "github.com/lpuig/batec/stockmanagement/src/backend/manager"
	"github.com/lpuig/batec/stockmanagement/src/backend/persist"
	"net/http"
)

func GetRecordsArchive(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetRecordsArchive").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	vars := mux.Vars(r)
	recordType := vars["recordtype"]
	var container persist.ArchivableRecordContainer
	switch recordType {
	case "actors":
		container = mgr.Actors
	case "users":
		container = mgr.Users
	default:
		AddError(w, logmsg, "unsupported archive type '"+recordType+"'", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", persist.ArchiveName(container)))
	w.Header().Set("Content-Type", "application/zip")

	err := persist.CreateRecordsArchive(w, container)
	if err != nil {
		AddError(w, logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.AddInfoResponse(fmt.Sprintf("%s archive produced", recordType), http.StatusOK)
}

func GetSaveArchive(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetSaveArchive").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	err := mgr.SaveArchive()
	if err != nil {
		logmsg.LogError(err)
	}
}
