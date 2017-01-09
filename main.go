package main

import (
	"github.com/codahale/charlie"
	"log"
	"net/http"
	"time"
)

var fakeKeyParams = charlie.New([]byte("oursecrettokengoeshere"))
var cookieName = "testappAuth"

type user struct {
	Name     string
	Password string
	Token    string
}

// TODO: use cookie to hold username as well as token
// use a map[string]string to hold any desired user data
var testUser = user{
	Name:     "walle",
	Password: "abc123",
	Token:    "",
}

func isAuthorized(name, password string) bool {
	if name == testUser.Name && password == testUser.Password {
		return true
	}
	return false
}

// our cookie in this case will expire in 1 hour
func getExpiration() time.Time {
	// idTimeSuffix := time.Now().Format(time.RFC1123)
	return time.Now().Add(1 * time.Hour)
}

func servePublic(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving public.")
	http.ServeFile(w, r, "index.html")
}

func setSession(name string, response http.ResponseWriter) {
	testUser.Token = fakeKeyParams.Generate(name) // here we might send this to something like Redis or whatever
	cookie := http.Cookie{Name: cookieName, Value: testUser.Token, Expires: getExpiration(), Path: "/"}
	log.Println("Setting cookie.")
	http.SetCookie(response, &cookie)
}

func destroySession(name string, response http.ResponseWriter) {
	if testUser.Token == "" {
		return
	}
	uncookie := http.Cookie{Name: cookieName, Expires: getExpiration(), MaxAge: -1, Path: "/"}
	http.SetCookie(response, &uncookie)
}

func notAuthorized(w http.ResponseWriter, r *http.Request) {
	log.Println("not authorized Redirecting")
	http.Redirect(w, r, "/", http.StatusUnauthorized)
}

func serveSecret(w http.ResponseWriter, r *http.Request) {
	gotCookie, e := r.Cookie(cookieName) // where cookieName relates to our app's session, and value relates to user-specific info
	if e != nil {
		log.Println("No cookie or error getting it.")
		log.Println(e)
		notAuthorized(w, r)
		return
	}
	if gotCookie != nil && gotCookie.Value != "" {
		sessionToken := gotCookie.Value
		if err := fakeKeyParams.Validate(testUser.Name, sessionToken); err != nil {
			notAuthorized(w, r)
			return
		}
		log.Println("Serving secret.")
		http.ServeFile(w, r, "secret.html")
		return
	}
}

func handleAuthentication(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling auth.")

	postedUsername := r.FormValue("username")
	postedPassword := r.FormValue("password")
	log.Println("Got uname : ", postedUsername)
	log.Println("Got pword : ", postedPassword)

	redirectTarget := "/" // assume invalid

	if postedUsername != "" && postedPassword != "" {
		if isAuthorized(postedUsername, postedPassword) {
			setSession(postedUsername, w)
			redirectTarget = "/secure"
		}
	}
	http.Redirect(w, r, redirectTarget, 302)
}

func handleDeauthentication(w http.ResponseWriter, r *http.Request) {
	log.Println("Got deauth request.")
	redirectTarget := "/secure" // assume the worst
	cookie, err := r.Cookie(cookieName)
	if err != nil || cookie.Value == "" {
		return
	}

	formCSRFToken := r.FormValue("csrftoken")
	log.Println("logging out got csrf token from form: ", formCSRFToken)
	if formCSRFToken == cookie.Value {
		destroySession(testUser.Name, w)
		redirectTarget = "/"
	}
	http.Redirect(w, r, redirectTarget, 302)
}

func main() {

	http.HandleFunc("/", servePublic)
	http.HandleFunc("/login", handleAuthentication)
	http.HandleFunc("/logout", handleDeauthentication)
	http.HandleFunc("/secure", serveSecret)

	log.Fatal(http.ListenAndServe(":8081", nil))

}
