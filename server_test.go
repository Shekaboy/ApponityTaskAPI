package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

type GetUserId struct {
	userID   string
	expected string
}
type PostTest struct {
	UserId   string
	Caption  string
	ImageUrl string
	Time     string
}
type UserCreate struct {
	UserId   string
	Name     string
	Email    string
	Password string
}

func Test_GetUserUsingId(t *testing.T) {
	UserID := []GetUserId{
		{"/users/ABC12dxxaaa345", `{"Error" :"Incorrect User ID" , "message": "mongo: no documents in result" }`},
		{"/users/Ball", `{"_id":"61619d94d604783918e324b5","UserId":"Ball","Name":"Apple","Email":"Cat","Password":"0eb129bf94594aaeee66e38361d7be212cd927c3df4dd92e3ded2e0da0c7ad88"}
`},
		{"/users/ABC12dxa345", `{"Error" :"Incorrect User ID" , "message": "mongo: no documents in result" }`},
	}

	for _, testcase := range UserID {
		req, err := http.NewRequest("GET", testcase.userID, nil)
		if err != nil {
			t.Fatalf("Couldn't create request: %v\n", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetUserUsingId)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		if rr.Body.String() != testcase.expected {
			t.Errorf("got :%v want :%v",
				rr.Body.String(), testcase.expected)
		}
	}
}

func TestCreateUser(t *testing.T) {
	var jsonStr = []byte(`{"UserId":"Ball","Name":"Apple","Email":"Cat","Password":"Hello"}`)
	expected := `{"Error" :"User ID Taken"}`
	sample(jsonStr, expected, t)

	jsonStr = []byte(`{"UserId":"Balsl1f255345","Name":"Apple","Email":"Cat","Password":"Hello"}`)
	expected = `{"Sucess" :"User Created" }`

	sample(jsonStr, expected, t)
}

func sample(data []byte, expected string, t *testing.T) {
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(data))
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateUser)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if rr.Body.String() != expected {
		t.Errorf("got :%v want :%v",
			rr.Body.String(), expected)
	}
}

func TestCreatePost(t *testing.T) {
	var jsonStr = []byte(`{"UserId":"Ball","Caption":"Apple","ImageUrl":"Cat","time":"10:30"}`)
	expected := `{"Sucess" :"Post Created"}`
	sample_1(jsonStr, expected, t)

	jsonStr = []byte(`{"UserId":"Balsssl","Caption":"Asdspple","ImageUrl":"Casdt","time":"11:30"}`)
	expected = `{"Sucess" :"Post Created"}`

	sample_1(jsonStr, expected, t)
}

func sample_1(data []byte, expected string, t *testing.T) {
	req, err := http.NewRequest("POST", "/posts", bytes.NewBuffer(data))
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreatePost)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	if rr.Body.String() != expected {
		t.Errorf("got :%v want :%v",
			rr.Body.String(), expected)
	}
}

func TestGetAllPosts(t *testing.T) {
	req, err := http.NewRequest("GET", "/posts/users/Balsssl", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllPosts)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"Post_Details":[{"_id":"6161bb42e1d621b4f8dd9f20","UserId":"Balsssl","Caption":"Asdspple","ImageUrl":"Casdt","time":"11:30"},{"_id":"6161bb5986e3c684e0baca12","UserId":"Balsssl","Caption":"Asdspple","ImageUrl":"Casdt","time":"11:30"},{"_id":"6161bb67e92c3b30b3cffca0","UserId":"Balsssl","Caption":"Asdspple","ImageUrl":"Casdt","time":"11:30"},{"_id":"6161bb874a37312983c4b6c3","UserId":"Balsssl","Caption":"Asdspple","ImageUrl":"Casdt","time":"11:30"},{"_id":"6161d36ff1c3022b3fbf1420","UserId":"Balsssl","Caption":"Asdspple","ImageUrl":"Casdt","time":"11:30"},{"_id":"6161d8c7177017526ffac55f","UserId":"Balsssl","Caption":"Asdspple","ImageUrl":"Casdt","time":"11:30"}],"total":6,"CurrentPost":1,"LastPost":6}
`
	if rr.Body.String() != expected {
		t.Errorf("got :%v want :%v",
			rr.Body.String(), expected)
	}
}
