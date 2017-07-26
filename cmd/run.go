package cmd

import (
	"os"

	"chkit-v2/chlib"
	"chkit-v2/helpers"

	"regexp"

	"github.com/spf13/cobra"
)

var runCmdName string

var runCmd = &cobra.Command{
	Use:   "run NAME (--image -i IMAGE | --configure)",
	Short: "Run deployment NAME and generate json file",
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			np.FEEDBACK.Println("NAME must be specified")
			cmd.Usage()
			os.Exit(1)
		}
		runCmdName = args[0]
		if !regexp.MustCompile(chlib.DeployRegex).MatchString(runCmdName) {
			np.FEEDBACK.Println("Invalid deployment name")
			cmd.Usage()
			os.Exit(1)
		}
		if !cmd.Flag("image").Changed && !cmd.Flag("configure").Changed {
			np.FEEDBACK.Println("Image or configure parameter must be specified")
			cmd.Usage()
			os.Exit(1)
		}
		if cmd.Flag("image").Changed && cmd.Flag("configure").Changed {
			np.FEEDBACK.Println("Only image or configured must be specified")
			cmd.Usage()
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		var params chlib.ConfigureParams
		if cmd.Flag("configure").Changed {
			params = chlib.PromptParams(np)
		} else {
			params = chlib.ParamsFromArgs(np, cmd.Flags())
		}
		client, err := chlib.NewClient(db, helpers.CurrentClientVersion, helpers.UuidV4(), np)
		if err != nil {
			np.ERROR.Println(err)
			return
		}
		ns, _ := cmd.Flags().GetString("namespace")
		np.FEEDBACK.Print("run...")
		_, err = client.Run(runCmdName, params, ns)
		if err != nil {
			np.FEEDBACK.Println("ERROR")
			np.ERROR.Println(err)
		} else {
			np.FEEDBACK.Println("OK")
		}
	},
}

func init() {
	runCmd.PersistentFlags().Bool("configure", false, "Run interactive configurator")
	runCmd.PersistentFlags().StringP("image", "i", "", "Image name")
	runCmd.PersistentFlags().IntSliceP("port", "p", []int{}, "Ports which will be opened.Format: 8080 ... 4556")
	runCmd.PersistentFlags().StringSliceP("labels", "l", []string{}, "Tags for deployment. Format: key1=value1 ... keyN=valueN")
	runCmd.PersistentFlags().StringSliceP("commands", "C", []string{}, "Commands executed on container start. Format: command1 ... commandN")
	runCmd.PersistentFlags().StringSliceP("env", "e", []string{}, "Environment variables. Format: key1=value1 ... keyN=valueN")
	runCmd.PersistentFlags().StringP("cpu", "c", chlib.DefaultCPURequest, "CPU cores. Format: (number)[m]")
	runCmd.PersistentFlags().StringP("memory", "m", chlib.DefaultMemoryRequest, "Memory size. Format: (number)[Mi|Gi]")
	runCmd.PersistentFlags().IntP("replicas", "r", chlib.DefaultReplicas, "Replicas count")
	runCmd.PersistentFlags().StringP("namespace", "n", "", "Namespace")
	cobra.MarkFlagCustom(runCmd.PersistentFlags(), "namespace", "__chkit_namespaces_list")
	RootCmd.AddCommand(runCmd)
}