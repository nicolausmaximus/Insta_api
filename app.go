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
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
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

//Mutex has been used to make the server thread safe
var lock sync.Mutex

//hashing blocksize
const BlockSize = 16

//we are passing the passphrase to create a 32 byte key that will be used to hash the password
func createHash(key string) string {
	lock.Lock()
	defer lock.Unlock()
	hasher := md5.New()
	hasher.Write([]byte(key))
	time.Sleep(1 * time.Second)
	return hex.EncodeToString(hasher.Sum(nil))
}

//hashing function -> used to hash the password
func hash_password(data []byte, passphrase string) []byte {
	lock.Lock()
	defer lock.Unlock()
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
	time.Sleep(1 * time.Second)
	return ciphertext

}

//function to create a new user
func CreateNewUser(response http.ResponseWriter, request *http.Request) {
	lock.Lock()
	defer lock.Unlock()
	response.Header().Add("content-type", "application/json")
	var user User
	json.NewDecoder(request.Body).Decode(&user)
	hash := hash_password([]byte(user.Password), "ckdjjekk29i2")
	user.Password = string(hash)
	collection := client.Database("instagram_api").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, user)
	json.NewEncoder(response).Encode(result)
	time.Sleep(1 * time.Second)
}

//function to create a new post
func CreateNewPosts(response http.ResponseWriter, request *http.Request) {
	lock.Lock()
	defer lock.Unlock()
	response.Header().Add("content-type", "application/json")
	var post Post
	json.NewDecoder(request.Body).Decode(&post)
	post.Time_stamp = time.Now().Format("2001-10-26 20:59:20")
	collection := client.Database("instagram_api").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, post)
	json.NewEncoder(response).Encode(result)
	time.Sleep(1 * time.Second)
}

//function to get user according to his is
func GetUser(response http.ResponseWriter, request *http.Request) {
	lock.Lock()
	defer lock.Unlock()
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var user User
	collection := client.Database("instagram_api").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, User{ID: id}).Decode(&user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(user)
	time.Sleep(1 * time.Second)
}

//get a post using id
func GetPostUsingID(response http.ResponseWriter, request *http.Request) {
	lock.Lock()
	defer lock.Unlock()
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var post Post
	collection := client.Database("instagram_api").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, Post{ID: id}).Decode(&post)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(post)
	time.Sleep(1 * time.Second)
}

//Get all Posts of the User
func GetAllPosts(response http.ResponseWriter, request *http.Request) {
	var postlist []Post
	lock.Lock()
	defer lock.Unlock()
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	limit, _ := strconv.Atoi(params["limit"])
	var post Post
	collection := client.Database("instagram_api").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cur, err := collection.Find(ctx, bson.M{"user": id})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	for cur.Next(ctx) {
		err := cur.Decode(&post)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{"message": "` + err.Error() + `"}`))
			return
		}
		postlist = append(postlist, post)
	}

	for _, item := range postlist {
		if item.User == id {
			if limit > 0 {
				limit--
				json.NewEncoder(response).Encode(item)
			}
		}
	}
	json.NewEncoder(response).Encode(&Post{})
	time.Sleep(1 * time.Second)
}

//main function to handle requests and call all other functions
func main() {
	fmt.Println("Starting the application")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	router := mux.NewRouter()
	router.HandleFunc("/users", CreateNewUser).Methods("POST")
	router.HandleFunc("/users/{id}", GetUser).Methods("GET")
	router.HandleFunc("/posts", CreateNewPosts).Methods("POST")
	router.HandleFunc("/posts/{id}", GetPostUsingID).Methods("GET")
	router.HandleFunc("/posts/users/{id}&limit={limit}", GetAllPosts).Methods("GET")
	time.Sleep(3 * time.Second)
	http.ListenAndServe(":1211", router)
}
