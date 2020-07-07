package httplib

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// 请求类型
const (
	METHOD_GET  = "GET"
	METHOD_POST = "POST"
)

// post类型
const (
	POST_FORM = iota + 1
	POST_JSON
	POST_XML
)

// 默认超时时间 ms
const (
	TIMEOUT_CONN = 1000
	TIMEOUT_RW   = 5000
)

// 请求结果状态
const (
	STATUS_WRONG   = -2
	STATUS_FAIL    = -1
	STATUS_SUCCESS = 200
)

const (
	ERR_400 = "Bad Request"
	ERR_401 = "Unauthorized"
	ERR_403 = "Forbidden"
	ERR_404 = "Not Found"
	ERR_405 = "Method Not Allowed"
	ERR_500 = "Internal Server Error"
	ERR_502 = "Bad Gateway"
	ERR_503 = "Service Unavailable"
	ERR_504 = "Gateway Time-out"
	ERR_505 = "HTTP Version not supported"
)

type Request struct {
	Method      string
	Url         string
	Headers     map[string]string
	Cookie      string
	Ssl         bool
	ConnTimeout int
	RWTimeout   int
	PostType    int
	PostParams  interface{}
}

type Response struct {
	Status int
	ErrMsg string
	Body   string
}

// Get请求, arg 为超时时间，默认5s
func Get(url string, arg ...int) (string, error) {
	timeout := TIMEOUT_RW
	if arg != nil {
		timeout = arg[0]
	}
	req := Request{
		Method:    METHOD_GET,
		Url:       url,
		RWTimeout: timeout,
	}
	resp := Call(req)
	if resp.Status != STATUS_SUCCESS {
		return "", errors.New(resp.ErrMsg)
	}
	return resp.Body, nil
}

// POST请求, arg 为超时时间，默认5s
func Post(url string, params interface{}, arg ...int) (string, error) {
	timeout := TIMEOUT_RW
	if arg != nil {
		timeout = arg[0]
	}
	req := Request{
		Method:     METHOD_POST,
		Url:        url,
		RWTimeout:  timeout,
		PostType:   POST_FORM,
		PostParams: params,
	}
	resp := Call(req)
	if resp.Status != STATUS_SUCCESS {
		return "", errors.New(resp.ErrMsg)
	}
	return resp.Body, nil
}

// 单次请求
func Call(r Request) (resp *Response) {
	resp = &Response{}
	if req, msg, ok := check(r); !ok {
		resp.Status = STATUS_WRONG
		resp.ErrMsg = msg
	} else {
		send(req, resp, nil, nil)
	}
	return
}

// 并发请求
var max_concurrency = 500
var default_concurrency = 10

// arg 为并发数量，默认10， 最大500
func MultiCall(reqs map[interface{}]Request, arg ...int) map[interface{}]*Response {
	resp := make(map[interface{}]*Response)
	if len(reqs) == 0 {
		return resp
	}

	concurrency := default_concurrency
	if arg != nil {
		num := arg[0]
		if num > 0 {
			concurrency = num
		}
		if num > max_concurrency {
			concurrency = max_concurrency
		}
	}

	ch := make(chan int, concurrency)
	wg := &sync.WaitGroup{}
	for idx, req := range reqs {
		if req, msg, ok := check(req); ok {
			resp[idx] = &Response{}
			ch <- 1
			wg.Add(1)
			go send(req, resp[idx], ch, wg)
		} else {
			resp[idx] = &Response{Status: STATUS_WRONG, ErrMsg: msg}
		}
	}
	wg.Wait()
	close(ch)
	return resp
}

// 发送请求
func send(req *Request, resp *Response, ch chan int, wg *sync.WaitGroup) {

	defer func() {
		if ch != nil {
			<-ch
		}
		if wg != nil {
			wg.Done()
		}
	}()

	request, err := http.NewRequest(req.Method, req.Url, nil)
	if err != nil {
		resp.Status = STATUS_FAIL
		resp.ErrMsg = err.Error()
		return
	}

	// 设置请求头
	for idx, value := range req.Headers {
		request.Header.Set(idx, value)
	}

	if len(req.Cookie) > 0 {
		request.Header.Set("Cookie", req.Cookie)
	}

	if req.Method == METHOD_POST {
		switch req.PostType {
		case POST_FORM:
			reqBody := &bytes.Buffer{}
			writer := multipart.NewWriter(reqBody)
			params, _ := req.PostParams.(map[string]string)
			for idx, val := range params {
				if len(val) > 0 && val[0] == '@' {
					file, err := os.Open(val[1:])
					if err != nil {
						resp.Status = STATUS_FAIL
						resp.ErrMsg = err.Error()
						return
					}
					part, _ := writer.CreateFormFile(idx, filepath.Base(val[1:]))
					_, err = io.Copy(part, file)
					if err != nil {
						resp.Status = STATUS_FAIL
						resp.ErrMsg = err.Error()
						return
					}
					file.Close()
				} else {
					writer.WriteField(idx, val)
				}
			}
			writer.Close()
			request.ContentLength = int64(reqBody.Len())
			request.Body = ioutil.NopCloser(reqBody)
			request.Header.Set("content-type", writer.FormDataContentType())

		case POST_XML:
			param, _ := req.PostParams.(string)
			request.ContentLength = int64(len(param))
			request.Header.Set("content-type", "text/xml")
			request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(param)))

		case POST_JSON:
			param, _ := req.PostParams.(string)
			request.ContentLength = int64(len(param))
			request.Header.Set("content-type", "application/json")
			request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(param)))
		}
	}

	// 超时设置
	tr := &http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			conn, err := net.DialTimeout(network, addr, time.Millisecond*time.Duration(req.ConnTimeout))
			if err != nil {
				return nil, err
			}
			conn.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(req.RWTimeout)))
			return conn, nil
		},
	}

	if req.Ssl {
		tr.TLSHandshakeTimeout = time.Millisecond * time.Duration(req.ConnTimeout)
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	client := &http.Client{
		Transport: tr,
	}
	response, err := client.Do(request)
	if err != nil {
		resp.Status = STATUS_FAIL
		resp.ErrMsg = err.Error()
		return
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		resp.Status = STATUS_FAIL
		resp.ErrMsg = err.Error()
		return
	}

	resp.Status = response.StatusCode
	switch response.StatusCode {
	case 400:
		resp.ErrMsg = ERR_400
	case 401:
		resp.ErrMsg = ERR_401
	case 403:
		resp.ErrMsg = ERR_403
	case 404:
		resp.ErrMsg = ERR_404
	case 405:
		resp.ErrMsg = ERR_405
	case 500:
		resp.ErrMsg = ERR_500
	case 502:
		resp.ErrMsg = ERR_502
	case 503:
		resp.ErrMsg = ERR_503
	case 504:
		resp.ErrMsg = ERR_504
	case 505:
		resp.ErrMsg = ERR_505
	}

	if resp.Status == STATUS_SUCCESS {
		resp.Body = string(body)
	}
	return
}

// 请求参数验证
func check(req Request) (*Request, string, bool) {
	request := &Request{
		Headers: req.Headers,
		Url:     req.Url,
		Cookie:  req.Cookie,
	}

	if req.Method != METHOD_GET && req.Method != METHOD_POST {
		return nil, "method not valid", false
	}
	request.Method = req.Method

	if request.Method == METHOD_POST {
		switch req.PostType {
		case POST_FORM:
		case POST_JSON:
		case POST_XML:
		default:
			return nil, "posttype not valid", false
		}
		request.PostType = req.PostType
		request.PostParams = req.PostParams
	}

	request.Ssl = req.Ssl
	if strings.Contains(req.Url, "https://") {
		request.Ssl = true
	}

	request.ConnTimeout = req.ConnTimeout
	if req.ConnTimeout <= 0 {
		request.ConnTimeout = TIMEOUT_CONN
	}

	request.RWTimeout = req.RWTimeout
	if req.RWTimeout <= 0 {
		request.RWTimeout = TIMEOUT_RW
	}

	return request, "", true
}

// 对response删除指定的cookie
func DeleteCookies(w http.ResponseWriter, cookieNames ...string) {
	for _, n := range cookieNames {
		c := &http.Cookie{
			Name:    n,
			Value:   "",
			Path:    "/",
			MaxAge:  -1,
			Expires: time.Time{},
		}
		http.SetCookie(w, c)
	}
}
