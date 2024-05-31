package hub_test

import (
	"context"
	"strings"
	"testing"

	"github.com/greenplum-db/gp-common-go-libs/testhelper"
	"github.com/greenplum-db/gpdb/gpservice/idl"
	"github.com/greenplum-db/gpdb/gpservice/internal/hub"
	"github.com/greenplum-db/gpdb/gpservice/internal/testutils"
)

func TestServer_GetAllHostNames(t *testing.T) {
	testhelper.SetupTestLogger()
	hubConfig := testutils.CreateDummyServiceConfig(t)

	t.Run("returns error when fails to load client credentials", func(t *testing.T) {
		testStr := "test error"
		hubServer := hub.New(hubConfig)
		request := idl.GetAllHostNamesRequest{HostList: []string{"sdw1", "sdw2"}}

		_, err := hubServer.GetAllHostNames(context.Background(), &request)
		if err == nil || !strings.Contains(err.Error(), testStr) {
			t.Fatalf("Got:%v, expected:%s", err, testStr)
		}
	})
}
