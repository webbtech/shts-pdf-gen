package model

// DbHandler interface
type DbHandler interface {
	Close()
	FetchEstimate(string) (*Estimate, error)
}
