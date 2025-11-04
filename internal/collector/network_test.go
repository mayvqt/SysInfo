package collector

import (
	"testing"
)

// TestCollectNetwork verifies basic network collection works
func TestCollectNetwork(t *testing.T) {
	data, err := CollectNetwork()
	if err != nil {
		t.Fatalf("CollectNetwork failed: %v", err)
	}

	if data == nil {
		t.Fatal("CollectNetwork returned nil data")
	}

	// Every system should have at least one network interface (even loopback)
	if len(data.Interfaces) == 0 {
		t.Error("No network interfaces found (expected at least loopback)")
	}

	// Verify interface data
	for i, iface := range data.Interfaces {
		if iface.Name == "" {
			t.Errorf("Interface[%d] has empty name", i)
		}

		// MTU should be positive
		if iface.MTU <= 0 {
			t.Logf("Warning: Interface[%d] (%s) has MTU <= 0: %d", i, iface.Name, iface.MTU)
		}

		// Counters are uint64 and can't be negative, just log them

		t.Logf("Interface: %s, MAC: %s, Addrs: %v, MTU: %d",
			iface.Name, iface.HardwareAddr, iface.Addresses, iface.MTU)

		if len(iface.Flags) > 0 {
			t.Logf("  Flags: %v", iface.Flags)
		}

		if iface.BytesSent > 0 || iface.BytesRecv > 0 {
			t.Logf("  Traffic: Sent=%d, Recv=%d", iface.BytesSent, iface.BytesRecv)
		}
	}

	// Connection count may be available
	if data.Connections > 0 {
		t.Logf("Active connections: %d", data.Connections)
	}
}

func TestCollectNetworkHasLoopback(t *testing.T) {
	data, err := CollectNetwork()
	if err != nil {
		t.Fatalf("CollectNetwork failed: %v", err)
	}

	// Check for loopback interface (common names: lo, lo0, Loopback)
	hasLoopback := false
	for _, iface := range data.Interfaces {
		if iface.Name == "lo" || iface.Name == "lo0" ||
			iface.Name == "Loopback Pseudo-Interface 1" ||
			containsFlag(iface.Flags, "loopback") {
			hasLoopback = true
			t.Logf("Found loopback interface: %s", iface.Name)
			break
		}
	}

	if !hasLoopback {
		t.Log("Warning: No obvious loopback interface found (may use different naming)")
	}
}

func TestCollectNetworkAddresses(t *testing.T) {
	data, err := CollectNetwork()
	if err != nil {
		t.Fatalf("CollectNetwork failed: %v", err)
	}

	foundAddress := false
	for _, iface := range data.Interfaces {
		if len(iface.Addresses) > 0 {
			foundAddress = true
			for _, addr := range iface.Addresses {
				if addr == "" {
					t.Errorf("Interface %s has empty address", iface.Name)
				}
			}
		}
	}

	if !foundAddress {
		t.Log("Warning: No interface has addresses (unusual but possible)")
	}
}

func BenchmarkCollectNetwork(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = CollectNetwork()
	}
}

// Helper function
func containsFlag(flags []string, target string) bool {
	for _, flag := range flags {
		if flag == target {
			return true
		}
	}
	return false
}
