package chromite

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

type Attachment struct {
	Name         string
	SourceURL    string
	TotalBytes   float64
	RecivedBytes float64
}

type TabProduct struct {
	Requests  map[string]*network.Request
	Responses map[string]*network.Response
	Downloads map[string]*Attachment
	Logs      []string
}

type Chrome struct {
	ctx       context.Context
	cancel    context.CancelFunc
	cachePath string
}

func (c *Chrome) NewTab(url *url.URL, timeout time.Duration, f func(ev any, product *TabProduct) error, actions ...chromedp.Action) (*TabProduct, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	if timeout > 0 {
		ctx, _ = context.WithTimeout(c.ctx, timeout)
	}
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()
	product := &TabProduct{
		Requests:  map[string]*network.Request{},
		Responses: map[string]*network.Response{},
		Downloads: map[string]*Attachment{},
		Logs:      []string{},
	}

	var err error
	l := func(ev interface{}) {

		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("%s", r)
			}
		}()

		switch ev := ev.(type) {
		case *network.EventRequestWillBeSent:
			product.Requests[ev.RequestID.String()] = ev.Request
		case *network.EventResponseReceived:
			product.Responses[ev.RequestID.String()] = ev.Response
		case *page.EventJavascriptDialogOpening:
			go func() {
				chromedp.Run(ctx,
					page.HandleJavaScriptDialog(true),
				)
			}()
		case *browser.EventDownloadWillBegin:
			product.Downloads[ev.GUID] = &Attachment{
				Name:      ev.SuggestedFilename,
				SourceURL: ev.URL,
			}
		case *browser.EventDownloadProgress:
			switch ev.State {
			case browser.DownloadProgressStateInProgress:
				product.Downloads[ev.GUID].TotalBytes = ev.TotalBytes
				product.Downloads[ev.GUID].RecivedBytes = ev.ReceivedBytes
			case browser.DownloadProgressStateCompleted:
			case browser.DownloadProgressStateCanceled:
			}
		case *runtime.EventConsoleAPICalled:
			var log string
			if ev.Type == "log" {
				for _, a := range ev.Args {
					switch a.Type {
					case "string":
						log += string(a.Value)[1 : len(a.Value)-1]
					case "number":
						log += string(a.Value)
					case "object":
						v, _ := a.MarshalJSON()
						log += string(v)
					default:
						v, _ := a.MarshalJSON()
						log += "{@" + "UnknownType " + string(a.Type) + "@Value " + string(v) + "}"
					}
				}
				product.Logs = append(product.Logs, log)
			} else {
				break
			}
		}
		if f != nil {
			err = f(ev, product)
		}
	}

	chromedp.ListenTarget(ctx, l)

	actions = append([]chromedp.Action{
		chromedp.Navigate(url.String()),
	}, actions...)

	chromedp.Run(ctx, actions...)

	return product, err
}

func (c *Chrome) CachePath() string {
	return c.cachePath
}

func (c *Chrome) Close() {
	c.cancel()
	os.RemoveAll(c.cachePath)
}

func NewChrome(ctx context.Context, cache string, opts ...func(*chromedp.ExecAllocator)) (*Chrome, error) {
	c := &Chrome{}

	if cache == "" || !path.IsAbs(cache) {
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		c.cachePath = path.Join(wd, cache, "chrome-cache")
	} else {
		c.cachePath = path.Join(cache, "chrome-cache")
	}

	opts = append(opts, chromedp.UserDataDir(c.cachePath))
	opts = append(opts, chromedp.Flag("disk-cache-dir", c.cachePath))
	opts = append(opts, chromedp.NoFirstRun)
	ctx, _ = chromedp.NewExecAllocator(ctx, opts...)
	ctx, cancel := chromedp.NewContext(ctx)
	c.ctx = ctx
	c.cancel = cancel
	return c, chromedp.Run(ctx)
}
