package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type Employee struct {
	Id           int
	Name         string
	Salary       int
	Destignation string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "employee"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		fmt.Println("db error", err)
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM Employee ORDER BY id ASC")
	if err != nil {
		fmt.Println("index db error", err)
	}
	emp := Employee{}
	res := []Employee{}
	for selDB.Next() {
		var id, salary int
		var name, destignation string
		err = selDB.Scan(&name, &id, &salary, &destignation)
		if err != nil {
			fmt.Println("index error", err)
		}
		emp.Id = id
		emp.Name = name
		emp.Salary = salary
		emp.Destignation = destignation
		res = append(res, emp)
	}
	tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM Employee WHERE id=?", nId)
	if err != nil {
		fmt.Println("show db error", err)
	}
	emp := Employee{}
	for selDB.Next() {
		var id, salary int
		var name, destignation string
		err = selDB.Scan(&name, &id, &salary, &destignation)
		if err != nil {
			fmt.Println("show error", err)
		}
		emp.Id = id
		emp.Name = name
		emp.Salary = salary
		emp.Destignation = destignation
	}
	tmpl.ExecuteTemplate(w, "Show", emp)
	defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM Employee WHERE id=?", nId)
	if err != nil {
		fmt.Println("edit db error", err)
	}
	emp := Employee{}
	for selDB.Next() {
		var id, salary int
		var name, destignation string
		err = selDB.Scan(&name, &id, &salary, &destignation)
		if err != nil {
			fmt.Println("edit db error", err)
		}
		emp.Id = id
		emp.Name = name
		emp.Salary = salary
		emp.Destignation = destignation
	}
	tmpl.ExecuteTemplate(w, "Edit", emp)
	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		salary := r.FormValue("salary")
		destignation := r.FormValue("destignation")

		insForm, err := db.Prepare("INSERT INTO Employee(name, salary,destignation) VALUES(?,?,?)")
		if err != nil {
			fmt.Println("Insert error", err)
		}
		insForm.Exec(name, salary, destignation)
		log.Println("INSERT: Name: " + name + " | salary: " + salary + " | destignation: " + destignation)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		salary := r.FormValue("salary")
		destignation := r.FormValue("destignation")

		id := r.FormValue("id")
		insForm, err := db.Prepare("UPDATE Employee SET name=?, salary=? , destignation =? WHERE id=?")
		if err != nil {
			fmt.Println("Update error", err)
		}
		insForm.Exec(name, salary, id, destignation)
		log.Println("UPDATE: Name: " + name + " | Salary: " + salary + " | destignation: " + destignation)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	emp := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM Employee WHERE id=?")
	if err != nil {
		fmt.Println("delete db error", err)
	}
	delForm.Exec(emp)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func main() {
	log.Println("Server started on: http://localhost:8080")
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	http.ListenAndServe(":8080", nil)
}
