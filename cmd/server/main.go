package main

import (
	"log"
	"net/http"

	"projanalyzer/internal/infrastructure"
	"projanalyzer/internal/interfaces"
)

func main() {
	log.Println("Сервер запускается...")

	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/analyze", interfaces.AnalyzeHandler)

	go infrastructure.OpenBrowser("http://localhost:55555")

	log.Println("✅ Сервер запущен: http://localhost:55555")
	log.Fatal(http.ListenAndServe(":55555", nil))
}
