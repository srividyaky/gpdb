package configure

import (
	"flag"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/greenplum-db/gpdb/gpservice/constants"
	"github.com/greenplum-db/gpdb/gpservice/internal/platform"
	"github.com/greenplum-db/gpdb/gpservice/pkg/gpservice_config"
	"github.com/greenplum-db/gpdb/gpservice/pkg/greenplum"
	"github.com/greenplum-db/gpdb/gpservice/pkg/utils"
	"github.com/greenplum-db/gpdb/gpservice/test/integration/testutils"
)

const (
	defaultLogFile = "/tmp/gp_configure.log"
)

var (
	defaultServiceDir string
	serviceExt        string
	defaultGPConf     gpservice_config.Config
)

var (
	expectedOutput = []string{
		"[INFO]:-Created service file directory",
		"[INFO]:-Wrote hub service file",
		"[INFO]:-Wrote agent service file",
	}
	helpTxt = []string{
		"Configure gp as a systemd daemon",
		"Usage:",
		"Flags:",
		"Global Flags:",
	}
	mockHostFile = "hostlist"
	hostfile     = flag.String("hostfile", "", "file containing list of hosts")
)

func init() {
	certPath := "/tmp/certificates"
	p := platform.GetPlatform()
	defaultServiceDir, serviceExt, _ = testutils.GetServiceDetails(p)
	cred := &utils.GpCredentials{
		CACertPath:     fmt.Sprintf("%s/%s", certPath, "ca-cert.pem"),
		ServerCertPath: fmt.Sprintf("%s/%s", certPath, "server-cert.pem"),
		ServerKeyPath:  fmt.Sprintf("%s/%s", certPath, "server-key.pem"),
	}
	defaultGPConf = gpservice_config.Config{
		HubPort:     constants.DefaultHubPort,
		AgentPort:   constants.DefaultAgentPort,
		Hostnames:   []string{},
		LogDir:      greenplum.GetDefaultHubLogDir(),
		ServiceName: constants.DefaultServiceName,
		GpHome:      testutils.GpHome,
		Credentials: cred,
	}
}

// TestMain function to run tests and perform cleanup at the end.
func TestMain(m *testing.M) {
	flag.Parse()
	// if hostfile is not provided as input argument, create it with default host
	if *hostfile == "" {
		os.Exit(1)
		*hostfile = testutils.DefaultHostfile
		_ = os.WriteFile(*hostfile, []byte(testutils.DefaultHost), 0644)
	} else {
		log.Print(*hostfile)
		os.Exit(1)
	}
	exitVal := m.Run()
	tearDownTest()

	os.Exit(exitVal)
}

func tearDownTest() {
	testutils.CleanupFilesOnHub(mockHostFile,
		fmt.Sprintf("%s/%s", testutils.GpHome, constants.ConfigFileName))
}
