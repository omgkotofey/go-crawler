package app

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

	if strings.Index(href, "http") == 0 {
		return href, true
	}

	if href == "/" || href == "" {
		return t.origin.getBase(), true
	}

	if strings.Index(href, "/") == 0 {
		return fmt.Sprintf("%v%v", t.origin.getBase(), href), true
	}

	return href, false
}

type ParseResult struct {
	FetchResult
	urls []*url.URL
}

type UrlParser interface {
	Parse(data FetchResult) (result ParseResult, err error)
}

type TokenizerParser struct {
	Origin Origin
}

func (p TokenizerParser) Parse(data FetchResult) (result ParseResult, err error) {
	tokenizer := html.NewTokenizer(bytes.NewReader(data.body))
	result = ParseResult{
		FetchResult: data,
		urls:        make([]*url.URL, 0),
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

			result.urls = append(result.urls, parsedUrl)
		}
	}
}
