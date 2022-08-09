package main

import ( 
	"log"
	"net/http"
	"individual/handler"
)



func main() {
	
	mux := http.NewServeMux()

	mux.HandleFunc("/", handler.HomeForm)
	mux.HandleFunc("/addForm", handler.AddForm)
	mux.HandleFunc("/processAddForm", handler.ProcessAddForm)
	mux.HandleFunc("/task", handler.GetTask)
	mux.HandleFunc("/processEditForm", handler.ProcessEditForm)
	mux.HandleFunc("/status", handler.ChangeStatus)
	mux.HandleFunc("/delete", handler.DeleteTask)

	fileServer := http.FileServer(http.Dir("assets"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Println("starting web in port 3000")

	err := http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}

