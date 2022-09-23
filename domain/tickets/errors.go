package tickets

import "errors"

var (
	TicketNotFoundErr          = errors.New("no ticket found")
	InvalidTicketParametersErr = errors.New("invalid parameters")
	InvalidPlayerParameters    = errors.New("invalid player parameters")
)
