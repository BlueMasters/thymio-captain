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
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/kidstuff/mongostore"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
)

const (
	dbName     = "thymio_captain"
	sessionC   = "sessions"
	sessionKey = "session-key"
	prefix     = "/v1"
	cardC      = "cards"
	robotC     = "robots"
)

type Info struct {
	CardId  string `json:"cardId" bson:"cardId"`
	IsAdmin bool   `json:"isAdmin" bson:"isAdmin"`
}

type Robot struct {
	Name   string `json:"name" bson:"name"`
	URL    string `json:"url" bson:"url"`
	CardId string `json:"cardId" bson:"cardId"`
}

type Card struct {
	CardId  string `json:"cardId" bson:"cardId"`
	Program []byte `json:"program" bson:"program"`
}

type JsonError struct {
	ErrorDescription string `json:"errorDescription"`
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
		http.Error(w, "session error", 500)
	}
	return
}

func report(w http.ResponseWriter, err error) error {
	if err != nil {
		errorDesc, _ := json.Marshal(JsonError{err.Error()})
		log.Info(err)
		http.Error(w, string(errorDesc), 400)
	}
	return err
}

func GetCard(w http.ResponseWriter, r *http.Request) {
	vars, _, err := initSession(w, r)
	if report(w, err) != nil {
		return
	}

	var card Card
	err = database.C(cardC).Find(bson.M{"cardId": vars["cardId"]}).One(&card)
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
	if report(w, err) != nil {
		return
	}

	var payload struct {
		Program []byte `json:"program"`
	}

	err = json.NewDecoder(r.Body).Decode(&payload)
	if report(w, err) != nil {
		return
	}

	var card Card
	card.CardId = vars["cardId"]
	card.Program = payload.Program

	_, err = database.C(cardC).Upsert(bson.M{"cardId": vars["cardId"]}, card)
	report(w, err)
	if err != nil {
		return
	}
	json.NewEncoder(w).Encode(JsonOK{"done"})
}

func GetInfo(w http.ResponseWriter, r *http.Request) {
	_, session, err := initSession(w, r)
	if report(w, err) != nil {
		return
	}

	var info Info
	if cardId, ok := session.Values["cardId"]; ok {
		info.CardId = cardId.(string)
	}

	if admin, ok := session.Values["admin"]; ok {
		info.IsAdmin = admin == "1"
	}

	json.NewEncoder(w).Encode(info)
}

func AssociateRobot(w http.ResponseWriter, r *http.Request) {
	vars, _, err := initSession(w, r)
	if report(w, err) != nil {
		return
	}

	n, err := database.C(cardC).Find(bson.M{"cardId": vars["cardId"]}).Count()
	if report(w, err) != nil {
		return
	} else if n != 1 {
		report(w, errors.New("Card not found"))
		return
	}

	// remove the card from all robots (to prevent duplicate)
	_, err = database.C(robotC).UpdateAll(
		bson.M{"cardId": vars["cardId"]},
		bson.M{"$set": bson.M{"cardId": ""}})
	report(w, err)
	if err != nil {
		return
	}

	// associate the robot with the card
	err = database.C(robotC).Update(
		bson.M{"name": vars["robotName"]},
		bson.M{"$set": bson.M{"cardId": vars["cardId"]}})

	report(w, err)
	if err != nil {
		return
	}
	json.NewEncoder(w).Encode(JsonOK{"done"})
}

func DissociateRobot(w http.ResponseWriter, r *http.Request) {
	vars, _, err := initSession(w, r)
	if report(w, err) != nil {
		return
	}

	err = database.C(robotC).Update(
		bson.M{"name": vars["robotName"]},
		bson.M{"$set": bson.M{"cardId": ""}})

	if report(w, err) != nil {
		return
	}
	json.NewEncoder(w).Encode(JsonOK{"done"})

}

func GetRobot(w http.ResponseWriter, r *http.Request) {
	vars, _, err := initSession(w, r)
	if report(w, err) != nil {
		return
	}

	var robot Robot
	err = database.C(robotC).Find(bson.M{"name": vars["robotName"]}).One(&robot)
	if report(w, err) != nil {
		return
	}
	json.NewEncoder(w).Encode(robot)

}

func PutRobot(w http.ResponseWriter, r *http.Request) {
	vars, _, err := initSession(w, r)
	if report(w, err) != nil {
		return
	}

	var payload struct {
		URL string `json:"url"`
	}

	err = json.NewDecoder(r.Body).Decode(&payload)
	if report(w, err) != nil {
		return
	}

	var robot Robot
	robot.Name = vars["robotName"]
	robot.URL = payload.URL
	robot.CardId = ""

	_, err = database.C(robotC).Upsert(
		bson.M{"name": vars["robotName"]},
		bson.M{
			"$set":         bson.M{"url": robot.URL},
			"$setOnInsert": bson.M{"cardId": ""}})
	if report(w, err) != nil {
		return
	}
	json.NewEncoder(w).Encode(JsonOK{"done"})

}

func DelRobot(w http.ResponseWriter, r *http.Request) {
	vars, _, err := initSession(w, r)
	if report(w, err) != nil {
		return
	}

	_, err = database.C(robotC).RemoveAll(bson.M{"name": vars["robotName"]})
	if report(w, err) != nil {
		return
	}
	json.NewEncoder(w).Encode(JsonOK{"done"})
}

func GetRobots(w http.ResponseWriter, r *http.Request) {
	_, _, err := initSession(w, r)
	if report(w, err) != nil {
		return
	}

	var robots []Robot
	err = database.C(robotC).Find(nil).All(&robots)
	if report(w, err) != nil {
		return
	}
	json.NewEncoder(w).Encode(robots)
}

func Ping(w http.ResponseWriter, r *http.Request) {
	vars, _, err := initSession(w, r)
	if report(w, err) != nil {
		return
	}

	var robot Robot
	err = database.C(robotC).Find(bson.M{"cardId": vars["cardId"]}).One(&robot)
	if report(w, err) != nil {
		return
	}
	u, _ := url.Parse(robot.URL)
	u.Path = filepath.Join(u.Path, "/ping")

	log.Infof("Sending ping command to robot: %v", u)
	var client http.Client
	res, err := client.Get(u.String())
	if report(w, err) != nil {
		return
	}
	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
}

func Run(w http.ResponseWriter, r *http.Request) {
	vars, _, err := initSession(w, r)
	if report(w, err) != nil {
		return
	}

	var robot Robot
	err = database.C(robotC).Find(bson.M{"cardId": vars["cardId"]}).One(&robot)
	if report(w, err) != nil {
		return
	}
	log.Infof("Sending run command to robot %v", robot.URL)
	u, _ := url.Parse(robot.URL)
	u.Path = filepath.Join(u.Path, "/run")
	var client http.Client
	res, err := client.Get(u.String())
	if report(w, err) != nil {
		return
	}
	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
}

func Stop(w http.ResponseWriter, r *http.Request) {
	vars, _, err := initSession(w, r)
	if report(w, err) != nil {
		return
	}

	var robot Robot
	err = database.C(robotC).Find(bson.M{"cardId": vars["cardId"]}).One(&robot)
	if report(w, err) != nil {
		return
	}
	log.Infof("Sending stop command to robot %v", robot.URL)
	u, _ := url.Parse(robot.URL)
	u.Path = filepath.Join(u.Path, "/stop")
	var client http.Client
	res, err := client.Get(u.String())
	if report(w, err) != nil {
		return
	}
	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
}

func Upload(w http.ResponseWriter, r *http.Request) {
	vars, _, err := initSession(w, r)
	if report(w, err) != nil {
		return
	}

	var robot Robot
	err = database.C(robotC).Find(bson.M{"cardId": vars["cardId"]}).One(&robot)
	if report(w, err) != nil {
		return
	}

	var card Card
	err = database.C(cardC).Find(bson.M{"cardId": vars["cardId"]}).One(&card)
	if report(w, err) != nil {
		return
	}

	u, _ := url.Parse(robot.URL)
	u.Path = filepath.Join(u.Path, "/stop")
	var client http.Client

	cardJ, err := json.Marshal(card)
	if report(w, err) != nil {
		return
	}
	cReq, err := http.NewRequest("PUT", u.String(), bytes.NewReader(cardJ))
	if report(w, err) != nil {
		return
	}
	res, err := client.Do(cReq)
	if report(w, err) != nil {
		return
	}

	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
}

type CorsServer struct {
	r *mux.Router
}

func (s *CorsServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if req.Method == "OPTIONS" {
		return
	}
	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Lets Gorilla work
	s.r.ServeHTTP(rw, req)
}

func main() {
	var port = flag.Int("port", 8081, "port")
	var debug = flag.Bool("debug", false, "run in debug mode")
	var domain = flag.String("domain", "thymio.tk", "Domain name (for the cookie)")
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

	r.HandleFunc(prefix+"/card/{cardId}", GetCard).Methods("GET")
	r.HandleFunc(prefix+"/card/{cardId}", PutCard).Methods("PUT", "POST")

	r.HandleFunc(prefix+"/robot/{robotName}", GetRobot).Methods("GET")
	r.HandleFunc(prefix+"/robot/{robotName}", PutRobot).Methods("PUT", "POST")
	r.HandleFunc(prefix+"/robot/{robotName}", DelRobot).Methods("DELETE")

	r.HandleFunc(prefix+"/robot/{robotName}/card/{cardId}", AssociateRobot).Methods("PUT", "POST")
	r.HandleFunc(prefix+"/robot/{robotName}/card", DissociateRobot).Methods("DELETE")

	r.HandleFunc(prefix+"/robots", GetRobots).Methods("GET")
	r.HandleFunc(prefix+"/card/{cardId}/ping", Ping).Methods("GET")
	r.HandleFunc(prefix+"/card/{cardId}/run", Run).Methods("GET")
	r.HandleFunc(prefix+"/card/{cardId}/stop", Stop).Methods("GET")
	r.HandleFunc(prefix+"/card/{cardId}/upload", Upload).Methods("GET")

	http.Handle("/", &CorsServer{r})

	log.Infof("Ready, listening on port %d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
