package main

import (
	_ "fmt"; _ "net/http"; _ "strings"; _ "log";
	_ "os";
	_ "io"
	"log"
	"net/http"
	"fmt"
	_ "google.golang.org/appengine"
	"encoding/json"
	"os"
	"flag"
	"path/filepath"
	"io"
	"bytes"
	_ "encoding/xml"
	"time"
)
var profArray = []profile{}
const profLocation = "profs"

func main() {
	flag.Parse()
	switch flag.Arg(2) {
	/*case "-t":
		c := make(chan []byte)
		genCard(0, c)
		readstream()*/

	default:
	//testprof := &profile{"test2","testbio", map[string]float32{"shirt":10.50,"tie":6.00}, []string{"this","them","theOther"}}
	//writeJson("profs/test2.json", testprof)
		readProfs(profLocation)
		http.HandleFunc("/hit", handleIt)
		http.HandleFunc("/getProf", prepHTML)
		http.ListenAndServe(":8080",nil)
	}
}

func handleIt(w http.ResponseWriter, r *http.Request)  {
	log.Println("hit")
	fmt.Fprintln(w, string("<p>hello world</p>"))
}


func readProfs(filename string){
	dir, _ := filepath.Abs(filepath.Dir(flag.Arg(0)))
	fmt.Println(dir)
	err := filepath.Walk(dir+`\`+filename, readJsonProf)
	if err != nil {
		fmt.Printf("err: %v",err)
	}

}

func readJsonProf(path string, info os.FileInfo, err error) error {
	fmt.Println("profile search visiting:     " + path)
	var p profile
	file, _ := os.Open(path)
	if path[len(path)-5:] == ".json" {
		ret := json.Unmarshal(readstream(file), &p)
		if ret != nil {
			return ret
		}
		profArray = append(profArray, p)
	}
	return nil
}

func readstream(stream io.Reader) []byte {
	buf:= new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}




func prepHTML(w http.ResponseWriter, r *http.Request){
	var request map[string] int
	var stringSum string
	c := make(chan string)
	defer close(c)
	for k, _ := range r.URL.Query(){
		fmt.Println(k)
		err := json.Unmarshal([]byte(k), &request)
		if err != nil {
			fmt.Println("error in json convert: ",err)
		}
		fmt.Println(request)

	}
	index := int(request["index"])
	count := int(request["count"])
	top := index + count
		for i := index; i <= top; i++{
			go genCard(i, c)
		}
	for x := index; x <= top; x++ {
		timeout := time.Second *1
		select {
		case resp := <-c:
			fmt.Println("card returned: ",x)
			stringSum += resp
		case <-time.After(timeout):
			fmt.Println("timeout reached: ", timeout)
			goto fin
		}
	}
	fin:w.Header().Add("Content-Type","test/html")
	w.Header().Add("charset", "uft-8")
	io.WriteString(w, stringSum)


}

func genCard(index int, c chan string) {
	if index >= len(profArray) {return }
	var retString = ""
	retString += fmt.Sprintf("<div class='profContainer' id='card%d'><div><span class='profName'>%s</span><div class='profRow'><div class='profItemContainer'><span class='profItemTitle'>Items</span>", index, profArray[index].Name)
	var i = 0
	for k, v := range profArray[index].Items {
		retString += fmt.Sprintf("<span class='itemEntry'><input type='checkbox' id='card%ditem%d' /> <label for='card%ditem%d'><span></span>%s  $%v</label></span>",index, i,index, i, k, v)
		i++
	}
	retString += fmt.Sprint("</div>")
	retString += fmt.Sprintf("<div class='profPhotoContainer' id='card%dPhoto'> <input type='button' id='card%dPrev'><label for='card%dPrev'><div class='photoButton'><span>&lt;</span></div></label>",index,index,index)
	retString += fmt.Sprintf("<input type='button' id='card%dNext'><label for='card%dNext'> <div class='photoButton'>&gt;</div></label></div>",index,index)
	retString += fmt.Sprint("</div>")
	retString += fmt.Sprintf("<span class='bioTitle'>Bio:</span><div class='bioText'>%s</div></div></div>",profArray[index].Bio)
	c<-retString
}

func writeJsonResponse(w io.Writer, a *response) {
	cont, _ := json.Marshal(a)
	w.Write(cont)

}

type profile struct{
	Name string
	Bio string
	Items map[string] float32
	Pics []string
}

type response struct {
	Content string
	Pics [][]string
}


