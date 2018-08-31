package spider

import "github.com/weihanchen/tw-currency-tool/spider/rter"

func NewRterSpider() *rter.RterSpider {
	return rter.CreateRterSpider()
}
