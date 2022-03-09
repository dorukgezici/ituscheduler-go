package app

import (
	"encoding/json"
	"fmt"
	"github.com/imdario/mergo"
	"gorm.io/gorm/clause"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func LoadUserFixtures(filename string) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		log.Printf("failed to open json file: %s, error: %v", filename, err)
	}

	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Printf("failed to read json file, error: %v", err)
	}

	var users []User
	if err = json.Unmarshal(jsonData, &users); err != nil {
		log.Printf("failed to unmarshal json file, error: %v\n", err)
		log.Printf("failed to close jsonFile, error: %s", jsonFile.Close().Error())
	}

	DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&users)
	for i, user := range users {
		log.Printf("USER#%d: ID: %d Username: %s", i, user.ID, user.Username)
	}
}

func LoadPostFixtures(filename string) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		log.Printf("failed to open json file: %s, error: %v", filename, err)
	}

	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Printf("failed to read json file, error: %v", err)
	}

	var posts []Post
	if err = json.Unmarshal(jsonData, &posts); err != nil {
		log.Printf("failed to unmarshal json file, error: %v\n", err)
		log.Printf("failed to close jsonFile, error: %s", jsonFile.Close().Error())
	}

	DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&posts)
	for i, post := range posts {
		log.Printf("POST#%d: Author: %s Date: %s", i, post.Author, post.Date)
	}
}

func render(filename string, w http.ResponseWriter, r *http.Request, data map[string]interface{}) {
	fm := template.FuncMap{
		"safe": func(value interface{}) template.HTML {
			return template.HTML(fmt.Sprint(value))
		},
		"date": func(date time.Time) template.HTML {
			return template.HTML(date.Format("Jan 2, 2006, 3:04 PM"))
		},
		"slugify": func(value interface{}) template.HTML {
			return template.HTML(fmt.Sprint(value))
		},
		"pathContains": func(value interface{}) bool {
			return strings.Contains(r.URL.Path, fmt.Sprint(value))
		},
	}
	tpl := template.Must(template.New(filename).Funcs(fm).ParseFiles("templates/base.gohtml", "templates/"+filename))

	var user User
	if cookie, err := r.Cookie("session"); err == nil {
		DB.Joins("JOIN sessions ON sessions.user_id = users.id").Find(&user, "sessions.token = ? AND sessions.deleted_at IS NULL", cookie.Value)
	}

	initialData := map[string]interface{}{
		"Time": time.Now(),
		"Path": r.URL.Path,
		"User": user,
	}
	if err := mergo.Merge(&data, initialData); err != nil {
		panic(err)
	}

	if err := tpl.Execute(w, data); err != nil {
		log.Printf("failed to render template: %s, error: %v", filename, err)
		panic(err)
	}
}
