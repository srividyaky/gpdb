package stop

import (
	"strings"
	"testing"

	"github.com/greenplum-db/gpdb/gpservice/internal/platform"
	"github.com/greenplum-db/gpdb/gpservice/test/integration/testutils"
)

func TestStopSuccess(t *testing.T) {
	hosts := testutils.GetHostListFromFile(*hostfile)

	const (
		gpserviceHub   = "gpservice_hub"
		gpserviceAgent = "gpservice_agent"
	)
	t.Run("stop services successfully", func(t *testing.T) {
		testutils.InitService(*hostfile, testutils.CertificateParams)
		_, _ = testutils.RunStart()

		expectedOut := []string{
			"Agent service stopped successfully",
			"Hub service stopped successfully",
		}

		// Running the gpservice stop command for services
		result, err := testutils.RunStop()
		// check for command result
		if err != nil {
			t.Errorf("\nUnexpected error: %#v", err)
		}
		if result.ExitCode != 0 {
			t.Errorf("\nExpected: %v \nGot: %v", 0, result.ExitCode)
		}
		for _, item := range expectedOut {
			if !strings.Contains(result.OutputMsg, item) {
				t.Errorf("\nExpected string: %#v \nNot found in: %#v", item, result.OutputMsg)
			}
		}

		// check if service is not running
		for _, svc := range []string{gpserviceHub, gpserviceAgent} {
			hostList := hosts[:1]
			if svc == gpserviceAgent {
				hostList = hosts
			}
			for _, host := range hostList {
				status, _ := testutils.GetSvcStatusOnHost(p.(platform.GpPlatform), svc, host)
				testutils.VerifySvcNotRunning(t, status.OutputMsg)
			}
		}
	})

	t.Run("stop hub successfully", func(t *testing.T) {
		testutils.InitService(*hostfile, testutils.CertificateParams)
		_, _ = testutils.RunStart("--hub")

		expectedOut := "Hub service stopped successfully"

		// Running the gp stop command for hub
		result, err := testutils.RunStop("--hub")
		// check for command result
		if err != nil {
			t.Errorf("\nUnexpected error: %#v", err)
		}
		if result.ExitCode != 0 {
			t.Errorf("\nExpected: %v \nGot: %v", 0, result.ExitCode)
		}
		if !strings.Contains(result.OutputMsg, expectedOut) {
			t.Errorf("\nExpected string: %#v \nNot found in: %#v", expectedOut, result.OutputMsg)
		}

		// check if service is not running
		status, _ := testutils.GetSvcStatusOnHost(p.(platform.GpPlatform), "gp_hub", hosts[0])
		testutils.VerifySvcNotRunning(t, status.OutputMsg)
	})

	t.Run("stop agents successfully", func(t *testing.T) {
		testutils.InitService(*hostfile, testutils.CertificateParams)
		_, _ = testutils.RunStart()

		expectedOut := "Agent service stopped successfully"

		// Running the gp stop command for agents
		result, err := testutils.RunStop("--agent")
		// check for command result
		if err != nil {
			t.Errorf("\nUnexpected error: %#v", err)
		}
		if result.ExitCode != 0 {
			t.Errorf("\nExpected: %v \nGot: %v", 0, result.ExitCode)
		}
		if !strings.Contains(result.OutputMsg, expectedOut) {
			t.Errorf("\nExpected string: %#v \nNot found in: %#v", expectedOut, result.OutputMsg)
		}

		// check if service is not running
		for _, host := range hosts {
			status, _ := testutils.GetSvcStatusOnHost(p.(platform.GpPlatform), "gp_agent", host)
			testutils.VerifySvcNotRunning(t, status.OutputMsg)
		}

		_, _ = testutils.RunStop("--hub")
	})

	t.Run("stop services command with --verbose shows status details", func(t *testing.T) {
		testutils.InitService(*hostfile, testutils.CertificateParams)
		_, _ = testutils.RunStart()

		cliParams := []string{
			"--verbose",
		}
		expectedOut := []string{
			"Agent service stopped successfull",
			"Hub service stopped successfully",
		}

		result, err := testutils.RunStop(cliParams...)
		// check for command result
		if err != nil {
			t.Errorf("\nUnexpected error: %#v", err)
		}
		if result.ExitCode != 0 {
			t.Errorf("\nExpected: %v \nGot: %v", 0, result.ExitCode)
		}
		for _, item := range expectedOut {
			if !strings.Contains(result.OutputMsg, item) {
				t.Errorf("\nExpected string: %#v \nNot found in: %#v", item, result.OutputMsg)
			}
		}

		// check if service is not running
		for _, svc := range []string{"gp_hub", "gp_agent"} {
			hostList := hosts[:1]
			if svc == "gp_agent" {
				hostList = hosts
			}
			for _, host := range hostList {
				status, _ := testutils.GetSvcStatusOnHost(p.(platform.GpPlatform), svc, host)
				testutils.VerifySvcNotRunning(t, status.OutputMsg)
			}
		}
	})
}

func TestStopSuccessHelp(t *testing.T) {
	TestCases := []struct {
		name        string
		cliParams   []string
		expectedOut []string
	}{
		//this is failing, bug is raised for error message
		{
			name: "stop command with invalid param shows help",
			cliParams: []string{
				"invalid",
			},
			expectedOut: append([]string{
				"Stop processes",
			}, testutils.CommonHelpText...),
		},
		{
			name: "stop command with --help shows help",
			cliParams: []string{
				"--help",
			},
			expectedOut: append([]string{
				"Stop hub and agent services",
			}, testutils.CommonHelpText...),
		},
		{
			name: "stop command with -h shows help",
			cliParams: []string{
				"-h",
			},
			expectedOut: append([]string{
				"Stop hub and agent services",
			}, testutils.CommonHelpText...),
		},
	}
	for _, tc := range TestCases {
		t.Run(tc.name, func(t *testing.T) {
			testutils.InitService(*hostfile, testutils.CertificateParams)
			testutils.RunStart()

			result, err := testutils.RunStop(tc.cliParams...)
			// check for command result
			if err != nil {
				t.Errorf("\nUnexpected error: %#v", err)
			}
			if result.ExitCode != 0 {
				t.Errorf("\nExpected: %v \nGot: %v", 0, result.ExitCode)
			}
			for _, item := range tc.expectedOut {
				if !strings.Contains(result.OutputMsg, item) {
					t.Errorf("\nExpected string: %#v \nNot found in: %#v", item, result.OutputMsg)
				}
			}
		})
	}
}
