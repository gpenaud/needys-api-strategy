package internal

import (
  fmt      "fmt"
  http     "net/http"
  json     "encoding/json"
  log      "github.com/sirupsen/logrus"
  strategy "github.com/gpenaud/needys-api-strategy/internal/strategy"
)

var handlerLog *log.Entry

func init() {
  handlerLog = log.WithFields(log.Fields{
    "_file": "internal/handler.go",
    "_type": "user",
  })
}

// -------------------------------------------------------------------------- //
// Common functions for handlers

func respondHTTPCodeOnly(w http.ResponseWriter, code int) {
  w.WriteHeader(code)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
  handlerLog.Error(message)
  respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
  response, _ := json.Marshal(payload)
  handlerLog.Debug(fmt.Sprintf("JSON response: %s", response))

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(code)
  w.Write(response)
}

// -------------------------------------------------------------------------- //
// Probe handlers

func (a *Application) isHealthy(w http.ResponseWriter, _ *http.Request) {
  payload := map[string]bool{
    "healthy": true,
  }

  if a.Config.LogHealthcheck {
    handlerLog.Debug("sent a GET request on /healthy")
    respondWithJSON(w, http.StatusOK, payload)
  } else {
    respondHTTPCodeOnly(w, http.StatusOK)
  }
}

func (a *Application) isReady(w http.ResponseWriter, _ *http.Request) {
  if a.Config.LogHealthcheck {
    handlerLog.Debug("sent a GET request on /ready")

    payload := map[string]interface{}{
      "ready": true,
    }

    respondWithJSON(w, http.StatusOK, payload)
  } else {
    respondHTTPCodeOnly(w, http.StatusOK)
  }
}

// -------------------------------------------------------------------------- //
// Strategy handlers

func (a *Application) getStrategies(w http.ResponseWriter, r *http.Request) {
  handlerLog.Info("sent a GET query on /strategies")
  var strategy strategy.Strategy

  strategies := strategy.GetStrategies(w, r)
  respondWithJSON(w, http.StatusOK, strategies)
}

func (a *Application) createStrategy(w http.ResponseWriter, r *http.Request) {
  handlerLog.Info("sent a POST query on /strategy to create a new strategy")
  var strategy strategy.Strategy

  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&strategy)

  if err != nil {
    respondWithError(w, http.StatusBadRequest, "The payload is invalid")
    return
  }

  defer r.Body.Close()

  respondWithJSON(w, http.StatusCreated, strategy.CreateStrategy(strategy))
}
