package route

import (
	"encoding/json"
	"fmt"
	"github.com/lpuig/batec/stockmanagement/src/backend/logger"
	mgr "github.com/lpuig/batec/stockmanagement/src/backend/manager"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/article"
	"net/http"
	"strings"
)

func GetArticles(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetArticles").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	w.Header().Set("Content-Type", "application/json")

	err := mgr.GetArticles(w)
	if err != nil {
		AddError(w, logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.Response = http.StatusOK
}

func UpdateArticles(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	logmsg := logger.TimedEntry("Route").AddRequest("UpdateArticles").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	w.Header().Set("Content-Type", "application/json")
	if r.Body == nil {
		AddError(w, logmsg, "request Articles missing", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	updatedArticles := []*article.Article{}
	err := json.NewDecoder(r.Body).Decode(&updatedArticles)
	if err != nil {
		AddError(w, logmsg, fmt.Sprintf("misformatted request body :%v", err.Error()), http.StatusBadRequest)
		return
	}

	err = mgr.UpdateArticles(updatedArticles)
	if err != nil {
		AddError(w, logmsg, fmt.Sprintf("error updating article:%v", err.Error()), http.StatusInternalServerError)
		return
	}

	logmsg.Response = http.StatusOK
}

// GetArticlesExport
func GetArticlesExport(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("GetArticlesExport").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", mgr.Articles.ExportName()))
	w.Header().Set("Content-Type", "vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	err := mgr.Articles.XLSExport(w)
	if err != nil {
		AddError(w, logmsg, "could not generate articles XLSx export file. "+err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg.AddInfoResponse(fmt.Sprintf("Export XLS produced for %d Articles", mgr.Articles.NbArticles()), http.StatusOK)
}

// PostArticlesImport
func PostArticlesImport(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := logger.TimedEntry("Route").AddRequest("PostArticlesImport").AddUser(mgr.CurrentUser.Name)
	defer logmsg.Log()

	w.Header().Set("Content-Type", "application/json")

	// Parse our multipart form, 30 << 20 specifies a maximum
	// upload of 30 MB files.
	if r.ParseMultipartForm(30<<20) != nil {
		AddError(w, logmsg, "file info missing", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		AddError(w, logmsg, "error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if !strings.HasSuffix(strings.ToUpper(handler.Filename), ".XLSX") {
		AddError(w, logmsg, "uploaded file is not a XLSx file", http.StatusBadRequest)
		return
	}

	loadedArticles, err := article.LoadArticlesFromXlsx(file)
	if err != nil {
		AddError(w, logmsg, fmt.Sprintf("error while reading articles file: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(loadedArticles)

	logmsg.AddInfoResponse(fmt.Sprintf("%d articles imported", len(loadedArticles)), http.StatusOK)
}
