// Code generated by MockGen. DO NOT EDIT.
// Source: agent.pb.go

// Package mock_idl is a generated GoMock package.
package mock_idl

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	idl "github.com/greenplum-db/gpdb/gp/idl"
	grpc "google.golang.org/grpc"
)

// MockAgentClient is a mock of AgentClient interface.
type MockAgentClient struct {
	ctrl     *gomock.Controller
	recorder *MockAgentClientMockRecorder
}

// MockAgentClientMockRecorder is the mock recorder for MockAgentClient.
type MockAgentClientMockRecorder struct {
	mock *MockAgentClient
}

// NewMockAgentClient creates a new mock instance.
func NewMockAgentClient(ctrl *gomock.Controller) *MockAgentClient {
	mock := &MockAgentClient{ctrl: ctrl}
	mock.recorder = &MockAgentClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAgentClient) EXPECT() *MockAgentClientMockRecorder {
	return m.recorder
}

// GetInterfaceAddrs mocks base method.
func (m *MockAgentClient) GetInterfaceAddrs(ctx context.Context, in *idl.GetInterfaceAddrsRequest, opts ...grpc.CallOption) (*idl.GetInterfaceAddrsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetInterfaceAddrs", varargs...)
	ret0, _ := ret[0].(*idl.GetInterfaceAddrsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInterfaceAddrs indicates an expected call of GetInterfaceAddrs.
func (mr *MockAgentClientMockRecorder) GetInterfaceAddrs(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInterfaceAddrs", reflect.TypeOf((*MockAgentClient)(nil).GetInterfaceAddrs), varargs...)
}

// MakeSegment mocks base method.
func (m *MockAgentClient) MakeSegment(ctx context.Context, in *idl.MakeSegmentRequest, opts ...grpc.CallOption) (*idl.MakeSegmentReply, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "MakeSegment", varargs...)
	ret0, _ := ret[0].(*idl.MakeSegmentReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MakeSegment indicates an expected call of MakeSegment.
func (mr *MockAgentClientMockRecorder) MakeSegment(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MakeSegment", reflect.TypeOf((*MockAgentClient)(nil).MakeSegment), varargs...)
}

// PgBasebackup mocks base method.
func (m *MockAgentClient) PgBasebackup(ctx context.Context, in *idl.PgBasebackupRequest, opts ...grpc.CallOption) (*idl.PgBasebackupResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "PgBasebackup", varargs...)
	ret0, _ := ret[0].(*idl.PgBasebackupResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PgBasebackup indicates an expected call of PgBasebackup.
func (mr *MockAgentClientMockRecorder) PgBasebackup(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PgBasebackup", reflect.TypeOf((*MockAgentClient)(nil).PgBasebackup), varargs...)
}

// PgControlData mocks base method.
func (m *MockAgentClient) PgControlData(ctx context.Context, in *idl.PgControlDataRequest, opts ...grpc.CallOption) (*idl.PgControlDataResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "PgControlData", varargs...)
	ret0, _ := ret[0].(*idl.PgControlDataResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PgControlData indicates an expected call of PgControlData.
func (mr *MockAgentClientMockRecorder) PgControlData(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PgControlData", reflect.TypeOf((*MockAgentClient)(nil).PgControlData), varargs...)
}

// StartSegment mocks base method.
func (m *MockAgentClient) StartSegment(ctx context.Context, in *idl.StartSegmentRequest, opts ...grpc.CallOption) (*idl.StartSegmentReply, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "StartSegment", varargs...)
	ret0, _ := ret[0].(*idl.StartSegmentReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StartSegment indicates an expected call of StartSegment.
func (mr *MockAgentClientMockRecorder) StartSegment(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartSegment", reflect.TypeOf((*MockAgentClient)(nil).StartSegment), varargs...)
}

// Status mocks base method.
func (m *MockAgentClient) Status(ctx context.Context, in *idl.StatusAgentRequest, opts ...grpc.CallOption) (*idl.StatusAgentReply, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Status", varargs...)
	ret0, _ := ret[0].(*idl.StatusAgentReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Status indicates an expected call of Status.
func (mr *MockAgentClientMockRecorder) Status(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Status", reflect.TypeOf((*MockAgentClient)(nil).Status), varargs...)
}

// Stop mocks base method.
func (m *MockAgentClient) Stop(ctx context.Context, in *idl.StopAgentRequest, opts ...grpc.CallOption) (*idl.StopAgentReply, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Stop", varargs...)
	ret0, _ := ret[0].(*idl.StopAgentReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Stop indicates an expected call of Stop.
func (mr *MockAgentClientMockRecorder) Stop(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockAgentClient)(nil).Stop), varargs...)
}

// UpdatePgConf mocks base method.
func (m *MockAgentClient) UpdatePgConf(ctx context.Context, in *idl.UpdatePgConfRequest, opts ...grpc.CallOption) (*idl.UpdatePgConfRespoonse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdatePgConf", varargs...)
	ret0, _ := ret[0].(*idl.UpdatePgConfRespoonse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdatePgConf indicates an expected call of UpdatePgConf.
func (mr *MockAgentClientMockRecorder) UpdatePgConf(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePgConf", reflect.TypeOf((*MockAgentClient)(nil).UpdatePgConf), varargs...)
}

// UpdatePgHbaConf mocks base method.
func (m *MockAgentClient) UpdatePgHbaConf(ctx context.Context, in *idl.UpdatePgHbaConfRequest, opts ...grpc.CallOption) (*idl.UpdatePgHbaConfResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdatePgHbaConf", varargs...)
	ret0, _ := ret[0].(*idl.UpdatePgHbaConfResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdatePgHbaConf indicates an expected call of UpdatePgHbaConf.
func (mr *MockAgentClientMockRecorder) UpdatePgHbaConf(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePgHbaConf", reflect.TypeOf((*MockAgentClient)(nil).UpdatePgHbaConf), varargs...)
}

// ValidateHostEnv mocks base method.
func (m *MockAgentClient) ValidateHostEnv(ctx context.Context, in *idl.ValidateHostEnvRequest, opts ...grpc.CallOption) (*idl.ValidateHostEnvReply, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ValidateHostEnv", varargs...)
	ret0, _ := ret[0].(*idl.ValidateHostEnvReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateHostEnv indicates an expected call of ValidateHostEnv.
func (mr *MockAgentClientMockRecorder) ValidateHostEnv(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateHostEnv", reflect.TypeOf((*MockAgentClient)(nil).ValidateHostEnv), varargs...)
}

// MockAgentServer is a mock of AgentServer interface.
type MockAgentServer struct {
	ctrl     *gomock.Controller
	recorder *MockAgentServerMockRecorder
}

// MockAgentServerMockRecorder is the mock recorder for MockAgentServer.
type MockAgentServerMockRecorder struct {
	mock *MockAgentServer
}

// NewMockAgentServer creates a new mock instance.
func NewMockAgentServer(ctrl *gomock.Controller) *MockAgentServer {
	mock := &MockAgentServer{ctrl: ctrl}
	mock.recorder = &MockAgentServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAgentServer) EXPECT() *MockAgentServerMockRecorder {
	return m.recorder
}

// GetInterfaceAddrs mocks base method.
func (m *MockAgentServer) GetInterfaceAddrs(arg0 context.Context, arg1 *idl.GetInterfaceAddrsRequest) (*idl.GetInterfaceAddrsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInterfaceAddrs", arg0, arg1)
	ret0, _ := ret[0].(*idl.GetInterfaceAddrsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInterfaceAddrs indicates an expected call of GetInterfaceAddrs.
func (mr *MockAgentServerMockRecorder) GetInterfaceAddrs(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInterfaceAddrs", reflect.TypeOf((*MockAgentServer)(nil).GetInterfaceAddrs), arg0, arg1)
}

// MakeSegment mocks base method.
func (m *MockAgentServer) MakeSegment(arg0 context.Context, arg1 *idl.MakeSegmentRequest) (*idl.MakeSegmentReply, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MakeSegment", arg0, arg1)
	ret0, _ := ret[0].(*idl.MakeSegmentReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MakeSegment indicates an expected call of MakeSegment.
func (mr *MockAgentServerMockRecorder) MakeSegment(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MakeSegment", reflect.TypeOf((*MockAgentServer)(nil).MakeSegment), arg0, arg1)
}

// PgBasebackup mocks base method.
func (m *MockAgentServer) PgBasebackup(arg0 context.Context, arg1 *idl.PgBasebackupRequest) (*idl.PgBasebackupResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PgBasebackup", arg0, arg1)
	ret0, _ := ret[0].(*idl.PgBasebackupResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PgBasebackup indicates an expected call of PgBasebackup.
func (mr *MockAgentServerMockRecorder) PgBasebackup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PgBasebackup", reflect.TypeOf((*MockAgentServer)(nil).PgBasebackup), arg0, arg1)
}

// PgControlData mocks base method.
func (m *MockAgentServer) PgControlData(arg0 context.Context, arg1 *idl.PgControlDataRequest) (*idl.PgControlDataResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PgControlData", arg0, arg1)
	ret0, _ := ret[0].(*idl.PgControlDataResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PgControlData indicates an expected call of PgControlData.
func (mr *MockAgentServerMockRecorder) PgControlData(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PgControlData", reflect.TypeOf((*MockAgentServer)(nil).PgControlData), arg0, arg1)
}

// StartSegment mocks base method.
func (m *MockAgentServer) StartSegment(arg0 context.Context, arg1 *idl.StartSegmentRequest) (*idl.StartSegmentReply, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartSegment", arg0, arg1)
	ret0, _ := ret[0].(*idl.StartSegmentReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StartSegment indicates an expected call of StartSegment.
func (mr *MockAgentServerMockRecorder) StartSegment(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartSegment", reflect.TypeOf((*MockAgentServer)(nil).StartSegment), arg0, arg1)
}

// Status mocks base method.
func (m *MockAgentServer) Status(arg0 context.Context, arg1 *idl.StatusAgentRequest) (*idl.StatusAgentReply, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Status", arg0, arg1)
	ret0, _ := ret[0].(*idl.StatusAgentReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Status indicates an expected call of Status.
func (mr *MockAgentServerMockRecorder) Status(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Status", reflect.TypeOf((*MockAgentServer)(nil).Status), arg0, arg1)
}

// Stop mocks base method.
func (m *MockAgentServer) Stop(arg0 context.Context, arg1 *idl.StopAgentRequest) (*idl.StopAgentReply, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop", arg0, arg1)
	ret0, _ := ret[0].(*idl.StopAgentReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Stop indicates an expected call of Stop.
func (mr *MockAgentServerMockRecorder) Stop(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockAgentServer)(nil).Stop), arg0, arg1)
}

// UpdatePgConf mocks base method.
func (m *MockAgentServer) UpdatePgConf(arg0 context.Context, arg1 *idl.UpdatePgConfRequest) (*idl.UpdatePgConfRespoonse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePgConf", arg0, arg1)
	ret0, _ := ret[0].(*idl.UpdatePgConfRespoonse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdatePgConf indicates an expected call of UpdatePgConf.
func (mr *MockAgentServerMockRecorder) UpdatePgConf(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePgConf", reflect.TypeOf((*MockAgentServer)(nil).UpdatePgConf), arg0, arg1)
}

// UpdatePgHbaConf mocks base method.
func (m *MockAgentServer) UpdatePgHbaConf(arg0 context.Context, arg1 *idl.UpdatePgHbaConfRequest) (*idl.UpdatePgHbaConfResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePgHbaConf", arg0, arg1)
	ret0, _ := ret[0].(*idl.UpdatePgHbaConfResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdatePgHbaConf indicates an expected call of UpdatePgHbaConf.
func (mr *MockAgentServerMockRecorder) UpdatePgHbaConf(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePgHbaConf", reflect.TypeOf((*MockAgentServer)(nil).UpdatePgHbaConf), arg0, arg1)
}

// ValidateHostEnv mocks base method.
func (m *MockAgentServer) ValidateHostEnv(arg0 context.Context, arg1 *idl.ValidateHostEnvRequest) (*idl.ValidateHostEnvReply, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateHostEnv", arg0, arg1)
	ret0, _ := ret[0].(*idl.ValidateHostEnvReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateHostEnv indicates an expected call of ValidateHostEnv.
func (mr *MockAgentServerMockRecorder) ValidateHostEnv(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateHostEnv", reflect.TypeOf((*MockAgentServer)(nil).ValidateHostEnv), arg0, arg1)
}
