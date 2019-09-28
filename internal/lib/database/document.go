package database

// Document ...
type Document interface {
	Collection() string
	Where() FilterList
}
