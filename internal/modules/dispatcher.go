// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package modules

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strings"

	log "unknwon.dev/clog/v2"
)

// DoRequest modify the request if the module was invoked.
func (m *Module) DoRequest(req *http.Request, body []byte) {
	if !m.isRequestHit(req, body) {
		return
	}

	if m.Req.TransmitURL != nil {
		req.URL = m.Req.TransmitURL
	}

	for k, v := range m.Req.Header {
		req.Header.Set(k, v)
	}

	// Modify request body with the local file.
	if filePath, ok := m.Req.Body["file"].(string); ok && filePath != "" {
		newBody, err := os.ReadFile(filePath)
		if err != nil {
			log.Error("Failed to load body file: %v", err)
			return
		}
		body = newBody
	}

	if replacement, ok := m.Resp.Body["replace"].(map[interface{}]interface{}); ok { // Replace body's keywords.
		bodyString := string(body)
		for k, v := range replacement {
			key, keyOk := k.(string)
			value, valueOk := v.(string)
			if keyOk && valueOk {
				bodyString = strings.ReplaceAll(bodyString, key, value)
			}
		}
		body = []byte(bodyString)
	} else if newBody, ok := m.Req.Body["replace"].(string); ok { // Replace the whole body.
		body = []byte(newBody)
	}

	// Write back the request body.
	req.Body = io.NopCloser(bytes.NewBuffer(body))
}

func (m *Module) isRequestHit(req *http.Request, body []byte) bool {
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
	case bool: // Condition.
		return v
	case string: // Specific URL.
		return v == req.URL.String()
	default:
		return false
	}
}

func (m *Module) DoResponse(resp *http.Response, body []byte) {
	if !m.isResponseHit(resp, body) {
		return
	}

	// Replace response status code.
	if m.Resp.StatusCode != 0 {
		resp.StatusCode = m.Resp.StatusCode
	}

	// Replace response header.
	for k, v := range m.Resp.Header {
		resp.Header.Set(k, v)
	}

	// Modify response body with the local file.
	if filePath, ok := m.Resp.Body["file"].(string); ok && filePath != "" {
		newBody, err := os.ReadFile(filePath)
		if err != nil {
			log.Error("Failed to load local file: %v", err)
			return
		}
		body = newBody
	}

	if replacement, ok := m.Resp.Body["replace"].(map[interface{}]interface{}); ok { // Replace body's keywords
		bodyString := string(body)
		for k, v := range replacement {
			key, keyOk := k.(string)
			value, valueOk := v.(string)
			if keyOk && valueOk {
				bodyString = strings.ReplaceAll(bodyString, key, value)
			}
		}
		body = []byte(bodyString)
	} else if newContent, ok := m.Resp.Body["replace"].(string); ok { // Replace the whole body.
		body = []byte(newContent)
	}

	// Write back the response body.
	resp.Body = io.NopCloser(bytes.NewBuffer(body))
}

func (m *Module) isResponseHit(resp *http.Response, body []byte) bool {
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
