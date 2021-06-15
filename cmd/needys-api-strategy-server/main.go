package main

import(
  config  "github.com/gpenaud/needys-api-strategy/internal/config"
  context "context"
  core    "github.com/gpenaud/needys-api-strategy/internal/core"
  cmdline "github.com/galdor/go-cmdline"
  fmt     "fmt"
  http    "net/http"
  log     "github.com/gpenaud/needys-api-strategy/pkg/log"
  mux     "github.com/gorilla/mux"
  os      "os"
  signal  "os/signal"
  syscall "syscall"
  time    "time"
  version "github.com/gpenaud/needys-api-strategy/build/version"
)

func health(w http.ResponseWriter, _ *http.Request) {
  w.WriteHeader(http.StatusOK)
  w.Header().Set("Content-Type", "application/json")
  w.Write([]byte("{\"state\": \"healthy\"}"))
}

func ready(w http.ResponseWriter, _ *http.Request) {
  w.WriteHeader(http.StatusOK)
  w.Header().Set("Content-Type", "application/json")
  w.Write([]byte("{\"state\": \"ready\"}"))
}

func main() {
  cmdline := cmdline.New()

  cmdline.AddOption("", "server.host", "HOST", "host of application")
  cmdline.SetOptionDefault("server.host", "localhost")
  cmdline.AddOption("", "server.port", "PORT", "port of application")
  cmdline.SetOptionDefault("server.port", "8011")

  cmdline.AddFlag("v", "verbose", "log more information")
  cmdline.Parse(os.Args)

  config.HttpServer.Host = cmdline.OptionValue("server.host")
  config.HttpServer.Port = cmdline.OptionValue("server.port")

  // mocked data
  core.Strategies = append(core.Strategies,
    core.Strategy{ID: "1", Description: "Aller me promener dans la nature", NeedID: "3"},
    core.Strategy{ID: "2", Description: "Faire une séance de cohérence cardiaque", NeedID: "3"},
    core.Strategy{ID: "3", Description: "Une bonne b...", NeedID: "3"},
  )

  router := mux.NewRouter()

  // probes
  router.HandleFunc("/health", health).Methods("GET")
  router.HandleFunc("/ready", ready).Methods("GET")

  // logic
  router.HandleFunc("/strategy", core.GetStrategies).Methods("GET")
  router.HandleFunc("/strategy", core.CreateStrategy).Methods("POST")

  interrupt := make(chan os.Signal, 1)
  signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

  fmt.Printf(
    "Starting needys-api-strategy on %s:%s...\n > build time: %s\n > release: %s\n > commit: %s\n",
    config.HttpServer.Host,
    config.HttpServer.Port,
    version.BuildTime,
    version.Release,
    version.Commit,
  )

  server_address := fmt.Sprintf("%s:%s", config.HttpServer.Host, config.HttpServer.Port)
  server := &http.Server{
    Addr:           server_address,
    Handler:        router,
    ReadTimeout:    10 * time.Second,
    WriteTimeout:   10 * time.Second,
    MaxHeaderBytes: 1 << 20,
  }

  go func() {
    log.ErrorLogger.Fatalln(server.ListenAndServe())
  }()

  killSignal := <-interrupt

  switch killSignal {
  case os.Interrupt:
    log.WarningLogger.Println("Received SIGINT...")
  case syscall.SIGTERM:
    log.WarningLogger.Println("Received SIGTERM...")
  }

  log.InfoLogger.Println("The service is shutting down...")
  server.Shutdown(context.Background())
  log.InfoLogger.Println("...Shutting done !")
}
