package internal

import (
  context "context"
  fmt     "fmt"
  http    "net/http"
  log     "github.com/sirupsen/logrus"
  _       "github.com/lib/pq"
  mux     "github.com/gorilla/mux"
  time    "time"
)

var applicationLog *log.Entry

func init() {
  applicationLog = log.WithFields(log.Fields{
    "_file": "internal/application.go",
    "_type": "system",
  })
}

type Configuration struct {
  Environment    string
  Verbosity      string
  LogFormat      string
  LogHealthcheck bool
  Server struct {
    Host string
    Port string
  }
  Healthcheck struct {
    Timeout  int
  }
}

type Version struct {
  BuildTime string
  Commit    string
  Release   string
}

type Application struct {
  Router  *mux.Router
  Config  *Configuration
  Version *Version
}

func (a *Application) Initialize() {
  a.Router = mux.NewRouter()

  a.initializeLogger()
  a.initializeRoutes()

  applicationLog.Info("application is initialized")
}

func (a *Application) initializeLogger() {
  switch a.Config.Verbosity {
  case "fatal":
    log.SetLevel(log.FatalLevel)
  case "error":
    log.SetLevel(log.ErrorLevel)
  case "warning":
    log.SetLevel(log.WarnLevel)
  case "info":
    log.SetLevel(log.InfoLevel)
  case "debug":
    log.SetLevel(log.DebugLevel)
    log.SetReportCaller(false)
  default:
    log.WithFields(
      log.Fields{"verbosity": a.Config.Verbosity},
    ).Fatal("Unkown verbosity level")
  }

  switch a.Config.Environment {
  case "development":
    log.SetFormatter(&log.TextFormatter{})
  case "integration":
    log.SetFormatter(&log.JSONFormatter{})
  case "production":
    log.SetFormatter(&log.JSONFormatter{})
  default:
    log.WithFields(
      log.Fields{"environment": a.Config.Environment},
    ).Fatal("Unkown environment type")
  }

  if a.Config.LogFormat != "unset" {
    switch a.Config.LogFormat {
    case "text":
      log.SetFormatter(&log.TextFormatter{})
    case "json":
      log.SetFormatter(&log.JSONFormatter{})
    default:
      log.WithFields(
        log.Fields{"log_format": a.Config.LogFormat},
      ).Fatal("Unkown log format")
    }
  }
}

func (a *Application) initializeRoutes() {
  // application strategy-related routes
  a.Router.HandleFunc("/strategies", a.getStrategies).Methods("GET")
  a.Router.HandleFunc("/strategy", a.createStrategy).Methods("POST")
  // application probes routes
  a.Router.HandleFunc("/health", a.isHealthy).Methods("GET")
  a.Router.HandleFunc("/ready", a.isReady).Methods("GET")
}

func (a *Application) Run(ctx context.Context) {
  server_address :=
    fmt.Sprintf("%s:%s", a.Config.Server.Host, a.Config.Server.Port)

  server_message :=
    fmt.Sprintf(
  `

START INFOS
-----------
Listening needys-api-strategy on %s:%s...

BUILD INFOS
-----------
time: %s
release: %s
commit: %s

`,
      a.Config.Server.Host,
      a.Config.Server.Port,
      a.Version.BuildTime,
      a.Version.Release,
      a.Version.Commit,
    )

  httpServer := &http.Server{
		Addr:    server_address,
		Handler: a.Router,
	}

  go func() {
    // we keep this log on standard format
    log.Info(server_message)
    applicationLog.Fatal(httpServer.ListenAndServe())
  }()

  <-ctx.Done()
  applicationLog.Info("server stopped")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

  var err error

	if err = httpServer.Shutdown(ctxShutDown); err != nil {
    applicationLog.WithFields(log.Fields{
      "error": err,
    }).Fatal("server shutdown failed")
	}

  applicationLog.Info("server exited properly")

	if err == http.ErrServerClosed {
		err = nil
	}

	return
}
