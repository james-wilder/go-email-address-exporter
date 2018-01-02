package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

const CookieName = "accounts"

type Context struct {
	Accounts []Account
}

type Account struct {
	Name string
	Hostname string
}

func main() {
    http.HandleFunc("/", viewHandler)
//    http.HandleFunc("/edit/", editHandler)
//    http.HandleFunc("/save/", saveHandler)
    http.ListenAndServe(":8080", nil)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	context := getContext(r)

	r.ParseForm()
	fmt.Println("Method", r.Method)
	fmt.Println("Form", r.Form)
	fmt.Println("PostForm", r.PostForm)
	if r.Method == "POST" {
		fmt.Println("Action", r.Form["action"])
		if (r.Form["action"][0] == "add") {
			name := r.Form["name"][0]
			hostname := r.Form["hostname"][0]
			context.Accounts = append(context.Accounts, Account{Name:name,Hostname:hostname})
		}
		if (r.Form["action"][0] == "delete") {
			var newAccounts []Account
			name := r.Form["name"][0]
			for _, account := range context.Accounts {
				if account.Name != name {
					newAccounts = append(newAccounts, account)
				}
			}
			context.Accounts = newAccounts
		}
	}

	showView(w, r, context)
}

func getContext(r *http.Request) Context {
	var cookieValue string

	// default test data
	a1 := Account{Name: "Gmail", Hostname: "gmail.com"}
	a2 := Account{Name: "Fastmail", Hostname: "fastmail.com"}
	c := Context{Accounts:[]Account{a1, a2}}

	cookie, err := r.Cookie(CookieName)

	if err != nil {
		fmt.Println(err)
	} else {
		decoded, err := base64.StdEncoding.DecodeString(cookie.Value)
		if err != nil {
			fmt.Println(err)
		} else {
			cookieValue = string(decoded)
		}
	}

	if cookieValue != "" {
		var cookieContext Context
		err = json.Unmarshal([]byte(cookieValue), &cookieContext)
		if err != nil {
			fmt.Println(err)
		} else {
			c = cookieContext
		}
	}

	return c
}

func showView(w http.ResponseWriter, r *http.Request, c Context) {
	jsonAccounts, err := json.Marshal(c)
	if err != nil {
		fmt.Println(err)
	} else {
		expiration := time.Now().Add(365 * 24 * time.Hour)
  	cookie := http.Cookie{Name: CookieName,Value:base64.StdEncoding.EncodeToString(jsonAccounts),Expires:expiration}
		http.SetCookie(w, &cookie)
	}

	t, _ := template.ParseFiles("main.html")
	t.Execute(w, c)
}
