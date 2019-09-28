package database

import "go.mongodb.org/mongo-driver/bson"

// Filter ...
type Filter struct {
	Field string
	Value interface{}
}

// FilterList ...
type FilterList []Filter

func (fs FilterList) bson() bson.M {
	m := bson.M{}
	for _, f := range fs {
		m[f.Field] = f.Value
	}

	return m
}
