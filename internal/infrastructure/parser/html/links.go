package html

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	"experiments/internal/domain/crawler"
	"golang.org/x/net/html"
)

const LinksParserName crawler.ParserType = "links"

type Origin struct {
	Base *url.URL
}

func (o Origin) getBase() string {
	return fmt.Sprintf("%v://%v", o.Base.Scheme, o.Base.Host)
}

func (o Origin) getAbsolute(relativeURL string) string {
	absoluteURL, err := o.Base.Parse(relativeURL)
	if err != nil {
		return ""
	}

	return absoluteURL.String()
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

type LinksParser struct {
	Origin Origin
}

func NewLinksParser(origin *url.URL) LinksParser {
	return LinksParser{
		Origin: Origin{Base: origin},
	}
}

func (p LinksParser) GetType() crawler.ParserType {
	return LinksParserName
}

func (p LinksParser) Parse(resource crawler.FetchedResource) crawler.ParsedData {
	tokenizer := html.NewTokenizer(bytes.NewReader(resource.GetBody()))
	result := crawler.NewParsedData(p.GetType(), make([]string, 0), nil)

	for {
		token := tokenizer.Next()
		switch token {
		case html.ErrorToken:
			// End of the document, we're done
			return result
		case html.StartTagToken:
			t := tokenizer.Token()

			// Check if the token is an <a> tag
			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}

			// Extract the href value, if there is one
			atag := ATag{token: t, origin: p.Origin}
			rawURL, ok := atag.getHref()
			if !ok {
				continue
			}

			// Make sure the url begines in http**
			isHTTP := strings.Index(rawURL, "http") == 0
			if !isHTTP {
				continue
			}

			// check url has same base with origin
			if !strings.Contains(rawURL, p.Origin.getBase()) {
				continue
			}

			parsedURL, err := url.ParseRequestURI(rawURL)
			if err != nil {
				result.SetError(err)

				return result
			}

			result.AppendData([]string{parsedURL.String()})
		}
	}
}
