package gateway

import (
	"context"
	"strconv"
	"time"

	"google.golang.org/grpc/codes"
	"ntsc.ac.cn/ta-registry/pkg/pb"
	"ntsc.ac.cn/ta-registry/pkg/rpc"

	"github.com/sirupsen/logrus"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

const (
	IFX_DB = "TA-SNMP"
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
		logrus.WithField("prefix", "handler_report").
			Tracef("[%s] report oid [%s],value type [%s] value: %v",
				req.MachineID, req.Oid, req.ValueType, req.Value)
		var value interface{}
		switch req.ValueType {
		case "TimeTicks":
			continue
		case "OctetString":
			value = req.Value
		case "Integer", "Counter64", "Counter32":
			value, err = strconv.Atoi(req.Value)
			if err != nil {
				logrus.WithField("prefix", "handler_report").
					Warnf("failed to convert integer: %s", req.Value)
				continue
			}
		default:
			logrus.WithField("prefix", "handler_report").
				Warnf("unsupport oid type: %s", req.ValueType)
			value = nil
		}
		if value == nil {
			continue
		}
		point := influxdb2.NewPoint(IFX_DB,
			map[string]string{
				"unit": req.ValueType,
				"host": req.MachineID,
			},
			map[string]interface{}{req.Oid: value},
			time.Now())
		if err := s.ifdWriter.WritePoint(context.Background(), point); err != nil {
			logrus.WithField("prefix", "handler_report").
				Errorf("failed to write influxDB data: %v", err)
		}

	}
}
