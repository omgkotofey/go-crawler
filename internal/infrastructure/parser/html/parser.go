package html

import (
	"bytes"
	"experiments/internal/domain/crawler"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type HTMLParser struct {
	Origin Origin
}

func NewHTMLParser(origin *url.URL) HTMLParser {
	return HTMLParser{
		Origin: Origin{Base: origin},
	}
}

func (p HTMLParser) Parse(resource crawler.FetchedResource) (result crawler.ParsedData, err error) {
	tokenizer := html.NewTokenizer(bytes.NewReader(resource.GetBody()))
	result = crawler.NewParsedData(resource, make([]string, 0))

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

			result.AppendData([]string{parsedUrl.String()})
		}
	}
}
