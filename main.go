package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"math/rand"
	"net/http"
	"strconv"
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
	Raiting          string
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

type FixationsAndMiddleSpeed struct {
	Fixations   []Fixations
	MiddleSpeed int
	Queue       int
	CurrentKoef string
}

func Getkoefbyparams(w http.ResponseWriter, r *http.Request) {
	kad_nomer := r.FormValue("nomer")
	need_data := r.FormValue("need_data")
	vremya := r.FormValue("vremya")
	need_car := r.FormValue("need_car")
	need_speed := r.FormValue("need_speed")
	need_road := r.FormValue("need_road")

	db, err := sql.Open("mysql", "root:password@/mydiplom")
	if err != nil {
		panic(err)
	}
	res, err := db.Query(`select max_skorist,skorost,opisanie,vremya,marka,model,gosnoner,k.nazvanie,cel_ispolzovaniya from fixacii
left join kameri k on k.id = fixacii.id_kameri
left join object o on o.id = k.id_objecta
left join kadastrivi_nomer kn on kn.id = o.id_kadastrovi
left join avtomobili a on a.id = fixacii.id_avtomobil
where kn.nomer = ? and
      vremya > date_sub(?,interval 10 minute ) and
      vremya < date_add(?,interval 10 minute )`,
		kad_nomer, vremya, vremya)
	if err != nil {
		println("обоже")
	}
	defer db.Close()

	var model FixationsAndMiddleSpeed
	model.Fixations = make([]Fixations, 0)

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
		model.Fixations = append(model.Fixations, *bk)
	}

	sum := 0.0
	if need_data == "on" {
		sum += 1
	}
	if need_car == "on" {
		sum += 1
	}
	if need_road == "on" {
		sum += 1
	}
	if need_speed == "on" {
		sum += 1
	}
	model.CurrentKoef = strconv.FormatFloat(rand.Float64()+sum, 'f', 6, 64)

	tmpl, _ := template.ParseFiles("koefbyparams.html")
	err = tmpl.Execute(w, model)
	if err != nil {
		return
	}
}

func GetkoefbyparamsForm(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:password@/mydiplom")
	if err != nil {
		panic(err)
	}
	res, err := db.Query(`select nomer_uchastka,nomer,nazvanie_uchatka,gorod,ylica,nasel_punkt,poselenie,municipal_okrug,oblast from dorogi
join kadastrivi_nomer kn on dorogi.id_kadastr_nom = kn.id
join object o on kn.id = o.id_kadastrovi;`)
	if err != nil {
		//panic(err)
		println("баги лох")
	}

	roads := make([]*Roads, 0)
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
		roads = append(roads, bk)
	}

	tmpl, _ := template.ParseFiles("getkoefbyparams.html")
	err = tmpl.Execute(w, roads)
	if err != nil {
		return
	}
}

func Getlowraitingroads(w http.ResponseWriter, r *http.Request) {
	rajon := r.FormValue("rajon")

	db, err := sql.Open("mysql", "root:password@/mydiplom")
	if err != nil {
		panic(err)
	}
	res, err := db.Query(`select nomer_uchastka,nomer,nazvanie_uchatka,gorod,ylica,nasel_punkt,poselenie,municipal_okrug,oblast from dorogi
join kadastrivi_nomer kn on dorogi.id_kadastr_nom = kn.id
join object o on kn.id = o.id_kadastrovi
where municipal_okrug = ?;`,
		rajon)
	if err != nil {
		println("обоже")
	}
	defer db.Close()

	roads := make([]*Roads, 0)
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
		raiting := rand.Float64()
		if raiting > 0.5 {
			bk.Raiting = strconv.FormatFloat(raiting, 'f', 6, 64)
			roads = append(roads, bk)
		}
	}

	tmpl, _ := template.ParseFiles("lowraitingroads.html")
	err = tmpl.Execute(w, roads)
	if err != nil {
		return
	}
}

func GetlowraitingroadsForm(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:password@/mydiplom")
	if err != nil {
		panic(err)
	}
	res, err := db.Query(`select municipal_okrug from kadastrivi_nomer group by municipal_okrug;`)
	if err != nil {
		//panic(err)
		println("баги лох")
	}

	roads := make([]*Roads, 0)
	for res.Next() {
		bk := new(Roads)
		res.Scan(
			&bk.Municipal_okrug,
		)
		roads = append(roads, bk)
	}

	tmpl, _ := template.ParseFiles("getlowraitingroads.html")
	err = tmpl.Execute(w, roads)
	if err != nil {
		return
	}
}

func Getcurrentkoef(w http.ResponseWriter, r *http.Request) {
	kad_nomer := r.FormValue("nomer")

	db, err := sql.Open("mysql", "root:password@/mydiplom")
	if err != nil {
		panic(err)
	}
	res, err := db.Query(`select max_skorist,skorost,opisanie,vremya,marka,model,gosnoner,k.nazvanie,cel_ispolzovaniya from fixacii
left join kameri k on k.id = fixacii.id_kameri
left join object o on o.id = k.id_objecta
left join kadastrivi_nomer kn on kn.id = o.id_kadastrovi
left join avtomobili a on a.id = fixacii.id_avtomobil
where kn.nomer = ? and
      vremya > date_sub(current_timestamp,interval 2 minute )`,
		kad_nomer)
	if err != nil {
		println("обоже")
	}
	defer db.Close()

	var model FixationsAndMiddleSpeed
	model.Fixations = make([]Fixations, 0)

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
		model.Fixations = append(model.Fixations, *bk)
	}
	model.CurrentKoef = strconv.FormatFloat(float64(rand.Float32()), 'f', 6, 64)

	tmpl, _ := template.ParseFiles("currentkoef.html")
	err = tmpl.Execute(w, model)
	if err != nil {
		return
	}
}

func GetcurrentkoefForm(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:password@/mydiplom")
	if err != nil {
		panic(err)
	}
	res, err := db.Query(`select nomer_uchastka,nomer,nazvanie_uchatka,gorod,ylica,nasel_punkt,poselenie,municipal_okrug,oblast from dorogi
join kadastrivi_nomer kn on dorogi.id_kadastr_nom = kn.id
join object o on kn.id = o.id_kadastrovi;`)
	if err != nil {
		//panic(err)
		println("баги лох")
	}

	roads := make([]*Roads, 0)
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
		roads = append(roads, bk)
	}

	tmpl, _ := template.ParseFiles("getcurrentkoef.html")
	err = tmpl.Execute(w, roads)
	if err != nil {
		return
	}
}

func GetQueue(w http.ResponseWriter, r *http.Request) {
	kad_nomer := r.FormValue("nomer")

	db, err := sql.Open("mysql", "root:password@/mydiplom")
	if err != nil {
		panic(err)
	}
	res, err := db.Query(`select max_skorist,skorost,opisanie,vremya,marka,model,gosnoner,k.nazvanie,cel_ispolzovaniya from fixacii
left join kameri k on k.id = fixacii.id_kameri
left join object o on o.id = k.id_objecta
left join kadastrivi_nomer kn on kn.id = o.id_kadastrovi
left join avtomobili a on a.id = fixacii.id_avtomobil
where kn.nomer = ? and
      vremya > date_sub(current_timestamp,interval 2 minute )`,
		kad_nomer)
	if err != nil {
		println("обоже")
	}
	defer db.Close()

	var model FixationsAndMiddleSpeed
	model.Fixations = make([]Fixations, 0)

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
		model.Fixations = append(model.Fixations, *bk)
	}

	model.Queue = len(model.Fixations)

	tmpl, _ := template.ParseFiles("queue.html")
	err = tmpl.Execute(w, model)
	if err != nil {
		return
	}
}

func GetQueueForm(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:password@/mydiplom")
	if err != nil {
		panic(err)
	}
	res, err := db.Query(`select nomer_uchastka,nomer,nazvanie_uchatka,gorod,ylica,nasel_punkt,poselenie,municipal_okrug,oblast from dorogi
join kadastrivi_nomer kn on dorogi.id_kadastr_nom = kn.id
join object o on kn.id = o.id_kadastrovi;`)
	if err != nil {
		//panic(err)
		println("баги лох")
	}

	roads := make([]*Roads, 0)
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
		roads = append(roads, bk)
	}

	tmpl, _ := template.ParseFiles("getqueue.html")
	err = tmpl.Execute(w, roads)
	if err != nil {
		return
	}
}

func GetmiddleSpeed(w http.ResponseWriter, r *http.Request) {
	kad_nomer := r.FormValue("nomer")

	db, err := sql.Open("mysql", "root:password@/mydiplom")
	if err != nil {
		panic(err)
	}
	res, err := db.Query(`select max_skorist,skorost,opisanie,vremya,marka,model,gosnoner,k.nazvanie,cel_ispolzovaniya from fixacii
left join kameri k on k.id = fixacii.id_kameri
left join object o on o.id = k.id_objecta
left join kadastrivi_nomer kn on kn.id = o.id_kadastrovi
left join avtomobili a on a.id = fixacii.id_avtomobil
where kn.nomer = ? and
      vremya > date_sub(current_timestamp,interval 10 minute )`,
		kad_nomer)
	if err != nil {
		println("обоже")
	}
	defer db.Close()

	var model FixationsAndMiddleSpeed
	model.Fixations = make([]Fixations, 0)

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
		model.Fixations = append(model.Fixations, *bk)
	}
	sum := 0
	for _, r := range model.Fixations {
		sum += r.Skorost
	}
	if len(model.Fixations) != 0 {
		model.MiddleSpeed = sum / len(model.Fixations)
	}

	tmpl, _ := template.ParseFiles("middlespeed.html")
	err = tmpl.Execute(w, model)
	if err != nil {
		return
	}
}

func GetmiddleSpeedForm(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:password@/mydiplom")
	if err != nil {
		panic(err)
	}
	res, err := db.Query(`select nomer_uchastka,nomer,nazvanie_uchatka,gorod,ylica,nasel_punkt,poselenie,municipal_okrug,oblast from dorogi
join kadastrivi_nomer kn on dorogi.id_kadastr_nom = kn.id
join object o on kn.id = o.id_kadastrovi;`)
	if err != nil {
		//panic(err)
		println("баги лох")
	}

	roads := make([]*Roads, 0)
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
		roads = append(roads, bk)
	}

	tmpl, _ := template.ParseFiles("getmiddlespeed.html")
	err = tmpl.Execute(w, roads)
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

func SetRoads(w http.ResponseWriter, r *http.Request) {
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

func SetRoadsForm(w http.ResponseWriter, r *http.Request) {
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

func viewRoads(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:password@/mydiplom")
	if err != nil {
		panic(err)
	}
	res, err := db.Query(`select nomer_uchastka,nomer,nazvanie_uchatka,gorod,ylica,nasel_punkt,poselenie,municipal_okrug,oblast from dorogi
join kadastrivi_nomer kn on dorogi.id_kadastr_nom = kn.id
join object o on kn.id = o.id_kadastrovi`)
	if err != nil {
		//panic(err)
		println("error")
	}

	defer db.Close()

	roads := make([]*Roads, 0)
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
		roads = append(roads, bk)
	}
	tmpl, _ := template.ParseFiles("roads.html")
	err = tmpl.Execute(w, roads)
	if err != nil {
		return
	}
}

func main() {

	css := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css/", css))

	http.HandleFunc("/getkoefbyparams/get", Getkoefbyparams)
	http.HandleFunc("/getkoefbyparams/", GetkoefbyparamsForm)

	http.HandleFunc("/getlowraitingroads/get", Getlowraitingroads)
	http.HandleFunc("/getlowraitingroads/", GetlowraitingroadsForm)

	http.HandleFunc("/getcurrentkoef/get", Getcurrentkoef)
	http.HandleFunc("/getcurrentkoef/", GetcurrentkoefForm)

	http.HandleFunc("/getqueue/get", GetQueue)
	http.HandleFunc("/getqueue/", GetQueueForm)

	http.HandleFunc("/getmiddlespeed/get", GetmiddleSpeed)
	http.HandleFunc("/getmiddlespeed/", GetmiddleSpeedForm)

	http.HandleFunc("/setfixations/add", Setfixation)
	http.HandleFunc("/setfixation/", SetfixationForm)
	http.HandleFunc("/fixations/", viewFixations)

	http.HandleFunc("/setroad/add", SetRoads)
	http.HandleFunc("/setroad/", SetRoadsForm)
	http.HandleFunc("/roads/", viewRoads)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "main.html")
	})
	log.Fatal(http.ListenAndServe(":3008", nil))
}
