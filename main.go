package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type Page struct {
	Title string
	Body  []byte
}

type Nts struct {
	ID  string
	Ttl  string
	Mess string
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func selectData(){
	db, err := sql.Open("mysql", "root:password@/notes")
	if err != nil {
		panic(err)
	}
	res, err := db.Query("select * from notes")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	nnote := make([]*Nts, 0)
	for res.Next() {
		bk := new(Nts)
		res.Scan(&bk.ID, &bk.Ttl, &bk.Mess)
		nnote = append(nnote, bk)
	}
	http.HandleFunc("/view/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.ParseFiles("view.html")
		tmpl.Execute(w, nnote)
	})
}
func loading(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "add", p)
}

func saveArticle(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	titl := r.FormValue("tname")
	p := &Page{Title: titl, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	db, err := sql.Open("mysql", "root:password@/notes")
	if err != nil {
		panic(err)
	}
	result, err := db.Exec("insert into notes.notes (title, mess) values (?, ?)",
		titl, body)
	if err != nil{
		panic(err)
	}
	defer db.Close()
	fmt.Println(result.LastInsertId())
	http.Redirect(w, r, "/view/notes", http.StatusFound)
}


func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	var templates = template.Must(template.ParseFiles("add.html", "view.html", "main.html"))

	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}



func do(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	var validPath = regexp.MustCompile("^/(add|save|view|main)/([a-zA-Z0-9]+)$")
	return func(w http.ResponseWriter, r *http.Request) {
		path := validPath.FindStringSubmatch(r.URL.Path)
		if path == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, path[2])
	}
}

func main() {
	selectData()
	css := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css/", css))
	http.HandleFunc("/add/", do(loading))
	http.HandleFunc("/save/", do(saveArticle))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "main.html")
	})
	log.Fatal(http.ListenAndServe(":3000", nil))
}