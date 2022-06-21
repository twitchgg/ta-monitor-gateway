package gateway

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"ntsc.ac.cn/ta-registry/pkg/pb"
	"ntsc.ac.cn/ta-registry/pkg/rpc"
)

func (s *Server) Report(stream pb.MonitorService_ReportServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return rpc.GenerateError(codes.Canceled, err)
		}
		if err := rpc.CheckMachineID(stream.Context(), req.MachineID); err != nil {
			return rpc.GenerateArgumentError("machine id")
		}
		fmt.Println(req.MachineID, req.Oid, req.ValueType, req.Value)
	}
}
