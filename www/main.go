// Copyright 2016 Jacques Supcik <jacques.supcik@hefr.ch>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/kidstuff/mongostore"
	"gopkg.in/mgo.v2"
	"html/template"
	"net/http"
)

const (
	dbName     = "thymio_captain"
	sessionC   = "sessions"
	maxAge     = 24 * 3600
	sessionKey = "session-key"
)

var (
	database *mgo.Session
	store    *mongostore.MongoStore
)

func initSession(w http.ResponseWriter, r *http.Request) (vars map[string]string, session *sessions.Session, err error) {
	database.Refresh()
	vars = mux.Vars(r)
	session, err = store.Get(r, sessionKey)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Session Error", 500)
	}
	return
}

func CardLogin(w http.ResponseWriter, r *http.Request) {
	vars, session, err := initSession(w, r)
	if err != nil {
		return
	}

	if vars["CardId"] == "friendship" {
		session.Values["admin"] = "1"
		sessions.Save(r, w)
		http.ServeFile(w, r, "html/login-ok.html")
	} else {
		session.Values["admin"] = "0"
		sessions.Save(r, w)
		http.ServeFile(w, r, "html/login-failed.html")
	}

}

func Logout(w http.ResponseWriter, r *http.Request) {
	_, session, err := initSession(w, r)
	if err != nil {
		return
	}

	session.Values["admin"] = "0"
	sessions.Save(r, w)
	http.ServeFile(w, r, "html/logout.html")
}

func Debug(w http.ResponseWriter, r *http.Request) {
	vars, session, err := initSession(w, r)
	if err != nil {
		return
	}

	tmpl, err := template.ParseFiles("html/debug.html")
	if err != nil {
		log.Infof("Error 1 %s", err)
	}
	s := fmt.Sprintf("%v", session.Values)
	if err != nil {
		log.Infof("Error 2 %s", err)
	}
	err = tmpl.Execute(w, struct{ Session string }{string(s)})
}

func Start(w http.ResponseWriter, r *http.Request) {
	vars, session, err := initSession(w, r)
	if err != nil {
		return
	}

	session.Values["cardId"] = vars["CardId"]
	sessions.Save(r, w)

	if session.Values["admin"] == "1" {
		http.ServeFile(w, r, "html/admin.html")
	} else {
		http.ServeFile(w, r, "html/user.html")
	}
}

func main() {
	var port = flag.Int("port", 8080, "port")
	var debug = flag.Bool("debug", false, "run in debug mode")
	var domain = flag.String("domain", "telecom.tk", "Domain name (for the cookie)")
	var mongoServer = flag.String("mongo-server", "localhost", "MongoDB server URL")
	var secretKey = flag.String("secret-key", "not-so-secret", "Secret key (for secure cookies)")

	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
		log.Debug("Debug mode")
	} else {
		log.SetLevel(log.InfoLevel)
	}

	var err error
	database, err = mgo.Dial(*mongoServer)
	if err != nil {
		log.Fatal(err)
	}
	store = mongostore.NewMongoStore(
		database.DB(dbName).C(sessionC),
		maxAge, true, []byte(*secretKey))

	store.Options.Domain = *domain

	r := mux.NewRouter()
	r.HandleFunc("/start/{CardId}", Start)
	r.HandleFunc("/cardlogin/{CardId}", CardLogin)
	r.HandleFunc("/logout", Logout)
	r.HandleFunc("/debug", Debug)

	http.Handle("/", r)
	log.Infof("Ready, listening on port %d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
