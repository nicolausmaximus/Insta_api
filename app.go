package appointy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type Person struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
}

type Post struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	User    primitive.ObjectID `json:"user,omitempty" bson:"user,omitempty"`
	Caption string             `json:"caption,omitempty" bson:"caption,omitempty"`
	Image   string             `json:"image,omitempty" bson:"image,omitempty"`
	Time    string             `json:"time,omitempty" bson:"time,omitempty"`
}

var client *mongo.Client
var postt []Post

//create Users Endpoint
func CreateUsersEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var person Person
	json.NewDecoder(request.Body).Decode(&person)
	hash, err := bcrypt.GenerateFromPassword([]byte(person.Password), bcrypt.DefaultCost)
	person.Password = string(hash)
	if err != nil {
		fmt.Println(err)
	}
	collection := client.Database("instagramAPI").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, person)
	json.NewEncoder(response).Encode(result)
}

func CreatePostsEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var post Post
	json.NewDecoder(request.Body).Decode(&post)
	post.Time = time.Now().Format("2006-01-02 15:04:05")
	collection := client.Database("instagramAPI").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, post)
	postt = append(postt, post)
	json.NewEncoder(response).Encode(result)
}

func GetUsersEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var user Person
	collection := client.Database("instagramAPI").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, Person{ID: id}).Decode(&user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(user)
}

func main() {
	fmt.Println("Starting the application")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	router := mux.NewRouter()
	router.HandleFunc("/users", CreateUsersEndpoint).Methods("POST")
	router.HandleFunc("/users/{id}", GetUsersEndpoint).Methods("GET")
	router.HandleFunc("/posts", CreatePostsEndpoint).Methods("POST")
	router.HandleFunc("/posts/{id}", GetPostsEndpoint).Methods("GET")
	router.HandleFunc("/posts/users/{id}", GetAllPostsEndpoint).Methods("GET")

	http.ListenAndServe(":4000", router)
}
