// interfaces/handlers.go
package interfaces

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"projanalyzer/internal/entity"
	"projanalyzer/internal/usecase"
)

var (
	analyzer usecase.Analyzer = usecase.NewDefaultAnalyzer()
	splitter usecase.Splitter = usecase.NewDefaultSplitter()
)

func AnalyzeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var data entity.RequestData
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	req := parseRequest(&data)

	// Валидация: Проверяем, что путь существует и это директория
	info, err := os.Stat(req.FolderPath)
	if err != nil {
		http.Error(w, "Folder path does not exist or is inaccessible", http.StatusBadRequest)
		return
	}
	if !info.IsDir() {
		http.Error(w, "Provided path is not a directory", http.StatusBadRequest)
		return
	}

	result, err := analyzer.AnalyzeProject(req)
	if err != nil {
		http.Error(w, "Scan failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	chunks := splitter.SplitTextWithPrompt(result, req.ChunkSize, data.AddPrompt)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chunks)
}

func parseRequest(data *entity.RequestData) *entity.ParsedRequest {
	return &entity.ParsedRequest{
		FolderPath:    filepath.Clean(data.FolderPath),
		Mode:          data.Mode,
		SupportedExts: csvToExtSet(data.SupportedFiles), // расширения
		SupportedDirs: csvToNameSet(data.SupportedDirs), // имена папок
		ExcludedDirs:  csvToNameSet(data.ExcludeDirs),   // имена папок
		ExcludedFiles: csvToNameSet(data.ExcludeFiles),  // имена файлов
		ChunkSize:     parseInt(data.ChunkSize, 15000),
	}
}

func parseInt(s string, def int) int {
	if val, err := strconv.Atoi(s); err == nil && val > 0 {
		return val
	}
	return def
}

// Для расширений (с точкой в начале)
func csvToExtSet(csv string) map[string]bool {
	set := make(map[string]bool)
	for _, item := range strings.Split(csv, ",") {
		item = strings.TrimSpace(strings.ToLower(item))
		if item == "" {
			continue
		}
		if !strings.HasPrefix(item, ".") {
			item = "." + item
		}
		set[item] = true
	}
	return set
}

// Для имён файлов и папок (без добавления точки)
func csvToNameSet(csv string) map[string]bool {
	set := make(map[string]bool)
	for _, item := range strings.Split(csv, ",") {
		item = strings.TrimSpace(strings.ToLower(item))
		if item != "" {
			set[item] = true
		}
	}
	return set
}
