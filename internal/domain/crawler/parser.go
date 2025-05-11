package crawler

type Parser interface {
	Parse(resource FetchedResource) (result ParsedData, err error)
}
