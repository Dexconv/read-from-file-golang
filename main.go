package main

import (
	"crypto/sha1"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var tpl *template.Template

func init(){
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}
func main(){
	http.HandleFunc("/", index)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080",nil)
}
func index(w http.ResponseWriter, r *http.Request){
	var s string
	if r.Method == http.MethodPost{

		f, h, err := r.FormFile("q")

		if err != nil{
			http.Error(w, http.StatusText(405), 405)
		}

		defer f.Close()

		fmt.Println("\nfile:", f, "\nheader:", h, "\nerr", err)

		bs, err := ioutil.ReadAll(f)
		if err != nil{
			http.Error(w, http.StatusText(500), 500)
			log.Fatalln(err)
		}
		s = string(bs)

		ext := strings.Split(h.Filename, ".")[1]
		hash := sha1.New()
		io.Copy(hash, f)
		fname := fmt.Sprintf("%x", hash.Sum(nil))+"."+ext
		dst, err := os.Create(filepath.Join("./storage/",fname))
		if err != nil{
			http.Error(w, http.StatusText(500), 500)
			log.Fatalln(err)
		}
		defer dst.Close()
		_, err = dst.Write(bs)
		if err != nil{
			http.Error(w, http.StatusText(500), 500)
			log.Fatalln(err)
		}
	}
	tpl.ExecuteTemplate(w,"index.gohtml", s)
}