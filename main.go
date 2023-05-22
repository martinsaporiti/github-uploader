package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"text/template"
)

var (
	clientID     string
	clientSecret string
)

//go:embed files/static/*
var staticFiles embed.FS

func main() {

	clientID = os.Getenv("CLIENT_ID")
	if clientID == "" {
		log.Fatal("CLIENT_ID env variable is not set")
	}

	clientSecret = os.Getenv("CLIENT_SECRET")
	if clientSecret == "" {
		log.Fatal("CLIENT_SECRET env variable is not set")
	}

	var staticFS = fs.FS(staticFiles)
	static, err := fs.Sub(staticFS, "files/static")
	if err != nil {
		log.Fatal(err)
	}

	// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./files/static"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(static))))
	http.HandleFunc("/", index)
	http.HandleFunc("/get_access_token", getAccessToken)
	log.Printf("Listening on port %d", 9000)
	log.Fatal(http.ListenAndServe("0.0.0.0:9000", nil))

}

// index - render the index.html file
func index(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles(path.Join("./files/index.html"))
	t.Execute(w, "")
}

// getAccessToken - get access token from github, using the code from the frontend.
func getAccessToken(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	postURL := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", clientID, clientSecret, code)

	r, err := http.NewRequest("POST", postURL, nil)
	if err != nil {
		RespondJSON(w, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Accept", "application/json")
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		RespondJSON(w, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	defer res.Body.Close()

	fmt.Println("response Status:", res.Status)
	fmt.Println("response Status:", res.Body)

	type AccessToken struct {
		AccessToken string `json:"access_token"`
		Scope       string `json:"scope"`
		TokenType   string `json:"token_type"`
	}

	at := &AccessToken{}
	err = json.NewDecoder(res.Body).Decode(at)
	if err != nil {
		RespondJSON(w, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	RespondJSON(w, at, http.StatusOK)
}

// RespondJSON - translate an interface to json for response
func RespondJSON(w http.ResponseWriter, resp any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	retJSON, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
		message := fmt.Sprintf(`{"error marshaling response": "%s"}`, err.Error())
		retJSON = []byte(message)
	}
	_, _ = w.Write(retJSON)
}

// func upload() {
// 	TOKEN := ""
// 	ctx := context.Background()
// 	ts := oauth2.StaticTokenSource(
// 		&oauth2.Token{AccessToken: TOKEN},
// 	)

// 	tc := oauth2.NewClient(ctx, ts)
// 	client := github.NewClient(tc)
// 	contentResp, resp, err := client.Repositories.CreateFile(context.Background(), "martinsaporiti", "service_A", "schemas/abc3.txt", &github.RepositoryContentFileOptions{
// 		Message: github.String("my commit message"),
// 		Content: []byte("Hello World"),
// 		Author: &github.CommitAuthor{
// 			Name:  github.String("martinsaporiti"),
// 			Email: github.String("martinsaporiti@gmail.com"),
// 		},
// 	})

// 	if err != nil {
// 		panic(err)
// 	}

// 	defer resp.Body.Close()

// 	fmt.Println("response Status:", resp.Status)
// 	fmt.Println("Content Response:", contentResp.GetHTMLURL())
// }
