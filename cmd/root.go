package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"ntsc.ac.cn/ta/ta-monitor-gateway/internal/gateway"
	ccmd "ntsc.ac.cn/tas/tas-commons/pkg/cmd"
)

var envs struct {
	listener    string
	certPath    string
	ifxEndpoint string
	ifxToken    string
}

var rootCmd = &cobra.Command{
	Use:    "ta-monitor-gateway",
	Short:  "TA monitor gateway",
	PreRun: prerun,
	Run:    run,
}

func init() {
	cobra.OnInitialize(func() {})
	viper.AutomaticEnv()
	viper.SetEnvPrefix("TA")

	rootCmd.Flags().StringVar(&ccmd.GlobalEnvs.LoggerLevel,
		"logger-level", "DEBUG", "logger level")
	rootCmd.Flags().StringVar(&envs.listener,
		"rpc-listener", "tcp://0.0.0.0:1358", "grpc listener url")
	rootCmd.Flags().StringVar(&envs.certPath,
		"cert-path", "/etc/ntsc/ta/certs", "system certificates path")
	rootCmd.Flags().StringVar(&envs.ifxEndpoint,
		"ifx-endpoint", "http://localhost:8086", "influxdb endpoint")
	rootCmd.Flags().StringVar(&envs.ifxToken,
		"ifx-token", "", "influxdb token")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func prerun(cmd *cobra.Command, args []string) {
	ccmd.InitGlobalVars()
	var err error
	if err = ccmd.ValidateStringVar(&envs.certPath, "cert_path", true); err != nil {
		logrus.WithField("prefix", "cmd.root").
			Fatalf("check boot var failed: %s", err.Error())
	}
	if err = ccmd.ValidateStringVar(&envs.listener, "rpc_listener", true); err != nil {
		logrus.WithField("prefix", "cmd.root").
			Fatalf("check boot var failed: %s", err.Error())
	}
	go func() {
		ccmd.RunWithSysSignal(nil)
	}()
}

func run(cmd *cobra.Command, args []string) {
	serv, err := gateway.NewServer(&gateway.Config{
		Listener: envs.listener,
		CertPath: envs.certPath,
		IfxDBConf: &gateway.InfluxDBConfig{
			Endpoint: envs.ifxEndpoint,
			Token:    envs.ifxToken,
		},
	})
	if err != nil {
		logrus.WithField("prefix", "cmd.root").
			Fatalf("create monitor gateway failed: %s", err.Error())
	}
	logrus.WithField("prefix", "cmd.root").
		Fatalf("run registry server failed: %s", <-serv.Start())
}
