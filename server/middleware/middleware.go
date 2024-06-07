package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"go-react-todo/models"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection = nil

func init() {
	loadEnv()
	connectToDB()
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func connectToDB() {
	dbURI := os.Getenv("DB_URI")
	dbName := os.Getenv("DB_NAME")
	collectionName := os.Getenv("DB_COLLECTION_NAME")

	clientOptions := options.Client().ApplyURI(dbURI)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection = client.Database(dbName).Collection(collectionName)
}

func setHeaders(w http.ResponseWriter, method string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", method)
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	setHeaders(w, "GET")
	payload := getAllTasks()
	json.NewEncoder(w).Encode(payload)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	setHeaders(w, "POST")
	var task models.ToDoList
	fmt.Println(r.Body)
	json.NewDecoder(r.Body).Decode(&task)
	insertTask(task)
	json.NewEncoder(w).Encode(task)
}

func CompleteTask(w http.ResponseWriter, r *http.Request) {
	setHeaders(w, "PUT")
	params := mux.Vars(r)
	completeTask(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func UndoTask(w http.ResponseWriter, r *http.Request) {
	setHeaders(w, "PUT")
	params := mux.Vars(r)
	undoTask(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	setHeaders(w, "DELETE")
	params := mux.Vars(r)
	deleteTask(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteAllTasks(w http.ResponseWriter, r *http.Request) {
	setHeaders(w, "DELETE")
	deleteAllTasks()
	json.NewEncoder(w).Encode("All tasks deleted")
}

// Database interaction functions

func getAllTasks() []primitive.M {
	cusror, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	var results []primitive.M
	for cusror.Next(context.Background()) {
		var result bson.M
		e := cusror.Decode(&result)
		if e != nil {
			log.Fatal(e)
		}
		results = append(results, result)
	}
	if err := cusror.Err(); err != nil {
		log.Fatal(err)
	}
	cusror.Close(context.Background())
	return results
}

func insertTask(task models.ToDoList) {
	fmt.Println(`task: `, task)
	_, err := collection.InsertOne(context.Background(), task)
	if err != nil {
		log.Fatal(err)
	}
}

func undoTask(taskID string) {
	id, _ := primitive.ObjectIDFromHex(taskID)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"complete": false}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
}

func deleteTask(taskID string) {
	id, _ := primitive.ObjectIDFromHex(taskID)
	filter := bson.M{"_id": id}
	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
}

func completeTask(taskID string) {
	id, _ := primitive.ObjectIDFromHex(taskID)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"complete": true}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
}

func deleteAllTasks() {
	cusror, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	for cusror.Next(context.Background()) {
		var result bson.M
		e := cusror.Decode(&result)
		if e != nil {
			log.Fatal(e)
		}
		id := result["_id"].(primitive.ObjectID)
		filter := bson.M{"_id": id}
		_, err := collection.DeleteOne(context.Background(), filter)
		if err != nil {
			log.Fatal(err)
		}
	}
}
