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

// Frontend
package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
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
	root       = "internal_pages"
)

var (
	database       *mgo.Session
	store          *mongostore.MongoStore
	adminSecretKey *string
	startSecretKey *string
)

func initSession(w http.ResponseWriter, r *http.Request) (vars map[string]string, session *sessions.Session, err error) {
	database.Refresh()
	vars = mux.Vars(r)
	session, err = store.Get(r, sessionKey)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Session Error", 500)
	} else {
		w.Header().Set("Cache-Control", "max-age=0, no-cache, no-store")
		w.Header().Set("Pragma", "no-cache")
	}
	return
}

func isValidToken(token string, key string) bool {
	t, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return false
	}

	if len(t) != 40 { // TODO: replace MAGIC
		log.Infof("Invalid token length: %d", len(t))
		return false
	}
	data := t[0:20] // TODO: replace MAGIC
	sig := t[20:40] // TODO: replace MAGIC
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write(data)
	log.Debugf("HMAC: %v", mac.Sum(nil))
	log.Debugf("EXP : %v", sig)
	return hmac.Equal(mac.Sum(nil), sig)
}

func CardLogin(w http.ResponseWriter, r *http.Request) {
	vars, session, err := initSession(w, r)
	if err != nil {
		return
	}

	if isValidToken(vars["CardId"], *adminSecretKey) {
		log.Debugf("Valid Card Login: %v", vars["CardId"])
		session.Values["admin"] = "1"
		sessions.Save(r, w)
		http.ServeFile(w, r, root+"/login-ok.html")
	} else {
		log.Infof("Bad Card Login: %v", vars["CardId"])
		session.Values["admin"] = "0"
		sessions.Save(r, w)
		http.ServeFile(w, r, root+"/login-failed.html")
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	_, session, err := initSession(w, r)
	if err != nil {
		return
	}

	log.Debug("Logout")
	session.Values["admin"] = "0"
	sessions.Save(r, w)
	http.ServeFile(w, r, root+"/logout.html")
}

func Debug(w http.ResponseWriter, r *http.Request) {
	_, session, err := initSession(w, r)
	if err != nil {
		return
	}

	log.Debug("Debug page")
	tmpl, err := template.ParseFiles(root + "/debug.html")
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

	if *startSecretKey == "" || isValidToken(vars["CardId"], *startSecretKey) {
		log.Debugf("Valid Start page: %v", vars["CardId"])
		session.Values["cardId"] = vars["CardId"]
		sessions.Save(r, w)

		if session.Values["admin"] == "1" {
			http.ServeFile(w, r, root+"/admin.html")
		} else {
			http.ServeFile(w, r, root+"/public.html")
		}
	} else {
		log.Infof("Bad Start page: %v", vars["CardId"])
		http.ServeFile(w, r, root+"/bad-card.html")
	}
}

func main() {
	var port = flag.Int("port", 8080, "port")
	var debug = flag.Bool("debug", false, "run in debug mode")
	var domain = flag.String("domain", "thymio.tk", "Domain name (for the cookie)")
	var mongoServer = flag.String("mongo-server", "localhost", "MongoDB server URL")
	var cookieSecretKey = flag.String("cookie-secret-key", "not-so-secret", "Secret key (for secure cookies)")
	adminSecretKey = flag.String("admin-secret-key", "change-me", "Secret key (for admin card-login)")
	startSecretKey = flag.String("start-secret-key", "", "Secret key (for start ID)")

	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
		log.Debug("Debug mode")
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if *startSecretKey == "" {
		log.Warn("Running without start id validation")
	} else {
		log.Info("Start id validation enabled")
	}

	var err error
	database, err = mgo.Dial(*mongoServer)
	if err != nil {
		log.Fatal(err)
	}
	store = mongostore.NewMongoStore(
		database.DB(dbName).C(sessionC),
		maxAge, true, []byte(*cookieSecretKey))

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
