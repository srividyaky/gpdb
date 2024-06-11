package hub_test

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/greenplum-db/gpdb/gpservice/idl/mock_idl"
	"github.com/greenplum-db/gpdb/gpservice/testutils"
	"google.golang.org/grpc/credentials"
	"reflect"
	"strings"
	"testing"

	"github.com/greenplum-db/gp-common-go-libs/testhelper"
	"github.com/greenplum-db/gpdb/gpservice/idl"
	"github.com/greenplum-db/gpdb/gpservice/internal/hub"
)

func TestServer_GetAllHostNames(t *testing.T) {
	testhelper.SetupTestLogger()
	hubConfig := testutils.CreateDummyServiceConfig(t)

	t.Run("returns error when one host errors getting hostname", func(t *testing.T) {
		testStr := "test error"
		hubServer := hub.New(hubConfig)
		ctrl := gomock.NewController(t)
		sdw1 := mock_idl.NewMockAgentClient(ctrl)
		sdw2 := mock_idl.NewMockAgentClient(ctrl)
		sdw1.EXPECT().GetHostName(gomock.Any(), gomock.Any()).Return(&idl.GetHostNameReply{}, fmt.Errorf(testStr))
		sdw2.EXPECT().GetHostName(gomock.Any(), gomock.Any()).Return(&idl.GetHostNameReply{Hostname: "sdw2"}, nil)
		request := idl.GetAllHostNamesRequest{HostList: []string{"sdw1", "sdw2"}}

		agentConns := []*hub.Connection{
			{AgentClient: sdw1, Hostname: "sdw1"},
			{AgentClient: sdw2, Hostname: "sdw2"},
		}
		hub.GetConnectionOnHostList = func(credentials credentials.TransportCredentials, agentPort int, hostList []string) (map[string]idl.AgentClient, error) {
			return map[string]idl.AgentClient{"sdw1": sdw1, "sdw2": sdw2}, nil
		}
		defer func() { hub.GetConnectionOnHostList = hub.GetConnectionOnHostListFn }()

		hubServer.Conns = agentConns

		_, err := hubServer.GetAllHostNames(context.Background(), &request)
		if err == nil || !strings.Contains(err.Error(), testStr) {
			t.Fatalf("Got:%v, expected:%s", err, testStr)
		}
	})
	t.Run("returns hostname map correctly", func(t *testing.T) {
		hubServer := hub.New(hubConfig)
		ctrl := gomock.NewController(t)
		sdw1 := mock_idl.NewMockAgentClient(ctrl)
		sdw2 := mock_idl.NewMockAgentClient(ctrl)
		sdw1.EXPECT().GetHostName(gomock.Any(), gomock.Any()).Return(&idl.GetHostNameReply{Hostname: "sdw1"}, nil)
		sdw2.EXPECT().GetHostName(gomock.Any(), gomock.Any()).Return(&idl.GetHostNameReply{Hostname: "sdw2"}, nil)
		request := idl.GetAllHostNamesRequest{HostList: []string{"sdw1", "sdw2"}}

		agentConns := []*hub.Connection{
			{AgentClient: sdw1, Hostname: "sdw1"},
			{AgentClient: sdw2, Hostname: "sdw2"},
		}
		hub.GetConnectionOnHostList = func(credentials credentials.TransportCredentials, agentPort int, hostList []string) (map[string]idl.AgentClient, error) {
			return map[string]idl.AgentClient{"sdw1": sdw1, "sdw2": sdw2}, nil
		}
		defer func() { hub.GetConnectionOnHostList = hub.GetConnectionOnHostListFn }()

		hubServer.Conns = agentConns
		response, err := hubServer.GetAllHostNames(context.Background(), &request)
		if err != nil {
			t.Fatalf("Got:%v, expected no error", err)
		}

		expected := map[string]string{"sdw1": "sdw1", "sdw2": "sdw2"}
		if !reflect.DeepEqual(response.HostNameMap, expected) {
			t.Fatalf("got: %v, want: %v", response.HostNameMap, expected)
		}
	})
	t.Run("returns hostname map correctly", func(t *testing.T) {
		testStr := "test error"
		hubConfig.Credentials = &testutils.MockCredentials{
			Err: fmt.Errorf(testStr),
		}
		hubServer := hub.New(hubConfig)

		request := idl.GetAllHostNamesRequest{HostList: []string{"sdw1", "sdw2"}}

		_, err := hubServer.GetAllHostNames(context.Background(), &request)
		if err == nil || !strings.Contains(err.Error(), testStr) {
			t.Fatalf("Got:%v, expected:%s", err, testStr)
		}

	})
}
