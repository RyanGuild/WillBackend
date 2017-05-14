package main

import (
	_ "fmt"; _ "golang.org/x/net/http2"; _ "strings"; _ "log";
	_ "os";
	_ "io"
	"net/http"
	_ "google.golang.org/appengine"
	"google.golang.org/appengine/blobstore"
	"google.golang.org/appengine"
	"io"
	"bytes"
	_ "encoding/xml"
	"fmt"
	"strconv"
	"time"
)

var (
	profArray = []profile{}
)

const (
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

func main() {
	http.HandleFunc("/", serveStatic)
	http.HandleFunc("/cards.htm", prepCards)
	http.HandleFunc("/card/", prepCard)
	http.ListenAndServe(":8080", nil)
}



func serveStatic(w http.ResponseWriter, r *http.Request) {
	updateProfiles(r)
	ctx := appengine.NewContext(r)
	key,_ := blobstore.BlobKeyForFile(ctx, r.RequestURI)
	blobstore.Send(w,key)
}


func updateProfiles(r *http.Request) {
	ctx := appengine.NewContext(r)
	key, _ := blobstore.BlobKeyForFile(ctx, "profs/")
	reader := blobstore.NewReader(ctx,key)
	fmt.Println(readstream(reader))
}

func readstream(stream io.Reader) []byte {
	buf:= new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}

func  prepCard(w http.ResponseWriter, r *http.Request) {
	uri := r.RequestURI
	uri = uri[6:]
	index, err := strconv.ParseInt(uri,10,64)
	if err != nil {http.NotFound(w,r); return }
	ch := make(chan string)
	defer close(ch)
	timeout := time.After(time.Second)
	go genCard(int(index), ch)
	select {
	case data := <-ch:
		w.Header().Set("Content-Type", "text/html; charset=uft-8")
		io.WriteString(w,cardHead+data+cardBase)
	case <-timeout:
		http.NotFound(w,r)
	}
}


func prepCards(w http.ResponseWriter, r *http.Request) {
	c := make(chan string)
	var page string
	for k := range profArray {
		go genCard(k,c)
	}
	timeout := time.After(time.Second*1)
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
