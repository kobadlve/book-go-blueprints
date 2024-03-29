package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/kobadlve/book-go-blueprints/ch4/thesaurus"
)

func main() {
	apiKey := os.Getenv("BHT_APIKEY")
	thesaurus := &thesaurus.BigHuge{APIKey: apiKey}
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		word := s.Text()
		syns, err := thesaurus.Synonyms(word)
		if err != nil {
			log.Fatalf("Failed to search synonyms %q: %v\n", word, err)
		}
		if len(syns) == 0 {
			log.Fatalf("synonyms of %q is none\n", word)
		}
		for _, syn := range syns {
			fmt.Println(syn)
		}
	}
}
