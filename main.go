package main

import (
	"context"
	"encoding/json"
	"github.com/coldze/memzie/engines/logic/mongo"
	"github.com/coldze/memzie/engines/store"
	"github.com/coldze/memzie/engines/store/mongo"
	"github.com/coldze/memzie/engines/store/mongo/impls"
	mgo "github.com/coldze/mongo-go-driver/mongo"
	"github.com/coldze/primitives/custom_error"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"log"
	"time"
	"math/rand"
)

type BoardData struct {
	ID           objectid.ObjectID  `bson:"_id"`
	Name         string             `bson:"Name"`
	Description  string             `bson:"Description"`
	LastQuestion *objectid.ObjectID `bson:"asked_word_id,omitempty"`
}

type WordData struct {
	ID        objectid.ObjectID `bson:"_id"`
	Shown     uint64            `bson:"shown_times"`
	Valid     uint64            `bson:"valid_times"`
	LastShown time.Time         `bson:"last_shown"`
}

type BoardAssignedData struct {
	BoardID objectid.ObjectID `bson:"board_id"`
	Data    interface{}       `bson:"data"`
}

type Question interface {
}

type Board interface {
	GetName() string
	GetDescription() string
	NextQuestion() Question
	CurrentQuestion() Question
	AddQuestion(q Question) custom_error.CustomError
}

type BoardImpl struct {
	id           objectid.ObjectID
	name         string
	description  string
	lastQuestion Question
	questions    []Question
}

func (b *BoardImpl) GetName() string {
	return b.name
}

func (b *BoardImpl) GetDescription() string {
	return b.description
}

func (b *BoardImpl) NextQuestion() Question {
	return nil
}

func (b *BoardImpl) CurrentQuestion() Question {
	return b.lastQuestion
}

func (b *BoardImpl) AddQuestion(q Question) custom_error.CustomError {
	b.questions = append(b.questions, q)
	return nil
}

func (b *BoardImpl) Save(data interface{}) custom_error.CustomError {
	d := BoardAssignedData{
		BoardID: b.id,
		Data:    data,
	}
	v, err := json.MarshalIndent(d, "", "   ")
	if err != nil {
		return custom_error.MakeErrorf("Failed to marshal. Error: %v", err)
	}
	log.Print(string(v))
	return nil
}

func main() {
	ctx := context.Background()
	client, err := mgo.Connect(ctx, "mongodb://localhost:27030", nil)
	if err != nil {
		log.Fatalf("Failed to connect to mongo-db. Error: %v", err)
	}
	defer client.Disconnect(ctx)

	collFactory := mongo.NewCollectionFactory(client, "memzie")

	wordFactory := impls.NewWordFactory(collFactory)
	boardFactory := impls.NewBoardFactory(wordFactory, collFactory)
	boardsFactory := impls.NewBoardsFactory(boardFactory, collFactory)
	root, customErr := impls.NewRoot(collFactory, boardsFactory)
	if customErr != nil {
		log.Fatalf("Failed to create root. Error: %v", customErr)
	}
	clientID := "5b025428ea5e7904880aa3be" //objectid.New().Hex()
	log.Printf("ClientID: '%v'", clientID)
	boards, customErr := root.GetBoards(clientID)
	if customErr != nil {
		log.Fatalf("Failed to get boards wrap. Error: %v", customErr)
	}
	if boards == nil {
		log.Fatalf("Boards are nil.")
	}
	boardID := "5b0254287126f07bd05e369f"
	board, customErr := boards.Get(boardID) // boards.Create("TEST_BOARD", "TEST BOARD DESCRIPTION")
	if customErr != nil {
		log.Fatalf("Failed to create board. Error: %v", customErr)
	}
	customErr = boards.List(func(board store.Board) (bool, custom_error.CustomError) {
		log.Printf("BOARD: %v. Name: %v", board.GetID(), board.GetName())
		return true, nil
	})
	if customErr != nil {
		log.Fatalf("Failed to list boards. Error: %v", customErr)
	}

	dummies := map[string]string{
		"test_01": "тест_01",
		"test_02": "тест_02",
		"test_03": "тест_03",
		"test_04": "тест_04",
		"test_05": "тест_05",
		"test_06": "тест_06",
		"test_07": "тест_07",
		"test_08": "тест_08",
		"test_09": "тест_09",
	}

	for k, v := range dummies {
		w, customErr := board.Create(&store.WordCreateParams{
			Text: k,
			Translations: []string{
				v,
			},
		})
		if customErr != nil {
			log.Fatalf("Failed to create word. %v, %v. Error: %v", k, v, customErr)
		}
		log.Printf("Created word: %v", w.GetID())
	}

	/*wordsCollection, customErr := collFactory(impls.WORDS_COLLECTION)
	if customErr != nil {
		log.Fatalf("Failed to get words collection. Error: %v", customErr)
	}*/

	logic, customErr := mongo2.NewLogic(client, "memzie", impls.WORDS_COLLECTION, board.GetID(), clientID, wordFactory)
	if customErr != nil {
		log.Fatalf("Failed to create logic. Error: %v", customErr)
	}


	for {
		word, customErr := logic.Next()
		if customErr != nil {
			log.Fatalf("Failed to get next word. Error: %v", customErr)
		}
		log.Printf("WordID: %v.          ID: %v", word.GetText(), word.GetID())
		customErr = word.Update(10 + (rand.Int63() % 10) - 5, false)
		if customErr != nil {
			log.Fatalf("Failed to update word. Error: %v", customErr)
		}
		time.Sleep(500 * time.Millisecond)
	}

	/*log.Printf("BoardID: '%v'", board.GetID())
	newWord := &store.WordCreateParams{
		Text:         "test",
		Translations: []string{"test", "Test", "tEstT"},
	}
	word, customErr := board.Create(newWord)
	if customErr != nil {
		log.Fatalf("Failed to create word. Error: %v", customErr)
	}
	log.Printf("WordID: '%v'", word.GetID())*/
	return

	/*logger := logs.NewStdLogger()
	app := cli.NewCliApp(logger)
	err := app.Run()
	if err == nil {
		return
	}
	logger.Fatalf("App run failed with error: %v", err)
	os.Exit(1)

	return
	board := BoardImpl{
		id:           objectid.New(),
		name:         "name_Test",
		description:  "desc_test",
		lastQuestion: nil,
		questions:    []Question{},
	}
	word := WordData{
		ID:        objectid.New(),
		Shown:     0,
		Valid:     0,
		LastShown: time.Time{},
	}
	board.Save(word)*/
}
