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
	"os"
	"path/filepath"
	"time"
)

const (
	dbName       = "thymio_captain"
	sessionC     = "sessions"
	maxAge       = 24 * 3600
	sessionKey   = "session-key"
	root         = "internal_pages"
	tokenRndLen  = 20
	tokenSignLen = 20
)

var (
	database       *mgo.Session
	store          *mongostore.MongoStore
	adminSecretKey *string
	startSecretKey *string
	templates      = make(map[string]*template.Template)
)

func initSession(w http.ResponseWriter, r *http.Request) (vars map[string]string, session *sessions.Session, err error) {
	database.Refresh()
	vars = mux.Vars(r)
	session, err = store.Get(r, sessionKey)
	log.Debugf("Session ID = %v", session.ID)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Session Error", 500)
	} else {
		log.Debug("Session OK")
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

	if len(t) != tokenRndLen+tokenSignLen {
		log.Infof("Invalid token length: %d", len(t))
		return false
	}
	data := t[0:tokenRndLen]
	sig := t[tokenRndLen : tokenRndLen+tokenSignLen]
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
		err = templates["login-ok.html"].Execute(w, nil)
	} else {
		log.Infof("Bad Card Login: %v", vars["CardId"])
		session.Values["admin"] = "0"
		sessions.Save(r, w)
		err = templates["login-failed.html"].Execute(w, nil)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	err = templates["logout.html"].Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	err := templates["index.html"].Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Help(w http.ResponseWriter, r *http.Request) {
	err := templates["help.html"].Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func About(w http.ResponseWriter, r *http.Request) {
	err := templates["about.html"].Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func notFound(w http.ResponseWriter, r *http.Request) {
	err := templates["404.html"].Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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

		var fileName string

		if session.Values["admin"] == "1" {
			log.Debug("Sending Admin UI")
			fileName = root + "/admin.html"
		} else {
			log.Debug("Sending User UI")
			fileName = root + "/public.html"
		}
		f, err := os.Open(fileName)
		if err == nil {
			http.ServeContent(w, r, fileName, time.Time{}, f)
		} else {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		log.Infof("Bad Start page: %v", vars["CardId"])
		err = templates["bad-cards.html"].Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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

	tpls := template.Must(template.ParseGlob("internal_pages/templates/*"))
	nameList, err := filepath.Glob("internal_pages/*.html")
	if err != nil {
		panic(err)
	}
	for _, name := range nameList {
		log.Debugf("Reading %s", name)
		key := filepath.Base(name)
		t, _ := tpls.Clone()
		templates[key] = t
		_, err = templates[key].ParseFiles(name)
		if err != nil {
			panic(err)
		}
	}

	r := mux.NewRouter()
	r.HandleFunc("/start/{CardId}", Start)
	r.HandleFunc("/cardlogin/{CardId}", CardLogin)
	r.HandleFunc("/logout", Logout)
	r.HandleFunc("/debug", Debug)

	r.HandleFunc("/", Index)
	r.HandleFunc("/about", About)
	r.HandleFunc("/help", Help)

	r.NotFoundHandler = http.HandlerFunc(notFound)

	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/vendor/", http.StripPrefix("/vendor/", http.FileServer(http.Dir("vendor"))))

	http.Handle("/", r)
	log.Infof("Ready, listening on port %d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
