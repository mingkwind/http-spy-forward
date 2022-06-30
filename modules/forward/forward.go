package forward

import (
	"http-spy-forward/logger"
	"http-spy-forward/models"
	"http-spy-forward/vars"
	"net/http"
	"strings"
)

func Forward(httpReq *models.HttpReq) {
	header := httpReq.Header
	header.Add("X-Real-IP", "")
	header.Set("X-Real-IP", httpReq.Client)
	// 测试发包
	/*
		uri, _ := url.Parse("http://127.0.0.1:8080")
		client := &http.Client{
			Transport: &http.Transport{
				// 设置代理
				Proxy: http.ProxyURL(uri),
			},
		}*/
	client := &http.Client{}
	switch httpReq.Method {
	case "GET":
		// req, _ := http.NewRequest("GET", "http://httpbin.org/get", nil)
		req, _ := http.NewRequest("GET", vars.TargetURL+httpReq.RequestURI, nil)
		req.Header = header
		resp, err := client.Do(req)
		_ = resp
		if err != nil {
			logger.Log.Error("Forward GET error:", err)
		}
		//body, _ := ioutil.ReadAll(resp.Body)
		//fmt.Println(string(body))
	case "POST":
		req, _ := http.NewRequest("POST", vars.TargetURL+httpReq.RequestURI, strings.NewReader(httpReq.Body))
		req.Header = header
		resp, err := client.Do(req)
		_ = resp
		if err != nil {
			logger.Log.Error("Forward POST error:", err)
		}
		//body, _ := ioutil.ReadAll(resp.Body)
		//fmt.Println(string(body))
	}
}
