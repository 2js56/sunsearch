package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var wg sync.WaitGroup
var cpunum = runtime.NumCPU()

func main() {
	Banner()
	starttime := time.Now()
	defer testtime(starttime)
	baseurl := "http://2js56.top/"
	urls := make(chan string)
	file, err := os.Open("dic.txt")
	if err != nil {
		fmt.Println("字典打开失败: ", err)
		return
	}
	defer file.Close()
	//开启一个携程，用来向urls通道中写入完整的url
	go func() {
		scaner := bufio.NewScanner(file)
		for scaner.Scan() {
			urls <- fmt.Sprintf("%s/%s", baseurl, scaner.Text())
		}
		if err := scaner.Err(); err != nil {
			fmt.Println("读取字典错误: ", err)
		}
		close(urls)
		fmt.Println("字典读取完毕。")
	}()

	//开启多个携程，通过http的head请求方式，进行验证
	worknum := (cpunum - 1) * 20
	for i := 0; i < worknum; i++ {
		wg.Add(1)
		go Sendurl(urls)
	}
	wg.Wait()

}

//发送http请求，并检查返回状态，输出2xx，3xx状态码url
func Sendurl(urls chan string) {
	for {
		select {
		case url, ok := <-urls:
			if !ok {
				wg.Done()
				return
			}
			client := &http.Client{
				Timeout: time.Second * 10,
			}
			req, err := http.NewRequest("HEAD", url, nil)
			if err != nil {
				//fmt.Println(err)
				continue
			}
			req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.87 Safari/537.36")
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println(err)
				continue
			}
			match, err := regexp.MatchString(`^(2|3)`, strconv.Itoa(resp.StatusCode))
			if err != nil {
				fmt.Println(err)
				continue
			}
			if match {
				fmt.Printf("%s 状态码: %v\n", url, resp.StatusCode)
			}
			resp.Body.Close()
		case <-time.After(time.Second * 60):
			wg.Done()
			return
		}
	}
}
func testtime(start time.Time) {
	fmt.Println("运行时间：", time.Since(start))
}

func Banner() {
	banner := `
	######   ##     ## ##    ##  ######  ########    ###    ########   ######  ##     ## 
	##    ## ##     ## ###   ## ##    ## ##         ## ##   ##     ## ##    ## ##     ## 
	##       ##     ## ####  ## ##       ##        ##   ##  ##     ## ##       ##     ## 
	 ######  ##     ## ## ## ##  ######  ######   ##     ## ########  ##       ######### 
		  ## ##     ## ##  ####       ## ##       ######### ##   ##   ##       ##     ## 
	##    ## ##     ## ##   ### ##    ## ##       ##     ## ##    ##  ##    ## ##     ## 
	 ######   #######  ##    ##  ######  ######## ##     ## ##     ##  ######  ##     ## 
	 																--version: 1.0
																	--author : 2JS56
`
	fmt.Println(banner)
}
