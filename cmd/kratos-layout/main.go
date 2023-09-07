package main

import (
	"flag"
	"net/url"
	"os"

	"kratos-layout/internal/conf"
	"kratos-layout/pkg/bootstrap"
	"kratos-layout/pkg/servicename"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	Flags   = bootstrap.NewCommandFlags()

	id, _ = os.Hostname()
)

func init() {
	Flags.Init()
	Name = servicename.AppService
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	)
}

func main() {
	flag.Parse()
	logger := bootstrap.NewLoggerProvider(id, Name, Version)
	var confKey string
	consul_url, err := url.Parse(Flags.Consul)
	if err == nil {
		Flags.Consul = consul_url.Host
		if consul_url.Path != "" {
			confKey = consul_url.Path
		} else {
			confKey = Flags.Env + "." + Name
		}
	}
	c := bootstrap.NewConfigProvider(Flags.Conf, Flags.Consul, confKey)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	app, cleanup, err := wireApp(bc.Server, bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
