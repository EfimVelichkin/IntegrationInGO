package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type User struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Friends []string `json:"friends"`
}

func (u *User) toString() string {
	return fmt.Sprintf("Name is %s, age %d, friends %v\n", u.Name, u.Age, u.Friends)
}

type servie struct {
	store map[string]*User
}

func main() {
	mux := http.NewServeMux()
	srv := servie{make(map[string]*User)}
	mux.HandleFunc("/create", srv.Create)
	mux.HandleFunc("/get", srv.GetALL)
	mux.HandleFunc("/addFriend", srv.AddFriend)
	mux.HandleFunc("/delete", srv.DeleteUser)
	mux.HandleFunc("/user/", srv.GetUserFriends)

	http.ListenAndServe(":8080", mux)
}

func (s *servie) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()

		var u User
		if err := json.Unmarshal(content, &u); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		u.ID = fmt.Sprintf("%x", md5.Sum([]byte(u.Name+strconv.Itoa(u.Age))))
		s.store[u.ID] = &u

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User was created " + u.Name + "\n"))
		w.Write([]byte("Your ID " + u.ID))
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *servie) GetALL(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		response := ""
		for _, user := range s.store {
			response += user.toString()
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *servie) AddFriend(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()

		var request struct {
			InitiatorID string `json:"initiatorId"`
			ReceiverID  string `json:"receiverId"`
		}
		if err := json.Unmarshal(content, &request); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		initiator, ok := s.store[request.InitiatorID]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Initiator not found"))
			return
		}

		receiver, ok := s.store[request.ReceiverID]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Receiver not found"))
			return
		}

		initiator.Friends = append(initiator.Friends, receiver.Name)
		receiver.Friends = append(receiver.Friends, initiator.Name)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%s и %s теперь друзья", initiator.Name, receiver.Name)))
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *servie) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()

		var request struct {
			TargetID string `json:"target_id"`
		}
		if err := json.Unmarshal(content, &request); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		user, ok := s.store[request.TargetID]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("User not found"))
			return
		}

		delete(s.store, request.TargetID)

		for _, friend := range user.Friends {
			for _, friendUser := range s.store {
				if friendUser.Name == friend {
					friendUser.Friends = remove(friendUser.Friends, user.Name)
				}
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User was deleted"))
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func remove(slice []string, item string) []string {
	for i, v := range slice {
		if v == item {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func (s *servie) GetUserFriends(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Извлечение userId из пути URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 || pathParts[1] != "user" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Invalid URL"))
		return
	}
	userId := pathParts[2]

	user, ok := s.store[userId]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user.Friends)
}

func (s *servie) UpdateUserAge(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Извлечение userId из пути URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 || pathParts[1] != "user" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Invalid URL"))
		return
	}
	userId := pathParts[2]

	var request struct {
		NewAge int `json:"new age"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request body"))
		return
	}

	user, ok := s.store[userId]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
		return
	}

	user.Age = request.NewAge

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("возраст пользователя успешно обновлён"))
}
