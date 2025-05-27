package crawler

type ParserType string

type Parser interface {
	GetType() ParserType
	Parse(resource FetchedResource) ParsedData
}

type ParsedResource struct {
	resource      FetchedResource
	parcingResult []ParsedData
}

func (r ParsedResource) GetResource() FetchedResource {
	return r.resource
}

func (p *ParsedResource) AddResults(results []ParsedData) {
	p.parcingResult = append(p.parcingResult, results...)
}

func (p *ParsedResource) GetResults() []ParsedData {
	return p.parcingResult
}

func NewParsedResource(resource FetchedResource, results []ParsedData) ParsedResource {
	return ParsedResource{resource: resource, parcingResult: results}
}

type ParsedData struct {
	parser ParserType
	data   []string
	err    error
}

func NewParsedData(parser ParserType, data []string, err error) ParsedData {
	return ParsedData{parser: parser, data: data, err: err}
}

func (p *ParsedData) IsSuccess() bool {
	return p.err == nil
}

func (p *ParsedData) GetParserType() ParserType {
	return p.parser
}

func (p *ParsedData) GetData() []string {
	return p.data
}

func (p *ParsedData) AppendData(rows []string) {
	p.data = append(p.data, rows...)
}

func (p *ParsedData) SetError(err error) {
	p.err = err
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
