package rter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/olekukonko/tablewriter"
	"github.com/weihanchen/tw-currency-tool/httpRequest"
	"github.com/weihanchen/tw-currency-tool/util"
)

type RterSpider struct {
	bankReg         *regexp.Regexp
	currencyCodeReg *regexp.Regexp
	endpoint        string
	dateFmt         string
}

func CreateRterSpider() *RterSpider {
	var (
		bankReg         = regexp.MustCompile("<[^>]*>")
		currencyCodeReg = regexp.MustCompile(`currency/([A-Z]+)/`)
		endpoint        = "https://tw.rter.info"
	)
	return &RterSpider{
		bankReg:         bankReg,
		currencyCodeReg: currencyCodeReg,
		endpoint:        endpoint,
		dateFmt:         "2006-01-02T15:04:05",
	}
}

func (r *RterSpider) getCurrencyCodes() []CurrencyCode {
	response := httpRequest.Get(r.endpoint, "", false)
	if response == nil {
		return nil
	}
	defer response.Body.Close()
	doc, err := goquery.NewDocumentFromResponse(response)
	util.ErrorHandler(err)
	codes := make([]CurrencyCode, 0)
	doc.Find(".dropdown-menu > li a").Each(func(i int, s *goquery.Selection) {
		href, exist := s.Attr("href")
		if !exist {
			return
		}
		countryMatches := r.currencyCodeReg.FindStringSubmatch(href)
		if len(countryMatches) > 0 {
			code := countryMatches[1]
			title := s.Text()
			codes = append(codes, CurrencyCode{Code: code, Title: title})
		}
	})
	return codes
}

func (r *RterSpider) getNews() []News {
	response := httpRequest.Get(r.endpoint, "", false)
	if response == nil {
		return nil
	}
	defer response.Body.Close()
	doc, err := goquery.NewDocumentFromResponse(response)
	util.ErrorHandler(err)
	news := make([]News, 0)
	doc.Find(".media-list > li").Each(func(i int, s *goquery.Selection) {
		hrefElm := s.Find("a")
		href, exist := hrefElm.Attr("href")
		if !exist {
			return
		}
		dateText := strings.TrimSpace(s.Find(".media-body > font").Text())
		date, _ := time.Parse(r.dateFmt, dateText)
		title := strings.TrimSpace(hrefElm.Text())
		shortContent := strings.TrimSpace(s.Find(".media-body > p.text-muted").Text())
		news = append(news, News{
			Date:         date,
			Href:         href,
			Title:        title,
			ShortContent: shortContent,
		})
	})
	return news
}

func (r *RterSpider) getRates(rateType string, code string) []CurrencyRate {
	requestURL := fmt.Sprintf("%s/json.php?t=currency&q=%s&iso=%s", r.endpoint, rateType, code)

	response := httpRequest.Get(requestURL, "", false)
	if response == nil {
		return nil
	}
	defer response.Body.Close()
	b, _ := ioutil.ReadAll(response.Body)
	var (
		rateBody CurrencyRateBody
		rates    []CurrencyRate
	)
	err := json.Unmarshal(b, &rateBody)
	if err != nil {
		util.ErrorHandler(err)
		return nil
	}
	for _, data := range rateBody.Data {
		bank := r.bankReg.ReplaceAllString(data[0], "")
		in, err := strconv.ParseFloat(data[1], 64)
		if err != nil {
			break
		}
		out, err := strconv.ParseFloat(data[2], 64)
		if err != nil {
			break
		}
		updateTime := data[3]
		rates = append(rates, CurrencyRate{
			Bank:       bank,
			In:         in,
			Out:        out,
			UpdateTime: updateTime,
		})
	}
	return rates
}

func (r *RterSpider) markLineBreak(text string, lenToBreak int) string {
	res := ""
	for i, c := range text {
		res = res + string(c)
		if i > 0 && i%lenToBreak == 0 {
			res = res + "\n"
		}
	}
	return res
}

func (r *RterSpider) PrintCodes() {
	currencyCodes := r.getCurrencyCodes()
	header := []string{"Code", "title"}
	var data [][]string
	for _, currencyCode := range currencyCodes {
		data = append(data, []string{currencyCode.Code, currencyCode.Title})
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.AppendBulk(data)
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold, tablewriter.FgWhiteColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgWhiteColor})
	table.SetColumnColor(tablewriter.Colors{tablewriter.Bold, tablewriter.FgGreenColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiBlackColor})
	table.Render()
}

func (r *RterSpider) PrintNews() {
	news := r.getNews()
	header := []string{"Title", "ShortContent", "Link"}
	var data [][]string
	for _, n := range news {
		shortContent := r.markLineBreak(n.ShortContent, 20)
		data = append(data, []string{n.Title, shortContent, n.Href})
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowLine(true)
	table.SetHeader(header)
	table.AppendBulk(data)
	table.Render()
}

type RateSorter []CurrencyRate

func (rate RateSorter) Len() int           { return len(rate) }
func (rate RateSorter) Swap(i, j int)      { rate[i], rate[j] = rate[j], rate[i] }
func (rate RateSorter) Less(i, j int) bool { return rate[i].Out < rate[j].Out }

func (r *RterSpider) PrintRate(rateType string, code string, sortable bool) {
	rates := r.getRates(rateType, code)
	header := []string{"銀行", "買進", "賣出", "更新日期"}
	var data [][]string
	if sortable {
		sort.Sort(RateSorter(rates))
	}
	for _, rate := range rates {
		in := strconv.FormatFloat(rate.In, 'f', 6, 64)
		out := strconv.FormatFloat(rate.Out, 'f', 6, 64)
		data = append(data, []string{rate.Bank, in, out, rate.UpdateTime})
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowLine(true)
	table.SetHeader(header)
	table.AppendBulk(data)
	table.Render()
}
