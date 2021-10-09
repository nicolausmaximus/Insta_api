package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Person struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson: "_id,omitempty"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
}

var client *mongo.Client

func NewRouter() *Router {
	return &Router{namedRoutes: make(map[string]*Route)}
}

type Router struct {
	NotFoundHandler         http.Handler
	MethodNotAllowedHandler http.Handler
	routes                  []*Route
	namedRoutes             map[string]*Route
	KeepContext             bool
	middlewares             []middleware
	routeConf
}
type routeConf struct {
	useEncodedPath bool
	strictSlash    bool
	skipClean      bool
	regexp         routeRegexpGroup
	matchers       []matcher
	buildScheme    string
	buildVarsFunc  BuildVarsFunc
}

func CreatePersonEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var person Person
	json.NewDecoder(request.Body).Decode(&person)
	collection := client.Database("instagramAPI").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, person)
	json.NewEncoder(response).Encode(result)
}
func NewRouter() *Router {
	return &Router{namedRoutes: make(map[string]*Route)}
}
func GetPeopleEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var people []Person
	collection := client.Database("instagramAPI").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
		//fmt.Println(err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person Person
		cursor.Decode(&person)
		people = append(people, person)
	}

	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
		//fmt.Println(err)
	}
	json.NewEncoder(response).Encode(people)

}
func main() {
	fmt.Println("Starting the application")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	resp, err := http.PostForm("http://localhost:8081/person", CreatePersonEndpoint)
	if err != nil {
		log.Printf("Request Failed: %s", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body) // Log the request body
	bodyString := string(body)
	log.Print(bodyString) // Unmarshal result
	post := Person{}
	err = json.Unmarshal(body, &post)
	if err != nil {
		log.Printf("Reading body failed: %s", err)
		return
	}

	router := NewRouter()
	router.HandleFunc("/person", CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
	http.ListenAndServe(":3101", router)
}
