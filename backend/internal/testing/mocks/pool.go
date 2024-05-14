// Code generated by MockGen. DO NOT EDIT.
// Source: ./db/connection/pool.go
//
// Generated by this command:
//
//	mockgen -source=./db/connection/pool.go -destination=./internal/testing/mocks/pool.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	pgx "github.com/jackc/pgx/v5"
	pgxpool "github.com/jackc/pgx/v5/pgxpool"
	gomock "go.uber.org/mock/gomock"
)

// MockPool is a mock of Pool interface.
type MockPool struct {
	ctrl     *gomock.Controller
	recorder *MockPoolMockRecorder
}

// MockPoolMockRecorder is the mock recorder for MockPool.
type MockPoolMockRecorder struct {
	mock *MockPool
}

// NewMockPool creates a new mock instance.
func NewMockPool(ctrl *gomock.Controller) *MockPool {
	mock := &MockPool{ctrl: ctrl}
	mock.recorder = &MockPoolMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPool) EXPECT() *MockPoolMockRecorder {
	return m.recorder
}

// AcquireFunc mocks base method.
func (m *MockPool) AcquireFunc(ctx context.Context, fn func(*pgxpool.Conn) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AcquireFunc", ctx, fn)
	ret0, _ := ret[0].(error)
	return ret0
}

// AcquireFunc indicates an expected call of AcquireFunc.
func (mr *MockPoolMockRecorder) AcquireFunc(ctx, fn any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AcquireFunc", reflect.TypeOf((*MockPool)(nil).AcquireFunc), ctx, fn)
}

// BeginTxFunc mocks base method.
func (m *MockPool) BeginTxFunc(ctx context.Context, txOptions pgx.TxOptions, fn func(pgx.Tx) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BeginTxFunc", ctx, txOptions, fn)
	ret0, _ := ret[0].(error)
	return ret0
}

// BeginTxFunc indicates an expected call of BeginTxFunc.
func (mr *MockPoolMockRecorder) BeginTxFunc(ctx, txOptions, fn any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeginTxFunc", reflect.TypeOf((*MockPool)(nil).BeginTxFunc), ctx, txOptions, fn)
}