package WillBackend

import (
	_ "fmt"; _ "net/http"; _ "strings"; _ "log";
	_ "os";
	"net/http"
	"fmt"
)

func main() {
	http.HandleFunc("/", handleIt)
}

func handleIt(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()
	fmt.Println(r.Form)
}
