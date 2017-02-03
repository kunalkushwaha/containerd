package main

import (
	"fmt"
	"os"

	gocontext "context"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"time"

	"github.com/docker/containerd/api/debugger"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

const memTemplate = `
-------- Memory Details --------
MemAlloc	: {{.Alloc}}
TotalAlloc	: {{.TotalAlloc}}
SysAlloc	: {{.Sys}}

HeapAlloc	: {{.HeapAlloc}}
HeapSys		: {{.HeapSys}}

StackInUse	: {{.StackInuse}}
StackSys	: {{.StackSys}}
--------------------------------
`

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
			fmt.Println(err)
			return err
		}
		fmt.Println("Containerd Version : ", response.Version)
		fmt.Println("GitCommit : ", response.GitCommit)
		t := template.New("Mem Template")
		t, _ = t.Parse(memTemplate)
		err = t.Execute(os.Stdout, response.MemInfo)
		if err != nil {
			fmt.Println(err)
			return err
		}
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
