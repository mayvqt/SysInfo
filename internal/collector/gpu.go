package collector

import (
	"fmt"

	"github.com/mayvqt/sysinfo/internal/types"
)

// CollectGPU gathers GPU information
func CollectGPU() (*types.GPUData, error) {
	gpus := collectGPUPlatform()

	if len(gpus) == 0 {
		return nil, fmt.Errorf("no GPU information available")
	}

	data := &types.GPUData{
		GPUs: gpus,
	}

	return data, nil
}
