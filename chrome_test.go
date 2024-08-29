package chromite

import (
	"context"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/B9O2/chromite/actions"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func TestChrome(t *testing.T) {
	c, err := NewChrome(context.Background(), "", chromedp.Headless)
	if err != nil {
		fmt.Println("New>", err)
		return
	}
	defer c.Close()
	target := "http://127.0.0.1:8888/click.html"
	//target := "https://www.baidu.com"
	u, err := url.Parse(target)
	if err != nil {
		fmt.Println("Parse>", err)
		return
	}
	// network.Enable()
	// network.SetExtraHTTPHeaders(t.headers)
	// browser.SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorDeny).WithDownloadPath(t.chrome.GetTmpPath())
	p, err := c.NewTab(u, 10*time.Second, func(ev any, product *TabProduct) error {
		switch ev := ev.(type) {
		case *page.EventJavascriptDialogOpening:
			fmt.Println("Dialog", ev.Message)
		}
		return nil
	}, actions.AutoClick())
	if err != nil {
		fmt.Println("NewTab>", err)
		return
	}
	fmt.Println("Requests:")
	for id, req := range p.Requests {
		fmt.Println(" ", id, req.URL)
	}
	fmt.Println("Total:", len(p.Requests))
	fmt.Println("Responses:")
	for id, resp := range p.Responses {
		fmt.Println(" ", id, resp.Status)
	}
	fmt.Println("Total:", len(p.Responses))
	fmt.Println("Downloads:")
	for id, file := range p.Downloads {
		fmt.Printf("%s %s(%s) %f/%f", id, file.Name, file.SourceURL, file.RecivedBytes, file.TotalBytes)
	}

	fmt.Println("Logs:")
	for id, log := range p.Logs {
		fmt.Println(id, log)
	}

}
