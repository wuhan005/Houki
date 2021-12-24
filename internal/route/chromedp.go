// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package route

import (
	"context"

	"github.com/chromedp/chromedp"
	log "unknwon.dev/clog/v2"
)

type ChromeDpHandler struct{}

func NewChromeDpHandler() *ChromeDpHandler {
	return &ChromeDpHandler{}
}

func (*ChromeDpHandler) New() {
	opts := []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.ProxyServer("localhost:8080"),
		chromedp.Flag("ignore-certificate-errors", "1"),
		func(a *chromedp.ExecAllocator) {
			chromedp.Flag("headless", false)(a)
		},
	}
	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	// Set up custom logger.
	taskCtx, _ := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Trace))

	go func() {
		_ = chromedp.Run(taskCtx)
	}()
}
