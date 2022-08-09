package handler

import (
	"net/http"
	"html/template"
	"path"
	"log"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
	"strconv"

)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "IndividualProject"
)

type taskStruct struct {
	ID int
	Task string
	Asignee string
	Deadline time.Time
	IsDone string
}

func (t taskStruct) TaskStatus() string {
	var status string
	if t.IsDone	== "N" {
		status = "Mark as Done"
	} else {
		status = "Done"
	}
	return status
}

func CreateConnection() *sql.DB {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)
	return db
	// defer db.Close()	

}

func CheckError(err error) {
	if err != nil {
			panic(err)
	}
}

func HomeForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		db := CreateConnection() 
		rows, err := db.Query(`SELECT * FROM "task" order by "deadline" ASC`)
		CheckError(err)
		defer rows.Close()
		var Result []taskStruct
		for rows.Next() {
		var each = taskStruct{}

		err = rows.Scan(&each.ID, &each.Task, &each.Asignee, &each.Deadline, &each.IsDone)
		Result = append(Result, each)

		CheckError(err)
		
		}

		tmpl, err := template.ParseFiles(path.Join("views", "index.html"), path.Join("views", "layout.html"))
		if err != nil {
			log.Println(err) 
			http.Error(w, "error is happening", http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, Result)
		if err != nil {
			log.Println(err)
			http.Error(w, "error is happening", http.StatusInternalServerError)
			return
		}
		return
	}
	http.Error(w, "error is happening", http.StatusBadRequest)
}

func AddForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles(path.Join("views", "addForm.html"), path.Join("views", "layout.html"))
		if err != nil {
			log.Println(err) 
			http.Error(w, "error is happening", http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Println(err)
			http.Error(w, "error is happening", http.StatusInternalServerError)
			return
		}
		return
	}
	http.Error(w, "error is happening", http.StatusBadRequest)
}

func ProcessAddForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
			http.Error(w, "error is happening", http.StatusInternalServerError)
			return
		}
	}

	task := r.Form.Get("task")
	asignee := r.Form.Get("asignee")
	deadline := r.Form.Get("deadline")
	isDone := "N"
	db := CreateConnection() 

	insertStmt := `insert into "task"("task", "asignee", "deadline", "isDone") values($1, $2, $3, $4)`
	_, e := db.Exec(insertStmt, task, asignee, deadline, isDone)
  CheckError(e)

	http.Redirect(w, r, "/", http.StatusFound)
	return
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	idNumb, err := strconv.Atoi(id)
	if err != nil || idNumb < 1 {
		http.NotFound(w, r)
		return
	}
	db := CreateConnection()
	defer db.Close()
	sqlStatement := `select * from "task" where "id" = $1`
	rows := db.QueryRow(sqlStatement, idNumb)
	var Result = taskStruct{}
	err = rows.Scan(&Result.ID, &Result.Task, &Result.Asignee, &Result.Deadline, &Result.IsDone)

	tmpl, err := template.ParseFiles(path.Join("views", "editForm.html"), path.Join("views", "layout.html"))
	if err != nil {
		log.Println(err)
		http.Error(w, "error is happening", http.StatusInternalServerError)
		return
	}
	log.Println(Result)
	err = tmpl.Execute(w, Result)
	if err != nil {
		log.Println(err)
		http.Error(w, "error is hapening", http.StatusInternalServerError)
		return
	}
	// http.Error(w, "error is happening", http.StatusBadRequest)
}

func ProcessEditForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
			http.Error(w, "error is happening", http.StatusInternalServerError)
			return
		}
	}

	task := r.Form.Get("task")
	asignee := r.Form.Get("asignee")
	deadline := r.Form.Get("deadline")
	id := r.Form.Get("id")
	db := CreateConnection()

	updateStmt := `update "task" set "task"=$1, "asignee"=$2, "deadline"=$3 where "id"=$4`
	_, e := db.Exec(updateStmt, task, asignee, deadline, id)
	CheckError(e)

	http.Redirect(w, r, "/", http.StatusFound)
	return
}

func ChangeStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
			http.Error(w, "error is happening", http.StatusInternalServerError)
			return
		}
	}
	id := r.URL.Query().Get("id")
	idNumb, err := strconv.Atoi(id)
	if err != nil || idNumb < 1 {
		http.NotFound(w, r)
		return
	}
	db := CreateConnection()
	defer db.Close()
	sqlStatement := `select "isDone" from "task" where "id" = $1`
	rows := db.QueryRow(sqlStatement, idNumb)
	var Result = taskStruct{}
	err = rows.Scan(&Result.IsDone)

	if Result.IsDone == "N" {
		Result.IsDone = "Y"
	} else {
		Result.IsDone = "N"
	}

	updateStmt := `update "task" set "isDone"=$1 where "id"=$2`
	_, e := db.Exec(updateStmt, Result.IsDone, idNumb)
	CheckError(e)

	http.Redirect(w, r, "/", http.StatusFound)
	return
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	idNumb, err := strconv.Atoi(id)
	if err != nil || idNumb < 1 {
		http.NotFound(w, r)
		return
	}
	db := CreateConnection()
	defer db.Close()
	deleteStmt := `delete from "task" where "id"=$1`
	_, e := db.Exec(deleteStmt, idNumb)
	CheckError(e)

	http.Redirect(w, r, "/", http.StatusFound)
	return
}