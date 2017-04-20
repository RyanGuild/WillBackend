package main

import (
	_ "fmt"; _ "net/http"; _ "strings";
	_ "os";
	_ "io"
	"net/http"
	"fmt"
	"encoding/json"
	"io"
	"bytes"
	_ "encoding/xml"
	"time"
	"google.golang.org/appengine/blobstore"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"golang.org/x/net/context"
)
var profArray = []profile{}
var c context.Context
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
    <script type="text/javascript" src="../js/pics.js"></script>
</html>`

)

func init() {
	/*http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(blobstore.BlobKeyForFile())))
	http.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.Dir("/html"))))
	http.Handle("/resourses/", http.StripPrefix("/resourses/", http.FileServer(http.Dir("/resourses"))))
	http.Handle("/stylesheets/", http.StripPrefix("/stylesheets/", http.FileServer(http.Dir("/stylesheets"))))
	*/
	http.HandleFunc("/", serveStatic)
	http.HandleFunc("/cards.htm", prepHTML)
}
func serveStatic(w http.ResponseWriter, r *http.Request) {
	c = appengine.NewContext(r)
	key,_ := blobstore.BlobKeyForFile(c, r.RequestURI)
	blobstore.Send(w,key)
}

func readProfs() {
	var p profile
	i := 1
	for true{
		key, _ :=blobstore.BlobKeyForFile(c,"/profs/prof"+string(i)+".json")
		reader := blobstore.NewReader(c,key)
		err := json.Unmarshal(readstream(reader), &p)
		if err != nil{goto read}
		log.Infof(c,string(p),nil)
		profArray = append(profArray, p)
		i++
	}
	read:
}

func readstream(stream io.Reader) []byte {
	buf:= new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}




func prepHTML(w http.ResponseWriter, r *http.Request) {
	readProfs()
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
	i = 0
	for _, url:= range profArray[index].Pics{
		//fmt.Println(url)
		Payload += fmt.Sprintf(`<img class="slide" src="../resourses/prof/%s" id="%dimg%d"/>`,url,index,i)
		i++
	}

	Payload += fmt.Sprintf("<input type='button' id='cardNext%d'><label for='cardNext%d'> <div class='photoButton'>&gt;</div></label></div>",index,index)
	Payload += fmt.Sprint("</div>")
	Payload += fmt.Sprintf("<span class='bioTitle'>Bio:</span><div class='bioText'>%s</div></div></div>",profArray[index].Bio)
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
