package debugger

import (
	"runtime"

	context "golang.org/x/net/context"

	"github.com/docker/containerd"
	dapi "github.com/docker/containerd/api/debugger"
)

//Service struct
type Service struct {
}

// DebugResponse struct to build json response.
type DebugResponse struct {
	Version   string
	GitCommit string
	Stack     string
	MemInfo   MemStats
}

//MemStats struct
type MemStats struct {
	Alloc      uint64 // bytes allocated and not yet freed
	TotalAlloc uint64 // bytes allocated (even if freed)
	Sys        uint64 // bytes obtained from system (sum of XxxSys below)

	// Main allocation heap statistics.
	HeapAlloc uint64 // bytes allocated and not yet freed (same as Alloc above)
	HeapSys   uint64 // bytes obtained from system

	// Low-level fixed-size structure allocator statistics.
	//	Inuse is bytes used now.
	//	Sys is bytes obtained from system.
	StackInuse uint64 // bytes used by stack allocator
	StackSys   uint64
}

// NewService returns a new shim service that can be used via GRPC
func NewService() *Service {
	return &Service{}
}

//DumpDebugInfo builds the debug context to be passed to client.
func (s *Service) DumpDebugInfo(ctx context.Context, r *dapi.CreateDebugRequest) (*dapi.DebugResponse, error) {
	//TODO:
	// Build a JSON string with
	// Containerd version.
	// {
	//    "version": "1.0.0",
	//    "commit": "blablalblabla",
	//    "stack: [
	//         {},{}
	//     ]
	// }
	response := DebugResponse{}
	response.Stack, _ = s.buildStackInfo(ctx)
	response.MemInfo, _ = s.buildMemInfo(ctx)
	response.Version = containerd.Version
	response.GitCommit = containerd.GitCommit
	return &dapi.DebugResponse{
		StackDump: response.Stack,
		//	MemInfo:   response.MemInfo,
		Version:   response.Version,
		GitCommit: response.GitCommit,
	}, nil
}

func (s *Service) buildStackInfo(ctx context.Context) (string, error) {
	var (
		buf       []byte
		stackSize int
	)
	bufferLen := 16384
	for stackSize == len(buf) {
		buf = make([]byte, bufferLen)
		stackSize = runtime.Stack(buf, true)
		bufferLen *= 2
	}
	buf = buf[:stackSize]
	return string(buf), nil
}

func (s *Service) buildMemInfo(ctx context.Context) (MemStats, error) {
	memInfo := runtime.MemStats{}
	memStats := MemStats{}
	runtime.ReadMemStats(&memInfo)

	memStats.Alloc = memInfo.Alloc
	memStats.TotalAlloc = memInfo.TotalAlloc
	memStats.Sys = memInfo.Sys

	memStats.HeapAlloc = memInfo.HeapAlloc
	memStats.HeapSys = memInfo.HeapSys

	memStats.StackInuse = memInfo.StackInuse
	memStats.StackSys = memInfo.StackSys

	return memStats, nil
}
