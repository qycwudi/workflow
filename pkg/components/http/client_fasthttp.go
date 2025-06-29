package http

import (
	"fmt"
	"net"
	"net/url"
	"strings"
	"sync"
	"time"
	
	"github.com/bytedance/sonic"
	"github.com/valyala/fasthttp"
	"github.com/zeromicro/go-zero/core/logx"
	
	"workflow/pkg/errors"
)

// FastHTTPStrategy FastHTTP实现策略
type FastHTTPStrategy struct {
	client     *fasthttp.Client
	timeout    time.Duration
	maxRetries int
}

// NewFastHTTPStrategy 创建FastHTTP策略
func NewFastHTTPStrategy(timeout time.Duration, maxRetries int) *FastHTTPStrategy {
	client := &fasthttp.Client{
		MaxConnsPerHost:         100,
		MaxIdleConnDuration:     90 * time.Second,
		MaxConnDuration:         10 * time.Minute,
		MaxConnWaitTimeout:      30 * time.Second,
		ReadTimeout:             timeout,
		WriteTimeout:            timeout,
		DisableHeaderNamesNormalizing: true,
		NoDefaultUserAgentHeader:      true,
		DisablePathNormalizing:        true,
		Dial: func(addr string) (net.Conn, error) {
			return fasthttp.DialTimeout(addr, 5*time.Second)
		},
	}
	
	return &FastHTTPStrategy{
		client:     client,
		timeout:    timeout,
		maxRetries: maxRetries,
	}
}

// Name 策略名称
func (s *FastHTTPStrategy) Name() string {
	return "fasthttp"
}

// DoRequest 执行HTTP请求
func (s *FastHTTPStrategy) DoRequest(opts RequestOptions) (int, []byte, error) {
	req := requestPool.Get().(*fasthttp.Request)
	resp := responsePool.Get().(*fasthttp.Response)
	defer requestPool.Put(req)
	defer responsePool.Put(resp)
	
	req.Reset()
	resp.Reset()
	
	// 设置请求
	if err := s.prepareRequest(req, opts); err != nil {
		return 0, nil, err
	}
	
	// 执行请求
	return s.executeRequest(req, resp, opts)
}

// prepareRequest 准备请求
func (s *FastHTTPStrategy) prepareRequest(req *fasthttp.Request, opts RequestOptions) error {
	req.SetRequestURI(opts.URL)
	req.Header.SetMethod(opts.Method)
	
	// 设置默认请求头
	req.Header.Set("Connection", "close")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; WorkflowBot/1.0)")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")
	
	// 设置自定义请求头
	for key, value := range opts.Headers {
		req.Header.Set(key, value)
	}
	
	// 设置请求体
	if opts.Body != nil {
		return s.setRequestBody(req, opts)
	}
	
	return nil
}

// setRequestBody 设置请求体
func (s *FastHTTPStrategy) setRequestBody(req *fasthttp.Request, opts RequestOptions) error {
	var bodyBytes []byte
	var err error
	
	if opts.IsJSON {
		bodyBytes, err = sonic.Marshal(opts.Body)
		if err != nil {
			return errors.ValidationError("http", "failed to marshal JSON body").
				Cause(err).
				Build()
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		if formData, ok := opts.Body.(map[string]any); ok {
			data := url.Values{}
			for k, v := range formData {
				data.Set(k, fmt.Sprintf("%v", v))
			}
			bodyBytes = []byte(data.Encode())
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		} else {
			return errors.ValidationError("http", "invalid form body format").Build()
		}
	}
	
	req.SetBody(bodyBytes)
	return nil
}

// executeRequest 执行请求
func (s *FastHTTPStrategy) executeRequest(req *fasthttp.Request, resp *fasthttp.Response, opts RequestOptions) (int, []byte, error) {
	var lastErr error
	
	for i := 0; i <= s.maxRetries; i++ {
		logx.Debugf("[HTTP] FastHTTP request attempt %d, timeout: %v", i+1, s.timeout)
		
		resp.Reset()
		
		err := s.client.DoTimeout(req, resp, s.timeout)
		if err != nil {
			lastErr = err
			logx.Errorw("[HTTP] FastHTTP request failed",
				logx.Field("url", opts.URL),
				logx.Field("method", opts.Method),
				logx.Field("attempt", i+1),
				logx.Field("error", err))
			
			// 连接错误直接返回
			if s.isConnectionError(err) {
				return 0, nil, errors.NetworkError("http", "connection failed").
					Context("url", opts.URL).
					Context("attempt", fmt.Sprintf("%d", i+1)).
					Cause(err).
					Build()
			}
			
			// 重试前等待
			if i < s.maxRetries {
				time.Sleep(time.Duration(i+1) * 300 * time.Millisecond)
			}
			continue
		}
		
		// 检查响应状态码
		statusCode := resp.StatusCode()
		if statusCode >= 200 && statusCode < 300 {
			bodyBytes := make([]byte, len(resp.Body()))
			copy(bodyBytes, resp.Body())
			return statusCode, bodyBytes, nil
		} else if statusCode >= 400 && statusCode < 500 {
			bodyBytes := make([]byte, len(resp.Body()))
			copy(bodyBytes, resp.Body())
			return statusCode, bodyBytes, errors.ExecutionError("http", "client error").
				Context("status_code", fmt.Sprintf("%d", statusCode)).
				Context("url", opts.URL).
				Build()
		} else {
			lastErr = errors.ExecutionError("http", "server error").
				Context("status_code", fmt.Sprintf("%d", statusCode)).
				Context("url", opts.URL).
				Build()
			
			if i < s.maxRetries {
				time.Sleep(time.Duration(i+1) * 200 * time.Millisecond)
			}
		}
	}
	
	return 0, nil, errors.ExecutionError("http", "request failed after retries").
		Context("max_retries", fmt.Sprintf("%d", s.maxRetries)).
		Context("url", opts.URL).
		Cause(lastErr).
		Build()
}

// isConnectionError 判断是否为连接错误
func (s *FastHTTPStrategy) isConnectionError(err error) bool {
	errStr := err.Error()
	return strings.Contains(errStr, "connection") ||
		strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "closed") ||
		strings.Contains(errStr, "first response byte")
}

// 全局对象池
var (
	requestPool = sync.Pool{
		New: func() interface{} {
			return &fasthttp.Request{}
		},
	}
	
	responsePool = sync.Pool{
		New: func() interface{} {
			return &fasthttp.Response{}
		},
	}
)