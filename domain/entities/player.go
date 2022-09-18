package entities

type Player struct {
	ID     string `json:"id"`
	League int64  `json:"league"`
	Table  int64  `json:"table"`
}
