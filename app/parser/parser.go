package parser

import (
	"bytes"
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

type ResourcePresenter interface {
	GetUrl() *url.URL
	GetBody() []byte
	GetSpent() float64
}

type ParseResult struct {
	Resource ResourcePresenter
	Urls     []*url.URL
}

type UrlParser interface {
	Parse(resource ResourcePresenter) (result ParseResult, err error)
}

type TokenizerParser struct {
	Origin Origin
}

func (p TokenizerParser) Parse(resource ResourcePresenter) (result ParseResult, err error) {
	tokenizer := html.NewTokenizer(bytes.NewReader(resource.GetBody()))
	result = ParseResult{
		Resource: resource,
		Urls:     make([]*url.URL, 0),
	}

	for {
		token := tokenizer.Next()

		switch {
		case token == html.ErrorToken:
			// End of the document, we're done
			return result, nil
		case token == html.StartTagToken:
			t := tokenizer.Token()

			// Check if the token is an <a> tag
			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}

			// Extract the href value, if there is one
			atag := ATag{token: t, origin: p.Origin}
			rawUrl, ok := atag.getHref()
			if !ok {
				continue
			}

			// Make sure the url begines in http**
			isHTTP := strings.Index(rawUrl, "http") == 0
			if !isHTTP {
				continue
			}

			// check url has same base with origin
			if !strings.Contains(rawUrl, p.Origin.getBase()) {
				continue
			}

			parsedUrl, err := url.ParseRequestURI(rawUrl)
			if err != nil {
				return result, err
			}

			result.Urls = append(result.Urls, parsedUrl)
		}
	}
}
