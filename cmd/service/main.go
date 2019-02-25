package main

import (
	"context"
	"flag"
	"github.com/golang/glog"
	"google.golang.org/grpc"
	"grpc-gateway-test/proto"
	"net"
)

var (
	addr    = flag.String("addr", ":9090", "endpoint of the gRPC service")
	network = flag.String("network", "tcp", "a valid network type which is consistent to -addr")
)

type YourService struct {
}

func NewYourService() *YourService {
	return &YourService{}
}

func (*YourService) Echo(ctx context.Context, msg *proto.StringMessage) (*proto.StringMessage, error) {
	return &proto.StringMessage{Value: "it is work"}, nil
}

func run(ctx context.Context, network, address string) error {
	l, err := net.Listen(network, address)
	if err != nil {
		return err
	}
	defer func() {
		if err := l.Close(); err != nil {
			glog.Errorf("Failed to close %s %s: %v", network, address, err)
		}
	}()
	s := grpc.NewServer()
	proto.RegisterYourServiceServer(s, NewYourService())
	go func() {
		defer s.GracefulStop()
		<-ctx.Done()
	}()
	return s.Serve(l)
}

func main() {
	flag.Parse()
	defer glog.Flush()

	ctx := context.Background()
	if err := run(ctx, *network, *addr); err != nil {
		glog.Fatal(err)
	}
}
