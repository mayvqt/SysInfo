package collector

import (
	"fmt"
	"net"

	"github.com/mayvqt/sysinfo/src/internal/types"
	psnet "github.com/shirou/gopsutil/v3/net"
)

// CollectNetwork gathers network interface information
func CollectNetwork() (*types.NetworkData, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get network interfaces: %w", err)
	}

	data := &types.NetworkData{
		Interfaces: make([]types.NetworkInterface, 0),
	}

	// Get I/O counters
	ioCounters, _ := psnet.IOCounters(true)
	ioMap := make(map[string]psnet.IOCountersStat)
	for _, io := range ioCounters {
		ioMap[io.Name] = io
	}

	for _, iface := range interfaces {
		addrs, _ := iface.Addrs()
		addrStrings := make([]string, 0)
		for _, addr := range addrs {
			addrStrings = append(addrStrings, addr.String())
		}

		flags := make([]string, 0)
		if iface.Flags&net.FlagUp != 0 {
			flags = append(flags, "UP")
		}
		if iface.Flags&net.FlagBroadcast != 0 {
			flags = append(flags, "BROADCAST")
		}
		if iface.Flags&net.FlagLoopback != 0 {
			flags = append(flags, "LOOPBACK")
		}
		if iface.Flags&net.FlagMulticast != 0 {
			flags = append(flags, "MULTICAST")
		}

		netInterface := types.NetworkInterface{
			Name:         iface.Name,
			HardwareAddr: iface.HardwareAddr.String(),
			Addresses:    addrStrings,
			Flags:        flags,
			MTU:          iface.MTU,
		}

		// Add I/O statistics if available
		if io, ok := ioMap[iface.Name]; ok {
			netInterface.BytesSent = io.BytesSent
			netInterface.BytesRecv = io.BytesRecv
			netInterface.PacketsSent = io.PacketsSent
			netInterface.PacketsRecv = io.PacketsRecv
			netInterface.ErrorsIn = io.Errin
			netInterface.ErrorsOut = io.Errout
			netInterface.DropsIn = io.Dropin
			netInterface.DropsOut = io.Dropout
		}

		data.Interfaces = append(data.Interfaces, netInterface)
	}

	// Get connection count
	connections, err := psnet.Connections("all")
	if err == nil {
		data.Connections = len(connections)
	}

	return data, nil
}
