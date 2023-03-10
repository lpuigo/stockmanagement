package main

import (
	"crypto/tls"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/lpuig/batec/stockmanagement/src/backend/config"
	"github.com/lpuig/batec/stockmanagement/src/backend/logger"
	"github.com/lpuig/batec/stockmanagement/src/backend/manager"
	"github.com/lpuig/batec/stockmanagement/src/backend/route"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Conf struct {
	manager.ManagerConfig

	AssetsDir  string
	AssetsRoot string
	RootDir    string

	LogFile     string
	ServiceHost string
	ServicePort string

	LaunchWebBrowser    bool
	InProduction        bool
	RedirectHTTPToHTTPS bool
}

const (
	ConfigFile = `./config.json`

	UsersDir       = `C:\Users\Laurent\Golang\src\github.com\lpuig\batec\stockmanagement\Ressources\Users`
	ActorsDir      = `C:\Users\Laurent\Golang\src\github.com\lpuig\batec\stockmanagement\Ressources\Actors`
	ArticlesDir    = `C:\Users\Laurent\Golang\src\github.com\lpuig\batec\stockmanagement\Ressources\Articles`
	MovementsDir   = `C:\Users\Laurent\Golang\src\github.com\lpuig\batec\stockmanagement\Ressources\Movements`
	WorksitesDir   = `C:\Users\Laurent\Golang\src\github.com\lpuig\batec\stockmanagement\Ressources\Worksites`
	StocksDir      = `C:\Users\Laurent\Golang\src\github.com\lpuig\batec\stockmanagement\Ressources\Stocks`
	SaveArchiveDir = `C:\Users\Laurent\Golang\src\github.com\lpuig\batec\stockmanagement\SaveArchive`
	SessionKey     = "SECRET_KEY"

	AssetsDir  = `./Assets`
	AssetsRoot = `/Assets/`
	RootDir    = `./Dist`

	LogFile     = `./server.log`
	ServiceHost = "vpsXXXXX.ovh.net"
	ServicePort = ":8088"

	LaunchWebBrowser    = true
	InProduction        = false
	RedirectHTTPToHTTPS = false
)

func LaunchPageInBrowser(c *Conf) error {
	if !c.LaunchWebBrowser {
		return nil
	}
	cmd := exec.Command("cmd", "/c", "start", "http://localhost"+c.ServicePort)
	return cmd.Start()
}

// createRouter sets a router with all functional route  using given configuration
func createRouter(mgr *manager.Manager, conf *Conf) http.Handler {
	withManager := func(hf route.MgrHandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			hf(mgr.Clone(), w, r)
		}
	}

	withUserManager := func(request string, hf route.MgrHandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			m := mgr.Clone()
			if !m.CheckSessionUser(r) {
				logmsg := logger.Entry("Route").AddRequest(request)
				route.AddError(w, logmsg, "User not connected or not authorized", http.StatusUnauthorized)
				logmsg.Log()
				return
			}
			hf(m, w, r)
		}
	}

	router := mux.NewRouter()
	// attach pprof route from defaultServeMux
	router.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)
	// session management
	router.HandleFunc("/api/login", withManager(route.GetUser)).Methods("GET")
	router.HandleFunc("/api/login", withUserManager("Logout", route.Logout)).Methods("DELETE")
	router.HandleFunc("/api/login", withManager(route.Login)).Methods("POST")

	// Archives methods
	router.HandleFunc("/api/{recordtype}/archive", withUserManager("GetRecordsArchive", route.GetRecordsArchive)).Methods("GET")
	router.HandleFunc("/api/archive", withUserManager("GetSaveArchive", route.GetSaveArchive)).Methods("GET")

	// Users methods
	router.HandleFunc("/api/users", withUserManager("GetUsers", route.GetUsers)).Methods("GET")
	router.HandleFunc("/api/users", withUserManager("UpdateUsers", route.UpdateUsers)).Methods("PUT")

	// Actors methods
	router.HandleFunc("/api/actors", withUserManager("GetActors", route.GetActors)).Methods("GET")
	router.HandleFunc("/api/actors", withUserManager("UpdateActors", route.UpdateActors)).Methods("PUT")

	// Articles methods
	router.HandleFunc("/api/articles", withUserManager("GetArticles", route.GetArticles)).Methods("GET")
	router.HandleFunc("/api/articles", withUserManager("UpdateArticles", route.UpdateArticles)).Methods("PUT")
	router.HandleFunc("/api/articles/export", withUserManager("GetArticlesExport", route.GetArticlesExport)).Methods("GET")
	router.HandleFunc("/api/articles/import", withUserManager("PostArticlesImport", route.PostArticlesImport)).Methods("POST")

	// Movements methods
	router.HandleFunc("/api/movements", withUserManager("GetMovements", route.GetMovements)).Methods("GET")
	router.HandleFunc("/api/movements/stock/{sid:[0-9]+}", withUserManager("GetMovementsForStock", route.GetMovementsForStock)).Methods("GET")
	router.HandleFunc("/api/movements", withUserManager("UpdateMovements", route.UpdateMovements)).Methods("PUT")

	// Worksites methods
	router.HandleFunc("/api/worksites", withUserManager("GetWorksites", route.GetWorksites)).Methods("GET")
	router.HandleFunc("/api/worksites", withUserManager("UpdateWorksites", route.UpdateWorksites)).Methods("PUT")

	// Stocks methods
	router.HandleFunc("/api/stocks", withUserManager("GetStocks", route.GetStocks)).Methods("GET")
	router.HandleFunc("/api/stocks/{sid:[0-9]+}", withUserManager("GetStockById", route.GetStockById)).Methods("GET")
	router.HandleFunc("/api/stocks", withUserManager("UpdateStocks", route.UpdateStocks)).Methods("PUT")

	// Administration methods
	router.HandleFunc("/api/admin/reload", withUserManager("ReloadPersister", route.ReloadPersister)).Methods("GET")

	// Static Files serving
	router.PathPrefix(conf.AssetsRoot).Handler(http.StripPrefix(conf.AssetsRoot, http.FileServer(http.Dir(conf.AssetsDir))))
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(conf.RootDir)))

	gzipedrouter := handlers.CompressHandler(router)

	return gzipedrouter
}

func makeServerFromMux(mux http.Handler) *http.Server {
	// set timeouts so that a slow or malicious client doesn't
	// hold resources forever
	return &http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
}

func catchArchiveRequest(mgr *manager.Manager) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.Signal(10)) // notify SIGUSR1

	go func() {
		for {
			<-c // wait for SIGUSR1 to happen
			logmsg := logger.TimedEntry("Signal")
			err := mgr.SaveArchive()
			if err == nil {
				logmsg.LogInfo("SaveArchive")
			} else {
				logmsg.LogError(err)
			}
			logmsg.Log()
		}
	}()
}

func main() {
	conf := &Conf{
		ManagerConfig: manager.ManagerConfig{
			UsersDir:       UsersDir,
			ActorsDir:      ActorsDir,
			ArticlesDir:    ArticlesDir,
			MovementsDir:   MovementsDir,
			WorksitesDir:   WorksitesDir,
			StocksDir:      StocksDir,
			SaveArchiveDir: SaveArchiveDir,
			SessionKey:     SessionKey,
		},

		AssetsDir:  AssetsDir,
		AssetsRoot: AssetsRoot,
		RootDir:    RootDir,

		LogFile:     LogFile,
		ServiceHost: ServiceHost,
		ServicePort: ServicePort,

		LaunchWebBrowser:    LaunchWebBrowser,
		InProduction:        InProduction,
		RedirectHTTPToHTTPS: RedirectHTTPToHTTPS,
	}

	if err := config.SetFromFile(ConfigFile, conf); err != nil {
		logger.Entry("Server").Fatal(err)
	}

	logFile, err := logger.StartLog(conf.LogFile)
	if err != nil {
		log.SetOutput(os.Stderr)
		log.Fatalf("could not init logger: %s\n", err.Error())
	}
	defer logFile.Close()
	logger.Entry("Server").LogInfo("============================= SERVER STARTING ==================================")

	mgr, err := manager.NewManager(conf.ManagerConfig)
	if err != nil {
		logger.Entry("Server").Fatal(err)
	}
	router := createRouter(mgr, conf)

	catchArchiveRequest(mgr)

	wg := sync.WaitGroup{}

	if conf.InProduction {
		logger.Entry("Server").LogInfo("Init Production setup")

		dataDir := "."
		certManager := &autocert.Manager{
			Prompt: autocert.AcceptTOS,
			//HostPolicy: hostPolicy,
			Cache: autocert.DirCache(dataDir),
		}

		httpsSrv := makeServerFromMux(router)
		httpsSrv.Addr = ":443"
		httpsSrv.TLSConfig = &tls.Config{GetCertificate: certManager.GetCertificate}

		wg.Add(2)
		go func() {
			logger.Entry("Server").LogInfo("listening HTTPS on " + httpsSrv.Addr)
			logger.Entry("Server").LogError(httpsSrv.ListenAndServeTLS("", ""))
			wg.Done()
		}()
		go func() {
			logger.Entry("Server").LogInfo("listening HTTP on :80 for certification handshake")
			logger.Entry("Server").LogError(http.ListenAndServe(":80", certManager.HTTPHandler(nil)))
			wg.Done()
		}()
	} else {
		logger.Entry("Server").LogInfo("Init Non Production setup")
		httpSrv := makeServerFromMux(router)
		httpSrv.Addr = conf.ServicePort
		wg.Add(1)
		go func() {
			logger.Entry("Server").LogInfo("listening HTTP on " + httpSrv.Addr)
			logger.Entry("Server").LogError(httpSrv.ListenAndServe())
			wg.Done()
		}()
	}

	logger.Entry("Server").LogInfo("============================== SERVER READY ====================================")
	LaunchPageInBrowser(conf)
	wg.Wait()
}
