package storage

type QuotesStorage interface {
	RandomQuote() (string, error)
}
