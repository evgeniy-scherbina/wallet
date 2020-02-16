package main

import (
	"github.com/jessevdk/go-flags"
)

const (
	defaultRpcListen  = "localhost:9090"
	defaultHTTPListen = "localhost:8080"

	defaultUser   = "postgres"
	defaultPass   = "mysecretpassword"
	defaultAddr   = "localhost"
	defaultDBName = "wallet"
)

type RpcOptions struct {
	Listen string `long:"listen"`
}

type HttpOptions struct {
	Listen string `long:"listen"`
}

type DBOptions struct {
	User string `long:"user"`
	Pass string `long:"pass"`
	Addr string `long:"addr"`
	Name string `long:"name"`
}

type Config struct {
	Rpc  *RpcOptions  `group:"rpc" namespace:"rpc"`
	Http *HttpOptions `group:"http" namespace:"http"`
	DB   *DBOptions   `group:"db" namespace:"db"`
}

func getConf() (*Config, error) {
	cfg := Config{
		Rpc: &RpcOptions{
			Listen: defaultRpcListen,
		},
		Http: &HttpOptions{
			Listen: defaultHTTPListen,
		},
		DB: &DBOptions{
			User: defaultUser,
			Pass: defaultPass,
			Addr: defaultAddr,
			Name: defaultDBName,
		},
	}
	_, err := flags.Parse(&cfg)
	return &cfg, err
}
