package agent

import (
	"context"
	"fmt"

	"github.com/greenplum-db/gpdb/gp/idl"
	"github.com/greenplum-db/gpdb/gp/utils"
	"github.com/greenplum-db/gpdb/gp/utils/postgres"
)

func (s *Server) UpdatePgHbaConf(ctx context.Context, req *idl.UpdatePgHbaConfRequest) (*idl.UpdatePgHbaConfResponse, error) {
	err := postgres.UpdateSegmentPgHbaConf(req.Pgdata, req.Addrs, req.Replication)
	if err != nil {
		return &idl.UpdatePgHbaConfResponse{}, fmt.Errorf("updating pg_hba.conf: %w", err)
	}

	pgCtlReloadCmd := &postgres.PgCtlReload{
		PgData: req.Pgdata,
	}
	out, err := utils.RunGpCommand(pgCtlReloadCmd, s.GpHome)
	if err != nil {
		return &idl.UpdatePgHbaConfResponse{}, fmt.Errorf("executing pg_ctl reload: %s, %w", out, err)
	}

	return &idl.UpdatePgHbaConfResponse{}, nil
}
