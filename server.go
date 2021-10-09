package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var ctx context.Context
var lock sync.Mutex

type UserID struct {
	UserId string `json:"UserId,omitempty" bson:"UserId,omitempty"`
}
type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId   string             `json:"UserId,omitempty" bson:"UserId,omitempty"`
	Name     string             `json:"Name,omitempty" bson:"Name,omitempty"`
	Email    string             `json:"Email,omitempty" bson:"Email,omitempty"`
	Password string             `json:"Password,omitempty" bson:"Password,omitempty"`
}

type Post struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId   string             `json:"UserId,omitempty" bson:"UserId,omitempty"`
	Caption  string             `json:"Caption,omitempty" bson:"Caption,omitempty"`
	ImageUrl string             `json:"ImageUrl,omitempty" bson:"ImageUrl,omitempty"`
	Time     string             `json:"time,omitempty" bson:"time,omitempty"`
}

type Posts struct {
	Posts       []Post `json:"Post_Details,omitempty"`
	Total       int    `json:"total,omitempty"`
	CurrentPost int    `json:"CurrentPost,omitempty"`
	LastPost    int    `json:"LastPost,omitempty"`
}
type UserPosts struct {
	PostID string `json:"PostID,omitempty" bson:"PostID,omitempty"`
}

func ConnectToDB() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb+srv://User1:Apple@cluster0.18prw.mongodb.net/test")
	client, _ = mongo.Connect(ctx, clientOptions)
}

func handleRequests() {
	http.HandleFunc("/users", CreateUser)
	http.HandleFunc("/users/", GetUserUsingId)
	http.HandleFunc("/posts", CreatePost)
	http.HandleFunc("/posts/users/", GetAllPosts)
	log.Fatal(http.ListenAndServe(":10000", nil))
}
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	lock.Lock()
	switch r.Method {
	case "POST":
		ConnectToDB()
		var user User
		_ = json.NewDecoder(r.Body).Decode(&user)
		h := sha256.New()
		h.Write([]byte(user.Password))
		user.Password = hex.EncodeToString(h.Sum(nil))
		collection := client.Database("INSTA").Collection("users")
		result, err := collection.InsertOne(ctx, user)
		if err != nil {
			fmt.Fprintf(w, `{"Error" :"User ID Taken"}`)
			fmt.Println(err.Error())
		} else {
			fmt.Fprintf(w, `{"Sucess" :"User Created" }`)
			fmt.Println(result)
		}
	case "GET":
		fmt.Fprintf(w, `{"Error" :"Used GET expected POST"}`)
	}
	fmt.Println("EndPoint Complete : Create a User")
	defer lock.Unlock()
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	lock.Lock()
	switch r.Method {
	case "POST":
		ConnectToDB()
		var post Post
		_ = json.NewDecoder(r.Body).Decode(&post)
		collection := client.Database("INSTA").Collection("Posts")
		_, err := collection.InsertOne(ctx, post)
		if err != nil {
			fmt.Fprintf(w, `{"Error" :"Try Again Later"}`)

		} else {
			fmt.Fprintf(w, `{"Sucess" :"Post Created"}`)
		}
	case "GET":
		fmt.Fprintf(w, `{"Error" :"Used GET expected POST"}`)
	}
	fmt.Println("EndPoint Hit : Create a Post")
	defer lock.Unlock()
}

func GetAllPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	lock.Lock()
	switch r.Method {
	case "POST":
		fmt.Fprintf(w, `{"Error" :"Used POST expected GET"}`)
	case "GET":
		ConnectToDB()
		var searchUser UserID
		searchUser.UserId = r.URL.Path[strings.LastIndex(r.URL.Path[1:], "/")+2:]

		fmt.Println(searchUser.UserId)
		filter := bson.M{
			"UserId": searchUser.UserId,
		}
		var posts []Post

		findOptions := options.Find()
		collection := client.Database("INSTA").Collection("Posts")
		page := 1
		var perPage int64 = 9
		total, _ := collection.CountDocuments(ctx, filter)

		findOptions.SetSkip((int64(page) - 1) * perPage)
		findOptions.SetLimit(perPage)
		count := 0
		cursor, _ := collection.Find(ctx, filter, findOptions)
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var post Post
			cursor.Decode(&post)
			posts = append(posts, post)
			count = count + 1
		}
		var Allposts Posts
		Allposts.Posts = posts
		Allposts.LastPost = count
		Allposts.CurrentPost = page
		Allposts.Total = int(total)
		fmt.Print(int(float64(total / perPage)))
		json.NewEncoder(w).Encode(Allposts)

	}
	fmt.Println("EndPoint Hit : Get All Posts User")
	defer lock.Unlock()
}

func GetUserUsingId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	lock.Lock()
	switch r.Method {
	case "POST":
		fmt.Fprintf(w, `{"Error" :"Used POST expected GET"}`)
	case "GET":
		ConnectToDB()
		var searchUser UserID
		searchUser.UserId = r.URL.Path[strings.LastIndex(r.URL.Path[1:], "/")+2:]
		var user User
		collection := client.Database("INSTA").Collection("users")
		err := collection.FindOne(context.TODO(), searchUser).Decode(&user)
		if err != nil {
			fmt.Fprintf(w, `{"Error" :"Incorrect User ID" , "message": "`+err.Error()+`" }`)
		} else {
			json.NewEncoder(w).Encode(user)
		}
	}
	defer lock.Unlock()
}

func main() {
	handleRequests()
}
