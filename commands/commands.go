// Package commands is where all cli logic is, including starting hoarder as a server.
package commands

import (
	"flag"
	"runtime"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

type commandsBuilder struct {
	cmd *cobra.Command
}

func (c *commandsBuilder) addCommand(cmd *cobra.Command) *commandsBuilder {
	c.cmd.AddCommand(cmd)
	return c
}

func (c *commandsBuilder) build() *cobra.Command {
	return c.cmd
}

func newCommandsBuilder(cmd *cobra.Command) *commandsBuilder {
	return &commandsBuilder{cmd: cmd}
}

var mongoAddr string

func Execute(args []string) error {

	flag.Set("alsologtostderr", "true")
	flag.CommandLine.Parse([]string{})

	spiderCmd := &cobra.Command{
		Use:   "tw-curremcy-tool",
		Short: "貨幣資訊輔助查詢工具",
		Long:  ``,
		SuggestionsMinimumDistance: 1,
	}
	cmd := newCommandsBuilder(spiderCmd).
		addCommand(newRterCommand()).
		build()
	cmd.SetArgs(args)
	// cli flags
	cmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs + 1)
	glog.Infof("CPU nums: %d", +numCPUs)
	return cmd.Execute()
}
