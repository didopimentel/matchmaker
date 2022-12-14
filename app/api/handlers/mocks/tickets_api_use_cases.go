// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"context"
	"github.com/didopimentel/matchmaker/app/api/handlers"
	"github.com/didopimentel/matchmaker/domain/tickets"
	"sync"
)

// Ensure, that TicketsAPIUseCasesMock does implement handlers.TicketsAPIUseCases.
// If this is not the case, regenerate this file with moq.
var _ handlers.TicketsAPIUseCases = &TicketsAPIUseCasesMock{}

// TicketsAPIUseCasesMock is a mock implementation of handlers.TicketsAPIUseCases.
//
// 	func TestSomethingThatUsesTicketsAPIUseCases(t *testing.T) {
//
// 		// make and configure a mocked handlers.TicketsAPIUseCases
// 		mockedTicketsAPIUseCases := &TicketsAPIUseCasesMock{
// 			CreateTicketFunc: func(ctx context.Context, input tickets.CreateTicketInput) (tickets.CreateTicketOutput, error) {
// 				panic("mock out the CreateTicket method")
// 			},
// 			GetTicketFunc: func(ctx context.Context, input tickets.GetTicketInput) (tickets.GetTicketOutput, error) {
// 				panic("mock out the GetTicket method")
// 			},
// 		}
//
// 		// use mockedTicketsAPIUseCases in code that requires handlers.TicketsAPIUseCases
// 		// and then make assertions.
//
// 	}
type TicketsAPIUseCasesMock struct {
	// CreateTicketFunc mocks the CreateTicket method.
	CreateTicketFunc func(ctx context.Context, input tickets.CreateTicketInput) (tickets.CreateTicketOutput, error)

	// GetTicketFunc mocks the GetTicket method.
	GetTicketFunc func(ctx context.Context, input tickets.GetTicketInput) (tickets.GetTicketOutput, error)

	// calls tracks calls to the methods.
	calls struct {
		// CreateTicket holds details about calls to the CreateTicket method.
		CreateTicket []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Input is the input argument value.
			Input tickets.CreateTicketInput
		}
		// GetTicket holds details about calls to the GetTicket method.
		GetTicket []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Input is the input argument value.
			Input tickets.GetTicketInput
		}
	}
	lockCreateTicket sync.RWMutex
	lockGetTicket    sync.RWMutex
}

// CreateTicket calls CreateTicketFunc.
func (mock *TicketsAPIUseCasesMock) CreateTicket(ctx context.Context, input tickets.CreateTicketInput) (tickets.CreateTicketOutput, error) {
	callInfo := struct {
		Ctx   context.Context
		Input tickets.CreateTicketInput
	}{
		Ctx:   ctx,
		Input: input,
	}
	mock.lockCreateTicket.Lock()
	mock.calls.CreateTicket = append(mock.calls.CreateTicket, callInfo)
	mock.lockCreateTicket.Unlock()
	if mock.CreateTicketFunc == nil {
		var (
			createTicketOutputOut tickets.CreateTicketOutput
			errOut                error
		)
		return createTicketOutputOut, errOut
	}
	return mock.CreateTicketFunc(ctx, input)
}

// CreateTicketCalls gets all the calls that were made to CreateTicket.
// Check the length with:
//     len(mockedTicketsAPIUseCases.CreateTicketCalls())
func (mock *TicketsAPIUseCasesMock) CreateTicketCalls() []struct {
	Ctx   context.Context
	Input tickets.CreateTicketInput
} {
	var calls []struct {
		Ctx   context.Context
		Input tickets.CreateTicketInput
	}
	mock.lockCreateTicket.RLock()
	calls = mock.calls.CreateTicket
	mock.lockCreateTicket.RUnlock()
	return calls
}

// GetTicket calls GetTicketFunc.
func (mock *TicketsAPIUseCasesMock) GetTicket(ctx context.Context, input tickets.GetTicketInput) (tickets.GetTicketOutput, error) {
	callInfo := struct {
		Ctx   context.Context
		Input tickets.GetTicketInput
	}{
		Ctx:   ctx,
		Input: input,
	}
	mock.lockGetTicket.Lock()
	mock.calls.GetTicket = append(mock.calls.GetTicket, callInfo)
	mock.lockGetTicket.Unlock()
	if mock.GetTicketFunc == nil {
		var (
			getTicketOutputOut tickets.GetTicketOutput
			errOut             error
		)
		return getTicketOutputOut, errOut
	}
	return mock.GetTicketFunc(ctx, input)
}

// GetTicketCalls gets all the calls that were made to GetTicket.
// Check the length with:
//     len(mockedTicketsAPIUseCases.GetTicketCalls())
func (mock *TicketsAPIUseCasesMock) GetTicketCalls() []struct {
	Ctx   context.Context
	Input tickets.GetTicketInput
} {
	var calls []struct {
		Ctx   context.Context
		Input tickets.GetTicketInput
	}
	mock.lockGetTicket.RLock()
	calls = mock.calls.GetTicket
	mock.lockGetTicket.RUnlock()
	return calls
}
