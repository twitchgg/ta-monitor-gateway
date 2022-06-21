package gateway

import (
	"fmt"

	"google.golang.org/grpc"
	"ntsc.ac.cn/ta-registry/pkg/pb"
	"ntsc.ac.cn/ta-registry/pkg/rpc"
)

type Server struct {
	conf      *Config
	rpcServer *rpc.Server
}

func NewServer(conf *Config) (*Server, error) {
	if conf == nil {
		return nil, fmt.Errorf("server config not define")
	}
	if err := conf.Check(); err != nil {
		return nil, fmt.Errorf("check server config failed: %s", err.Error())
	}
	rpcConf, err := conf.RPCConfig()
	if err != nil {
		return nil, fmt.Errorf("generate rpc config failed: %s", err.Error())
	}
	server := Server{
		conf: conf,
	}
	rpcServ, err := rpc.NewServer(rpcConf, []grpc.ServerOption{
		grpc.StreamInterceptor(
			rpc.StreamServerInterceptor(rpc.CertCheckFunc)),
		grpc.UnaryInterceptor(
			rpc.UnaryServerInterceptor(rpc.CertCheckFunc)),
	}, func(g *grpc.Server) {
		pb.RegisterMonitorServiceServer(g, &server)
	})
	if err != nil {
		return nil, fmt.Errorf("create grpc server failed: %s", err.Error())
	}
	server.rpcServer = rpcServ
	return &server, nil
}

func (s *Server) Start() chan error {
	errChan := make(chan error, 1)
	go func() {
		err := <-s.rpcServer.Start()
		errChan <- err
	}()
	return errChan
}
