package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	// "github.com/mongodb/mongo-go-driver/bson/primitive"
	// "github.com/mongodb/mongo-go-driver/mongo/options"
	// "github.com/mongodb/mongo-go-driver/mongo/readpref"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	//	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	//"go.mongodb.org/mongo-driver/*"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Data Structure
type Person struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
}

var client *mongo.Client
var collection *mongo.Collection

// func CreatePersonEndpoint(response http.ResponseWriter, request *http.Request) {
// 	response.Header().Set("content-type", "application/json")
// 	var person Person
// 	_ = json.NewDecoder(request.Body).Decode(&person)
// 	collection := client.Database("api").Collection("people")
// 	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
// 	result, _ := collection.InsertOne(ctx, person)
// 	json.NewEncoder(response).Encode(result)
// }

func CreatePersonEndpoint(response http.ResponseWriter, request *http.Request) {

	response.Header().Set("content-type", "application/json")
	var person Person
	_ = json.NewDecoder(request.Body).Decode(&person)
	//fmt.Printf("%s", json.NewDecoder(request.Body).Decode(&person))
	fmt.Println(" pt1")
	//	fmt.Printf("%#v\n", collection)
	//	fmt.Printf("%#v\n", person)
	fmt.Printf("%s\n", person.ID)

	//fmt.Println(json.NewDecoder(request.Body).Decode(&person))
	collection := client.Database("api").Collection("people")
	fmt.Printf("%#v\n", collection)
	//collection := client.Collection("people")

	fmt.Println(" pt2")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := collection.InsertOne(ctx, person)
	if err != nil {
		fmt.Println("!!insert! ")
		log.Fatal(err)
	}
	fmt.Println(result)
	_ = json.NewDecoder(request.Body).Decode(&person)

}

func GetPeopleEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var people []Person
	collection := client.Database("api").Collection("people")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var person Person
		cursor.Decode(&person)
		people = append(people, person)

	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return

	}
	json.NewEncoder(response).Encode(people)

}

//
//
//
//

func GetPersonEndpoint(response http.ResponseWriter, request *http.Request) {

	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var person Person
	collection := client.Database("api").Collection("people")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, Person{ID: id}).Decode(&person)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return

	}
	_ = json.NewEncoder(response).Encode(person)

}

func main() {
	fmt.Println("Starting ... ")

	//ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//var client *mongo.Client

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	mongoURI := os.Getenv("MONGO_URI")
	fmt.Println(mongoURI)

	//	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	client, err = mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		fmt.Println("!!client.connect! ")
		log.Fatal(err)
	}

	//And connect it to your running MongoDB server:

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		fmt.Println("!!client.connect! ")
		log.Fatal(err)
	}
	collection := client.Database("api").Collection("people")
	//fmt.Printf("%#v\n", collection)
	//fmt.Printf("%#v\n", client)
	_ = collection
	// To do this in a single step, you can use the Connect function:
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	// client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://clickingmouse:Aa123654@cluster0-5pcao.azure.mongodb.net/test?retryWrites=true&w=majority"))
	// if err != nil {
	// 	fmt.Println("!!client.connect! ")
	// 	log.Fatal(err)
	// }
	// Calling Connect does not block for server discovery. If you wish to know if a MongoDB server has been found and connected to, use the Ping method:

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Println("!!ping! ")
		log.Fatal(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/person", CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/person/{id}", GetPersonEndpoint).Methods("GET")

	log.Fatal(http.ListenAndServe(":3030", router))

}
