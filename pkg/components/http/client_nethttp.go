package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
	
	"github.com/bytedance/sonic"
	"github.com/zeromicro/go-zero/core/logx"
	
	"workflow/pkg/errors"
)

// NetHTTPStrategy 标准net/http实现策略
type NetHTTPStrategy struct {
	client     *http.Client
	timeout    time.Duration
	maxRetries int
}

// NewNetHTTPStrategy 创建NetHTTP策略
func NewNetHTTPStrategy(timeout time.Duration, maxRetries int) *NetHTTPStrategy {
	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			DisableKeepAlives:   true,
			MaxIdleConns:        1,
			MaxIdleConnsPerHost: 1,
			IdleConnTimeout:     1 * time.Second,
		},
	}
	
	return &NetHTTPStrategy{
		client:     client,
		timeout:    timeout,
		maxRetries: maxRetries,
	}
}

// Name 策略名称
func (s *NetHTTPStrategy) Name() string {
	return "net/http"
}

// DoRequest 执行HTTP请求
func (s *NetHTTPStrategy) DoRequest(opts RequestOptions) (int, []byte, error) {
	req, err := s.createRequest(opts)
	if err != nil {
		return 0, nil, err
	}
	
	return s.executeRequest(req, opts)
}

// createRequest 创建HTTP请求
func (s *NetHTTPStrategy) createRequest(opts RequestOptions) (*http.Request, error) {
	var body *bytes.Buffer
	
	if opts.Body != nil {
		if opts.IsJSON {
			jsonData, err := sonic.Marshal(opts.Body)
			if err != nil {
				return nil, errors.ValidationError("http", "failed to marshal JSON body").
					Cause(err).
					Build()
			}
			body = bytes.NewBuffer(jsonData)
		} else {
			if formData, ok := opts.Body.(map[string]any); ok {
				data := url.Values{}
				for k, v := range formData {
					data.Set(k, fmt.Sprintf("%v", v))
				}
				body = bytes.NewBufferString(data.Encode())
			} else {
				return nil, errors.ValidationError("http", "invalid form body format").Build()
			}
		}
	}
	
	var req *http.Request
	var err error
	
	if body != nil {
		req, err = http.NewRequest(opts.Method, opts.URL, body)
	} else {
		req, err = http.NewRequest(opts.Method, opts.URL, nil)
	}
	
	if err != nil {
		return nil, errors.ValidationError("http", "failed to create request").
			Cause(err).
			Build()
	}
	
	// 设置请求头
	s.setRequestHeaders(req, opts)
	
	return req, nil
}

// setRequestHeaders 设置请求头
func (s *NetHTTPStrategy) setRequestHeaders(req *http.Request, opts RequestOptions) {
	// 设置默认请求头
	req.Header.Set("Connection", "close")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; WorkflowBot/1.0)")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Cache-Control", "no-cache")
	
	// 设置Content-Type
	if opts.Body != nil {
		if opts.IsJSON {
			req.Header.Set("Content-Type", "application/json")
		} else {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	}
	
	// 设置自定义请求头
	for key, value := range opts.Headers {
		req.Header.Set(key, value)
	}
}

// executeRequest 执行请求
func (s *NetHTTPStrategy) executeRequest(req *http.Request, opts RequestOptions) (int, []byte, error) {
	var lastErr error
	
	for i := 0; i <= s.maxRetries; i++ {
		logx.Debugf("[HTTP] NetHTTP request attempt %d", i+1)
		
		resp, err := s.client.Do(req)
		if err != nil {
			lastErr = err
			logx.Errorw("[HTTP] NetHTTP request failed",
				logx.Field("url", opts.URL),
				logx.Field("method", opts.Method),
				logx.Field("attempt", i+1),
				logx.Field("error", err))
			
			if i < s.maxRetries {
				time.Sleep(time.Duration(i+1) * 500 * time.Millisecond)
			}
			continue
		}
		
		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = err
			continue
		}
		
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return resp.StatusCode, respBody, nil
		} else if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			return resp.StatusCode, respBody, errors.ExecutionError("http", "client error").
				Context("status_code", fmt.Sprintf("%d", resp.StatusCode)).
				Context("url", opts.URL).
				Build()
		} else {
			lastErr = errors.ExecutionError("http", "server error").
				Context("status_code", fmt.Sprintf("%d", resp.StatusCode)).
				Context("url", opts.URL).
				Build()
			
			if i < s.maxRetries {
				time.Sleep(time.Duration(i+1) * 500 * time.Millisecond)
			}
		}
	}
	
	return 0, nil, errors.ExecutionError("http", "request failed after retries").
		Context("max_retries", fmt.Sprintf("%d", s.maxRetries)).
		Context("url", opts.URL).
		Cause(lastErr).
		Build()
}