package html

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type Origin struct {
	Base *url.URL
}

func (o Origin) getBase() string {
	return fmt.Sprintf("%v://%v", o.Base.Scheme, o.Base.Host)
}

func (o Origin) getAbsolute(relativeUrl string) string {
	absoluteUrl, err := o.Base.Parse(relativeUrl)
	if err != nil {
		return ""
	}

	return absoluteUrl.String()
}

type ATag struct {
	token  html.Token
	origin Origin
}

func (t ATag) getHref() (href string, ok bool) {
	for _, a := range t.token.Attr {
		if a.Key == "href" {
			href = a.Val
			break
		}
	}

	// For http(s) urls
	if strings.Index(href, "http") == 0 {
		return href, true
	}

	// For "/" and "" shortcuts
	if href == "/" || href == "" {
		return t.origin.getBase(), true
	}

	// For relative urls starts with slash
	if strings.Index(href, "/") == 0 {
		return fmt.Sprintf("%v%v", t.origin.getBase(), href), true
	}

	// For relative urls starts with dot
	if strings.Index(href, ".") == 0 {
		return t.origin.getAbsolute(href), true
	}

	return href, false
}
