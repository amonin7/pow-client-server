package manager

import "math/rand"

func GetResource() string {
	quotes := GetWordOfWisdomQuotes()
	quotesAmount := len(quotes)
	return quotes[rand.Intn(quotesAmount)]
}
