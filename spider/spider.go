package spider

import "gitlab.com/will.chen/tw-currency-tool/spider/rter"

func NewRterSpider() *rter.RterSpider {
	return rter.CreateRterSpider()
}
