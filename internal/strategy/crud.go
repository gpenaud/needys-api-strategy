package strategy

import(
  http    "net/http"
  // json    "encoding/json"
)

type Strategy struct {
  ID          int    `json:"id"`
  Description string `json:"description"`
  NeedID      int    `json:"needId"`
}

var Strategies []Strategy

func init() {
  strategy := Strategy{ID: 1, Description: "Aller manger au restaurant", NeedID: 2}
  Strategies = append(Strategies, strategy)
}

func (s *Strategy) GetStrategies(w http.ResponseWriter, r *http.Request) []Strategy {
  return Strategies
}

func (s *Strategy) CreateStrategy(strategy Strategy) Strategy {
  strategy.ID = len(Strategies) + 1

  Strategies = append(Strategies, strategy)
  return strategy
}
