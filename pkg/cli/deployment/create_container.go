package clideployment

import (
	"fmt"
	"os"

	"github.com/containerum/chkit/pkg/context"
	containerControl "github.com/containerum/chkit/pkg/controls/container"
	"github.com/containerum/chkit/pkg/export"
	"github.com/containerum/chkit/pkg/model/configmap"
	"github.com/containerum/chkit/pkg/model/deployment"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
)

func CreateContainer(ctx *context.Context) *cobra.Command {
	var flags struct {
		Force         bool   `flag:"force f" desc:"suppress confirmation"`
		ContainerName string `flag:"container" desc:"container name, required on --force"`
		Deployment    string `desc:"deployment name, required on --force"`
		containerControl.Flags
	}
	command := &cobra.Command{
		Use:     "deployment-container",
		Aliases: []string{"depl-cont", "container", "dc"},
		Short:   "create deployment container.",
		Long: "Create deployment container.\n" +
			"Runs in one-line mode, suitable for integration with other tools, and in interactive wizard mode.",
		Run: func(cmd *cobra.Command, args []string) {
			var logger = ctx.Log.Command("create deployment container")
			logger.Debugf("START")
			defer logger.Debugf("END")

			if flags.Deployment == "" && flags.Force {
				ferr.Printf("deployment name must be provided as --deployment while using --force")
				os.Exit(1)
			} else if flags.Deployment == "" {
				logger.Debugf("getting deployment list from namespace %q", ctx.Namespace)
				var depl, err = ctx.Client.GetDeploymentList(ctx.Namespace.ID)
				if err != nil {
					logger.WithError(err).Errorf("unable to get deployment list from namespace %q", ctx.Namespace)
					ferr.Println(err)
					os.Exit(1)
				}
				logger.Debugf("selecting deployment")
				(&activekit.Menu{
					Title: "Select deployment",
					Items: activekit.StringSelector(depl.Names(), func(s string) error {
						logger.Debugf("deployment %q selected", s)
						flags.Deployment = s
						return nil
					}),
				}).Run()
			}

			if flags.ContainerName == "" && flags.Force {
				ferr.Printf("container name must be provided as --container while using --force")
				os.Exit(1)
			}

			logger.Debugf("building container from flags")
			var cont, err = flags.Container()
			if err != nil {
				logger.WithError(err).Errorf("unable to build container from flags")
				ferr.Println(err)
				os.Exit(1)
			}
			cont.Name = flags.ContainerName
			fmt.Println(cont.RenderTable())

			if flags.Force {
				logger.Debugf("running in --force mode")
				logger.Debugf("validating changed container %q", cont.Name)
				if err := cont.Validate(); err != nil {
					logger.WithError(err).Errorf("invalid container %q", cont.Name)
					fmt.Println(err)
					os.Exit(1)
				}
				logger.Debugf("creating container %q", cont.Name)
				if err := ctx.Client.CreateDeploymentContainer(ctx.Namespace.ID, flags.Deployment, cont); err != nil {
					logger.WithError(err).Errorf("unable to replace container %q", cont.Name)
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Println("Ok")
			}

			//	volumes, err := ctx.Client.GetVolumeList(ctx.Namespace.ID)
			//	if err != nil {
			//		ferr.Println(err)
			//		os.Exit(1)
			//	}

			var deployments = make(chan deployment.DeploymentList)
			go func() {
				logger := logger.Component("getting deployment list")
				logger.Debugf("START")
				defer logger.Debugf("END")
				defer close(deployments)
				deplList, err := ctx.Client.GetDeploymentList(ctx.Namespace.ID)
				if err != nil {
					logger.WithError(err).Errorf("unable to get deployment list from namespace %q", ctx.Namespace)
					ferr.Println(err)
					os.Exit(1)
				}
				deployments <- deplList
			}()

			var configs = make(chan configmap.ConfigMapList)
			go func() {
				logger := logger.Component("getting configmap list")
				logger.Debugf("START")
				defer logger.Debugf("END")
				defer close(configs)
				configList, err := ctx.Client.GetConfigmapList(ctx.Namespace.ID)
				if err != nil {
					logger.WithError(err).Errorf("unable to get configmap list")
					ferr.Println(err)
					os.Exit(1)
				}
				configs <- configList
			}()
			logger.Debugf("running wizard")
			containerControl.Wizard{
				EditName:   true,
				Container:  cont,
				Deployment: flags.Deployment,
				//		Volumes:     volumes.Names(),
				Configs:     (<-configs).Names(),
				Deployments: (<-deployments).Names(),
			}.Run()

			if activekit.YesNo("Are you sure you want to create container %q in deployment %q?", cont.Name, flags.Deployment) {
				logger.Debugf("creating container %q in deployment %q", cont.Name, flags.Deployment)
				if err := ctx.Client.CreateDeploymentContainer(ctx.Namespace.ID, flags.Deployment, cont); err != nil {
					logger.WithError(err).Errorf("unable to replace container %q in deployment %q", cont.Name, flags.Deployment)
					ferr.Println(err)
				}
			}
			(&activekit.Menu{
				Items: activekit.MenuItems{
					{
						Label: "Save container to file",
						Action: func() error {
							for {
								var fname = activekit.Promt("Type output filename: ")
								export.ExportData(cont, export.ExportConfig{
									Filename: fname,
									Format:   export.YAML,
								})
							}

							return nil
						},
					},
				},
			}).Run()
		},
	}
	if err := gpflag.ParseTo(&flags, command.PersistentFlags()); err != nil {
		panic(err)
	}
	return command
}
