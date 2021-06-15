package core

type Strategy struct {
  ID          string `json:"id"`
  Description string `json:"description"`
  NeedID      string `json:"needId"`
}

var Strategies []Strategy
