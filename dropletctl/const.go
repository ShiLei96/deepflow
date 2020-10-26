package dropletctl

import "gitlab.x.lan/yunshan/droplet-libs/debug"

const (
	DEBUG_LISTEN_IP   = "127.0.0.1"
	DEBUG_LISTEN_PORT = 9527
)

const (
	DROPLETCTL_ADAPTER debug.ModuleId = iota
	DROPLETCTL_QUEUE
	DROPLETCTL_LABELER
	DROPLETCTL_RPC
	DROPLETCTL_LOGLEVEL
	DROPLETCTL_CONFIG
	DROPLETCTL_ROZE_QUEUE
	DROPLETCTL_STREAM_QUEUE

	DROPLETCTL_MAX
)

const (
	DEBUG_MESSAGE_LEN = 4096
)

var ConfigPath string
