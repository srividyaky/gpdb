package init_cluster

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/greenplum-db/gp-common-go-libs/cluster"
	"github.com/greenplum-db/gp-common-go-libs/dbconn"
	"github.com/greenplum-db/gpdb/gp/cli"
	"github.com/greenplum-db/gpdb/gp/constants"
	"github.com/greenplum-db/gpdb/gp/test/integration/testutils"
)

func TestLocaleValidation(t *testing.T) {
	localTypes := []string{"LC_COLLATE", "LC_CTYPE", "LC_MESSAGES", "LC_MONETARY", "LC_NUMERIC", "LC_TIME"}

	t.Run("when LC_ALL is given, it sets the locale for all the types", func(t *testing.T) {
		expected := testutils.GetRandomLocale(t)

		configFile := testutils.GetTempFile(t, "config.json")
		SetConfigKey(t, configFile, "locale", cli.Locale{
			LcAll: expected,
		}, true)

		result, err := testutils.RunInitCluster(configFile)
		if err != nil {
			t.Fatalf("unexpected error: %s, %v", result.OutputMsg, err)
		}

		for _, localType := range localTypes {
			testutils.AssertPgConfig(t, localType, expected)
		}

		_, err = testutils.DeleteCluster()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("individual locale type takes precedence over LC_ALL", func(t *testing.T) {
		expected := testutils.GetRandomLocale(t)
		expectedLcCtype := testutils.GetRandomLocale(t)

		configFile := testutils.GetTempFile(t, "config.json")
		SetConfigKey(t, configFile, "locale", cli.Locale{
			LcAll:   expected,
			LcCtype: expectedLcCtype,
		}, true)

		result, err := testutils.RunInitCluster(configFile)
		if err != nil {
			t.Fatalf("unexpected error: %s, %v", result.OutputMsg, err)
		}

		for _, localType := range localTypes {
			if localType == "LC_CTYPE" {
				testutils.AssertPgConfig(t, localType, expectedLcCtype)
			} else {
				testutils.AssertPgConfig(t, localType, expected)
			}
		}

		_, err = testutils.DeleteCluster()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("when no locale value is provided, inherits the locale from the environment", func(t *testing.T) {
		// TODO: on macos launchd does not inherit the system locale value
		// so skip it for now until we find a way to test it.
		if runtime.GOOS == constants.PlatformDarwin {
			t.Skip()
		}

		configFile := testutils.GetTempFile(t, "config.json")
		UnsetConfigKey(t, configFile, "locale", true)

		result, err := testutils.RunInitCluster(configFile)
		if err != nil {
			t.Fatalf("unexpected error: %s, %v", result.OutputMsg, err)
		}

		for _, localType := range localTypes {
			testutils.AssertPgConfig(t, localType, testutils.GetSystemLocale(t, localType))
		}

		_, err = testutils.DeleteCluster()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("when invalid locale is given", func(t *testing.T) {
		configFile := testutils.GetTempFile(t, "config.json")
		SetConfigKey(t, configFile, "locale", cli.Locale{
			LcAll: "invalid.locale",
		}, true)

		result, err := testutils.RunInitCluster(configFile)
		if e, ok := err.(*exec.ExitError); !ok || e.ExitCode() != 1 {
			t.Fatalf("got %v, want exit status 1", err)
		}

		expectedOut := `\[ERROR\]:-validating hosts: host: (\S+), locale value 'invalid.locale' is not a valid locale`

		match, err := regexp.MatchString(expectedOut, result.OutputMsg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !match {
			t.Fatalf("got %q, want %q", result.OutputMsg, expectedOut)
		}
	})
}

func TestPgConfig(t *testing.T) {
	t.Run("sets the correct config values as provided by the user", func(t *testing.T) {
		configFile := testutils.GetTempFile(t, "config.json")
		SetConfigKey(t, configFile, "coordinator-config", map[string]string{
			"max_connections": "15",
		}, true)
		SetConfigKey(t, configFile, "segment-config", map[string]string{
			"max_connections": "10",
		})
		SetConfigKey(t, configFile, "common-config", map[string]string{
			"max_wal_senders": "5",
		})

		result, err := testutils.RunInitCluster(configFile)
		if err != nil {
			t.Fatalf("unexpected error: %s, %v", result.OutputMsg, err)
		}

		testutils.AssertPgConfig(t, "max_connections", "15", -1)
		testutils.AssertPgConfig(t, "max_connections", "10", 0)
		testutils.AssertPgConfig(t, "max_wal_senders", "5", -1)
		testutils.AssertPgConfig(t, "max_wal_senders", "5", 0)

		_, err = testutils.DeleteCluster()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("segment-config and coordinator-config take precedence over the common-config", func(t *testing.T) {
		configFile := testutils.GetTempFile(t, "config.json")
		SetConfigKey(t, configFile, "coordinator-config", map[string]string{
			"max_connections": "15",
		}, true)
		SetConfigKey(t, configFile, "segment-config", map[string]string{
			"max_connections": "10",
		})
		SetConfigKey(t, configFile, "common-config", map[string]string{
			"max_connections": "25",
		})

		result, err := testutils.RunInitCluster(configFile)
		if err != nil {
			t.Fatalf("unexpected error: %s, %v", result.OutputMsg, err)
		}

		testutils.AssertPgConfig(t, "max_connections", "15", -1)
		testutils.AssertPgConfig(t, "max_connections", "10", 0)

		_, err = testutils.DeleteCluster()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("check if the gp_segment_configuration table has the correct value", func(t *testing.T) {
		var value cli.Segment
		var ok bool
		configFile := testutils.GetTempFile(t, "config.json")
		config := GetDefaultConfig(t)

		err := config.WriteConfigAs(configFile)
		if err != nil {
			t.Fatalf("unexpected error: %#v", err)
		}

		coordinator := config.Get("coordinator")
		if value, ok = coordinator.(cli.Segment); !ok {
			t.Fatalf("unexpected data type for coordinator %T", value)
		}

		primarySegs := config.Get("segment-array")
		valueSegPair, ok := primarySegs.([]cli.SegmentPair)
		if !ok {
			t.Fatalf("unexpected data type for segment-array %T", primarySegs)
		}

		var primarySegments []cli.Segment
		primarySegments = append(primarySegments, value)

		for _, segPair := range valueSegPair {
			primarySegments = append(primarySegments, *segPair.Primary)
		}

		result, err := testutils.RunInitCluster(configFile)
		if err != nil {
			t.Fatalf("unexpected error: %s, %v", result.OutputMsg, err)
		}

		conn := dbconn.NewDBConnFromEnvironment("postgres")
		if err := conn.Connect(1); err != nil {
			t.Fatalf("Error connecting to the database: %v", err)
		}
		defer conn.Close()

		segConfigs, err := cluster.GetSegmentConfiguration(conn, false)
		if err != nil {
			t.Fatalf("Error getting segment configuration: %v", err)
		}

		resultSegs := make([]cli.Segment, len(segConfigs))
		for i, seg := range segConfigs {
			resultSegs[i] = cli.Segment{
				Hostname:      seg.Hostname,
				Port:          seg.Port,
				DataDirectory: seg.DataDir,
				Address:       seg.Hostname,
			}
		}

		if !reflect.DeepEqual(resultSegs, primarySegments) {
			t.Fatalf("got %+v, want %+v", resultSegs, primarySegments)
		}

		_, err = testutils.DeleteCluster()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("initialize cluster with default config and verify default values used correctly", func(t *testing.T) {
		var expectedOut string
		configFile := testutils.GetTempFile(t, "config.json")
		config := GetDefaultConfig(t)

		err := config.WriteConfigAs(configFile)
		if err != nil {
			t.Fatalf("unexpected error: %#v", err)
		}

		result, err := testutils.RunInitCluster(configFile)
		if err != nil {
			t.Fatalf("unexpected error: %s, %v", result.OutputMsg, err)
		}
		expectedOutput := result.OutputMsg

		expectedOut = "[INFO]:-Could not find encoding in cluster config, defaulting to UTF-8"
		if !strings.Contains(expectedOutput, expectedOut) {
			t.Fatalf("Output does not contain the expected string.\nExpected: %q\nGot: %q", expectedOut, expectedOutput)
		}

		expectedOut = "[INFO]:-Coordinator max_connections not set, will set to value 150 from CommonConfig"
		if !strings.Contains(expectedOutput, expectedOut) {
			t.Fatalf("Output does not contain the expected string.\nExpected: %q\nGot: %q", expectedOut, expectedOutput)
		}

		expectedOut = "[INFO]:-shared_buffers is not set in CommonConfig, will set to default value 128000kB"
		if !strings.Contains(expectedOutput, expectedOut) {
			t.Fatalf("Output does not contain the expected string.\nExpected: %q\nGot: %q", expectedOut, expectedOutput)
		}

		testutils.AssertPgConfig(t, "max_connections", "150", -1)
		testutils.AssertPgConfig(t, "shared_buffers", "125MB", -1)
		testutils.AssertPgConfig(t, "client_encoding", "UTF8", -1)

		_, err = testutils.DeleteCluster()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestCollations(t *testing.T) {
	t.Run("collations are imported successfully", func(t *testing.T) {
		configFile := testutils.GetTempFile(t, "config.json")
		config := GetDefaultConfig(t)

		err := config.WriteConfigAs(configFile)
		if err != nil {
			t.Fatalf("unexpected error: %#v", err)
		}

		result, err := testutils.RunInitCluster(configFile)
		if err != nil {
			t.Fatalf("unexpected error: %s, %v", result.OutputMsg, err)
		}

		expectedOut := "[INFO]:-Importing system collations"
		if !strings.Contains(result.OutputMsg, expectedOut) {
			t.Errorf("got %q, want %q", result.OutputMsg, expectedOut)
		}

		// before importing collations
		testutils.ExecQuery(t, "", "CREATE TABLE collationimport1 AS SELECT * FROM pg_collation WHERE collnamespace = 'pg_catalog'::regnamespace")

		// importing collations
		rows := testutils.ExecQuery(t, "", "SELECT pg_import_system_collations('pg_catalog')")
		testutils.AssertRowCount(t, rows, 1)

		// after importing collations
		testutils.ExecQuery(t, "", "CREATE TABLE collationimport2 AS SELECT * FROM pg_collation WHERE collnamespace = 'pg_catalog'::regnamespace")

		// there should be no difference before and after
		rows = testutils.ExecQuery(t, "", "SELECT * FROM collationimport1 EXCEPT SELECT * FROM collationimport2")
		testutils.AssertRowCount(t, rows, 0)

		_, err = testutils.DeleteCluster()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestDbCreationValidation(t *testing.T) {
	testDatabaseCreation := func(t *testing.T, dbName string) {
		configFile := testutils.GetTempFile(t, "config.json")
		config := GetDefaultConfig(t)

		err := config.WriteConfigAs(configFile)
		if err != nil {
			t.Fatalf("unexpected error: %#v", err)
		}

		SetConfigKey(t, configFile, "db-name", dbName, true)

		InitClusterResult, err := testutils.RunInitCluster(configFile)
		if err != nil {
			t.Fatalf("unexpected error: %s, %v", InitClusterResult.OutputMsg, err)
		}

		rows := testutils.ExecQuery(t, "", "SELECT datname FROM pg_database")
		defer rows.Close()
		foundDB := false
		for rows.Next() {
			var db string
			if err := rows.Scan(&db); err != nil {
				t.Fatalf("unexpected error scanning result: %v", err)
			}
			if db == dbName {
				foundDB = true
				break
			}
		}
		if !foundDB {
			t.Fatalf("Database %s should exist after creating it", dbName)
		}

		_, err = testutils.DeleteCluster()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	t.Run("validate cluster creation by specifying db name", func(t *testing.T) {
		testDatabaseCreation(t, "testdb")
	})

	t.Run("validate cluster creation by specifying hyphen in db name", func(t *testing.T) {
		testDatabaseCreation(t, "test-db")
	})

	t.Run("validate default databases creation when no db name is specified", func(t *testing.T) {
		configFile := testutils.GetTempFile(t, "config.json")
		config := GetDefaultConfig(t)

		err := config.WriteConfigAs(configFile)
		if err != nil {
			t.Fatalf("unexpected error: %#v", err)
		}

		result, err := testutils.RunInitCluster(configFile)
		if err != nil {
			t.Fatalf("unexpected error: %s, %v", result.OutputMsg, err)
		}

		foundDBs := make(map[string]bool)

		rows := testutils.ExecQuery(t, "", "SELECT datname from pg_database")
		for rows.Next() {
			var dbName string
			if err := rows.Scan(&dbName); err != nil {
				t.Fatalf("unexpected error scanning result: %v", err)
			}
			foundDBs[dbName] = true
		}

		expectedDBs := []string{"postgres", "template1", "template0"}
		for _, db := range expectedDBs {
			if !foundDBs[db] {
				t.Fatalf("Default database %s should exist after creating it", db)
			}
		}

		if len(foundDBs) != len(expectedDBs) {
			t.Fatalf("Unexpected number of databases found: expected %d, found %d", len(expectedDBs), len(foundDBs))
		}

		_, err = testutils.DeleteCluster()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

}

func TestGpToolKitValidation(t *testing.T) {
	t.Run("check if the gp_toolkit extension is created", func(t *testing.T) {
		configFile := testutils.GetTempFile(t, "config.json")
		config := GetDefaultConfig(t)

		err := config.WriteConfigAs(configFile)
		if err != nil {
			t.Fatalf("unexpected error: %#v", err)
		}

		result, err := testutils.RunInitCluster(configFile)
		if err != nil {
			t.Fatalf("unexpected error: %s, %v", result.OutputMsg, err)
		}

		QueryResult := testutils.ExecQuery(t, "", "select extname from pg_extension ")
		foundGpToolkit := false
		for QueryResult.Next() {
			var extName string
			err := QueryResult.Scan(&extName)
			if err != nil {
				t.Fatalf("unexpected error scanning result: %v", err)
			}
			if extName == "gp_toolkit" {
				foundGpToolkit = true
				break
			}
		}

		// Validate that "gp_toolkit" is present
		if !foundGpToolkit {
			t.Fatalf("Extension 'gp_toolkit' should exist in pg_extension")
		}

		_, err = testutils.DeleteCluster()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestPgHbaConfValidation(t *testing.T) {
	/* FIXME:concurse is failing to resolve ip to hostname*/
	/*t.Run("pghba config file validation when hbahostname is true", func(t *testing.T) {
		var value cli.Segment
		var ok bool
		var valueSeg []cli.Segment
		var okSeg bool
		configFile := testutils.GetTempFile(t, "config.json")
		config := GetDefaultConfig(t)

		err := config.WriteConfigAs(configFile)
		if err != nil {
			t.Fatalf("unexpected error: %#v", err)
		}

		SetConfigKey(t, configFile, "hba-hostnames", true, true)

		result, err := testutils.RunInitCluster(configFile)
		if err != nil {
			t.Fatalf("unexpected error: %s, %v", result.OutputMsg, err)
		}

		coordinator := config.Get("coordinator")
		if value, ok = coordinator.(cli.Segment); !ok {
			t.Fatalf("unexpected data type for coordinator %T", value)
		}

		filePathCord := filepath.Join(coordinator.(cli.Segment).DataDirectory, "pg_hba.conf")
		hostCord := coordinator.(cli.Segment).Hostname
		cmdStr := "whoami"
		cmd := exec.Command("ssh", hostCord, cmdStr)
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("unexpected error : %v", err)
		}

		resultCord := strings.TrimSpace(string(output))
		pgHbaLine := fmt.Sprintf("host\tall\t%s\t%s\ttrust", resultCord, coordinator.(cli.Segment).Hostname)
		cmdStrCord := fmt.Sprintf("/bin/bash -c 'cat %s | grep \"%s\"'", filePathCord, pgHbaLine)
		cmdCord := exec.Command("ssh", hostCord, cmdStrCord)
		_, err = cmdCord.CombinedOutput()
		if err != nil {
			t.Fatalf("unexpected error : %v", err)
		}

		primarySegs := config.Get("segment-array")
		valueSegPair, ok := primarySegs.([]cli.SegmentPair)

		if !ok {
			t.Fatalf("unexpected data type for segment-array %T", primarySegs)
		}

		pgHbaLineSeg := fmt.Sprintf("host\tall\tall\t%s\ttrust", primarySegs.([]cli.Segment)[0].Hostname)
		filePathSeg := filepath.Join(primarySegs.([]cli.Segment)[0].DataDirectory, "pg_hba.conf")
		cmdStr_seg := fmt.Sprintf("/bin/bash -c 'cat %s | grep \"%s\"'", filePathSeg, pgHbaLineSeg)
		hostSeg := primarySegs.([]cli.Segment)[0].Hostname
		cmdSeg := exec.Command("ssh", hostSeg, cmdStr_seg)
		_, err = cmdSeg.CombinedOutput()
		if err != nil {
			t.Fatalf("unexpected error : %v", err)
		}

		_, err = testutils.DeleteCluster()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})*/

	t.Run("pghba config file validation when hbahostname is false", func(t *testing.T) {
		var value cli.Segment
		var ok bool

		configFile := testutils.GetTempFile(t, "config.json")
		config := GetDefaultConfig(t)

		err := config.WriteConfigAs(configFile)
		if err != nil {
			t.Fatalf("unexpected error: %#v", err)
		}

		SetConfigKey(t, configFile, "hba-hostnames", false, true)

		result, err := testutils.RunInitCluster(configFile)
		if err != nil {
			t.Fatalf("unexpected error: %s, %v", result.OutputMsg, err)
		}

		coordinator := config.Get("coordinator")
		if value, ok = coordinator.(cli.Segment); !ok {
			t.Fatalf("unexpected data type for coordinator %T", value)
		}

		filePathCord := filepath.Join(coordinator.(cli.Segment).DataDirectory, "pg_hba.conf")
		hostCord := coordinator.(cli.Segment).Hostname
		cmdStr := "whoami"
		cmd := exec.Command("ssh", hostCord, cmdStr)
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("unexpected error : %v", err)
		}

		resultCord := strings.TrimSpace(string(output))
		cmdStrCord := "ip -4 addr show | grep inet | grep -v 127.0.0.1/8 | awk '{print $2}'"
		cmdCord := exec.Command("ssh", hostCord, cmdStrCord)
		outputCord, err := cmdCord.Output()
		if err != nil {
			t.Fatalf("unexpected error : %v", err)
		}

		resultCordValue := string(outputCord)
		firstCordValue := strings.Split(resultCordValue, "\n")[0]
		pgHbaLine := fmt.Sprintf("host\tall\t%s\t%s\ttrust", resultCord, firstCordValue)
		cmdStrCordValue := fmt.Sprintf("/bin/bash -c 'cat %s | grep \"%s\"'", filePathCord, pgHbaLine)
		cmdCordValue := exec.Command("ssh", hostCord, cmdStrCordValue)
		_, err = cmdCordValue.CombinedOutput()
		if err != nil {
			t.Fatalf("unexpected error : %v", err)
		}

		primarySegs := config.Get("segment-array")
		valueSegPair, ok := primarySegs.([]cli.SegmentPair)

		if !ok {
			t.Fatalf("unexpected data type for segment-array %T", primarySegs)
		}

		filePathSeg := filepath.Join(valueSegPair[0].Primary.DataDirectory, "pg_hba.conf")
		hostSegValue := valueSegPair[0].Primary.Hostname
		cmdStrSegValue := "whoami"
		cmdSegvalue := exec.Command("ssh", hostSegValue, cmdStrSegValue)
		outputSeg, errSeg := cmdSegvalue.Output()
		if errSeg != nil {
			t.Fatalf("unexpected error : %v", errSeg)
		}

		resultSeg := strings.TrimSpace(string(outputSeg))
		cmdStrSeg := "ip -4 addr show | grep inet | grep -v 127.0.0.1/8 | awk '{print $2}'"
		cmdSegValueNew := exec.Command("ssh", hostSegValue, cmdStrSeg)
		outputSegNew, err := cmdSegValueNew.Output()
		if err != nil {
			t.Fatalf("unexpected error : %v", err)
		}

		resultSegValue := string(outputSegNew)
		firstValueNew := strings.Split(resultSegValue, "\n")[0]
		pgHbaLineNew := fmt.Sprintf("host\tall\t%s\t%s\ttrust", resultSeg, firstValueNew)
		cmdStrSegNew := fmt.Sprintf("/bin/bash -c 'cat %s | grep \"%s\"'", filePathSeg, pgHbaLineNew)
		cmdSegNew := exec.Command("ssh", hostSegValue, cmdStrSegNew)
		_, err = cmdSegNew.CombinedOutput()
		if err != nil {
			t.Fatalf("unexpected error : %v", err)
		}

		_, err = testutils.DeleteCluster()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

}

// expansion cases:
func TestExpansionValidation(t *testing.T) {
	t.Run("check if primary ports are adjusted as per coordinator port in all the hosts when not specified", func(t *testing.T) {
		var value cli.Segment
		var ok bool
		configFile := testutils.GetTempFile(t, "config.json")
		config := GetDefaultConfig(t, true)

		err := config.WriteConfigAs(configFile)
		if err != nil {
			t.Fatalf("unexpected error: %#v", err)
		}

		coordinator := config.Get("coordinator")
		if value, ok = coordinator.(cli.Segment); !ok {
			t.Fatalf("unexpected data type for coordinator %T", value)
		}
		coordinatorPort := value.Port

		configMap := config.AllSettings()
		delete(configMap, "primary-base-port")

		encodedConfig, err := json.MarshalIndent(configMap, "", " ")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		err = os.WriteFile(configFile, encodedConfig, 0777)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		content, err := ioutil.ReadFile(configFile)
		if err != nil {
			fmt.Printf("Error reading configuration file: %v\n", err)
			return
		}

		// Print the content of the configuration file
		fmt.Println("Configuration file content:")
		fmt.Println(string(content))

		result, err := testutils.RunInitCluster(configFile)
		if err != nil {
			t.Fatalf("unexpected error: %s, %v", result.OutputMsg, err)
		}

		//below warning message needs to be corrected once bug is fixed
		expectedWarning := fmt.Sprintf("[WARNING]:-primary-base-port value not specified. Setting default to: %d", coordinatorPort+2)
		if !strings.Contains(result.OutputMsg, expectedWarning) {
			t.Fatalf("got %q, want %q", result.OutputMsg, expectedWarning)
		}

		expectedOut := "[INFO]:-Cluster initialized successfully"
		if !strings.Contains(result.OutputMsg, expectedOut) {
			t.Fatalf("got %q, want %q", result.OutputMsg, expectedOut)
		}

		conn := dbconn.NewDBConnFromEnvironment("postgres")
		if err := conn.Connect(1); err != nil {
			t.Fatalf("Error connecting to the database: %v", err)
		}
		defer conn.Close()

		segConfigs, err := cluster.GetSegmentConfiguration(conn, false)
		if err != nil {
			t.Fatalf("Error getting segment configuration: %v", err)
		}
		fmt.Println("segconfig")
		fmt.Println(segConfigs)

		var primarySegConfigs []cluster.SegConfig
		hostList := (config.GetStringSlice("hostlist"))

		for _, hostname := range hostList {
			for _, seg := range segConfigs {
				if seg.ContentID != -1 && seg.Role == "p" && seg.Hostname == hostname {
					primarySegConfigs = append(primarySegConfigs, seg)
				}
			}
			fmt.Printf("Printing primary segments for host %s:\n", hostname)
			fmt.Println(primarySegConfigs)

			for i, seg := range primarySegConfigs {
				fmt.Printf("%d. %s\n", i, seg.Hostname)
				expectedPrimaryPort := coordinatorPort + 2 + i
				fmt.Println("Expected primary port:", expectedPrimaryPort)

				if seg.Port != expectedPrimaryPort {
					t.Fatalf("Primary port mismatch for segment %s. Got: %d, Expected: %d", seg.Hostname, seg.Port, expectedPrimaryPort)
				}
			}

			// Clear primarySegConfigs for the next hostname
			primarySegConfigs = nil
		}

		_, err = testutils.DeleteCluster()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("check if mirror ports are adjusted as per coordinator port in all the hosts when not specified", func(t *testing.T) {
		var value cli.Segment
		var ok bool
		configFile := testutils.GetTempFile(t, "config.json")
		config := GetDefaultConfig(t, true)

		err := config.WriteConfigAs(configFile)
		if err != nil {
			t.Fatalf("unexpected error: %#v", err)
		}

		coordinator := config.Get("coordinator")
		if value, ok = coordinator.(cli.Segment); !ok {
			t.Fatalf("unexpected data type for coordinator %T", value)
		}
		coordinatorPort := value.Port

		configMap := config.AllSettings()
		delete(configMap, "mirror-base-port")

		encodedConfig, err := json.MarshalIndent(configMap, "", " ")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		err = os.WriteFile(configFile, encodedConfig, 0777)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		content, err := ioutil.ReadFile(configFile)
		if err != nil {
			fmt.Printf("Error reading configuration file: %v\n", err)
			return
		}

		// Print the content of the configuration file
		fmt.Println("Configuration file content:")
		fmt.Println(string(content))

		result, err := testutils.RunInitCluster(configFile)
		if err != nil {
			t.Fatalf("unexpected error: %s, %v", result.OutputMsg, err)
		}

		//below warning message needs to be corrected once bug is fixed
		expectedWarning := fmt.Sprintf("[WARNING]:-mirror-base-port value not specified. Setting default to: %d", coordinatorPort+1002)
		if !strings.Contains(result.OutputMsg, expectedWarning) {
			t.Fatalf("got %q, want %q", result.OutputMsg, expectedWarning)
		}

		expectedOut := "[INFO]:-Cluster initialized successfully"
		if !strings.Contains(result.OutputMsg, expectedOut) {
			t.Fatalf("got %q, want %q", result.OutputMsg, expectedOut)
		}

		conn := dbconn.NewDBConnFromEnvironment("postgres")
		if err := conn.Connect(1); err != nil {
			t.Fatalf("Error connecting to the database: %v", err)
		}
		defer conn.Close()

		segConfigs, err := cluster.GetSegmentConfiguration(conn, false, true)
		if err != nil {
			t.Fatalf("Error getting segment configuration: %v", err)
		}
		fmt.Println("segconfig")
		fmt.Println(segConfigs)

		var mirrorSegConfigs []cluster.SegConfig
		hostList := (config.GetStringSlice("hostlist"))

		for _, hostname := range hostList {
			for _, seg := range segConfigs {
				if seg.ContentID != -1 && seg.Role == "m" && seg.Hostname == hostname {
					mirrorSegConfigs = append(mirrorSegConfigs, seg)
				}
			}
			fmt.Printf("Printing mirror segments for host %s:\n", hostname)
			fmt.Println(mirrorSegConfigs)

			for i, seg := range mirrorSegConfigs {
				fmt.Printf("%d. %s\n", i, seg.Hostname)
				expectedMirrorPort := coordinatorPort + 1002 + i
				fmt.Println("Expected primary port:", expectedMirrorPort)

				if seg.Port != expectedMirrorPort {
					t.Fatalf("Mirror port mismatch for segment %s. Got: %d, Expected: %d", seg.Hostname, seg.Port, expectedMirrorPort)
				}
			}

			// Clear primarySegConfigs for the next hostname
			mirrorSegConfigs = nil
		}

		_, err = testutils.DeleteCluster()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	// t.Run("check if mirror ports are adjusted as per coordinator port when not specified", func(t *testing.T) {
	// 	var value cli.Segment
	// 	var ok bool
	// 	configFile := testutils.GetTempFile(t, "config.json")
	// 	config := GetDefaultConfig(t, true)

	// 	err := config.WriteConfigAs(configFile)
	// 	if err != nil {
	// 		t.Fatalf("unexpected error: %#v", err)
	// 	}

	// 	coordinator := config.Get("coordinator")
	// 	if value, ok = coordinator.(cli.Segment); !ok {
	// 		t.Fatalf("unexpected data type for coordinator %T", value)
	// 	}
	// 	coordinatorPort := value.Port

	// 	configMap := config.AllSettings()
	// 	delete(configMap, "mirror-base-port")

	// 	encodedConfig, err := json.MarshalIndent(configMap, "", " ")
	// 	if err != nil {
	// 		t.Fatalf("unexpected error: %v", err)
	// 	}

	// 	err = os.WriteFile(configFile, encodedConfig, 0777)
	// 	if err != nil {
	// 		t.Fatalf("unexpected error: %v", err)
	// 	}

	// 	result, err := testutils.RunInitCluster(configFile)
	// 	if err != nil {
	// 		t.Fatalf("unexpected error: %s, %v", result.OutputMsg, err)
	// 	}

	// 	//below warning message needs to be corrected once bug is fixed
	// 	expectedWarning := fmt.Sprintf("[WARNING]:-No mirror-base-port value provided. Setting default to: %d", coordinatorPort+1002)
	// 	if !strings.Contains(result.OutputMsg, expectedWarning) {
	// 		t.Fatalf("got %q, want %q", result.OutputMsg, expectedWarning)
	// 	}

	// 	expectedOut := "[INFO]:-Cluster initialized successfully"
	// 	if !strings.Contains(result.OutputMsg, expectedOut) {
	// 		t.Fatalf("got %q, want %q", result.OutputMsg, expectedOut)
	// 	}

	// 	conn := dbconn.NewDBConnFromEnvironment("postgres")
	// 	if err := conn.Connect(1); err != nil {
	// 		t.Fatalf("Error connecting to the database: %v", err)
	// 	}
	// 	defer conn.Close()

	// 	segConfigs, err := cluster.GetSegmentConfiguration(conn, false, true)
	// 	if err != nil {
	// 		t.Fatalf("Error getting segment configuration: %v", err)
	// 	}

	// 	var mirrorSegConfigs []cluster.SegConfig
	// 	for _, seg := range segConfigs {
	// 		if seg.Port != coordinatorPort {
	// 			mirrorSegConfigs = append(mirrorSegConfigs, seg)
	// 		}
	// 	}

	// 	for i, seg := range mirrorSegConfigs {
	// 		// Calculate the expected primary port based on the coordinator port
	// 		expectedMirrorPort := coordinatorPort + 1002 + i

	// 		// Check if the primary port matches the expected value
	// 		if seg.Port != expectedMirrorPort {
	// 			t.Fatalf("Mirror port mismatch for segment %s. Got: %d, Expected: %d", seg.Hostname, seg.Port, expectedMirrorPort)
	// 		}
	// 	}

	// 	_, err = testutils.DeleteCluster()
	// 	if err != nil {
	// 		t.Fatalf("unexpected error: %v", err)
	// 	}
	// })

	t.Run("verify expansion by initialize cluster with default config and verify default values used correctly", func(t *testing.T) {
		var expectedOut string
		configFile := testutils.GetTempFile(t, "config.json")
		config := GetDefaultConfig(t, true)

		err := config.WriteConfigAs(configFile)
		if err != nil {
			t.Fatalf("unexpected error: %#v", err)
		}

		result, err := testutils.RunInitCluster(configFile)
		if err != nil {
			t.Fatalf("unexpected error: %s, %v", result.OutputMsg, err)
		}
		expectedOutput := result.OutputMsg

		expectedOut = "[INFO]:-Could not find encoding in cluster config, defaulting to UTF-8"
		if !strings.Contains(expectedOutput, expectedOut) {
			t.Fatalf("Output does not contain the expected string.\nExpected: %q\nGot: %q", expectedOut, expectedOutput)
		}

		expectedOut = "[INFO]:-Coordinator max_connections not set, will set to value 150 from CommonConfig"
		if !strings.Contains(expectedOutput, expectedOut) {
			t.Fatalf("Output does not contain the expected string.\nExpected: %q\nGot: %q", expectedOut, expectedOutput)
		}

		expectedOut = "[INFO]:-shared_buffers is not set in CommonConfig, will set to default value 128000kB"
		if !strings.Contains(expectedOutput, expectedOut) {
			t.Fatalf("Output does not contain the expected string.\nExpected: %q\nGot: %q", expectedOut, expectedOutput)
		}

		testutils.AssertPgConfig(t, "max_connections", "150", -1)
		testutils.AssertPgConfig(t, "shared_buffers", "125MB", -1)
		testutils.AssertPgConfig(t, "client_encoding", "UTF8", -1)

		_, err = testutils.DeleteCluster()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("validate expansion that proper number of primary and mirror directories are created in each hosts", func(t *testing.T) {
		configFile := testutils.GetTempFile(t, "config.json")
		config := GetDefaultConfig(t, true)

		err := config.WriteConfigAs(configFile)
		if err != nil {
			t.Fatalf("unexpected error: %#v", err)
		}
		primaryDirs := len(config.GetStringSlice("primary-data-directories"))
		mirrorDirs := len(config.GetStringSlice("mirror-data-directories"))

		hostList := len(config.GetStringSlice("hostlist"))

		fmt.Println("primary dir")
		fmt.Println(primaryDirs)

		fmt.Println("mirror dir")
		fmt.Println(mirrorDirs)

		fmt.Println("hostlist dir")
		fmt.Println(hostList)

		//validate the no of dd should be host* no of host

		//uncomment below lines once code is ready
		result, err := testutils.RunInitCluster(configFile)
		if err != nil {
			t.Fatalf("unexpected error: %s, %v", result.OutputMsg, err)
		}

		expectedOut := "[INFO]:-Cluster initialized successfully"
		if !strings.Contains(result.OutputMsg, expectedOut) {
			t.Fatalf("got %q, want %q", result.OutputMsg, expectedOut)
		}

		conn := dbconn.NewDBConnFromEnvironment("postgres")
		if err := conn.Connect(1); err != nil {
			t.Fatalf("Error connecting to the database: %v", err)
		}
		defer conn.Close()

		segConfigs, err := cluster.GetSegmentConfiguration(conn, true)
		if err != nil {
			t.Fatalf("Error getting segment configuration: %v", err)
		}
		fmt.Printf("all primary segs")
		fmt.Println(segConfigs)

		var primaryDataDirs []string
		var mirrorDataDirs []string
		//hosts := make(map[string]bool) // Map to store unique hosts

		for _, seg := range segConfigs {
			if seg.ContentID == -1 {
				// Skip appending coordinator directories
				continue
			} else if seg.Role == "p" {
				primaryDataDirs = append(primaryDataDirs, seg.DataDir)
				fmt.Println("primary DD")
				fmt.Println(primaryDataDirs)
				//hosts[seg.Hostname] = true

			} else if seg.Role == "m" {
				mirrorDataDirs = append(mirrorDataDirs, seg.DataDir)
				//hosts[seg.Hostname] = true

			}
		}

		primaryCount := len(primaryDataDirs)
		mirrorCount := len(mirrorDataDirs)

		fmt.Printf("Primary Count: %d\n", primaryCount)
		fmt.Printf("Mirror Count: %d\n", mirrorCount)

		actualPrimaryCount := len(primaryDataDirs)
		actualMirrorCount := len(mirrorDataDirs)

		expectedPrimaryCount := primaryDirs * hostList
		expectedMirrorCount := mirrorDirs * hostList

		fmt.Println("actual primary count")
		fmt.Println(actualPrimaryCount)

		fmt.Println("expectedPrimaryCount")
		fmt.Println(expectedPrimaryCount)
		if actualPrimaryCount != expectedPrimaryCount {
			t.Fatalf("Error: Primary data directories count mismatch: expected %d, got %d", expectedPrimaryCount, primaryCount)
		}

		if actualMirrorCount != expectedMirrorCount {
			t.Fatalf("Error: Mirror data directories count mismatch: expected %d, got %d", expectedMirrorCount, mirrorCount)
		}

		_, err = testutils.DeleteCluster()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("validate the group mirroring and check if segments are distributed properly across hosts", func(t *testing.T) {
		fmt.Println("hostlist count")
		fmt.Println(len(hostList))
		if len(hostList) == 1 {
			t.Skip()
		}
		var value cli.Segment
		var ok bool

		configFile := testutils.GetTempFile(t, "config.json")
		config := GetDefaultConfig(t, true)

		err := config.WriteConfigAs(configFile)
		if err != nil {
			t.Fatalf("unexpected error: %#v", err)
		}

		coordinator := config.Get("coordinator")
		if value, ok = coordinator.(cli.Segment); !ok {
			t.Fatalf("unexpected data type for coordinator %T", value)
		}

		// uncomment below lines once code is ready
		result, err := testutils.RunInitCluster(configFile)
		if err != nil {
			t.Fatalf("unexpected error: %s, %v", result.OutputMsg, err)
		}

		conn := dbconn.NewDBConnFromEnvironment("postgres")
		if err := conn.Connect(1); err != nil {
			t.Fatalf("Error connecting to the database: %v", err)
		}
		defer conn.Close()

		segConfigs, err := cluster.GetSegmentConfiguration(conn, true)
		if err != nil {
			t.Fatalf("Error getting segment configuration: %v", err)
		}
		fmt.Printf("all primary segs")
		fmt.Println(segConfigs)

		// // Step 1: Fetch all primaries and their hosts
		// hostList := (config.GetStringSlice("hostlist"))
		// fmt.Println("hostlist")
		// fmt.Println(hostList)

		// primaries := make(map[int][]cluster.SegConfig)
		// for _, hostname := range hostList {
		// 	fmt.Printf("hostname prinitng %s\n", hostname)
		// 	for _, seg := range segConfigs {
		// 		if seg.ContentID != -1 && seg.Role == "p" && seg.Hostname == hostname {
		// 			fmt.Println("inside if 1")
		// 			primaries[seg.ContentID] = append(primaries[seg.ContentID], seg)
		// 		}
		// 	}
		// }

		// fmt.Println("all primaries")
		// fmt.Println(primaries)

		// // Step 2: Fetch corresponding mirrors and their hosts
		// mirrors := make(map[int][]cluster.SegConfig)
		// for _, seg := range segConfigs {
		// 	if seg.Role == "m" {
		// 		if primary, ok := primaries[seg.ContentID]; ok {
		// 			mirrors[primary[0].ContentID] = append(mirrors[primary[0].ContentID], seg)
		// 		}
		// 	}
		// }

		// fmt.Println("all mirrors")
		// fmt.Println(mirrors)

		// // Step 3: Validate spread and group mirroring
		// var mirrorHostnames []string
		// seen := make(map[string]bool)
		// var primaryHostnames []string

		// for _, configs := range mirrors {
		// 	for _, config := range configs {
		// 		mirrorHostnames = append(mirrorHostnames, config.Hostname)
		// 		seen[config.Hostname] = true
		// 	}
		// }

		// fmt.Println("seen")
		// fmt.Println(seen)

		// for _, configs := range primaries {
		// 	for _, config := range configs {
		// 		primaryHostnames = append(primaryHostnames, config.Hostname)
		// 	}
		// }

		// fmt.Println("primary hostname")
		// fmt.Println(primaryHostnames)

		// for _, mirrorHostname := range mirrorHostnames {
		// 	for _, primaryHostname := range primaryHostnames {
		// 		if mirrorHostname == primaryHostname {
		// 			t.Fatalf("Error: Mirrors are hosted on the same host as primary: %s", mirrorHostname)
		// 		}
		// 	}
		// }

		// if len(seen) > 1 {
		// 	t.Fatalf("Error: Group mirroring validation Failed: All hostnames are not same for mirrors: %v", mirrorHostnames)
		// }
		hostname := config.Get("hostlist").([]string)[0]
		primaries := make(map[int][]cluster.SegConfig)
		for _, seg := range segConfigs {
			if seg.ContentID != -1 && seg.Role == "p" && seg.Hostname == hostname { //include code to skip conetent id -1
				primaries[seg.ContentID] = append(primaries[seg.ContentID], seg)
			}
		}
		fmt.Println("all primaries")
		fmt.Println(primaries)

		// Step 3: Fetch corresponding mirrors and their hosts
		mirrors := make(map[int][]cluster.SegConfig)
		for _, seg := range segConfigs {
			if seg.Role == "m" {
				if primary, ok := primaries[seg.ContentID]; ok {
					mirrors[primary[0].ContentID] = append(mirrors[primary[0].ContentID], seg)
				}
			}
		}

		fmt.Println("all mirrors")
		fmt.Println(mirrors)
		//uncomment below once code is ready not need becz by default it iakes spread

		// Step 4: Validate spread and group mirroring
		var mirrorHostnames []string
		seen := make(map[string]bool)
		var primaryHostnames []string

		for _, configs := range mirrors {
			for _, config := range configs {
				mirrorHostnames = append(mirrorHostnames, config.Hostname)
				seen[config.Hostname] = true
			}
		}

		fmt.Println("seen")
		fmt.Println(seen)

		for _, configs := range primaries {
			for _, config := range configs {
				primaryHostnames = append(primaryHostnames, config.Hostname)
			}
		}

		fmt.Println("primary hostname")
		fmt.Println(primaryHostnames)

		for _, mirrorHostname := range mirrorHostnames {
			for _, primaryHostname := range primaryHostnames {
				if mirrorHostname == primaryHostname {
					t.Fatalf("Error: Mirrors are hosted on the same host as primary: %s", mirrorHostname)
				}
			}
		}

		if len(seen) > 1 {
			t.Fatalf("Error: Group mirroing validation Failed: All hostnames are not same for mirrors: %v", mirrorHostnames)
		}

		_, err = testutils.DeleteCluster()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

	})

	t.Run("validate the spread mirroring and check if segments are distributed properly across hosts", func(t *testing.T) {
		if len(hostList) == 1 {
			t.Skip()
		}
		var value cli.Segment
		var ok bool

		configFile := testutils.GetTempFile(t, "config.json")
		config := GetDefaultConfig(t, true)

		err := config.WriteConfigAs(configFile)
		if err != nil {
			t.Fatalf("unexpected error: %#v", err)
		}

		coordinator := config.Get("coordinator")
		if value, ok = coordinator.(cli.Segment); !ok {
			t.Fatalf("unexpected data type for coordinator %T", value)
		}

		config.Set("mirroring-type", "spread")
		if err := config.WriteConfigAs(configFile); err != nil {
			t.Fatalf("failed to write config to file: %v", err)
		}

		// uncomment below lines once code is ready
		result, err := testutils.RunInitCluster(configFile)
		if err != nil {
			t.Fatalf("unexpected error: %s, %v", result.OutputMsg, err)
		}

		conn := dbconn.NewDBConnFromEnvironment("postgres")
		if err := conn.Connect(1); err != nil {
			t.Fatalf("Error connecting to the database: %v", err)
		}
		defer conn.Close()

		segConfigs, err := cluster.GetSegmentConfiguration(conn, true)
		if err != nil {
			t.Fatalf("Error getting segment configuration: %v", err)
		}
		fmt.Printf("all primary segs")
		fmt.Println(segConfigs)

		//decalare variable hostname and get it from config hostlist may be hostlist[0] ->u will get sdw0 and remove localhost and replace with hostna,e
		hostname := config.Get("hostlist").([]string)[0]
		primaries := make(map[int][]cluster.SegConfig)
		for _, seg := range segConfigs {
			if seg.ContentID != -1 && seg.Role == "p" && seg.Hostname == hostname { //include code to skip conetent id -1
				primaries[seg.ContentID] = append(primaries[seg.ContentID], seg)
			}
		}

		fmt.Println("all primaries")
		fmt.Println(primaries)

		// Step 3: Fetch corresponding mirrors and their hosts
		mirrors := make(map[int][]cluster.SegConfig)
		for _, seg := range segConfigs {
			if seg.Role == "m" {
				if primary, ok := primaries[seg.ContentID]; ok {
					mirrors[primary[0].ContentID] = append(mirrors[primary[0].ContentID], seg)
				}
			}
		}

		fmt.Println("all mirrors")
		fmt.Println(mirrors)

		//uncomment below once code is ready // may be not needed
		//mirroringType := config.Get("mirroring-type")
		//mirroringType := "group"

		// Step 4: Validate spread and group mirroring
		var mirrorHostnames []string
		seen := make(map[string]bool)
		var primaryHostnames []string

		for _, configs := range mirrors {
			for _, config := range configs {
				mirrorHostnames = append(mirrorHostnames, config.Hostname)
				seen[config.Hostname] = true
			}
		}

		for _, configs := range primaries {
			for _, config := range configs {
				primaryHostnames = append(primaryHostnames, config.Hostname)
			}
		}

		for _, mirrorHostname := range mirrorHostnames {
			for _, primaryHostname := range primaryHostnames {
				if mirrorHostname == primaryHostname {
					t.Fatalf("Error: Mirrors are hosted on the same host as primary: %s", mirrorHostname)
				}
			}
		}

		if len(seen) != len(mirrorHostnames) {
			t.Fatalf("Error: Spread mirroing Validation Failed, Hostnames are not different. %v", mirrorHostnames)
		}

		_, err = testutils.DeleteCluster()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

	})

}
