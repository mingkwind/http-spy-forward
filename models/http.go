package models

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

// HttpReq结构体，用于存储http请求的信息
type HttpReq struct {
	Host          string
	Ip            string
	Client        string
	Port          string
	URL           *url.URL
	Header        http.Header
	RequestURI    string
	Method        string
	ReqParameters url.Values
	Body          string
}

func NewHttpReq(req *http.Request, client string, ip string, port string) (httpReq *HttpReq, err error) {
	// 解析请求参数
	body, _ := ioutil.ReadAll(req.Body)
	// 不能同时使用ioutil.ReadAll(req.Body)和req.ParseForm()
	// err = req.ParseForm()
	// 以下复用req.ParseForm()的源码
	if req.PostForm == nil {
		if req.Method == "POST" || req.Method == "PUT" || req.Method == "PATCH" {
			req.PostForm, err = url.ParseQuery(string(body))
		}
		if req.PostForm == nil {
			req.PostForm = make(url.Values)
		}
	}
	if req.Form == nil {
		if len(req.PostForm) > 0 {
			req.Form = make(url.Values)
			copyValues(req.Form, req.PostForm)
		}
		var newValues url.Values
		if req.URL != nil {
			var e error
			newValues, e = url.ParseQuery(req.URL.RawQuery)
			if err == nil {
				err = e
			}
		}
		if newValues == nil {
			newValues = make(url.Values)
		}
		if req.Form == nil {
			req.Form = newValues
		} else {
			copyValues(req.Form, newValues)
		}
	}

	return &HttpReq{Host: req.Host, Client: client, Ip: ip, Port: port, URL: req.URL, Header: req.Header,
		RequestURI: req.RequestURI, Method: req.Method, ReqParameters: req.Form, Body: string(body)}, err
}

func copyValues(dst, src url.Values) {
	for k, vs := range src {
		dst[k] = append(dst[k], vs...)
	}
}
