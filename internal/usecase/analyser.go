package usecase

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"projanalyzer/internal/entity"
)

type DefaultAnalyzer struct{}

func NewDefaultAnalyzer() *DefaultAnalyzer {
	return &DefaultAnalyzer{}
}

func (a *DefaultAnalyzer) AnalyzeProject(req *entity.ParsedRequest) (string, error) {
	builder := &strings.Builder{}
	builder.WriteString(fmt.Sprintf("Проект: %s\n\n", req.FolderPath))

	err := filepath.Walk(req.FolderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Просто игнорируем ошибку и продолжаем обход
			return nil
		}

		relPath, _ := filepath.Rel(req.FolderPath, path)
		name := strings.ToLower(info.Name())

		if info.IsDir() {
			if req.ExcludedDirs[name] {
				return filepath.SkipDir
			}
			if len(req.SupportedDirs) > 0 && !req.SupportedDirs[name] {
				return filepath.SkipDir
			}
			if req.Mode == "paths" {
				builder.WriteString("[Папка] " + relPath + "\n")
			}
			return nil
		}

		// Фильтрация файлов
		ext := strings.ToLower(filepath.Ext(name))
		if !req.SupportedExts[ext] {
			return nil
		}

		if req.ExcludedFiles[name] {
			return nil
		}

		if req.Mode == "paths" {
			builder.WriteString("[Файл] " + relPath + "\n")
		} else {
			content, err := os.ReadFile(path)
			if err != nil {
				// Игнорируем ошибку чтения файла, но продолжаем обход
				return nil
			}
			builder.WriteString(fmt.Sprintf("%s:\n%s\n\n", relPath, content))
		}

		return nil
	})

	return builder.String(), err
}
