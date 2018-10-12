package main

import (
	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/micro/cli"
	"github.com/micro/go-log"
	micro "github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/broker/kafka"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	_ "github.com/micro/go-plugins/broker/nats"
	_ "github.com/micro/go-plugins/registry/nats"
	_ "github.com/micro/go-plugins/transport/nats"
	_ "github.com/micro/go-plugins/selector/static"
	"github.com/jerryhoio/cloudflare/pkg/proto"
	"github.com/jerryhoio/cloudflare/internal/handler"
)

const (
	serviceName = "io.jerryho.svr.cloudflare"
)

var (
	version = "latest"
)

func main() {

	var cloudFlareToken = cli.StringFlag{
		Name:   "cloudflare_token",
		Usage:  "cloudflare token",
		EnvVar: "CLOUDFLARE_TOKEN",
	}
	var cloudFlareEmail = cli.StringFlag{
		Name:   "cloudflare_email",
		Usage:  "cloudflare email",
		EnvVar: "CLOUDFLARE_EMAIL",
	}
	var cloudFlareCl *cloudflare.API
	var err error
	// New Service
	service := micro.NewService(
		micro.Name(serviceName),
		micro.Version(version),
		micro.Flags(cloudFlareToken, cloudFlareEmail),
	)
	log.Logf("initialize %s service", serviceName)
	service.Init(
		micro.Action(func(c *cli.Context) {
			cloudFlareCl, err = cloudflare.New(c.String("cloudflare_token"), c.String("cloudflare_email"))
			if err != nil {
				log.Fatal(err)
			}
		}),
	)

	domainHandler := handler.NewDomainsHandler(cloudFlareCl)
	proto.RegisterDomainsHandler(service.Server(), domainHandler)
	if err := service.Run(); err != nil {
		log.Log(err)
	}
}
