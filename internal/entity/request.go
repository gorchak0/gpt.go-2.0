package entity

// RequestData — DTO из JSON-запроса от клиента
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

// ParsedRequest — структура, используемая в бизнес-логике после валидации и преобразования
type ParsedRequest struct {
	FolderPath    string
	Mode          string
	SupportedExts map[string]bool
	SupportedDirs map[string]bool
	ExcludedDirs  map[string]bool
	ExcludedFiles map[string]bool
	ChunkSize     int
}
