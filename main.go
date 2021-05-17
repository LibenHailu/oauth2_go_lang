package main

import (
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// OAuth2 Config
var (
	conf = &oauth2.Config{
		ClientID:     "856828906157-5l5tcd0h8ebjfn98dktjeseauuhb1b6j.apps.googleusercontent.com",
		ClientSecret: "nUbaCmxc_CPZhc75j-d6kMID",
		RedirectURL:  "http://localhost:8080/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/bigquery",
			"https://www.googleapis.com/auth/blogger",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
	// TODO: make the string random this will be used as csrf token
	randomState = "random"
)

func main() {

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/callback", handleCallback)
	http.ListenAndServe(":8080", nil)

}

func handleHome(w http.ResponseWriter, r *http.Request) {
	var html = `<html><body><a href="/login">Google Login</body></html>`
	fmt.Fprint(w, html)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	// Redirect user to Google's consent page to ask for permission
	// this will send our randomState/our csrf token
	// autho injects some code as an interface to login
	url := conf.AuthCodeURL(randomState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)

}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	// makes sure it is the right client
	if r.FormValue("state") != randomState {
		fmt.Println("state is not valid")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	token, err := conf.Exchange(oauth2.NoContext, r.FormValue("code"))
	fmt.Println(token)
	if err != nil {
		fmt.Printf("could not get token: %s/n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//  sends http.client with some configurations preconfigured
	client := conf.Client(oauth2.NoContext, token)

	resp, err := client.Get("https://google.com")
	http.Redirect(w, r, "https://google.com", http.StatusTemporaryRedirect)
	fmt.Printf("Response %s", resp.Status)

}
