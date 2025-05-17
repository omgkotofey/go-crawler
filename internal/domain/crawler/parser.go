package crawler

type ParserType string

type Parser interface {
	GetType() ParserType
	Parse(resource FetchedResource) ParsedData
}

type ParserSet struct {
	parsers map[ParserType]Parser
}

func (p *ParserSet) Parse(resource FetchedResource) ParsedResource {
	parsedData := make([]ParsedData, 0)
	for parserName := range p.parsers {
		result := p.parsers[parserName].Parse(resource)
		parsedData = append(parsedData, result)
	}

	return NewParsedResource(resource, parsedData)
}

func (p *ParserSet) AddParser(parser Parser) {
	if p.parsers == nil {
		p.parsers = make(map[ParserType]Parser)
	}
	p.parsers[parser.GetType()] = parser
}

func (p *ParserSet) GetParsers() []Parser {
	result := make([]Parser, 0)

	for parserName := range p.parsers {
		result = append(result, p.parsers[parserName])
	}

	return result
}
