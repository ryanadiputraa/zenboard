// Code generated by mockery v2.31.1. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/ryanadiputraa/zenboard/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

// BoardRepository is an autogenerated mock type for the BoardRepository type
type BoardRepository struct {
	mock.Mock
}

// FetchByOwnerID provides a mock function with given fields: ctx, ownerID
func (_m *BoardRepository) FetchByOwnerID(ctx context.Context, ownerID string) ([]domain.Board, error) {
	ret := _m.Called(ctx, ownerID)

	var r0 []domain.Board
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]domain.Board, error)); ok {
		return rf(ctx, ownerID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []domain.Board); ok {
		r0 = rf(ctx, ownerID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Board)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, ownerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Init provides a mock function with given fields: ctx, board, task1, task2, task3
func (_m *BoardRepository) Init(ctx context.Context, board domain.InitBoardDTO, task1 domain.InitTaskDTO, task2 domain.InitTaskDTO, task3 domain.InitTaskDTO) (domain.Board, error) {
	ret := _m.Called(ctx, board, task1, task2, task3)

	var r0 domain.Board
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.InitBoardDTO, domain.InitTaskDTO, domain.InitTaskDTO, domain.InitTaskDTO) (domain.Board, error)); ok {
		return rf(ctx, board, task1, task2, task3)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.InitBoardDTO, domain.InitTaskDTO, domain.InitTaskDTO, domain.InitTaskDTO) domain.Board); ok {
		r0 = rf(ctx, board, task1, task2, task3)
	} else {
		r0 = ret.Get(0).(domain.Board)
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.InitBoardDTO, domain.InitTaskDTO, domain.InitTaskDTO, domain.InitTaskDTO) error); ok {
		r1 = rf(ctx, board, task1, task2, task3)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewBoardRepository creates a new instance of BoardRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBoardRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *BoardRepository {
	mock := &BoardRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
