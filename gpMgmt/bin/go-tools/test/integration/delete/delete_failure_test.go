package delete

import (
	"github.com/greenplum-db/gpdb/gp/test/integration/testutils"
	"os"
	"strings"
	"testing"
)

func TestDeleteFailure(t *testing.T) {
	//hosts := testutils.GetHostListFromFile(*hostfile)

	t.Run("delete services fails when services are not configured", func(t *testing.T) {
		expectedOut := []string{"could not open config file"}

		result, err := testutils.RunDelete("services")
		if err == nil {
			t.Errorf("\nExpected error Got: %#v", err)
		}
		if result.ExitCode != testutils.ExitCode1 {
			t.Errorf("\nExpected: %#v \nGot: %v", testutils.ExitCode1, result.ExitCode)
		}
		for _, item := range expectedOut {
			if !strings.Contains(result.OutputMsg, item) {
				t.Errorf("\nExpected string: %#v \nNot found in: %#v", item, result.OutputMsg)
			}
		}
	})

	t.Run("delete services fails when config file does not exist", func(t *testing.T) {
		testutils.InitService(*hostfile, testutils.CertificateParams)

		expectedOut := []string{"could not open config file"}

		//Remove configuration file
		_ = testutils.CopyFile(testutils.DefaultConfigurationFile, "/tmp/config.conf")
		_ = os.RemoveAll(testutils.DefaultConfigurationFile)

		result, err := testutils.RunDelete("services")
		if err == nil {
			t.Errorf("\nExpected error Got: %#v", err)
		}
		if result.ExitCode != testutils.ExitCode1 {
			t.Errorf("\nExpected: %#v \nGot: %v", testutils.ExitCode1, result.ExitCode)
		}
		for _, item := range expectedOut {
			if !strings.Contains(result.OutputMsg, item) {
				t.Errorf("\nExpected string: %#v \nNot found in: %#v", item, result.OutputMsg)
			}
		}
	})

	t.Run("delete services fails without hub service file", func(t *testing.T) {
		testutils.InitService(*hostfile, testutils.CertificateParams)

		expectedOut := []string{"could not remove hub service"}

		testutils.DisableandDeleteHubServiceFile(p, "gp_hub")

		result, err := testutils.RunDelete("services")
		if err == nil {
			t.Errorf("\nExpected error Got: %#v", err)
		}
		if result.ExitCode != testutils.ExitCode1 {
			t.Errorf("\nExpected: %#v \nGot: %v", testutils.ExitCode1, result.ExitCode)
		}
		for _, item := range expectedOut {
			if !strings.Contains(result.OutputMsg, item) {
				t.Errorf("\nExpected string: %#v \nNot found in: %#v", item, result.OutputMsg)
			}
		}
	})

	t.Run("delete services fails without agent service file", func(t *testing.T) {
		testutils.InitService(*hostfile, testutils.CertificateParams)

		expectedOut := []string{"could not remove agent service"}

		testutils.DisableandDeleteAgentServiceFile(p, "gp_agent")

		result, err := testutils.RunDelete("services")
		if err == nil {
			t.Errorf("\nExpected error Got: %#v", err)
		}
		if result.ExitCode != testutils.ExitCode1 {
			t.Errorf("\nExpected: %#v \nGot: %v", testutils.ExitCode1, result.ExitCode)
		}
		for _, item := range expectedOut {
			if !strings.Contains(result.OutputMsg, item) {
				t.Errorf("\nExpected string: %#v \nNot found in: %#v", item, result.OutputMsg)
			}
		}
	})
}
