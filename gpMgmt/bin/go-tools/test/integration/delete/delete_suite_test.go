package delete

import (
	"flag"
	"fmt"
	"github.com/greenplum-db/gpdb/gp/test/integration/testutils"
	"github.com/greenplum-db/gpdb/gp/utils"
	"os"
	"testing"
)

var (
	p          = utils.GetPlatform()
	configCopy = "config_copy.conf"
	hostfile   = flag.String("hostfile", "", "file containing list of hosts")
)

func TestMain(m *testing.M) {
	flag.Parse()
	// if hostfile is not provided as input argument, create it with default host
	// Hostfile is required to distinguish between single host and multi host testing
	if *hostfile == "" {
		file, err := os.CreateTemp("", "")
		if err != nil {
			fmt.Printf("could not create hostfile: %v, and no hostfile provided", err)
			os.Exit(1)
		}

		*hostfile = file.Name()
		hostname, err := os.Hostname()
		if err != nil {
			fmt.Printf("could not get hostname: %v", err)
			os.Exit(1)
		}

		err = os.WriteFile(*hostfile, []byte(hostname), 0777)
		if err != nil {
			fmt.Printf("could not create hostfile: %v, and no hostfile provided", err)
			os.Exit(1)
		}

	}
	exitVal := m.Run()
	tearDownTest()

	os.Exit(exitVal)
}

func tearDownTest() {
	testutils.CleanupFilesOnHub(configCopy)
	testutils.DisableandDeleteServiceFiles(p)
}
