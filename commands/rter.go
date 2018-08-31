package commands

import (
	"github.com/spf13/cobra"
	"github.com/weihanchen/tw-currency-tool/spider"
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

func (m *Mobile01Commander) code() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "code",
		Short: "查詢貨幣代碼",
		Run: func(ccmd *cobra.Command, args []string) {
			spider := spider.NewRterSpider()
			spider.PrintCodes()
		},
	}
	return cmd
}

func (m *Mobile01Commander) news() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "news",
		Short: "相關新聞",
		Run: func(ccmd *cobra.Command, args []string) {
			spider := spider.NewRterSpider()
			spider.PrintNews()
		},
	}
	return cmd
}

func (m *Mobile01Commander) rate() *cobra.Command {
	var (
		code     string
		rateType string
		sortable bool
	)
	cmd := &cobra.Command{
		Use:   "rate",
		Short: "匯率資訊",
		SuggestionsMinimumDistance: 1,
		Run: func(ccmd *cobra.Command, args []string) {
			spider := spider.NewRterSpider()
			spider.PrintRate(rateType, code, sortable)
		},
	}
	cmd.Flags().StringVarP(&code, "code", "c", "USD", "currency code, search by: tw-currency-tool rter code")
	cmd.Flags().StringVarP(&rateType, "type", "t", "check", "現金匯率(cach)、即期匯率(check)")
	cmd.Flags().BoolVarP(&sortable, "sortable", "s", false, "auto sort")
	cmd.MarkFlagRequired("code")
	return cmd
}

func newRterCommand() *cobra.Command {
	commander := &Mobile01Commander{}
	cmd := newCommandsBuilder(commander.basic()).
		addCommand(commander.code()).
		addCommand(commander.news()).
		addCommand(commander.rate()).
		build()
	return cmd
}
