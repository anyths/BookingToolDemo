package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/ini.v1"
)

type Res struct {
	Code int
	Data map[string]int
}

// 请求体必要参数
type Params struct {
	// api    string
	vapi   string
	cookie string
	token  string
	// period string
	// body   map[string]string
	// freq   int
}

func main() {
	var (
		request *http.Request // 请求对象
		path    string
		osGwd   string
		api     string
		res     Res
		params  Params
	)

	osGwd, _ = os.Getwd()
	path = osGwd + "/.ini"
	cfg, err := ini.Load(path)
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	params.vapi = cfg.Section("").Key("vapi").String()
	params.cookie = cfg.Section("").Key("cookie").String()
	params.token = cfg.Section("").Key("token").String()

	fmt.Println("请输入查询月份和日期, 格式为: 0506 代表5月6日")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("请输入： ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		input = strings.TrimSpace(input)
		if len(input) != 4 {
			fmt.Println("输入错误！请重试")
			continue
		}
		inputBytes := []byte(input)
		month := string(inputBytes[:2])
		day := string(inputBytes[2:])
		api = "http://gym.dazuiwl.cn/api/sport_schedule/booked/id/34?day=2022-" + month + "-" + day

		// 创建请求头
		request, err = http.NewRequest("GET", api, nil)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("初始化请求头错误...重试中...")
			continue
		}
		// 添加头内容
		HeaderAdd(request, &params)

		for {
			res, err = Req(*request, input)
			if err != nil {
				fmt.Println(err.Error())
				fmt.Println("请求错误...重试中...")
				continue
			}
			// if res.Code != 502 {
			// 	fmt.Println("服务器很垃圾，502堵塞了")
			// 	fmt.Println("5s后重试...")
			// 	time.Sleep(time.Second * 5)
			// 	continue
			// }
			// if res.Code != 404 {
			// 	fmt.Println("服务器很垃圾，竟然404了")
			// 	fmt.Println("5s后重试...")
			// 	time.Sleep(time.Second * 5)
			// 	continue
			// }
			break
		}
		Print(res, input)

	}

}

func Req(request http.Request, input string) (Res, error) {
	client := &http.Client{
		Timeout: time.Second * 30,
	}
	var res Res
	response, err := client.Do(&request)
	if err != nil {
		return res, err
	}
	if response.StatusCode != 200 {
		return res, errors.New(response.Status)
	}
	resBody, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal([]byte(resBody), &res)
	defer response.Body.Close()
	return res, nil
}

func Print(res Res, input string) {
	fmt.Println("*************************")
	fmt.Printf(">>>>> %10s <<<<<< \n", "2022"+input)
	fmt.Println("=====================================================================================================")
	fmt.Printf("%s %7s ", "时间↓↓", "场地>>  ")
	for x := 155; x <= 166; x++ {
		fmt.Printf("   #%3d", x)
	}
	fmt.Printf("\n")
	fmt.Println("-----------------------------------------------------------------------------------------------------")
	t := 7
	for n := int(328101); n <= 328107; n++ {
		fmt.Printf("[%6d] %2d-%2d:  ", n, t, t+1)
		for m := int(155); m <= 166; m++ {
			for k, v := range res.Data {
				arr := strings.Split(k, "-")
				m1, _ := strconv.Atoi(arr[0])
				n1, _ := strconv.Atoi(arr[1])
				if m1 == m && n1 == n {
					fmt.Printf("   %4d", v)
				}
			}
		}
		fmt.Printf("\n")
		t++
	}
	for n := int(328125); n <= 328132; n++ {
		fmt.Printf("[%6d] %2d-%2d:  ", n, t, t+1)
		for m := int(155); m <= 166; m++ {
			for k, v := range res.Data {
				arr := strings.Split(k, "-")
				m1, _ := strconv.Atoi(arr[0])
				n1, _ := strconv.Atoi(arr[1])
				if m1 == m && n1 == n {
					fmt.Printf("   %4d", v)
				}
			}
		}
		fmt.Printf("\n")
		t++
	}
	fmt.Println("=====================================================================================================")
}

func HeaderAdd(r *http.Request, p *Params) {
	r.Header.Add("Origin", "http://gym.dazuiwl.cn")
	r.Header.Add("Host", "gym.dazuiwl.cn")
	r.Header.Add("Referer", "http://gym.dazuiwl.cn/h5/")
	r.Header.Add("Accept", "*/*")
	r.Header.Add("Accept-Language", "en,en-US;q=0.9,zh-CN;q=0.8,zh;q=0.7")
	r.Header.Add("Cache-Control", "no-cache")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	// r.Header.Add("Cookie", p.cookie)
	r.Header.Add("Pragma", "no-cache")
	r.Header.Add("Proxy-Connection", "keep-alive")
	r.Header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.41 Mobile Safari/537.36 Edg/101.0.1210.32")
	r.Header.Add("X-Requested-With", "XMLHttpRequest")
	// r.Header.Add("token", p.token)
}
