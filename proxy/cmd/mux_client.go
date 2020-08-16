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

var muxClient mux.ClientConfig

var muxClientCmd = &cobra.Command{
	Use:   "client",
	Short: "proxy on mux client mode",
	Run: func(cmd *cobra.Command, args []string) {
		if forever {
			return
		}
		muxClient.SKCPConfig = kcpCfg

		srv := mux.NewClient(muxClient,
			mux.WithClientLogger(zap.S()),
			mux.WithClientGPool(sword.GPool))
		err := srv.Start()
		if err != nil {
			log.Fatalf("run service [%s],%s", cmd.Name(), err)
		}
		server = srv
	},
}

func init() {
	flags := muxClientCmd.Flags()

	flags.StringVarP(&muxClient.ParentType, "parent-type", "T", "tcp", "parent protocol type <tcp|tls|stcp|kcp>")
	flags.StringVarP(&muxClient.Parent, "parent", "P", "", "parent address, such as: \"23.32.32.19:28008\"")
	flags.StringVarP(&muxClient.CertFile, "cert", "C", "proxy.crt", "cert file for tls")
	flags.StringVarP(&muxClient.KeyFile, "key", "K", "proxy.key", "key file for tls")
	flags.StringVar(&muxClient.SecretKey, "sk", "default", "key same with server")
	flags.BoolVar(&muxClient.Compress, "compress", false, "compress data when tcp|tls|stcp mode")
	flags.StringVar(&muxClient.STCPMethod, "stcp-method", "aes-192-cfb", "method of local stcp's encrpyt/decrypt, these below are supported :\n"+strings.Join(encrypt.CipherMethods(), ","))
	flags.StringVar(&muxClient.STCPPassword, "stcp-password", "thinkgos's_goproxy", "password of local stcp's encrpyt/decrypt")
	flags.DurationVarP(&muxClient.Timeout, "timeout", "i", time.Second*2, "tcp timeout duration when connect to real server or parent proxy")
	flags.StringVar(&muxClient.RawProxyURL, "proxy", "", "https or socks5 proxies used when connecting to parent, only worked of -T is tls or tcp, format is https://username:password@host:port https://host:port or socks5://username:password@host:port socks5://host:port")

	rootCmd.AddCommand(muxClientCmd)
}
