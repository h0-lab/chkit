package cli

import (
	"fmt"

	"github.com/containerum/chkit/pkg/cli/configmap"
	"github.com/containerum/chkit/pkg/cli/ingress"
	"github.com/containerum/chkit/pkg/cli/postrun"
	"github.com/containerum/chkit/pkg/cli/prerun"
	"github.com/containerum/chkit/pkg/cli/service"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/spf13/cobra"
)

func Replace(ctx *context.Context) *cobra.Command {
	command := &cobra.Command{
		Use:   "replace",
		Short: "Replace deployment or service",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if err := prerun.PreRun(ctx); err != nil {
				angel.Angel(ctx, err)
				ctx.Exit(1)
			}
			if err := prerun.GetNamespaceByUserfriendlyID(ctx, cmd.Flags()); err != nil {
				fmt.Println(err)
				ctx.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
		PersistentPostRun: postrun.PostRunFunc(ctx),
	}
	command.PersistentFlags().
		StringP("namespace", "n", ctx.Namespace.ID, "")

	command.AddCommand(
		//clideployment.Replace(ctx),
		cliserv.Replace(ctx),
		clingress.Replace(ctx),
		cliconfigmap.Replace(ctx),
	)
	return command
}
