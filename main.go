package main

import(

	"net/http"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"fmt"
	"os"
	"io/ioutil"
	
)

var(

	googleOauthConfig *oauth2.Config
)

var randomState = "random"

func init(){

	googleOauthConfig = &oauth2.Config{
		RedirectURL: "http://localhost:8080/callback",
		ClientID: os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("rITGu1_08Lps0M9LuD1ZNFHY"),
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint: google.Endpoint,
	}


}

func main(){

	http.HandleFunc("/",handleHome)
	http.HandleFunc("/login",handleLogin)
	http.HandleFunc("/callback",handleCallback)
	http.ListenAndServe(":8080", nil)


}

func handleHome(w http.ResponseWriter, r *http.Request){

	var html = `<html><body><a href="/login">Google Log in</a></body></html>`
	fmt.Fprint(w, html)
}

func handleLogin(w http.ResponseWriter, r *http.Request){

	url:= googleOauthConfig.AuthCodeURL(randomState)
	http.Redirect(w,r,url, http.StatusTemporaryRedirect)
	fmt.Println("saindo login")
	
}

func handleCallback(w http.ResponseWriter, r *http.Request){
	fmt.Println("callback")
	fmt.Printf(r.FormValue("state"))
	if r.FormValue("state") != randomState{
		fmt.Println("state is not valid")
		http.Redirect(w,r,"/", http.StatusTemporaryRedirect)
		return
	}

	token, err:=googleOauthConfig.Exchange(oauth2.NoContext, r.FormValue("code"))
		if err != nil{
			fmt.Printf("could not get token: %s\n", err.Error())
			http.Redirect(w,r,"/", http.StatusTemporaryRedirect)
			return
		}

	resp, err:= http.Get("https://www.googleapis.com/oauth2/v2/userinfo?acces_token=" + token.AccessToken)
		if err != nil{
			fmt.Printf("could not get request: %s\n", err.Error())
			http.Redirect(w,r,"/", http.StatusTemporaryRedirect)
			return
		}

	defer resp.Body.Close()
	content, err:= ioutil.ReadAll(resp.Body)
	if err != nil{
		fmt.Printf("could not parse: %s\n", err.Error())
		http.Redirect(w,r,"/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Fprintf(w, "Response: %s", content)

}