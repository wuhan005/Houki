// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package route

import (
	"context"

	"github.com/chromedp/chromedp"
	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/Houki/internal/proxy"
)

type ChromeDpHandler struct{}

func NewChromeDpHandler() *ChromeDpHandler {
	return &ChromeDpHandler{}
}

func (*ChromeDpHandler) New() {
	log.Trace(proxy.Address())
	opts := []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.ProxyServer(proxy.Address()),
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
