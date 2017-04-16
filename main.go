package WillBackend

import (
	_ "fmt"; _ "net/http"; _ "strings"; _ "log";
	_ "os";
	"net/http"
	"fmt"
)

func main() {
	http.HandleFunc("/", handleIt)
	http.ListenAndServe("localhost:8080", nil)
}

func handleIt(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()
	fmt.Println(r.Form)
}
