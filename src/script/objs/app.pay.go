package objs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// 解析原始表单
func FormParse(form string) (map[string]string, error) {
	rst := map[string]string{}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(form))
	if err != nil {
		return rst, err
	}
	doc.Find("input[type=hidden]").Each(func(i int, s *goquery.Selection) {
		k, _ := s.Attr("name")
		v, _ := s.Attr("value")
		rst[k] = v
	})
	fmt.Println("原始解析Map", rst)
	return rst, nil
}

// 获取并解析 北理工支付表单
func FormBitPay(f map[string]string) (map[string]string, error) {
	var (
		client  *http.Client
		request *http.Request // 请求对象
	)

	body := url.Values{}
	for k, v := range f {
		body.Add(k, v)
	}
	api := "https://pay.info.bit.edu.cn/pay/prepay"

	request, _ = http.NewRequest("POST", api, strings.NewReader(body.Encode()))
	BitHeader(request)

	client = &http.Client{
		Timeout: time.Second * 18,
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	rst := map[string]string{}
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return rst, err
	}
	doc.Find("input[type=hidden]#tenantParams").Each(func(i int, s *goquery.Selection) {
		k, _ := s.Attr("name")
		v, _ := s.Attr("value")
		newV := strings.ReplaceAll(v, "&quot;", "\"")
		rst[k] = newV
	})
	rst["gateway"] = "WX"
	fmt.Println(rst)
	defer response.Body.Close()
	return rst, nil
}
func GetWxPay(f map[string]string) (string, error) {
	var (
		client  *http.Client
		request *http.Request // 请求对象
	)
	// 生成body
	body := url.Values{}
	for k, v := range f {
		body.Add(k, v)
	}
	api := "https://pay.info.bit.edu.cn/WXPay/pay/"

	request, _ = http.NewRequest("POST", api, strings.NewReader(body.Encode()))
	WxHeader(request)
	// request.Header.Add("Cookie", "JSESSIONID=usib33s6o0y5rbk3o1p25f5j")
	client = &http.Client{
		Timeout: time.Second * 18,
		// CheckRedirect: func(req *http.Request, via []*http.Request) error {
		// 	return http.ErrUseLastResponse
		// },
	}
	var payLink string
	response, err := client.Do(request)
	if err != nil {
		return payLink, err
	}
	// 如果想请求原跳转链接，这样操作
	// payLink = response.Header.Get("Location")
	// if err != nil {
	// 	fmt.Println(err)
	// 	resBody, _ := ioutil.ReadAll(response.Body)
	// 	fmt.Println(string(resBody))
	// 	rst := *payLink
	// 	return rst, nil
	// }
	resBody, _ := ioutil.ReadAll(response.Body)

	str := string(resBody)
	fmt.Println(str)
	preArr := strings.Split(str, "var url=\"")
	if len(preArr) <= 1 {
		doc, err := goquery.NewDocumentFromReader(response.Body)
		if err != nil {
			return payLink, err
		}
		wxMsg := ""
		doc.Find("p.weui-msg__desc").Each(func(i int, s *goquery.Selection) {
			wxMsg = s.Text()
		})
		wxMsg = strings.ReplaceAll(wxMsg, " ", "")
		wxMsg = strings.ReplaceAll(wxMsg, "\n", "")
		defer response.Body.Close()
		return payLink, errors.New(wxMsg)
	}
	str1 := preArr[1]
	str2 := strings.Split(str1, "top.location.href=url")[0]
	str3 := strings.ReplaceAll(str2, " ", "")
	str4 := str3[:len(str3)-3]
	arr := strings.Split(str4, "\";\nvarredirect_url=\"")
	payLink = arr[0] + "@" + arr[1]
	defer response.Body.Close()
	return payLink, nil
}

func BitHeader(request *http.Request) {
	request.Header.Add("Origin", "http://gym.dazuiwl.cn")
	request.Header.Add("Host", "pay.info.bit.edu.cn")
	request.Header.Add("Referer", "http://gym.dazuiwl.cn/")
	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	request.Header.Add("Accept-Language", "en,en-US;q=0.9,zh-CN;q=0.8,zh;q=0.7")
	request.Header.Add("Cache-Control", "no-cache")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Pragma", "no-cache")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.41 Mobile Safari/537.36 Edg/101.0.1210.32")
	request.Header.Add("X-Requested-With", "XMLHttpRequest")
}

func WxHeader(request *http.Request) {
	request.Header.Add("Origin", "https://pay.info.bit.edu.cn")
	request.Header.Add("Host", "pay.info.bit.edu.cn")
	request.Header.Add("Referer", "https://pay.info.bit.edu.cn/pay/prepay")
	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	request.Header.Add("Accept-Language", "en,en-US;q=0.9,zh-CN;q=0.8,zh;q=0.7")
	request.Header.Add("Cache-Control", "no-cache")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Pragma", "no-cache")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.41 Mobile Safari/537.36 Edg/101.0.1210.32")
}
