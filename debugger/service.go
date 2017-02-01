package debugger

import (
	"context"
	"runtime"
	"runtime/debug"

	dapi "github.com/docker/containerd/api/debugger"
)

//Service struct
type Service struct {
}

// NewService returns a new shim service that can be used via GRPC
func NewService() *Service {
	return &Service{}
}

//CreateDebugInfo builds the debug context to be passed to client.
func (s *Service) CreateDebugInfo(ctx context.Context, r *dapi.CreateDebugInfoRequest) (*dapi.CreateStackDumpResponse, error) {

	return nil, nil
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

func (s *Service) buildMemInfo(ctx context.Context) (string, error) {
	memInfo := runtime.MemStats{}
	runtime.ReadMemStats(&memInfo)

	//TODO: Parse MemInfo into human readable format.
	return "", nil
}

func (s *Service) buildGCStatsInfo(ctx context.Context) (string, error) {
	gcInfo := debug.GCStats{}
	debug.ReadGCStats(&gcInfo)
	//FIXME: Parse gcInfo in human readable format
	return "", nil
}
