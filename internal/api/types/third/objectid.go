package third

import (
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UnmarshalObjectID ...
func UnmarshalObjectID(v interface{}) (primitive.ObjectID, error) {
	switch value := v.(type) {
	case primitive.ObjectID:
		return value, nil
	case string:
		return primitive.ObjectIDFromHex(value)
	default:
		return [12]byte{}, fmt.Errorf("type assertion error!\t expected: primitive.ObjectID")
	}
}

// MarshalObjectID ...
func MarshalObjectID(objectID primitive.ObjectID) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.Quote(objectID.Hex()))
	})
}
