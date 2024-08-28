package chromite

import (
	"context"
	"fmt"
	"net/url"
	"testing"
	"time"
)

func TestChrome(t *testing.T) {
	c, err := NewChrome(context.Background(), "")
	if err != nil {
		fmt.Println("New>", err)
		return
	}
	defer c.Close()
	u, err := url.Parse("https://www.baidu.com")
	if err != nil {
		fmt.Println("Parse>", err)
		return
	}
	// network.Enable()
	// network.SetExtraHTTPHeaders(t.headers)
	// browser.SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorDeny).WithDownloadPath(t.chrome.GetTmpPath())
	p, err := c.NewTab(u, 5*time.Second, nil)
	if err != nil {
		fmt.Println("NewTab>", err)
		return
	}
	fmt.Println("Requests:")
	for id, req := range p.Requests {
		fmt.Println(" ", id, req.URL)
	}
	fmt.Println("Responses:")
	for id, resp := range p.Responses {
		fmt.Println(" ", id, resp.Status)
	}

	fmt.Println("Downloads:")
	for id, file := range p.Downloads {
		fmt.Printf("%s %s(%s) %f/%f", id, file.Name, file.SourceURL, file.RecivedBytes, file.TotalBytes)
	}

	time.Sleep(10 * time.Second)
}
