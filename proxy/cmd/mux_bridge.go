package cmd

import (
	"log"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/thinkgos/jocasta/lib/encrypt"

	"github.com/thinkgos/jocasta/pkg/sword"
	"github.com/thinkgos/jocasta/services/mux"
)

var muxBridge mux.BridgeConfig

var muxBridgeCmd = &cobra.Command{
	Use:   "bridge",
	Short: "proxy on mux bridge mode",
	Run: func(cmd *cobra.Command, args []string) {
		if forever {
			return
		}
		muxBridge.SKCPConfig = kcpCfg

		srv := mux.NewBridge(muxBridge,
			mux.WithBridgeLogger(zap.S()),
			mux.WithBridgeGPool(sword.GPool))
		err := srv.Start()
		if err != nil {
			log.Fatalf("run service [%s],%s", cmd.Name(), err)
		}
		server = srv
	},
}

func init() {
	flags := muxBridgeCmd.Flags()

	flags.StringVarP(&muxBridge.LocalType, "local-type", "t", "tcp", "local protocol type <tcp|tls|stcp|kcp>")
	flags.StringVarP(&muxBridge.Local, "local", "p", ":28080", "local ip:port to listen")
	flags.BoolVar(&muxBridge.Compress, "compress", false, "compress data when tcp|tls|stcp mode")
	flags.StringVarP(&muxBridge.CertFile, "cert", "C", "proxy.crt", "cert file for tls")
	flags.StringVarP(&muxBridge.KeyFile, "key", "K", "proxy.key", "key file for tls")
	flags.StringVar(&muxBridge.STCPConfig.Method, "stcp-method", "aes-192-cfb", "method of local stcp's encrpyt/decrypt, these below are supported :\n"+strings.Join(encrypt.CipherMethods(), ","))
	flags.StringVar(&muxBridge.STCPConfig.Password, "stcp-password", "thinkgos's_goproxy", "password of local stcp's encrpyt/decrypt")
	flags.DurationVarP(&muxBridge.Timeout, "timeout", "e", 2*time.Second, "tcp timeout duration when connect to real server or parent proxy")

	rootCmd.AddCommand(muxBridgeCmd)
}
