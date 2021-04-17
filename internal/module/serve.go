// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package module

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	log "unknwon.dev/clog/v2"
)

func (m *module) DoRequest(req *http.Request, body []byte) {
	if !m.isRequestHit(req, body) {
		return
	}

	if m.Req.TransmitURL != nil {
		req.URL = m.Req.TransmitURL
	}

	for k, v := range m.Req.Header {
		req.Header.Set(k, v)
	}

	// Load body from local file
	if filePath, ok := m.Req.Body["file"].(string); ok && filePath != "" {
		newBody, err := os.ReadFile(filePath)
		if err != nil {
			log.Error("Failed to load body file: %v", err)
			return
		}
		body = newBody
	}
	// Replace body
	if replacement, ok := m.Req.Body["replace"].(map[string]string); ok {
		bodyString := string(body)
		for k, v := range replacement {
			bodyString = strings.ReplaceAll(bodyString, k, v)
		}
		body = []byte(bodyString)
	} else if newBody, ok := m.Req.Body["replace"].(string); ok {
		body = []byte(newBody)
	}

	// Write back
	req.Body = io.NopCloser(bytes.NewBuffer(body))
}

func (m *module) isRequestHit(req *http.Request, body []byte) bool {
	result, _, err := m.Req.OnPrg.Eval(map[string]interface{}{
		"method":  req.Method,
		"url":     req.URL.String(),
		"headers": req.Header,
		"body":    string(body),
	})
	if err != nil {
		log.Error("Check module request active error: %v", err)
		return false
	}

	switch v := result.Value().(type) {
	case bool:
		return v
	case string:
		return v == req.URL.String()
	default:
		return false
	}
}

func (m *module) DoResponse(resp *http.Response, body []byte) {
	if !m.isResponseHit(resp, body) {
		return
	}

	// Status code
	if m.Resp.StatusCode != 0 {
		resp.StatusCode = m.Resp.StatusCode
	}

	// Response header
	for k, v := range m.Resp.Header {
		resp.Header.Set(k, v)
	}

	// Load body from local file
	if filePath, ok := m.Resp.Body["file"].(string); ok && filePath != "" {
		newBody, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Error("Failed to load local file: %v", err)
			return
		}
		body = newBody
	}

	// Replace body
	if replacement, ok := m.Resp.Body["replace"].(map[interface{}]interface{}); ok {
		bodyString := string(body)
		for k, v := range replacement {
			key, keyOk := k.(string)
			value, valueOk := v.(string)
			if keyOk && valueOk {
				bodyString = strings.ReplaceAll(bodyString, key, value)
			}
		}
		body = []byte(bodyString)
	} else if newContent, ok := m.Resp.Body["replace"].(string); ok {
		body = []byte(newContent)
	}

	// Write back
	resp.Body = io.NopCloser(bytes.NewBuffer(body))
}

func (m *module) isResponseHit(resp *http.Response, body []byte) bool {
	result, _, err := m.Resp.OnPrg.Eval(map[string]interface{}{
		"url":         resp.Request.URL.String(),
		"status_code": resp.StatusCode,
		"headers":     resp.Header,
		"body":        string(body),
	})
	if err != nil {
		log.Error("Check module response active error: %v", err)
		return false
	}

	switch v := result.Value().(type) {
	case bool:
		return v
	default:
		return false
	}
}
