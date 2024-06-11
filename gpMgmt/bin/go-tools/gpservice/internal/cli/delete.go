package cli

import (
	"github.com/greenplum-db/gp-common-go-libs/gplog"
	_ "github.com/greenplum-db/gpdb/gpservice/idl"
	"github.com/greenplum-db/gpdb/gpservice/pkg/gpservice_config"
	"github.com/spf13/cobra"
)

func DeleteCmd() *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "de-register services and cleanup related files",
		RunE: func(cmd *cobra.Command, args []string) error {
			return DeleteServices(conf)
		},
	}
	return deleteCmd
}

func DeleteServices(conf *gpservice_config.Config) error {
	// stop services if running, else ignore the error
	err := StopServices(conf)
	if err != nil {
		gplog.Debug("stop services returned error:%v. Probably already stopped", err)
	}

	err = platform.RemoveHubService(conf.ServiceName)
	if err != nil {
		return err
	}
	gplog.Info("Successfully deleted hub service")

	err = platform.RemoveAgentService(conf.GpHome, conf.ServiceName, conf.Hostnames)
	if err != nil {
		return err
	}
	gplog.Info("Successfully deleted agent service")

	err = platform.RemoveHubServiceFile(conf.ServiceName)
	if err != nil {
		return err
	}

	err = platform.RemoveAgentServiceFile(conf.GpHome, conf.ServiceName, conf.Hostnames)
	if err != nil {
		return err
	}

	// remove gpservice.conf from hub and agents using ssh
	err = conf.Remove(configFilepath)
	if err != nil {
		return err
	}

	return nil
}
