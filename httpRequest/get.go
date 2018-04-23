package httpRequest

import (
	"net/http"
	"time"

	"gitlab.com/will.chen/tw-currency-tool/util"
)

var usrAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36"

func Get(url, referer string, needWait bool) *http.Response {
	c := &http.Client{}

	reqst, err := http.NewRequest("GET", url, nil)
	if err != nil {
		util.ErrorHandler(err)
		return nil
	}

	if referer != "" {
		reqst.Header.Add("Referer", referer)
	}

	reqst.Header.Add("user-agent", usrAgent)

	resp, err := c.Do(reqst)
	if err != nil {
		util.ErrorHandler(err)
		return nil
	}
	if needWait {
		sTime := util.RandIntOverRange(1000, 500)
		time.Sleep(time.Millisecond * time.Duration(sTime))
	}

	return resp
}
