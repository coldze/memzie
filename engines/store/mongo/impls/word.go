package impls

import (
	"time"

	"github.com/coldze/memzie/engines/store"
	"github.com/coldze/memzie/engines/store/mongo"
	"github.com/coldze/memzie/engines/store/mongo/structs"
	"github.com/coldze/primitives/custom_error"
)

type wordImpl struct {
	data *structs.Word
	coll store.Collection
}

func (w *wordImpl) GetText() string {
	return w.data.Text
}

func (w *wordImpl) GetID() string {
	return w.data.ID.Hex()
}

func (w *wordImpl) GetWeight() float64 {
	return w.data.Weight
}

func (w *wordImpl) GetShownTimes() int64 {
	return w.data.ShownTimes
}

func (w *wordImpl) GetFails() int64 {
	return w.data.Fails
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

func (w *wordImpl) Update(weightChange float64, failed bool) custom_error.CustomError {
	w.data.ShownTimes++
	w.data.Weight -= 0.001 * weightChange
	w.data.LastAnswered = time.Now().UnixNano()
	if failed {
		w.data.Fails++
	}
	updated, err := w.coll.Update(w.data, map[string]interface{}{"_id": w.data.ID})
	if err != nil {
		return custom_error.NewErrorf(err, "Failed to save word.")
	}
	if !updated {
		return custom_error.MakeErrorf("Failed to update")
	}
	return nil
}

func NewWordFactory(collFactory mongo.CollectionFactory) mongo.WordFactory {
	return func(word *structs.Word) (store.Word, custom_error.CustomError) {
		coll, err := collFactory(WORDS_COLLECTION)
		if err != nil {
			return nil, custom_error.NewErrorf(err, "Failed to get word collection.")
		}
		if word == nil {
			return nil, custom_error.MakeErrorf("Failed to wrap word. Data is nil.")
		}
		return &wordImpl{
			data: word,
			coll: coll,
		}, nil
	}
}
