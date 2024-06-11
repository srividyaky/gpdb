package start

import (
	"strings"
	"testing"

	"github.com/greenplum-db/gpdb/gpservice/constants"
	"github.com/greenplum-db/gpdb/gpservice/internal/platform"
	"github.com/greenplum-db/gpdb/gpservice/test/integration/testutils"
)

func TestStartSuccess(t *testing.T) {

	const (
		gpserviceHub   = "gpservice_hub"
		gpserviceAgent = "gpservice_agent"
	)

	hosts := testutils.GetHostListFromFile(*hostfile)

	t.Run("start hub successfully", func(t *testing.T) {
		testutils.InitService(*hostfile, testutils.CertificateParams)
		cliParams := []string{"--hub"}
		expectedOut := []string{"[INFO]:-Hub service started successfully"}

		runStartCmdAndCheckOutput(t, cliParams, expectedOut)
		// check if service is running
		status, _ := testutils.GetSvcStatusOnHost(p.(platform.GpPlatform), gpserviceHub, hosts[0])
		testutils.VerifyServicePIDOnPort(t, status.OutputMsg, constants.DefaultHubPort, hosts[0])

		_, _ = testutils.RunStop("--hub")
	})

	t.Run("start hub and agents successfully", func(t *testing.T) {
		testutils.InitService(*hostfile, testutils.CertificateParams)

		cliParams := []string{""}
		expectedOut := []string{
			"[INFO]:-Hub service started successfully",
			"[INFO]:-Agent service started successfully",
		}
		runStartCmdAndCheckOutput(t, cliParams, expectedOut)
		// check if service is running
		for _, svc := range []string{gpserviceHub, gpserviceAgent} {
			listeningPort := constants.DefaultHubPort
			hostList := hosts[:1]
			if svc == gpserviceAgent {
				listeningPort = constants.DefaultAgentPort
				hostList = hosts
			}
			for _, host := range hostList {
				status, _ := testutils.GetSvcStatusOnHost(p.(platform.GpPlatform), svc, host)
				testutils.VerifyServicePIDOnPort(t, status.OutputMsg, listeningPort, host)
			}
		}
		_, _ = testutils.RunStop()
	})

	t.Run("start hub after gpservice configure with --service-name param", func(t *testing.T) {
		_, _ = testutils.RunConfigure(true, []string{
			"--hostfile", *hostfile,
			"--service-name", "dummySvc",
		}...)

		cliParams := []string{""}
		expectedOut := []string{
			"[INFO]:-Hub dummySvc started successfully",
			"[INFO]:-Agent dummySvc started successfully",
		}
		runStartCmdAndCheckOutput(t, cliParams, expectedOut)
		// check if service is running
		for _, svc := range []string{"dummySvc_hub", "dummySvc_agent"} {
			listeningPort := constants.DefaultHubPort
			hostList := hosts[:1]
			if svc == "dummySvc_agent" {
				listeningPort = constants.DefaultAgentPort
				hostList = hosts
			}
			for _, host := range hostList {
				status, _ := testutils.GetSvcStatusOnHost(p.(platform.GpPlatform), svc, host)
				testutils.VerifyServicePIDOnPort(t, status.OutputMsg, listeningPort, host)
			}
		}
		_, _ = testutils.RunStop()
	})

	t.Run("start services with --verbose param shows sevice status", func(t *testing.T) {
		testutils.InitService(*hostfile, testutils.CertificateParams)

		cliParams := []string{"--verbose"}
		expectedOut := []string{
			"[INFO]:-Hub service started successfully",
			"[INFO]:-Agent service started successfully",
			"ROLE", "HOST", "STATUS", "PID", "UPTIME",
		}

		runStartCmdAndCheckOutput(t, cliParams, expectedOut)
		// check if service is running
		for _, svc := range []string{gpserviceHub, gpserviceAgent} {
			listeningPort := constants.DefaultHubPort
			hostList := hosts[:1]
			if svc == gpserviceAgent {
				listeningPort = constants.DefaultAgentPort
				hostList = hosts
			}
			for _, host := range hostList {
				status, _ := testutils.GetSvcStatusOnHost(p.(platform.GpPlatform), svc, host)
				testutils.VerifyServicePIDOnPort(t, status.OutputMsg, listeningPort, host)
			}
		}
		_, _ = testutils.RunStop()
	})

	t.Run("start hub with --verbose param shows hub status", func(t *testing.T) {
		testutils.InitService(*hostfile, testutils.CertificateParams)

		cliParams := []string{"--hub", "--verbose"}
		expectedOut := []string{"[INFO]:-Hub service started successfully", "Hub", "running"}

		runStartCmdAndCheckOutput(t, cliParams, expectedOut)
		// check if service is running
		status, _ := testutils.GetSvcStatusOnHost(p.(platform.GpPlatform), gpserviceHub, hosts[0])
		testutils.VerifyServicePIDOnPort(t, status.OutputMsg, constants.DefaultHubPort, hosts[0])

		_, _ = testutils.RunStop("--hub")
	})
}

func TestStartSuccessHelp(t *testing.T) {
	TestCases := []struct {
		name        string
		cliParams   []string
		expectedOut []string
	}{

		{
			name: "start command with invalid params shows help",
			cliParams: []string{
				"invalid",
			},
			expectedOut: append([]string{
				"Start hub and agent services",
			}, testutils.CommonHelpText...),
		},
		{
			name: "start command with --help params shows help",
			cliParams: []string{
				"--help",
			},
			expectedOut: append([]string{
				"Start hub and agent services",
			}, testutils.CommonHelpText...),
		},
		{
			name: "start command with -h params shows help",
			cliParams: []string{
				"-h",
			},
			expectedOut: append([]string{
				"Start hub and agent services",
			}, testutils.CommonHelpText...),
		},
	}
	for _, tc := range TestCases {
		t.Run(tc.name, func(t *testing.T) {
			runStartCmdAndCheckOutput(t, tc.cliParams, tc.expectedOut)
		})
	}
}

func runStartCmdAndCheckOutput(t *testing.T, input []string, output []string) {
	testutils.InitService(*hostfile, testutils.CertificateParams)
	result, err := testutils.RunStart(input...)
	// check for command result
	if err != nil {
		t.Errorf("\nUnexpected error: %#v", err)
	}
	if result.ExitCode != 0 {
		t.Errorf("\nExpected: %v \nGot: %v", 0, result.ExitCode)
	}
	for _, item := range output {
		if !strings.Contains(result.OutputMsg, item) {
			t.Errorf("\nExpected string: %#v \nNot found in: %#v", item, result.OutputMsg)
		}
	}
}
