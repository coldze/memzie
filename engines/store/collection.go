package store

import "github.com/coldze/primitives/custom_error"

type Decoder func(func(interface{}) error) (interface{}, custom_error.CustomError)

type Collection interface {
	FindOne(decoder Decoder, filter map[string]interface{}) (interface{}, custom_error.CustomError)
	FindAll(decoder Decoder, filter map[string]interface{}, processor func(res interface{}) (bool, custom_error.CustomError)) custom_error.CustomError

	Create(object interface{}) (string, custom_error.CustomError)
	Delete(filters map[string]interface{}) (bool, custom_error.CustomError)

	Update(object interface{}, filter map[string]interface{}) (bool, custom_error.CustomError)
}
