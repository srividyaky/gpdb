package cli_test

import (
	"fmt"
	"github.com/greenplum-db/gp-common-go-libs/testhelper"
	"github.com/greenplum-db/gpdb/gpservice/internal/agent"
	"github.com/greenplum-db/gpdb/gpservice/pkg/gpservice_config"
	"os/exec"
	"strings"
	"testing"

	"github.com/greenplum-db/gpdb/gpservice/internal/cli"
	"github.com/greenplum-db/gpdb/gpservice/pkg/utils"
	"github.com/greenplum-db/gpdb/gpservice/testutils"
	"github.com/greenplum-db/gpdb/gpservice/testutils/exectest"
)

func TestDeleteServices(t *testing.T) {

	_, _, logfile := testhelper.SetupTestLogger()
	resetConf := cli.SetConf(testutils.CreateDummyServiceConfig(t))
	defer resetConf()

	platform := &testutils.MockPlatform{}
	agent.SetPlatform(platform)
	defer agent.ResetPlatform()

	t.Run("DeleteServices does not fail if StopServices fails", func(t *testing.T) {
		cli.StopServices = func(conf *gpservice_config.Config) error {
			return fmt.Errorf("error")
		}
		defer func() { cli.StopServices = cli.StopServicesFunc }()

		utils.System.ExecCommand = exectest.NewCommand(exectest.Success)
		utils.System.Remove = func(name string) error {
			return nil
		}
		defer utils.ResetSystemFunctions()

		_, err := testutils.ExecuteCobraCommand(t, cli.DeleteCmd(), "services")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		testutils.AssertLogMessage(t, logfile, `\[DEBUG\]:-stop services returned error:error. Probably already stopped`)
		testutils.AssertLogMessage(t, logfile, `\[INFO\]:-Successfully deleted hub service`)
		testutils.AssertLogMessage(t, logfile, `\[INFO\]:-Successfully deleted agent service`)
		testutils.AssertLogMessage(t, logfile, `\[INFO\]:-Removed hub service file .* from coordinator host`)
		testutils.AssertLogMessage(t, logfile, `\[INFO\]:-Successfully removed agent service file .* from segment hosts`)
	})
	t.Run("DeleteServices fails when removal of conf file fails", func(t *testing.T) {
		cli.StopServices = func(conf *gpservice_config.Config) error {
			return nil
		}
		defer func() { cli.StopServices = cli.StopServicesFunc }()

		utils.System.ExecCommand = func(name string, args ...string) *exec.Cmd {
			if strings.Contains(strings.Join(args, " "), "rm") {
				return exectest.NewCommand(exectest.Failure)(name, args...)
			}
			return exectest.NewCommand(exectest.Success)(name, args...)
		}
		utils.System.Remove = func(name string) error {
			return nil
		}

		defer utils.ResetSystemFunctions()

		_, err := testutils.ExecuteCobraCommand(t, cli.DeleteCmd())
		expectedErrPrefix := "could not delete agent service file"
		if !strings.HasPrefix(err.Error(), expectedErrPrefix) {
			t.Fatalf("got %v, want %s", err, expectedErrPrefix)
		}
	})

	t.Run("DeleteServices fails when hub service deletion fails", func(t *testing.T) {

		cli.StopServices = func(conf *gpservice_config.Config) error {
			return nil
		}
		defer func() { cli.StopServices = cli.StopServicesFunc }()

		utils.System.ExecCommand = func(name string, args ...string) *exec.Cmd {
			if strings.Contains(strings.Join(args, " "), "remove") {
				return exectest.NewCommand(exectest.Failure)(name, args...)
			}

			return exectest.NewCommand(exectest.Success)(name, args...)
		}

		defer utils.ResetSystemFunctions()

		_, err := testutils.ExecuteCobraCommand(t, cli.DeleteCmd())
		expectedErrPrefix := "could not remove hub service"
		if !strings.HasPrefix(err.Error(), expectedErrPrefix) {
			t.Fatalf("got %v, want %s", err, expectedErrPrefix)
		}
	})
}
