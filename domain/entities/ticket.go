package entities

import "encoding/json"

type MatchmakingStatus string

const (
	MatchmakingStatus_Pending MatchmakingStatus = "pending"
	MatchmakingStatus_Found   MatchmakingStatus = "found"
	MatchmakingStatus_Expired MatchmakingStatus = "expired"
)

type MatchmakingTicket struct {
	ID              string
	PlayerId        string
	CreatedAt       int64
	Status          MatchmakingStatus
	GameSessionId   string
	MatchParameters []MatchmakingTicketParameter
}

func (i MatchmakingTicket) MarshalBinary() (data []byte, err error) {
	bytes, err := json.Marshal(i)
	return bytes, err
}

// MatchmakingTicketParameterType allows us to define which parameters we will accept later on
type MatchmakingTicketParameterType string

type MatchmakingTicketParameterOperator string

const (
	MatchmakingTicketParameterOperator_Equal       MatchmakingTicketParameterOperator = "="
	MatchmakingTicketParameterOperator_NotEqual    MatchmakingTicketParameterOperator = "<>"
	MatchmakingTicketParameterOperator_SmallerThan MatchmakingTicketParameterOperator = "<"
	MatchmakingTicketParameterOperator_GreaterThan MatchmakingTicketParameterOperator = ">"
)

type MatchmakingTicketParameter struct {
	Type     MatchmakingTicketParameterType
	Operator MatchmakingTicketParameterOperator
	Value    float64
}
