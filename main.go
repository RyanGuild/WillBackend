package main

import (
	_ "fmt"; _ "net/http"; _ "strings"; _ "log";
	_ "os";
	_ "io"
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
const (
	profLocation = "profs"
	cardHead = `
<!DOCTYPE HTML5>
<html>
    <head>
        <link href="https://fonts.googleapis.com/css?family=Josefin+Slab" rel="stylesheet">
        <link type="text/css" rel="stylesheet" href="../stylesheets/main.css"/>
        <script type="text/javascript" src="../js/util/jquery-3.1.1.js"></script>
        <script type="text/javascript">
            parent.childElementById = function (id) {return document.getElementById(id);}
        </script>
    </head>
    <body id="contentBody">`
	cardBase = `</body>
    <script type="text/javascript" src="../js/content.js"></script>
</html>`

)

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
		http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
		http.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.Dir("html"))))
		http.Handle("/resourses/", http.StripPrefix("/resourses/", http.FileServer(http.Dir("resourses"))))
		http.Handle("/stylesheets/", http.StripPrefix("/stylesheets/", http.FileServer(http.Dir("stylesheets"))))
		http.HandleFunc("/cards.htm", prepHTML)
		http.ListenAndServe(":8080",nil)
	}
}

/*func serveMain(w http.ResponseWriter, r *http.Request)  {
	path, _ := os.Getwd()
	path += "/html/main.html"
	fmt.Println(path)
	file, _:= os.Open(path)

}*/


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




/*func prepHTML(w http.ResponseWriter, r *http.Request){
	var request map[string] int
	var Cards []*card
	c := make(chan *card)
	defer close(c)
	for k, _ := range r.URL.Query(){
		err := json.Unmarshal([]byte(k), &request)
		if err != nil {
			fmt.Println("error in json convert: ", err)
		}
	}
	index := int(request["index"])
	count := int(request["count"])
	top := index + count
	if top > len(profArray) {
		top = len(profArray)
	}
	for i := index; i <= top; i++{
		fmt.Println("creating card: ", i)
		go genCard(i, c)
	}
	for x := index; x <= top; x++ {
		timeout := time.Second *1
		select {
		case resp := <-c:
			Cards = append(Cards, resp)
		case <-time.After(timeout):
			fmt.Println("timeout reached: ", timeout)
			goto fin
		}
	}
	fin:w.Header().Add("Content-Type","test/html")
	w.Header().Add("charset", "uft-8")
	resp := response{Cards}
	jcoder := json.NewEncoder(w).Encode(resp)
	if jcoder != nil {
		fmt.Println(jcoder)
	}
}*/

func prepHTML(w http.ResponseWriter, r *http.Request) {
	c := make(chan string)
	var page string
	for k, _ := range profArray {
		go genCard(k,c)
	}
	timeout :=time.After(time.Second*1)
	for range profArray {
		select {
		case ret := <-c:
			page += ret

		case <-timeout:
			goto fin
		}
	}
	fin:
	w.Header().Add("Content-Type","text/html; charset=uft-8")
	io.WriteString(w,cardHead+page+cardBase)

}

func genCard(index int, c chan string) {
	if index >= len(profArray) {return }
	var Payload = ""
	Payload += fmt.Sprintf("<div class='profContainer' id='card%d'><div><span class='profName'>%s</span><div class='profRow'><div class='profItemContainer'><span class='profItemTitle'>Items</span>", index, profArray[index].Name)
	var i = 0
	for k, v := range profArray[index].Items {
		Payload += fmt.Sprintf("<span class='itemEntry'><input type='checkbox' id='carditem%d%d' /> <label for='carditem%d%d'><span></span>%s  $%v</label></span>",index, i,index, i, k, v)
		i++
	}
	Payload += fmt.Sprint("</div>")
	Payload += fmt.Sprintf("<div class='profPhotoContainer' id='cardPhoto%d'> <input type='button' id='cardPrev%d'><label for='cardPrev%d'><div class='photoButton'><span>&lt;</span></div></label>",index,index,index)
	Payload += fmt.Sprintf("<input type='button' id='cardNext%d'><label for='cardNext%d'> <div class='photoButton'>&gt;</div></label></div>",index,index)
	Payload += fmt.Sprint("</div>")
	Payload += fmt.Sprintf("<span class='bioTitle'>Bio:</span><div class='bioText'>%s</div></div></div>",profArray[index].Bio)
	fmt.Println("card",index,"served")
	c<- Payload
}


type profile struct{
	Name string
	Bio string
	Items map[string] float32
	Pics []string
}


type picContainer struct {
	Index int
	PhotoBucket []string
}
