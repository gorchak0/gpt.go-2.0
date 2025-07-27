package usecase

type Splitter interface {
	SplitTextWithPrompt(text string, chunkSize int, addPrompt bool) []string
}
