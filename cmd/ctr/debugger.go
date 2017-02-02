package main

import (
	"fmt"

	gocontext "context"
	"io/ioutil"
	"log"
	"net"
	"time"

	"github.com/docker/containerd/api/debugger"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

var debuggerCommand = cli.Command{
	Name:  "debug",
	Usage: "debug containerd",
	Action: func(context *cli.Context) error {
		debuggerService, err := getDebuggerService(context)
		if err != nil {
			return err
		}

		response, err := debuggerService.DumpDebugInfo(gocontext.Background(), &debugger.CreateDebugRequest{})
		if err != nil {
			fmt.Println("-- After Making request --")
			fmt.Println(err)
			return err
		}
		fmt.Println("Containerd Version : ", response.Version)
		fmt.Println("GitCommit : ", response.GitCommit)
		fmt.Println("Stack Dump : ")
		fmt.Println(response.StackDump)
		return nil
	},
}

func getDebuggerService(context *cli.Context) (debugger.DebuggerServiceClient, error) {
	//FIXME: Should not be hardcode.
	bindSocket := "/run/containerd/containerd-dapi.sock"

	// reset the logger for grpc to log to dev/null so that it does not mess with our stdio
	grpclog.SetLogger(log.New(ioutil.Discard, "", log.LstdFlags))
	dialOpts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithTimeout(100 * time.Second)}
	dialOpts = append(dialOpts,
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", bindSocket, timeout)
		},
		))
	conn, err := grpc.Dial(fmt.Sprintf("unix://%s", bindSocket), dialOpts...)
	if err != nil {
		return nil, err
	}
	return debugger.NewDebuggerServiceClient(conn), nil

}
