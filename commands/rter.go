package commands

import (
	"github.com/spf13/cobra"
	"gitlab.com/will.chen/tw-currency-tool/spider"
)

type Mobile01Commander struct {
}

func (m *Mobile01Commander) basic() *cobra.Command {
	return &cobra.Command{
		Use:   "rter",
		Short: "即匯站",
		SuggestionsMinimumDistance: 1,
		RunE: nil,
	}
}

func (m *Mobile01Commander) countries() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "countries",
		Short: "國家",
		Run: func(ccmd *cobra.Command, args []string) {

		},
	}
	return cmd
}

func (m *Mobile01Commander) news() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "news",
		Short: "相關新聞",
		Run: func(ccmd *cobra.Command, args []string) {

		},
	}
	return cmd
}

func (m *Mobile01Commander) rate() *cobra.Command {
	var boardID string
	cmd := &cobra.Command{
		Use:   "rate",
		Short: "匯率資訊",
		SuggestionsMinimumDistance: 1,
		Run: func(ccmd *cobra.Command, args []string) {
			spider := spider.NewRterSpider()
			spider.Start()
		},
	}
	cmd.Flags().StringVarP(&boardID, "board_id", "b", "291", "mobile01 board id")
	cmd.MarkFlagRequired("board_id")
	return cmd
}

func newRterCommand() *cobra.Command {
	commander := &Mobile01Commander{}
	cmd := newCommandsBuilder(commander.basic()).
		addCommand(commander.countries()).
		addCommand(commander.news()).
		addCommand(commander.rate()).
		build()
	return cmd
}
