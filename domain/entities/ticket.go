package entities

import "encoding/json"

type MatchmakingTicket struct {
	ID         string
	PlayerID   string
	League     int64
	Table      int64
	Parameters []MatchmakingTicketParameter
}

func (i MatchmakingTicket) MarshalBinary() (data []byte, err error) {
	bytes, err := json.Marshal(i)
	return bytes, err
}

type MatchmakingTicketParameterType string

const (
	MatchmakingTicketParameterType_League MatchmakingTicketParameterType = "league"
	MatchmakingTicketParameterType_Table  MatchmakingTicketParameterType = "table"
)

type MatchmakingTicketParameterOperator string

const (
	MatchmakingTicketParameterOperator_Equal       MatchmakingTicketParameterOperator = "="
	MatchmakingTicketParameterOperator_NotEqual    MatchmakingTicketParameterOperator = "<>"
	MatchmakingTicketParameterOperator_SmallerThan MatchmakingTicketParameterOperator = "<"
	MatchmakingTicketParameterOperator_GreaterThan MatchmakingTicketParameterOperator = ">"
)

type MatchmakingTicketParameter struct {
	Type     MatchmakingTicketParameterType     `json:"type"`
	Operator MatchmakingTicketParameterOperator `json:"operator"`
	Value    int64                              `json:"value"`
}
