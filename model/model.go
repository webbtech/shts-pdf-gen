package model

// DbHandler interface
type DbHandler interface {
	Close()
	FetchEstimate(int) (*Estimate, error)
}
