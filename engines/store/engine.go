package store

import "github.com/coldze/primitives/custom_error"

type Engine interface {
	Close()
	GetCollection(name string) (Collection, custom_error.CustomError)
}

type Root interface {
	GetFolders(clientID string) (Folders, custom_error.CustomError)
}

type FolderListHandle func(Folder) (bool, custom_error.CustomError)

type Folders interface {
	List(handle FolderListHandle) custom_error.CustomError
	Get(id string) (Folder, custom_error.CustomError)
	Create(name string, description string) (Folder, custom_error.CustomError)
	Delete(id string) custom_error.CustomError
}

type WordListHandle func(Word) (bool, custom_error.CustomError)

type WordCreateParams struct {
	Text         string
	Translations []string
}

type Folder interface {
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
	GetWeight() float64
	GetShownTimes() int64
	GetFails() int64
	List(handle TranslationHandle) custom_error.CustomError
	Update(update float64, failed bool) custom_error.CustomError
}
