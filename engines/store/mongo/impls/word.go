package impls

import (
	"github.com/coldze/memzie/engines/store"
	"github.com/coldze/memzie/engines/store/mongo"
	"github.com/coldze/memzie/engines/store/mongo/structs"
	"github.com/coldze/primitives/custom_error"
)

type wordImpl struct {
	data *structs.Word
}

func (w *wordImpl) GetText() string {
	return w.data.Text
}

func (w *wordImpl) GetID() string {
	return w.data.ID.Hex()
}

func (w *wordImpl) List(handle store.TranslationHandle) custom_error.CustomError {
	for i := range w.data.Translations {
		next, err := handle(w.data.Translations[i])
		if err != nil {
			return custom_error.NewErrorf(err, "Enumeration terminated.")
		}
		if !next {
			break
		}
	}
	return nil
}

func NewWordFactory() mongo.WordFactory {
	return func(word *structs.Word) (store.Word, custom_error.CustomError) {
		if word == nil {
			return nil, custom_error.MakeErrorf("Failed to wrap word. Data is nil.")
		}
		return &wordImpl{
			data: word,
		}, nil
	}
}
