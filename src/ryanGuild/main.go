package main

import (
	_ "fmt"; _ "net/http"; _ "strings"; _ "log";
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
)
var (
	profArray = []profile{}
)


func main() {

	http.HandleFunc("/", serveStatic)
	http.HandleFunc("/src.htm", prepCards)
	http.HandleFunc("/card/", prepCard)
	appengine.Main()
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
