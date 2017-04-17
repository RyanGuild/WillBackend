package WillBackend

import (
	_ "fmt"; _ "net/http"; _ "strings"; _ "log";
	_ "os";
	"io"
	"log"
	"net/http"
)

func init() {
	http.HandleFunc("/", handleIt)
}

func handleIt(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()
	log.Println("hit")
	io.WriteString(w, string("<p>hello world</p>"))
}