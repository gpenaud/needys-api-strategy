package core

import(
  http    "net/http"
  json    "encoding/json"
  _       "github.com/gpenaud/needys-api-strategy/pkg/log"
  strconv "strconv"
)

func GetStrategies(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(Strategies)
}

func CreateStrategy(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")

  var strategy Strategy
  json.NewDecoder(r.Body).Decode(&strategy)

  strategy.ID = strconv.Itoa(len(Strategies) + 1)

  Strategies = append(Strategies, strategy)
  json.NewEncoder(w).Encode(strategy)
}
