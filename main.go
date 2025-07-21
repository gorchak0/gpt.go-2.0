package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

type RequestData struct {
	FolderPath     string `json:"folderPath"`
	Mode           string `json:"mode"`
	SupportedFiles string `json:"supportedFiles"`
	SupportedDirs  string `json:"supportedDirs"`
	ExcludeFiles   string `json:"excludeFiles"`
	ExcludeDirs    string `json:"excludeDirs"`
	ChunkSize      string `json:"chunkSize"`
	AddPrompt      bool   `json:"addPrompt"`
}

func main() {
	http.Handle("/", http.FileServer(http.Dir(".")))

	http.HandleFunc("/analyze", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, _ := io.ReadAll(r.Body)

		var data RequestData
		if err := json.Unmarshal(body, &data); err != nil {
			http.Error(w, "invalid json", 400)
			return
		}

		// Преобразуем строку chunkSize → int
		chunkSize, err := strconv.Atoi(strings.TrimSpace(data.ChunkSize))
		if err != nil || chunkSize <= 0 {
			http.Error(w, "invalid chunkSize", 400)
			return
		}

		// Создание заглушки текста на 10 000 символов
		var builder strings.Builder
		i := 0
		for builder.Len() < 9996 {
			builder.WriteString(fmt.Sprintf("%d ", i*100))
			i++
		}
		text := builder.String()

		// Разбиваем текст на чанки
		var chunks []string
		for i := 0; i < len(text); i += chunkSize {
			end := i + chunkSize
			if end > len(text) {
				end = len(text)
			}
			chunks = append(chunks, text[i:end])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(chunks)
	})

	// Открытие браузера
	go func() {
		_ = exec.Command("xdg-open", "http://localhost:8080").Start()
		_ = exec.Command("open", "http://localhost:8080").Start()
		_ = exec.Command("rundll32", "url.dll,FileProtocolHandler", "http://localhost:8080").Start()
	}()

	fmt.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
