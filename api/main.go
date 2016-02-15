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
	"encoding/json"
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/kidstuff/mongostore"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

const (
	dbName     = "thymio_captain"
	sessionC   = "sessions"
	sessionKey = "session-key"
	prefix     = "/v1"
	cardC      = "card"
	cardId     = "cardid"
	robotC     = "robot"
	robotURL   = "roboturl"
)

type Info struct {
	CardId  string
	IsAdmin bool
}

type Robot struct {
	URL    string
	Name   string
	CardId string
}

type Card struct {
	CardId  string
	Program []byte
}

type JsonError struct {
	ErrorDescription string `json:"error-description"`
}

type JsonOK struct {
	Result string `json:"result"`
}

var (
	database *mgo.Database
	store    *mongostore.MongoStore
)

func initSession(w http.ResponseWriter, r *http.Request) (vars map[string]string, session *sessions.Session, err error) {
	database.Session.Refresh()
	vars = mux.Vars(r)
	session, err = store.Get(r, sessionKey)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Session Error", 500)
	}
	return
}

func GetCard(w http.ResponseWriter, r *http.Request) {
	vars, _, err := initSession(w, r)
	if err != nil {
		return
	}

	var card Card
	log.Debug(bson.M{cardId: vars[cardId]})
	err = database.C(cardC).Find(bson.M{cardId: vars[cardId]}).One(&card)
	if err != nil {
		errorDesc, _ := json.Marshal(JsonError{err.Error()})
		log.Info(err)
		http.Error(w, string(errorDesc), 400)
	} else {
		json.NewEncoder(w).Encode(card)
	}
}

func PutCard(w http.ResponseWriter, r *http.Request) {
	vars, _, err := initSession(w, r)
	if err != nil {
		return
	}

	var payload struct {
		Program []byte `json"program"`
	}

	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		errorDesc, _ := json.Marshal(JsonError{"Invalid payload"})
		http.Error(w, string(errorDesc), 400)
		return
	}

	var card Card
	card.CardId = vars[cardId]
	card.Program = payload.Program

	_, err = database.C(cardC).Upsert(bson.M{cardId: vars[cardId]}, card)
	if err != nil {
		errorDesc, _ := json.Marshal(JsonError{err.Error()})
		http.Error(w, string(errorDesc), 500)
	} else {
		json.NewEncoder(w).Encode(JsonOK{"Done"})
	}
}

func GetInfo(w http.ResponseWriter, r *http.Request) {
	_, session, err := initSession(w, r)
	if err != nil {
		return
	}

	var info Info
	if cardId, ok := session.Values[cardId]; ok {
		info.CardId = cardId.(string)
	}

	if admin, ok := session.Values["admin"]; ok {
		info.IsAdmin = admin == "1"
	}

	json.NewEncoder(w).Encode(info)
}

func AssociateRobot(w http.ResponseWriter, r *http.Request) {
	vars, _, err := initSession(w, r)
	if err != nil {
		return
	}

	n, err := database.C(cardC).Find(bson.M{cardId: vars[cardId]}).Count()
	if err != nil {
		errorDesc, _ := json.Marshal(JsonError{err.Error()})
		http.Error(w, string(errorDesc), 400)
		return
	} else if n != 1 {
		errorDesc, _ := json.Marshal(JsonError{"Card not found"})
		http.Error(w, string(errorDesc), 400)
		return
	}

	err = database.C(robotC).Update(
		bson.M{"url": vars[robotURL]},
		bson.M{"$set": bson.M{cardId: vars[cardId]}})

	if err != nil {
		errorDesc, _ := json.Marshal(JsonError{err.Error()})
		http.Error(w, string(errorDesc), 400)
	} else {
		json.NewEncoder(w).Encode(JsonOK{"Done"})
	}
}

func DissociateRobot(w http.ResponseWriter, r *http.Request) {
	vars, _, err := initSession(w, r)
	if err != nil {
		return
	}

	err = database.C(robotC).Update(
		bson.M{"url": vars[robotURL]},
		bson.M{"$set": bson.M{cardId: ""}})

	if err != nil {
		errorDesc, _ := json.Marshal(JsonError{err.Error()})
		http.Error(w, string(errorDesc), 400)
	} else {
		json.NewEncoder(w).Encode(JsonOK{"Done"})
	}
}

func GetRobot(w http.ResponseWriter, r *http.Request) {
	vars, _, err := initSession(w, r)
	if err != nil {
		return
	}

	var robot Robot
	err = database.C(robotC).Find(bson.M{"url": vars[robotURL]}).One(&robot)
	if err != nil {
		errorDesc, _ := json.Marshal(JsonError{err.Error()})
		http.Error(w, string(errorDesc), 400)
	} else {
		json.NewEncoder(w).Encode(robot)
	}
}

func PutRobot(w http.ResponseWriter, r *http.Request) {
	vars, _, err := initSession(w, r)
	if err != nil {
		return
	}

	var payload struct {
		Name string `json"name"`
	}

	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		errorDesc, _ := json.Marshal(JsonError{err.Error()})
		http.Error(w, string(errorDesc), 400)
		return
	}

	var robot Robot
	robot.URL = vars[robotURL]
	robot.Name = payload.Name
	robot.CardId = ""

	_, err = database.C(robotC).Upsert(bson.M{"url": vars[robotURL]}, robot)
	if err != nil {
		errorDesc, _ := json.Marshal(JsonError{err.Error()})
		http.Error(w, string(errorDesc), 500)
	} else {
		json.NewEncoder(w).Encode(JsonOK{"Done"})
	}
}

func DelRobot(w http.ResponseWriter, r *http.Request) {
	vars, _, err := initSession(w, r)
	if err != nil {
		return
	}

	_, err = database.C(robotC).RemoveAll(bson.M{"url": vars[robotURL]})
	if err != nil {
		errorDesc, _ := json.Marshal(JsonError{err.Error()})
		http.Error(w, string(errorDesc), 400)

	} else {
		json.NewEncoder(w).Encode(JsonOK{"Done"})
	}
}

func GetRobots(w http.ResponseWriter, r *http.Request) {
	_, _, err := initSession(w, r)
	if err != nil {
		return
	}

	var robots []Robot
	err = database.C(robotC).Find(nil).All(&robots)
	if err != nil {
		errorDesc, _ := json.Marshal(JsonError{err.Error()})
		http.Error(w, string(errorDesc), 400)
	} else {
		json.NewEncoder(w).Encode(robots)
	}
}

func Run(w http.ResponseWriter, r *http.Request) {
	vars, _, err := initSession(w, r)
	if err != nil {
		return
	}

	var robot Robot
	err = database.C(robotC).Find(bson.M{cardId: vars[cardId]}).One(&robot)
	if err != nil {
		errorDesc, _ := json.Marshal(JsonError{err.Error()})
		http.Error(w, string(errorDesc), 400)
	} else {
		log.Infof("Sending run command to robot %v", robot.URL)
		json.NewEncoder(w).Encode(JsonOK{"Done"})
	}
}

func Stop(w http.ResponseWriter, r *http.Request) {
	vars, _, err := initSession(w, r)
	if err != nil {
		return
	}

	var robot Robot
	err = database.C(robotC).Find(bson.M{cardId: vars[cardId]}).One(&robot)
	if err != nil {
		errorDesc, _ := json.Marshal(JsonError{err.Error()})
		http.Error(w, string(errorDesc), 400)
	} else {
		log.Infof("Sending stop command to robot %v", robot.URL)
		json.NewEncoder(w).Encode(JsonOK{"Done"})
	}
}

func Upload(w http.ResponseWriter, r *http.Request) {
	vars, _, err := initSession(w, r)
	if err != nil {
		return
	}

	var robot Robot
	err = database.C(robotC).Find(bson.M{cardId: vars[cardId]}).One(&robot)
	if err != nil {
		errorDesc, _ := json.Marshal(JsonError{err.Error()})
		http.Error(w, string(errorDesc), 400)
	} else {
		log.Infof("Uploading code to robot %v", robot.URL)
		json.NewEncoder(w).Encode(JsonOK{"Done"})
	}
}

func main() {
	var port = flag.Int("port", 8081, "port")
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
	mongoSession, err := mgo.Dial(*mongoServer)
	if err != nil {
		log.Fatal(err)
	}
	database = mongoSession.DB(dbName)
	store = mongostore.NewMongoStore(
		database.C(sessionC),
		0, true, []byte(*secretKey))

	store.Options.Domain = *domain

	r := mux.NewRouter()
	r.HandleFunc(prefix+"/info", GetInfo).Methods("GET")

	r.HandleFunc(prefix+"/card/{"+cardId+"}", GetCard).Methods("GET")
	r.HandleFunc(prefix+"/card/{"+cardId+"}", PutCard).Methods("PUT")

	r.HandleFunc(prefix+"/robot/{"+robotURL+"}", GetRobot).Methods("GET")
	r.HandleFunc(prefix+"/robot/{"+robotURL+"}", PutRobot).Methods("PUT")
	r.HandleFunc(prefix+"/robot/{"+robotURL+"}", DelRobot).Methods("DELETE")

	r.HandleFunc(prefix+"/robot/{"+robotURL+"}/card/{"+cardId+"}", AssociateRobot).Methods("PUT")
	r.HandleFunc(prefix+"/robot/{"+robotURL+"}/card", DissociateRobot).Methods("DELETE")

	r.HandleFunc(prefix+"/robots", GetRobots).Methods("GET")
	r.HandleFunc(prefix+"/card/{"+cardId+"}/run", Run).Methods("GET")
	r.HandleFunc(prefix+"/card/{"+cardId+"}/stop", Stop).Methods("GET")
	r.HandleFunc(prefix+"/card/{"+cardId+"}/upload", Upload).Methods("GET")

	http.Handle("/", r)
	log.Infof("Ready, listening on port %d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
