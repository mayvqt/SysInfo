//go:build windows

package collector

import (
	"fmt"
	"strings"

	"github.com/mayvqt/sysinfo/internal/types"
	"github.com/mayvqt/sysinfo/internal/utils"
	"github.com/yusufpapurcu/wmi"
)

// Win32DiskDrive represents Windows WMI disk drive information
type Win32DiskDrive struct {
	Caption           string
	DeviceID          string
	Model             string
	SerialNumber      string
	Size              uint64
	MediaType         string
	InterfaceType     string
	Partitions        uint32
	BytesPerSector    uint32
	SectorsPerTrack   uint32
	TotalCylinders    uint64
	TotalHeads        uint32
	TotalSectors      uint64
	TotalTracks       uint64
	TracksPerCylinder uint32
	Status            string
	Index             uint32
	FirmwareRevision  string
}

// MSFT_PhysicalDisk represents Windows Storage Spaces physical disk information (Windows 8+)
type MSFT_PhysicalDisk struct {
	FriendlyName      string
	SerialNumber      string
	MediaType         uint16 // 0=Unspecified, 3=HDD, 4=SSD, 5=SCM
	BusType           uint16 // 0=Unknown, 1=SCSI, 2=ATAPI, 3=ATA, 4=1394, 5=SSA, 6=FC, 7=USB, 8=RAID, 9=iSCSI, 10=SAS, 11=SATA, 12=SD, 13=MMC, 14=Virtual, 15=FileBackedVirtual, 16=Storage Spaces, 17=NVMe, 18=SCM
	CanPool           bool
	Size              uint64
	AllocatedSize     uint64
	OperationalStatus uint16
	HealthStatus      uint16
	Usage             uint16
	SpindleSpeed      uint32 // RPM
}

func collectPhysicalDisksPlatform() []types.PhysicalDisk {
	// Try modern MSFT_PhysicalDisk first (Windows 8+)
	modernDisks := collectDisksMSFT()
	if len(modernDisks) > 0 {
		return modernDisks
	}

	// Fallback to Win32_DiskDrive (compatible with older Windows)
	return collectDisksWMI()
}

// collectDisksMSFT uses MSFT_PhysicalDisk for modern Windows systems
func collectDisksMSFT() []types.PhysicalDisk {
	var wmiDisks []MSFT_PhysicalDisk
	query := "SELECT * FROM MSFT_PhysicalDisk"

	// Query root\Microsoft\Windows\Storage namespace
	err := wmi.QueryNamespace(query, &wmiDisks, `root\Microsoft\Windows\Storage`)
	if err != nil {
		// MSFT_PhysicalDisk not available (pre-Windows 8 or access denied)
		return nil
	}

	disks := make([]types.PhysicalDisk, 0)
	for _, wmiDisk := range wmiDisks {
		disk := types.PhysicalDisk{
			Name:          wmiDisk.FriendlyName,
			Model:         wmiDisk.FriendlyName,
			SerialNumber:  strings.TrimSpace(wmiDisk.SerialNumber),
			Size:          wmiDisk.Size,
			SizeFormatted: utils.FormatBytes(wmiDisk.Size),
		}

		// Map MediaType
		switch wmiDisk.MediaType {
		case 3:
			disk.Type = "HDD"
			if wmiDisk.SpindleSpeed > 0 {
				disk.RPM = wmiDisk.SpindleSpeed
			}
		case 4:
			disk.Type = "SSD"
		case 5:
			disk.Type = "SCM" // Storage Class Memory
		default:
			disk.Type = "Unknown"
		}

		// Map BusType to Interface
		disk.Interface = mapBusType(wmiDisk.BusType)

		// NVMe detection
		if wmiDisk.BusType == 17 {
			disk.Type = "NVMe"
			disk.Interface = "NVMe"
		}

		// Removable media check (USB, SD, MMC)
		if wmiDisk.BusType == 7 || wmiDisk.BusType == 12 || wmiDisk.BusType == 13 {
			disk.Removable = true
		}

		disks = append(disks, disk)
	}

	return disks
}

// collectDisksWMI uses Win32_DiskDrive for older Windows systems
func collectDisksWMI() []types.PhysicalDisk {
	var wmiDisks []Win32DiskDrive
	query := "SELECT * FROM Win32_DiskDrive"

	err := wmi.Query(query, &wmiDisks)
	if err != nil {
		return nil
	}

	disks := make([]types.PhysicalDisk, 0)
	for _, wmiDisk := range wmiDisks {
		disk := types.PhysicalDisk{
			Name:          wmiDisk.DeviceID,
			Model:         strings.TrimSpace(wmiDisk.Model),
			SerialNumber:  strings.TrimSpace(wmiDisk.SerialNumber),
			Size:          wmiDisk.Size,
			SizeFormatted: utils.FormatBytes(wmiDisk.Size),
			Interface:     strings.ToUpper(wmiDisk.InterfaceType),
		}

		// Try to determine disk type from media type or model
		mediaType := strings.ToLower(wmiDisk.MediaType)
		model := strings.ToLower(wmiDisk.Model)

		if strings.Contains(mediaType, "ssd") || strings.Contains(model, "ssd") || strings.Contains(model, "solid state") {
			disk.Type = "SSD"
		} else if strings.Contains(model, "nvme") {
			disk.Type = "NVMe"
			disk.Interface = "NVMe"
		} else if strings.Contains(mediaType, "fixed") || strings.Contains(mediaType, "hard disk") {
			disk.Type = "HDD"
		} else if strings.Contains(mediaType, "removable") {
			disk.Removable = true
			disk.Type = "Removable"
		}

		// Detect removable from interface
		if strings.Contains(strings.ToLower(wmiDisk.InterfaceType), "usb") {
			disk.Removable = true
		}

		disks = append(disks, disk)
	}

	return disks
}

// mapBusType converts Windows bus type code to interface name
func mapBusType(busType uint16) string {
	busTypes := map[uint16]string{
		0:  "Unknown",
		1:  "SCSI",
		2:  "ATAPI",
		3:  "ATA",
		4:  "1394",
		5:  "SSA",
		6:  "FC",
		7:  "USB",
		8:  "RAID",
		9:  "iSCSI",
		10: "SAS",
		11: "SATA",
		12: "SD",
		13: "MMC",
		14: "Virtual",
		15: "File Backed Virtual",
		16: "Storage Spaces",
		17: "NVMe",
		18: "SCM",
	}

	if name, ok := busTypes[busType]; ok {
		return name
	}
	return fmt.Sprintf("Unknown (%d)", busType)
}
