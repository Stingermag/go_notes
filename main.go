package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"text/template"
)

type Roads struct {
	Nomer_uchastka   int
	Nomer            string
	Nazvanie_uchatka string
	Gorod            string
	Ylica            string
	Nasel_punkt      string
	Poselenie        string
	Municipal_okrug  string
	Oblast           string
}

type Fixations struct {
	Max_skorist       int
	Skorost           int
	Opisanie          string
	Vremya            string
	Marka             string
	Model             string
	Gosnoner          string
	Nazvanie          string
	Cel_ispolzovaniya string
}

type Avto struct {
	Id       int
	Gosnoner string
}
type Kameri struct {
	Id              int
	Nazvanie        string
	Oblast          string
	Gorod           string
	Municipal_okrug string
	Ylica           string
	Zem_ychastok    string
}
type AvtoKameri struct {
	Avto   []Avto
	Kameri []Kameri
}

func GetmiddleSpeed(w http.ResponseWriter, r *http.Request) {
	vremya := r.FormValue("vremya")
	opisanie := r.FormValue("opisanie")
	skorost := r.FormValue("skorost")
	max_skorist := r.FormValue("max_skorist")
	id_avtomobil := r.FormValue("id_avtomobil")
	id_kameri := r.FormValue("id_kameri")

	db, err := sql.Open("mysql", "root:password@/mydiplom")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`insert into fixacii  (vremya, opisanie, skorost,max_skorist,id_avtomobil, id_kameri)
values (?,?,?,?,?,?);`,
		vremya, opisanie, skorost, max_skorist, id_avtomobil, id_kameri)
	if err != nil {
		//panic(err)
		println("обоже")
	}
	defer db.Close()
	http.Redirect(w, r, "/fixations", http.StatusFound)
}

func GetmiddleSpeedForm(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:password@/mydiplom")
	if err != nil {
		panic(err)
	}
	res, err := db.Query(`select id,gosnoner from avtomobili;`)
	if err != nil {
		//panic(err)
		println("баги лох")
	}
	res2, err := db.Query(`select kameri.id,kameri.nazvanie,oblast,gorod,municipal_okrug,ylica,zem_ychastok from kameri
left join object o on o.id = kameri.id_objecta
left join kadastrivi_nomer kn on kn.id = o.id_kadastrovi;`)
	if err != nil {
		//panic(err)
		println("баги лох")
	}

	defer db.Close()
	var model AvtoKameri
	model.Avto = make([]Avto, 0)
	model.Kameri = make([]Kameri, 0)
	for res.Next() {
		bk := new(Avto)
		res.Scan(&bk.Id, &bk.Gosnoner)
		model.Avto = append(model.Avto, *bk)
	}
	for res2.Next() {
		bk := new(Kameri)
		res2.Scan(&bk.Id, &bk.Nazvanie, &bk.Oblast, &bk.Gorod, &bk.Municipal_okrug, &bk.Ylica, &bk.Zem_ychastok)
		model.Kameri = append(model.Kameri, *bk)
	}

	tmpl, _ := template.ParseFiles("setfixations.html")
	err = tmpl.Execute(w, model)
	if err != nil {
		return
	}
}

func Setfixation(w http.ResponseWriter, r *http.Request) {
	vremya := r.FormValue("vremya")
	opisanie := r.FormValue("opisanie")
	skorost := r.FormValue("skorost")
	max_skorist := r.FormValue("max_skorist")
	id_avtomobil := r.FormValue("id_avtomobil")
	id_kameri := r.FormValue("id_kameri")

	db, err := sql.Open("mysql", "root:password@/mydiplom")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`insert into fixacii  (vremya, opisanie, skorost,max_skorist,id_avtomobil, id_kameri)
values (?,?,?,?,?,?);`,
		vremya, opisanie, skorost, max_skorist, id_avtomobil, id_kameri)
	if err != nil {
		//panic(err)
		println("обоже")
	}
	defer db.Close()
	http.Redirect(w, r, "/fixations", http.StatusFound)
}

func SetfixationForm(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:password@/mydiplom")
	if err != nil {
		panic(err)
	}
	res, err := db.Query(`select id,gosnoner from avtomobili;`)
	if err != nil {
		//panic(err)
		println("баги лох")
	}
	res2, err := db.Query(`select kameri.id,kameri.nazvanie,oblast,gorod,municipal_okrug,ylica,zem_ychastok from kameri
left join object o on o.id = kameri.id_objecta
left join kadastrivi_nomer kn on kn.id = o.id_kadastrovi;`)
	if err != nil {
		//panic(err)
		println("баги лох")
	}

	defer db.Close()
	var model AvtoKameri
	model.Avto = make([]Avto, 0)
	model.Kameri = make([]Kameri, 0)
	for res.Next() {
		bk := new(Avto)
		res.Scan(&bk.Id, &bk.Gosnoner)
		model.Avto = append(model.Avto, *bk)
	}
	for res2.Next() {
		bk := new(Kameri)
		res2.Scan(&bk.Id, &bk.Nazvanie, &bk.Oblast, &bk.Gorod, &bk.Municipal_okrug, &bk.Ylica, &bk.Zem_ychastok)
		model.Kameri = append(model.Kameri, *bk)
	}

	tmpl, _ := template.ParseFiles("setfixations.html")
	err = tmpl.Execute(w, model)
	if err != nil {
		return
	}
}

func viewFixations(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:password@/mydiplom")
	if err != nil {
		panic(err)
	}
	res, err := db.Query(`select max_skorist,skorost,opisanie,vremya,marka,model,gosnoner,nazvanie,cel_ispolzovaniya from mydiplom.fixacii
left join avtomobili a on a.id = fixacii.id_avtomobil
left join kameri k on fixacii.id_kameri = k.id order by vremya desc `)
	if err != nil {
		//panic(err)
		println("error")
	}

	defer db.Close()

	marshruti := make([]*Fixations, 0)
	for res.Next() {
		bk := new(Fixations)
		res.Scan(
			&bk.Max_skorist,
			&bk.Skorost,
			&bk.Opisanie,
			&bk.Vremya,
			&bk.Marka,
			&bk.Model,
			&bk.Gosnoner,
			&bk.Nazvanie,
			&bk.Cel_ispolzovaniya,
		)
		marshruti = append(marshruti, bk)
	}
	tmpl, _ := template.ParseFiles("fixations.html")
	err = tmpl.Execute(w, marshruti)
	if err != nil {
		return
	}
}

func viewRoads(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:password@/mydiplom")
	if err != nil {
		panic(err)
	}
	res, err := db.Query(`select nomer_uchastka,nomer,nazvanie_uchatka,gorod,ylica,nasel_punkt,poselenie,municipal_okrug,oblast from dorogi_haracter
join dorogi d on dorogi_haracter.id_dorogi = d.id
join kadastrivi_nomer kn on kn.id = d.id_kadastr_nom`)
	if err != nil {
		//panic(err)
		println("error")
	}

	defer db.Close()

	marshruti := make([]*Roads, 0)
	for res.Next() {
		bk := new(Roads)
		res.Scan(
			&bk.Nomer_uchastka,
			&bk.Nomer,
			&bk.Nazvanie_uchatka,
			&bk.Gorod,
			&bk.Ylica,
			&bk.Nasel_punkt,
			&bk.Poselenie,
			&bk.Municipal_okrug,
			&bk.Oblast,
		)
		marshruti = append(marshruti, bk)
	}
	tmpl, _ := template.ParseFiles("roads.html")
	err = tmpl.Execute(w, marshruti)
	if err != nil {
		return
	}
}

func main() {

	css := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css/", css))
	//http.HandleFunc("/view/", viewAvtotransport)
	//http.HandleFunc("/viewr/", viewRemonti)
	//http.HandleFunc("/add/", do(loading))
	http.HandleFunc("/getmiddlespeed/add", GetmiddleSpeed)
	http.HandleFunc("/getmiddlespeed/", GetmiddleSpeedForm)
	http.HandleFunc("/setfixations/add", Setfixation)
	http.HandleFunc("/setfixation/", SetfixationForm)
	http.HandleFunc("/fixations/", viewFixations)
	http.HandleFunc("/roads/", viewRoads)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "main.html")
	})
	log.Fatal(http.ListenAndServe(":3008", nil))
}
