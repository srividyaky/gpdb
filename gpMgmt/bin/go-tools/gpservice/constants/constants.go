package constants

import "time"

// gpservice configuration specific constants
const (
	DefaultHubPort     = 4242
	DefaultAgentPort   = 8000
	DefaultServiceName = "gpservice"
	ConfigFileName     = "gpservice.conf"
	PlatformDarwin     = "darwin"
	PlatformLinux      = "linux"
	DefaultGpCtlName   = "gpctl"
)

const (
	ShellPath               = "/bin/bash"
	GpSSH                   = "gpssh"
	MaxRetries              = 10
	DefaultQdMaxConnect     = 150
	QeConnectFactor         = 3
	DefaultBuffer           = "128000kB"
	OsOpenFiles             = 65535
	DefaultDatabase         = "template1"
	DefaultEncoding         = "UTF-8"
	EtcHostsFilepath        = "/etc/hosts"
	CleanFileName           = "ClusterInitCLeanup.txt"
	ReplicationSlotName     = "internal_wal_replication_slot"
	DefaultStartTimeout     = 600
	DefaultPostgresLogDir   = "log"
	GroupMirroring          = "group"
	SpreadMirroring         = "spread"
	DefaultSegName          = "gpseg"
	UserInputWaitDurtion    = 30
	CheckInterruptFrequency = 500 * time.Millisecond
)

// gp_segment_configuration specific constants
const (
	RolePrimary = "p"
	RoleMirror  = "m"
)

// Catalog tables
const (
	GpSegmentConfiguration = "gp_segment_configuration"
)
