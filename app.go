package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
}

type Post struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Caption    string             `json:"caption,omitempty" bson:"caption,omitempty"`
	Image_URL  string             `json:"image_url,omitempty" bson:"image_url,omitempty"`
	Time_stamp string             `json:"time_stamp,omitempty" bson:"time_stamp,omitempty"`
	User       primitive.ObjectID `json:"user,omitempty" bson:"user,omitempty"`
}

var client *mongo.Client
var postlist []Post

const BlockSize = 16

//we are passing the passphrase to create a 32 byte key that will be used to hash the password
func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

//hashing function -> used to hash the password
func hash_password(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

//function to create a new user
func CreateNewUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var user User
	json.NewDecoder(request.Body).Decode(&user)
	hash := hash_password([]byte(user.Password), "ckdjjekk29i2")
	user.Password = string(hash)
	collection := client.Database("instagramAPI").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, user)
	json.NewEncoder(response).Encode(result)
}

//function to create a new post
func CreateNewPosts(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var post Post
	json.NewDecoder(request.Body).Decode(&post)
	post.Time_stamp = time.Now().Format("2001-10-26 20:59:20")
	collection := client.Database("instagramAPI").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, post)
	postlist = append(postlist, post)
	json.NewEncoder(response).Encode(result)
}

//function to get user according to his is
func GetUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var user User
	collection := client.Database("instagramAPI").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, User{ID: id}).Decode(&user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(user)
}

//get a post using id
func GetPostUsingID(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var post Post
	collection := client.Database("instagramAPI").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, Post{ID: id}).Decode(&post)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(post)
}

//Get all Posts of the User
func GetAllPosts(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var post Post
	collection := client.Database("instagramAPI").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, Post{User: id}).Decode(&post)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	for _, item := range postlist {
		if item.User == id {
			fmt.Println(item.ID)
			fmt.Println(item.Caption)

			json.NewEncoder(response).Encode(item)
			//return
		}
	}
	json.NewEncoder(response).Encode(&Post{})
}
func main() {
	fmt.Println("Starting the application")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	router := mux.NewRouter()
	router.HandleFunc("/users", CreateNewUser).Methods("POST")
	router.HandleFunc("/users/{id}", GetUser).Methods("GET")
	router.HandleFunc("/posts", CreateNewPosts).Methods("POST")
	router.HandleFunc("/posts/{id}", GetPostUsingID).Methods("GET")
	router.HandleFunc("/posts/users/{id}", GetAllPosts).Methods("GET")

	http.ListenAndServe(":1211", router)
}
