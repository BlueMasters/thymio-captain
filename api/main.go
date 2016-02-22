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

//   _   _                     _                             _        _
//  | |_| |__  _   _ _ __ ___ (_) ___         ___ __ _ _ __ | |_ __ _(_)_ __
//  | __| '_ \| | | | '_ ` _ \| |/ _ \ _____ / __/ _` | '_ \| __/ _` | | '_ \
//  | |_| | | | |_| | | | | | | | (_) |_____| (_| (_| | |_) | || (_| | | | | |
//   \__|_| |_|\__, |_| |_| |_|_|\___/       \___\__,_| .__/ \__\__,_|_|_| |_|
//             |___/                                  |_|
//

// API
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/kidstuff/mongostore"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
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

// sessionValues extracts session info from the HTTP header. It first looks for a "Authorization" header and then
// it looks for a cookie. It returns a map of the session data.
func sessionValues(r *http.Request) (values map[interface{}]interface{}, err error) {
	err = nil
	// check authorization header
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(auth, "Cookie") {
		log.Debugf("Authorization header: %v", auth)
		t := strings.Split(auth, " ")
		if len(t) <= 1 {
			err = errors.New("Invalid Authorization header")
		}
		var sessionID string
		if err == nil {
			err = securecookie.DecodeMulti(sessionKey, t[len(t)-1], &sessionID, store.Codecs...)
			log.Debugf("Session ID = %v", sessionID)
		}
		if err == nil {
			if !bson.IsObjectIdHex(sessionID) {
				err = errors.New("Invalid session ID")
			}
		}
		var s mongostore.Session
		if err == nil {
			err = database.C(sessionC).FindId(bson.ObjectIdHex(sessionID)).One(&s)
		}
		if err == nil {
			err = securecookie.DecodeMulti(sessionKey, s.Data, &values, store.Codecs...)
		}
	} else {
		session, err := store.Get(r, sessionKey)
		if err == nil {
			log.Debugf("Cookie found / ID = %v", session.ID)
			values = session.Values
		}
	}
	if err != nil {
		log.Error(err.Error())
	}
	return
}

// initSession "bootstrap" the HTTP session. It configure the variables in the multiplexer and returns
// the session values. It also add cache control headers.
func initSession(w http.ResponseWriter, r *http.Request) (
	vars map[string]string, values map[interface{}]interface{}, err error) {

	log.Debugf("request: %v", r.URL.String())
	database.Session.Refresh()
	vars = mux.Vars(r)

	values, err = sessionValues(r)

	admin, ok := values["admin"]
	if ok {
		if admin == "1" {
			log.Debug("admin: yes")
		} else {
			log.Debugf("admin: no (%v)", admin)
		}
	} else {
		log.Debug("admin: UNKNOWN")
	}

	cardId, ok := values["cardId"]
	if ok {
		log.Debugf("cardId: %v", cardId)
	} else {
		log.Debug("cardId: UNKNOWN")
	}
	w.Header().Set("Cache-Control", "max-age=0, no-cache, no-store")
	w.Header().Set("Pragma", "no-cache")

	return
}

// checkAdmin returns nil if the session is an authorized admin.
func checkAdmin(w http.ResponseWriter, session map[interface{}]interface{}) error {
	value, ok := session["admin"]
	if ok && value == "1" {
		return nil
	} else {
		err := errors.New("Not authorized")
		http.Error(w, err.Error(), 401)
		return err
	}
}

// report check the err argument and if not nil, it logs the error and returns the error using HTTP
func report(w http.ResponseWriter, err error) error {
	if err != nil {
		errorDesc, _ := json.Marshal(JsonError{err.Error()})
		log.Info(err)
		http.Error(w, string(errorDesc), 400)
	}
	return err
}

// GetInfo is the handler for the "GET /info" method
func GetInfo(w http.ResponseWriter, r *http.Request) {
	_, session, err := initSession(w, r)
	if report(w, err) != nil {
		return
	}

	var info Info
	if cardId, ok := session["cardId"]; ok {
		info.CardId = cardId.(string)
	}

	if admin, ok := session["admin"]; ok {
		info.IsAdmin = admin == "1"
	}

	json.NewEncoder(w).Encode(info)
}

// GetCard is the handler for the "GET /card/{cardId}" method
func GetCard(w http.ResponseWriter, r *http.Request) {
	vars, _, err := initSession(w, r)
	if report(w, err) != nil {
		return
	}

	var card Card
	err = database.C(cardC).Find(bson.M{"cardId": vars["cardId"]}).One(&card)
	if report(w, err) != nil {
		return
	}
	json.NewEncoder(w).Encode(card)
}

// PutCard is the handler for the "PUT|POST /card/{cardId}" method
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

// GetRobots is the handler for the "GET /robots" method
func GetRobots(w http.ResponseWriter, r *http.Request) {
	_, session, err := initSession(w, r)
	if report(w, err) != nil {
		return
	}

	if checkAdmin(w, session) != nil {
		return
	}

	var robots []Robot
	err = database.C(robotC).Find(nil).All(&robots)
	if report(w, err) != nil {
		return
	}
	json.NewEncoder(w).Encode(robots)
}

// GetRobot is the handler for the "GET /robot/{robotName}" method
func GetRobot(w http.ResponseWriter, r *http.Request) {
	vars, session, err := initSession(w, r)
	if report(w, err) != nil {
		return
	}

	if checkAdmin(w, session) != nil {
		return
	}

	var robot Robot
	err = database.C(robotC).Find(bson.M{"name": vars["robotName"]}).One(&robot)
	if report(w, err) != nil {
		return
	}
	json.NewEncoder(w).Encode(robot)

}

// PutRobot is the handler for the "PUT|POST /robot/{robotName}" method
func PutRobot(w http.ResponseWriter, r *http.Request) {
	vars, session, err := initSession(w, r)
	if report(w, err) != nil {
		return
	}

	if checkAdmin(w, session) != nil {
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

// DelRobot is the handler for the "DELETE /robot/{robotName}" method
func DelRobot(w http.ResponseWriter, r *http.Request) {
	vars, session, err := initSession(w, r)
	if report(w, err) != nil {
		return
	}

	if checkAdmin(w, session) != nil {
		return
	}

	_, err = database.C(robotC).RemoveAll(bson.M{"name": vars["robotName"]})
	if report(w, err) != nil {
		return
	}
	json.NewEncoder(w).Encode(JsonOK{"done"})
}

// PingRobot is the handler for the "GET /robot/{robotName}/ping" method
func PingRobot(w http.ResponseWriter, r *http.Request) {
	vars, session, err := initSession(w, r)
	if report(w, err) != nil {
		return
	}

	if checkAdmin(w, session) != nil {
		return
	}

	var robot Robot
	err = database.C(robotC).Find(bson.M{"name": vars["robotName"]}).One(&robot)
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

// AssociateRobot is the handler for the "PUT|POST /robot/{robotName}/card/{cardId}" method
func AssociateRobot(w http.ResponseWriter, r *http.Request) {
	vars, session, err := initSession(w, r)
	if report(w, err) != nil {
		return
	}

	if checkAdmin(w, session) != nil {
		return
	}

	n, err := database.C(cardC).Find(bson.M{"cardId": vars["cardId"]}).Count()
	if report(w, err) != nil {
		return
	} else if n != 1 {
		report(w, errors.New("Card not found"))
		return
	}

	// check if there is already an association
	n, err = database.C(robotC).Find(bson.M{"cardId": vars["cardId"]}).Count()
	report(w, err)
	if err != nil {
		return
	} else if n > 0 {
		report(w, errors.New("Robot already associated"))
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

// DissociateRobot is the handler for the "DELETE /robot/{robotName}/card/" method
func DissociateRobot(w http.ResponseWriter, r *http.Request) {
	vars, session, err := initSession(w, r)
	if report(w, err) != nil {
		return
	}

	if checkAdmin(w, session) != nil {
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

// PingCardRobot is the handler for the "GET /card/{cardId}/ping" method
func PingCardRobot(w http.ResponseWriter, r *http.Request) {
	vars, _, err := initSession(w, r)
	if report(w, err) != nil {
		return
	}

	var robot Robot
	err = database.C(robotC).Find(bson.M{"cardId": vars["cardId"]}).One(&robot)
	if report(w, err) != nil {
		return
	}
	log.Infof("Received ping command from card: %v", vars["cardId"])
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

// RunCardRobot is the handler for the "GET /card/{cardId}/run" method
func RunCardRobot(w http.ResponseWriter, r *http.Request) {
	vars, _, err := initSession(w, r)
	if report(w, err) != nil {
		return
	}

	var robot Robot
	err = database.C(robotC).Find(bson.M{"cardId": vars["cardId"]}).One(&robot)
	if report(w, err) != nil {
		return
	}
	log.Infof("Received run command from card: %v", vars["cardId"])
	u, _ := url.Parse(robot.URL)
	u.Path = filepath.Join(u.Path, "/run")
	log.Infof("Sending run command to robot: %v", u)
	var client http.Client
	res, err := client.Get(u.String())
	if report(w, err) != nil {
		return
	}
	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
}

// StopCardRobot is the handler for the "GET /card/{cardId}/stop" method
func StopCardRobot(w http.ResponseWriter, r *http.Request) {
	vars, _, err := initSession(w, r)
	if report(w, err) != nil {
		return
	}

	var robot Robot
	err = database.C(robotC).Find(bson.M{"cardId": vars["cardId"]}).One(&robot)
	if report(w, err) != nil {
		return
	}
	log.Infof("Received stop command from card: %v", vars["cardId"])
	u, _ := url.Parse(robot.URL)
	u.Path = filepath.Join(u.Path, "/stop")
	log.Debugf("Sending stop command to robot: %v", u)
	var client http.Client
	res, err := client.Get(u.String())
	if report(w, err) != nil {
		return
	}
	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
}

// UploadCardRobot is the handler for the "GET|PUT|POST /card/{cardId}/upload" method
func UploadCardRobot(w http.ResponseWriter, r *http.Request) {
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

	log.Infof("Received upload command from card: %v", vars["cardId"])
	u, _ := url.Parse(robot.URL)
	u.Path = filepath.Join(u.Path, "/upload")
	var client http.Client

	cardJ, err := json.Marshal(card)
	if report(w, err) != nil {
		return
	}
	log.Debugf("Uploading card to %v: %v", u, string(cardJ))
	cReq, err := http.NewRequest("PUT", u.String(), bytes.NewReader(cardJ))
	if report(w, err) != nil {
		return
	}
	cReq.Header.Set("Content-Type", "application/json; charset=UTF-8")
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

// ServeHTTP is a HTTP handler that implements permissive CORS rules
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

	// Info
	r.HandleFunc(prefix+"/info", GetInfo).Methods("GET")

	// Card management
	r.HandleFunc(prefix+"/card/{cardId}", GetCard).Methods("GET")
	r.HandleFunc(prefix+"/card/{cardId}", PutCard).Methods("PUT", "POST")

	// Robot management
	r.HandleFunc(prefix+"/robots", GetRobots).Methods("GET")
	r.HandleFunc(prefix+"/robot/{robotName}", GetRobot).Methods("GET")
	r.HandleFunc(prefix+"/robot/{robotName}", PutRobot).Methods("PUT", "POST")
	r.HandleFunc(prefix+"/robot/{robotName}", DelRobot).Methods("DELETE")
	r.HandleFunc(prefix+"/robot/{robotName}/ping", PingRobot).Methods("GET")

	// Robot/Card associations
	r.HandleFunc(prefix+"/robot/{robotName}/card/{cardId}", AssociateRobot).Methods("PUT", "POST")
	r.HandleFunc(prefix+"/robot/{robotName}/card", DissociateRobot).Methods("DELETE")

	// Robot control
	r.HandleFunc(prefix+"/card/{cardId}/ping", PingCardRobot).Methods("GET")
	r.HandleFunc(prefix+"/card/{cardId}/run", RunCardRobot).Methods("GET")
	r.HandleFunc(prefix+"/card/{cardId}/stop", StopCardRobot).Methods("GET")
	r.HandleFunc(prefix+"/card/{cardId}/upload", UploadCardRobot).Methods("GET", "PUT", "POST")

	http.Handle("/", &CorsServer{r})

	log.Infof("Ready, listening on port %d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
