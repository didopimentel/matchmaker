package entities

import "encoding/json"

type MatchmakingStatus string

const (
	MatchmakingStatus_Pending MatchmakingStatus = "pending"
	MatchmakingStatus_Found   MatchmakingStatus = "found"
	MatchmakingStatus_Expired MatchmakingStatus = "expired"
)

type MatchmakingTicket struct {
	ID            string
	PlayerId      string
	League        int64
	Table         int64
	CreatedAt     int64
	Status        MatchmakingStatus
	GameSessionId string
	Parameters    []MatchmakingTicketParameter
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
	Type     MatchmakingTicketParameterType
	Operator MatchmakingTicketParameterOperator
	Value    int64
}
