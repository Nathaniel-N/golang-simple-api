package main

import (
    "fmt"
    "strconv"
    "log"
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"
    "io/ioutil"
)

type Article struct {
    Id      string `json:"id"`
    Title   string `json:"title"`
    Desc    string `json:"desc"`
    Content string `json:"content"`
}

var Articles map[int]Article
var TotalArticle int

func homePage(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Welcome to the HomePage!")
    fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
    /* http handler is replaced by mux
    http.HandleFunc("/", homePage)
    http.HandleFunc("/articles", returnAllArticles)
    log.Fatal(http.ListenAndServe(":10000", nil))
    */
    router := mux.NewRouter().StrictSlash(true)
    router.Use(commonMiddleware)
    router.HandleFunc("/", homePage)
    router.HandleFunc("/articles", createNewArticle).Methods("POST")
    router.HandleFunc("/articles", returnAllArticles).Methods("GET")
    router.HandleFunc("/articles/{id}", returnSingleArticle).Methods("GET")
    router.HandleFunc("/articles/{id}", updateArticle).Methods("PUT")
    router.HandleFunc("/articles/{id}", deleteArticle).Methods("DELETE")
    log.Fatal(http.ListenAndServe(":10000", router))
}

func main() {
    fmt.Println("Rest API Running")
    Articles = make(map[int]Article)
    Articles[1] = Article{Id: "1", Title: "Hello", Desc: "Article Description", Content: "Article Content"}
    Articles[2] = Article{Id: "2", Title: "Hello 2", Desc: "Article 2 Description", Content: "Article 2 Content"}
    TotalArticle = 2
    handleRequests()
}

func commonMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Add("Content-Type", "application/json")
        next.ServeHTTP(w, r)
    })
}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
    reqBody, _ := ioutil.ReadAll(r.Body)
    var article Article 
    json.Unmarshal(reqBody, &article)
    
    TotalArticle += 1

    article.Id = strconv.Itoa(TotalArticle)
    Articles[TotalArticle] = article

    json.NewEncoder(w).Encode(article)
}

func returnAllArticles(w http.ResponseWriter, r *http.Request){
    fmt.Println("Endpoint Hit: returnAllArticles")
    json.NewEncoder(w).Encode(Articles)}

func returnSingleArticle(w http.ResponseWriter, r *http.Request){
    vars := mux.Vars(r)
    key := vars["id"]

    id, err := strconv.Atoi(key)
    if err == nil {
        article, exist := Articles[id]
        if exist {
            json.NewEncoder(w).Encode(article)
        } else {
            fmt.Fprintf(w, "Article does not exist")
        }
    } else {
        fmt.Fprintf(w, "Invalid Id")
    }


}

func updateArticle(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    key := vars["id"]

    reqBody, _ := ioutil.ReadAll(r.Body)
    var updatedArticle Article 
    json.Unmarshal(reqBody, &updatedArticle)
    // memastikan id tidak dapat diubah
    updatedArticle.Id = key
    
    id, err := strconv.Atoi(key)
    if err == nil {
        _, exist := Articles[id]
        if exist {
            Articles[id] = updatedArticle
            json.NewEncoder(w).Encode(updatedArticle)
        } else {
            fmt.Fprintf(w, "Article does not exist")
        }
    } else {
        fmt.Fprintf(w, "Invalid Id")
    }
    
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    key := vars["id"]

    id, err := strconv.Atoi(key)
    if err == nil {
        _, exist := Articles[id]
        if exist {
            delete(Articles, id)
            fmt.Fprintf(w, "Article deleted")
        } else {
            fmt.Fprintf(w, "Article does not exist")
        }
    } else {
        fmt.Fprintf(w, "Invalid Id")
    }
}