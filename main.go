package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/coldze/memzie/engines/logic/mongo"
	"github.com/coldze/memzie/engines/store"
	"github.com/coldze/memzie/engines/store/mongo"
	"github.com/coldze/memzie/engines/store/mongo/impls"
	"github.com/coldze/memzie/interfaces/rpc"
	"github.com/coldze/primitives/custom_error"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	mgo "github.com/mongodb/mongo-go-driver/mongo"
	"net/http"
)

type FolderData struct {
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

type FolderAssignedData struct {
	FolderID objectid.ObjectID `bson:"folder_id"`
	Data     interface{}       `bson:"data"`
}

type Question interface {
}

type Folder interface {
	GetName() string
	GetDescription() string
	NextQuestion() Question
	CurrentQuestion() Question
	AddQuestion(q Question) custom_error.CustomError
}

type FolderImpl struct {
	id           objectid.ObjectID
	name         string
	description  string
	lastQuestion Question
	questions    []Question
}

func (b *FolderImpl) GetName() string {
	return b.name
}

func (b *FolderImpl) GetDescription() string {
	return b.description
}

func (b *FolderImpl) NextQuestion() Question {
	return nil
}

func (b *FolderImpl) CurrentQuestion() Question {
	return b.lastQuestion
}

func (b *FolderImpl) AddQuestion(q Question) custom_error.CustomError {
	b.questions = append(b.questions, q)
	return nil
}

func (b *FolderImpl) Save(data interface{}) custom_error.CustomError {
	d := FolderAssignedData{
		FolderID: b.id,
		Data:     data,
	}
	v, err := json.MarshalIndent(d, "", "   ")
	if err != nil {
		return custom_error.MakeErrorf("Failed to marshal. Error: %v", err)
	}
	log.Print(string(v))
	return nil
}

type TestSA struct {
	FieldStr string      `json:"field_str"`
	Field    interface{} `json:"field"`
}

type TestSubSB struct {
	FA string `json:"fa"`
	FB string `json:"fb"`
}

type TestSB struct {
	FieldStr string    `json:"field_str"`
	Field    TestSubSB `json:"field"`
}

func Create() interface{} {
	return &TestSubSB{}
}

func Test() {
	t := TestSB{
		FieldStr: "field_str_value",
		Field: TestSubSB{
			FA: "FA_value",
			FB: "FB_value",
		},
	}
	d, err := json.Marshal(t)
	if err != nil {
		log.Fatalf("Failed. Error: %v", err)
	}

	ta := TestSA{
		Field: Create(),
	}
	d2, err := json.Marshal(ta)
	if err != nil {
		log.Fatalf("Failed. Error: %v", err)
	}
	err = json.Unmarshal(d, &ta)
	if err != nil {
		log.Fatalf("Failed. Error: %v", err)
	}
	log.Printf("M1: %v", string(d))
	log.Printf("M2: %v", string(d2))

	log.Printf("V1: %+v", t)
	vx, ok := ta.Field.(*TestSubSB)
	if !ok {
		log.Fatalf("Failed. Non expected type.")
	}
	log.Printf("V2: %+v. V2f: %+v", ta, *vx)
}

func TestRPC() {
	rpc.Register(func(path string, handler rpc.HttpRequestHandler) custom_error.CustomError {
		http.HandleFunc("/"+path, handler)
		return nil
	})
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatalf("Failed. Error: %v", err)
	}
}

func main() {
	Test()
	TestRPC()
	return
	ctx := context.Background()
	client, err := mgo.Connect(ctx, "mongodb://localhost:27030", nil)
	if err != nil {
		log.Fatalf("Failed to connect to mongo-db. Error: %v", err)
	}
	defer client.Disconnect(ctx)

	collFactory := mongo.NewCollectionFactory(client, "memzie")

	wordFactory := impls.NewWordFactory(collFactory)
	folderFactory := impls.NewFolderFactory(wordFactory, collFactory)
	foldersFactory := impls.NewFoldersFactory(folderFactory, collFactory)
	root, customErr := impls.NewRoot(collFactory, foldersFactory)
	if customErr != nil {
		log.Fatalf("Failed to create root. Error: %v", customErr)
	}
	clientID := "5b025428ea5e7904880aa3be" //objectid.New().Hex()
	log.Printf("ClientID: '%v'", clientID)
	folders, customErr := root.GetFolders(clientID)
	if customErr != nil {
		log.Fatalf("Failed to get folders wrap. Error: %v", customErr)
	}
	if folders == nil {
		log.Fatalf("Folders are nil.")
	}
	folderID := "5b0254287126f07bd05e369f"
	folder, customErr := folders.Get(folderID) // folders.Create("TEST_FOLDER", "TEST FOLDER DESCRIPTION")
	if customErr != nil {
		log.Fatalf("Failed to create folder. Error: %v", customErr)
	}
	customErr = folders.List(func(folder store.Folder) (bool, custom_error.CustomError) {
		log.Printf("FOLDER: %v. Name: %v", folder.GetID(), folder.GetName())
		return true, nil
	})
	if customErr != nil {
		log.Fatalf("Failed to list folders. Error: %v", customErr)
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
	tests := map[string]string{}
	for k, v := range dummies {
		w, customErr := folder.Create(&store.WordCreateParams{
			Text: k,
			Translations: []string{
				v,
			},
		})
		if customErr != nil {
			log.Fatalf("Failed to create word. %v, %v. Error: %v", k, v, customErr)
		}
		log.Printf("Created word: %v", w.GetID())
		tests[k] = v
	}

	/*wordsCollection, customErr := collFactory(impls.WORDS_COLLECTION)
	if customErr != nil {
		log.Fatalf("Failed to get words collection. Error: %v", customErr)
	}*/

	logic, customErr := mongo2.NewLogic(client, "memzie", impls.WORDS_COLLECTION, folder.GetID(), clientID, wordFactory)
	if customErr != nil {
		log.Fatalf("Failed to create logic. Error: %v", customErr)
	}
	tests["test_01"] = "asd"

	for {
		word, customErr := logic.Next()
		if customErr != nil {
			log.Fatalf("Failed to get next word. Error: %v", customErr)
		}
		//log.Printf("WordID: %v.          ID: %v", word.GetText(), word.GetID())
		var translation string
		//fmt.Printf("New word: %v\n", word.GetText())
		//fmt.Scanf("%s\n", &translation)
		//log.Printf("%s", translation)
		translation = tests[word.GetText()]
		strings.ToLower(strings.Trim(translation, " "))
		var correct *bool = new(bool)
		*correct = false
		word.List(func(valid string) (bool, custom_error.CustomError) {
			*correct = strings.ToLower(strings.Trim(valid, " ")) == translation
			return !*correct, nil
		})
		word.GetWeight()
		//x0, x1
		//x1/x0

		shown := float64(word.GetShownTimes()) + 1
		fails := float64(word.GetFails())

		var res string
		if *correct {
			res = "Correct. Step: "
		} else {
			res = "Incorrect. Step: "
			fails += 1.0
		}
		step := float64(0.999) + rand.Float64()/10000.0 - fails/shown/1000
		if !*correct && (rand.Int63()%100) < 5 {
			step = 0.0
		}
		res += fmt.Sprintf("%v. Fails: %v, Shown: %v, Div: %v", step, int64(fails), int64(shown), fails/shown/1000)
		//fmt.Println(res)
		customErr = word.Update(step, !*correct)
		if customErr != nil {
			log.Fatalf("Failed to update word. Error: %v", customErr)
		}
		time.Sleep(100 * time.Millisecond)
	}

	/*log.Printf("FolderID: '%v'", folder.GetID())
	newWord := &store.WordCreateParams{
		Text:         "test",
		Translations: []string{"test", "Test", "tEstT"},
	}
	word, customErr := folder.Create(newWord)
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
	folder := FolderImpl{
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
	folder.Save(word)*/
}
