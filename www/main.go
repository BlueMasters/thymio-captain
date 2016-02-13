package main

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/kidstuff/mongostore"
	"gopkg.in/mgo.v2"
	"net/http"
	"time"
)

const (
	mongoServer           = "localhost"
	dbName                = "thymio_captain"
	sessionCollectionName = "sessions"
	maxAge                = int((time.Duration(24) * time.Hour).Seconds())
	sessionKey            = "session-key"
	secretKey             = "secret-key"
)

var (
	database *mgo.Session
	store    *mongostore.MongoStore
)

func TokenSigning(w http.ResponseWriter, r *http.Request) {
	database.Refresh()
	session, err := store.Get(r, sessionKey)
	if err != nil {
		log.Println(err.Error())
	}

	token := r.FormValue("idtoken")
	url := fmt.Sprintf(
		"https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=%s",
		token)

	resp, err := http.Get(url)
	if err != nil {
		session.Values["admin"] = "0"
		sessions.Save(r, w)
		log.Errorf("Failed to get URL: %s.", err)
		http.Error(w, "error", 500)
		return
	}
	if resp.StatusCode != 200 {
		session.Values["admin"] = "0"
		sessions.Save(r, w)
		log.Errorf("Status code is %d instead of 200.", resp.StatusCode)
		http.Error(w, "error", resp.StatusCode)
		return
	}

	var a AuthRequest
	err = json.NewDecoder(resp.Body).Decode(&a)

	query := struct {
		Iss string `json:"iss"`
		Sub string `json:"sub"`
	}{a.Iss, a.Sub}

	var u User
	err = database.DB(dbName).C("users").Find(query).One(&u)

	if err == mgo.ErrNotFound {
		u = User{
			Iss:   a.Iss,
			Sub:   a.Sub,
			Name:  a.Name,
			Email: a.Email,
			Admin: false,
		}
		database.DB(dbName).C("users").Insert(u)
	}

	if u.Admin {
		session.Values["admin"] = "1"
		fmt.Fprintln(w, "Hello Admin")
	} else {
		session.Values["admin"] = "0"
		fmt.Fprintln(w, "Hello User")
	}
	sessions.Save(r, w)

}

func AdminHandler(w http.ResponseWriter, r *http.Request) {
	database.Refresh()
	_, err := store.Get(r, sessionKey)
	if err != nil {
		log.Println(err.Error())
	}
	http.ServeFile(w, r, "admin.html")
}

func StartHandler(w http.ResponseWriter, r *http.Request) {
	database.Refresh()
	session, err := store.Get(r, sessionKey)
	if err != nil {
		log.Println(err.Error())
	}

	if session.Values["admin"] == "1" {
		fmt.Fprintln(w, "Admin ok")
	} else {
		fmt.Fprintln(w, "User ok")
	}
}

func main() {
	var err error
	database, err = mgo.Dial(mongoServer)
	if err != nil {
		log.Fatal(err)
	}
	store = mongostore.NewMongoStore(
		database.DB(dbName).C(sessionCollectionName),
		maxAge, true, []byte(secretKey))

	r := mux.NewRouter()
	r.HandleFunc("/start/{ID}", StartHandler)
	r.HandleFunc("/admin", AdminHandler)
	r.HandleFunc("/tokensignin", TokenSigning)
	http.Handle("/", r)
	log.Infoln("Ready.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
