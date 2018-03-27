package cliserv

import (
	"github.com/containerum/chkit/cmd/util"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/model/service"
	"gopkg.in/urfave/cli.v2"
)

var GetService = &cli.Command{
	Name:    "service",
	Aliases: []string{"srv", "services", "svc"},
	Action: func(ctx *cli.Context) error {
		client := util.GetClient(ctx)
		defer util.StoreClient(ctx, client)
		var show model.Renderer

		switch ctx.NArg() {
		case 0:
			namespace := util.GetNamespace(ctx)
			list, err := client.GetServiceList(namespace)
			if err != nil {
				return err
			}
			show = list
		default:
			namespace := util.GetNamespace(ctx)
			servicesNames := ctx.Args().Slice()
			var list service.ServiceList
			for _, servName := range servicesNames {
				serv, err := client.GetService(namespace, servName)
				if err != nil {
					return err
				}
				list = append(list, serv)
			}
			show = list
		}
		return util.WriteData(ctx, show)
	},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "file",
			Aliases: []string{"f"},
		},
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
		},
	},
}
