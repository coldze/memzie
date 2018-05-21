package main

import (
	"encoding/json"
	"github.com/coldze/memzie/engines/store/mongo"
	"github.com/coldze/memzie/engines/store/mongo/impls"
	"github.com/coldze/primitives/custom_error"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"log"
	"time"
	"github.com/coldze/memzie/engines/store"
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
	storeEngine, customErr := mongo.NewEngine("mongodb://localhost:27030", "memzie")
	if customErr != nil {
		log.Fatalf("Failed to create engine. Error: %v", customErr)
	}
	defer storeEngine.Close()

	wordFactory := mongo.NewWordFactory()
	boardFactory := impls.NewBoardFactory(wordFactory)
	boardsFactory := mongo.NewBoardsFactory(boardFactory)
	root, customErr := mongo.NewRoot(storeEngine, boardsFactory)
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
	//boardID := "5b0254287126f07bd05e369f"
	board, customErr := /*boards.Get(boardID) */boards.Create("TEST_BOARD", "TEST BOARD DESCRIPTION")
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
	log.Printf("BoardID: '%v'", board.GetID())
	newWord := &store.WordCreateParams{
		Text: "test",
		Translations: []string{"test", "Test", "tEstT"},
	}
	word, customErr := board.Create(newWord)
	if customErr != nil {
		log.Fatalf("Failed to create word. Error: %v", customErr)
	}
	log.Printf("WordID: '%v'", word.GetID())
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
