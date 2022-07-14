package gateway

import (
	"fmt"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"google.golang.org/grpc"
	"ntsc.ac.cn/tas/tas-commons/pkg/pb"
	"ntsc.ac.cn/tas/tas-commons/pkg/rpc"
)

type Server struct {
	conf      *Config
	rpcServer *rpc.Server
	ifdClient influxdb2.Client
	ifdWriter api.WriteAPIBlocking
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
		pb.RegisterHealthServer(g, &server)
	})
	if err != nil {
		return nil, fmt.Errorf("create grpc server failed: %s", err.Error())
	}
	server.rpcServer = rpcServ
	server.ifdClient = influxdb2.NewClient(
		conf.IfxDBConf.Endpoint, conf.IfxDBConf.Token)

	server.ifdWriter = server.ifdClient.WriteAPIBlocking("", IFX_DB)
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
