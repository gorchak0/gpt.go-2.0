package usecase

import "strconv"

type DefaultSplitter struct{}

func NewDefaultSplitter() *DefaultSplitter {
	return &DefaultSplitter{}
}

func (s *DefaultSplitter) SplitTextWithPrompt(text string, chunkSize int, addPrompt bool) []string {
	runes := []rune(text)

	if !addPrompt {
		var chunks []string
		for i := 0; i < len(runes); i += chunkSize {
			end := i + chunkSize
			if end > len(runes) {
				end = len(runes)
			}
			chunks = append(chunks, string(runes[i:end]))
		}
		return chunks
	}

	makeFirstPrompt := func() string {
		return `Я собираюсь передать вам большой объем текста, который не может быть отправлен одним сообщением. Я буду делить его на части и оформлять каждую по следующему шаблону:

[НАЧАЛО ЧАСТИ 1/10]
<содержимое части>
[КОНЕЦ ЧАСТИ 1/10]

После получения каждой части подтверждайте её простым сообщением:
«Получено, часть 1/10»
Не начинайте анализ или обработку, пока я не отправлю финальное сообщение:
«ВСЕ ЧАСТИ ОТПРАВЛЕНЫ»
После этого вы можете использовать полученные данные для выполнения моих запросов.
`
	}

	makeIntermediateStart := func(part, total int) string {
		return "Пока не отвечайте. Это всего лишь часть текста, который я хочу вам отправить. Просто получите его и подтвердите получение сообщением «Часть " +
			strconv.Itoa(part) + "/" + strconv.Itoa(total) + " получена», а затем ждите следующую часть.\n" +
			"[НАЧАЛО ЧАСТИ " + strconv.Itoa(part) + "/" + strconv.Itoa(total) + "]\n"
	}

	makeIntermediateEnd := func(part, total int) string {
		return "\n[КОНЕЦ ЧАСТИ " + strconv.Itoa(part) + "/" + strconv.Itoa(total) + "]\nНе отвечайте пока. Просто подтвердите получение этой части сообщением \n«Часть " +
			strconv.Itoa(part) + "/" + strconv.Itoa(total) + " получена», а затем ждите следующую часть.\n"
	}

	makeLastStart := func(part, total int) string {
		return "[НАЧАЛО ЧАСТИ " + strconv.Itoa(part) + "/" + strconv.Itoa(total) + "]\n"
	}

	makeLastEnd := func(part, total int) string {
		return "\n[КОНЕЦ ЧАСТИ " + strconv.Itoa(part) + "/" + strconv.Itoa(total) + "]\nТеперь все данные переданы. Вы можете приступить к их анализу и выполнить последующие инструкции, основанные на полученной информации.\n"
	}

	totalChunks := (len(runes) + chunkSize - 1) / chunkSize
	chunks := make([]string, 0, totalChunks+1)

	chunks = append(chunks, makeFirstPrompt())

	for i := 0; i < totalChunks; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > len(runes) {
			end = len(runes)
		}
		partNum := i + 1

		content := string(runes[start:end])
		var block string

		if partNum == totalChunks {
			block = makeLastStart(partNum, totalChunks) + content + makeLastEnd(partNum, totalChunks)
		} else {
			block = makeIntermediateStart(partNum, totalChunks) + content + makeIntermediateEnd(partNum, totalChunks)
		}

		chunks = append(chunks, block)
	}

	return chunks
}
