package status

import (
	"strings"
	"testing"

	"github.com/greenplum-db/gpdb/gpservice/pkg/gpservice_config"
	"github.com/greenplum-db/gpdb/gpservice/test/integration/testutils"
)

var gpCfg gpservice_config.Config

func TestStatusSuccess(t *testing.T) {
	const (
		gpserviceHub   = "gpservice_hub"
		gpserviceAgent = "gpservice_agent"
	)
	var StatusSuccessTestCases = []struct {
		name        string
		cliParams   []string
		expectedOut []string
		serviceName []string
	}{
		{
			name: "status services shows status of hub and agents",
			cliParams: []string{
				"",
			},
			expectedOut: []string{
				"ROLE", "HOST", "STATUS", "PID", "UPTIME",
				"Hub", "running",
				"Agent", "running",
			},
			serviceName: []string{
				gpserviceHub,
				gpserviceAgent,
			},
		},
		{
			name: "status services with --verbose cli param",
			cliParams: []string{
				"--verbose",
			},
			expectedOut: []string{
				"ROLE", "HOST", "STATUS", "PID", "UPTIME",
				"Hub", "running",
				"Agent", "running",
			},
			serviceName: []string{
				gpserviceHub,
				gpserviceAgent,
			},
		},
	}

	for _, tc := range StatusSuccessTestCases {
		t.Run(tc.name, func(t *testing.T) {
			testutils.InitService(*hostfile, testutils.CertificateParams)
			_, _ = testutils.RunStart()
			gpCfg = testutils.ParseConfig(testutils.DefaultConfigurationFile)

			// Running the gpservice status command
			result, err := testutils.RunStatus(tc.cliParams...)
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

			// verify the pid in status is listening on correct port
			statusMap := testutils.ExtractStatusData(result.OutputMsg)
			for _, svc := range tc.serviceName {
				listeningPort := gpCfg.HubPort
				hostPidMap := statusMap["Hub"]
				if svc == "gp_agent" {
					listeningPort = gpCfg.AgentPort
					hostPidMap = statusMap["Agent"]
				}
				for host, pid := range hostPidMap {
					testutils.VerifyServicePIDOnPort(t, pid, listeningPort, host)
				}
			}
			_, _ = testutils.RunStop()
		})
	}
}

func TestStatusSuccessWithoutDefaultService(t *testing.T) {
	t.Run("status services when gp installed with --service-name param", func(t *testing.T) {
		params := append(testutils.CertificateParams, []string{"--service-name", "dummySvc"}...)
		testutils.InitService(*hostfile, params)
		_, _ = testutils.RunStart()

		cliParams := []string{
			"--verbose",
		}
		expectedOut := []string{
			"ROLE", "HOST", "STATUS", "PID", "UPTIME",
			"Hub", "running",
			"Agent", "running",
		}

		gpCfg = testutils.ParseConfig(testutils.DefaultConfigurationFile)

		result, err := testutils.RunStatus(cliParams...)
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

		// verify the pid in status is listening on correct port
		statusMap := testutils.ExtractStatusData(result.OutputMsg)
		for _, svc := range []string{"dummySvc_hub", "dummySvc_agent"} {
			listeningPort := gpCfg.HubPort
			hostPidMap := statusMap["Hub"]
			if svc == "dummySvc_agent" {
				listeningPort = gpCfg.AgentPort
				hostPidMap = statusMap["Agent"]
			}
			for host, pid := range hostPidMap {
				testutils.VerifyServicePIDOnPort(t, pid, listeningPort, host)
			}
		}
		_, _ = testutils.RunStop()
	})

}

func TestStatusSuccessHelp(t *testing.T) {
	TestCases := []struct {
		name        string
		cliParams   []string
		expectedOut []string
	}{
		{
			name: "status command with invalid param shows help",
			cliParams: []string{
				"invalid",
			},
			expectedOut: append([]string{
				"Display status",
			}, testutils.CommonHelpText...),
		},
		{
			name: "status command with --help shows help",
			cliParams: []string{
				"--help",
			},
			expectedOut: append([]string{
				"Display status",
			}, testutils.CommonHelpText...),
		},
		{
			name: "status command with -h shows help",
			cliParams: []string{
				"-h",
			},
			expectedOut: append([]string{
				"Display status",
			}, testutils.CommonHelpText...),
		},
	}

	for _, tc := range TestCases {
		t.Run(tc.name, func(t *testing.T) {
			testutils.InitService(*hostfile, testutils.CertificateParams)
			_, _ = testutils.RunStart()
			result, err := testutils.RunStatus(tc.cliParams...)
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
