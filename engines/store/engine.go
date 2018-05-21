package store

import "github.com/coldze/primitives/custom_error"

type Engine interface {
	Close()
	GetCollection(name string) (Collection, custom_error.CustomError)
}

type Root interface {
	GetBoards(clientID string) (Boards, custom_error.CustomError)
}

type BoardListHandle func(Board) (bool, custom_error.CustomError)

type Boards interface {
	List(handle BoardListHandle) custom_error.CustomError
	Get(id string) (Board, custom_error.CustomError)
	Create(name string, description string) (Board, custom_error.CustomError)
	Delete(id string) custom_error.CustomError
}

type WordListHandle func(Word) (bool, custom_error.CustomError)

type WordCreateParams struct {
	Text         string
	Translations []string
}

type Board interface {
	GetName() string
	GetDescription() string
	GetID() string
	List(handle WordListHandle) custom_error.CustomError
	Get(id string) (Word, custom_error.CustomError)
	Create(params *WordCreateParams) (Word, custom_error.CustomError)
	Delete(id string) custom_error.CustomError
}

type TranslationHandle func(translation string) (bool, custom_error.CustomError)

type Word interface {
	GetText() string
	GetID() string
	List(handle TranslationHandle) custom_error.CustomError
}
