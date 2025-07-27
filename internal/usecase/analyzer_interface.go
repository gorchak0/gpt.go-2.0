package usecase

import "projanalyzer/internal/entity"

type Analyzer interface {
	AnalyzeProject(*entity.ParsedRequest) (string, error)
}
