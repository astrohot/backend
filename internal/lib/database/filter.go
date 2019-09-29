package database

import (
	"go.mongodb.org/mongo-driver/bson"
)

// Filter ...
type Filter struct {
	Field string
	Value interface{}
}

func (f Filter) bson() bson.M {
	m := bson.M{}

	switch filter := f.Value.(type) {
	case Filter:
		m[f.Field] = filter.bson()

	case FilterList:
		filters := []bson.M{}
		for _, f := range filter {
			filters = append(filters, f.bson())
		}

		m[f.Field] = filters

	default:
		m[f.Field] = filter
	}

	return m
}

// FilterList ...
type FilterList []Filter

func (fs FilterList) bson() bson.M {
	m := bson.M{}

	for _, f := range fs {
		m[f.Field] = f.bson()[f.Field]
	}

	return m
}
