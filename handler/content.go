package handler

// ContentHandler defines an interface for handling content sources that are
// represented as URLs. When the handler will finish its job it should notify
// the progress. Each error should be sent to the errors channel.
type ContentHandler interface {
	Process(url string, progress chan<- int, errors chan<- error)
}
