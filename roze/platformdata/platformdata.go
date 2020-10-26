package platformdata

import (
	"net"

	"gitlab.x.lan/yunshan/droplet-libs/debug"
	"gitlab.x.lan/yunshan/droplet-libs/grpc"
	"gitlab.x.lan/yunshan/droplet-libs/receiver"
)

var PlatformData *grpc.PlatformInfoTable

const (
	CMD_PLATFORMDATA = 33
)

func New(ips []net.IP, port int, processName string, receiver *receiver.Receiver) {
	PlatformData = grpc.NewPlatformInfoTable(ips, port, processName, receiver)
	debug.ServerRegisterSimple(CMD_PLATFORMDATA, PlatformData)
}

func Start() {
	PlatformData.Start()
}

func Close() {
	PlatformData.Close()
}
