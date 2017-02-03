package debugger

import (
	"bufio"
	"bytes"
	"runtime"

	context "golang.org/x/net/context"

	"github.com/docker/containerd"
	dapi "github.com/docker/containerd/api/debugger"
	pp "github.com/maruel/panicparse/stack"
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

	var memInfo dapi.MemInfo
	memInfo.Alloc = response.MemInfo.Alloc
	memInfo.TotalAlloc = response.MemInfo.TotalAlloc
	memInfo.Sys = response.MemInfo.Sys
	memInfo.HeapAlloc = response.MemInfo.HeapAlloc
	memInfo.HeapSys = response.MemInfo.HeapSys
	memInfo.StackInuse = response.MemInfo.StackInuse
	memInfo.StackSys = response.MemInfo.StackSys

	return &dapi.DebugResponse{
		StackDump: response.Stack,
		MemInfo:   &memInfo,
		Version:   response.Version,
		GitCommit: response.GitCommit,
	}, nil
}

func (s *Service) buildStackInfo(ctx context.Context) (string, error) {
	var (
		buf            []byte
		outBuffer      bytes.Buffer
		stackSize      int
		responseBuffer bytes.Buffer
	)
	bufferLen := 16384
	for stackSize == len(buf) {
		buf = make([]byte, bufferLen)
		stackSize = runtime.Stack(buf, true)
		bufferLen *= 2
	}
	buf = buf[:stackSize]
	goroutines, err := pp.ParseDump(bytes.NewReader(buf), bufio.NewWriter(&outBuffer))
	if err != nil {
		return "", err
	}
	p := &pp.Palette{}
	buckets := pp.SortBuckets(pp.Bucketize(goroutines, pp.AnyValue))
	srcLen, pkgLen := pp.CalcLengths(buckets, false)
	for _, bucket := range buckets {
		out1 := p.BucketHeader(&bucket, false, len(buckets) > 1)
		out2 := p.StackLines(&bucket.Signature, srcLen, pkgLen, false)
		responseBuffer.WriteString(out1)
		responseBuffer.WriteString(out2)
	}

	return responseBuffer.String(), nil

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
