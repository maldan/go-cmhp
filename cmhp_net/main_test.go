package cmhp_net_test

import (
	"testing"

	"github.com/maldan/go-cmhp/cmhp_net"
)

type Todo struct {
	UserId    int    `json:"userId"`
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type Post struct {
	UserId int    `json:"userId"`
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func TestGet(t *testing.T) {
	todoList := make([]Todo, 0)

	response := cmhp_net.Request(cmhp_net.HttpArgs{
		Url:        "https://jsonplaceholder.typicode.com/todos",
		Method:     "GET",
		OutputJSON: &todoList,
	})

	if response.StatusCode != 200 {
		t.Errorf("Can't request")
	}
	if len(todoList) == 0 {
		t.Errorf("No data")
	}
	if todoList[0].Id != 1 {
		t.Errorf("Can't parse json")
	}
}

func TestPost(t *testing.T) {

	response := cmhp_net.Request(cmhp_net.HttpArgs{
		Url:    "https://jsonplaceholder.typicode.com/posts",
		Method: "POST",
		InputJSON: &Post{
			UserId: 1,
			Id:     1,
			Title:  "A",
			Body:   "B",
		},
	})
	if response.StatusCode != 201 {
		t.Errorf("Can't create")
	}
}
