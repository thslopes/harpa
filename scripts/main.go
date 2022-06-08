package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type hino struct {
	Id    int
	Nome  string
	Letra string
	Linha string
	Busca string
}

func (h hino) String() string {
	// return fmt.Sprintf("%d %s\n", h.Id, h.Busca)
	return fmt.Sprintf("\n{\"id\": %d,\"nome\": \"%s\",\n\"letra\": \"%s\",\n\"linha\": \"%s\",\n\"busca\": \"%s\"},",
		h.Id, h.Nome,
		strings.ReplaceAll(h.Letra, "\"", "\\\""),
		strings.ReplaceAll(h.Linha, "\"", "\\\""),
		strings.ReplaceAll(h.Busca, "\"", "\\\""))
}

func main() {
	dat, err := os.ReadFile("hinos.ts")
	check(err)

	var hinos []hino
	p := bluemonday.StripTagsPolicy()

	err = json.Unmarshal(dat, &hinos)
	check(err)

	reAccents := regexp.MustCompile(`&(\w)\w+;`)
	reMarks := regexp.MustCompile(`[^\w\s]`)

	for i := range hinos {
		linha := hinos[i].Letra
		linha = strings.ReplaceAll(linha, "<br />", " ")
		linha = strings.ReplaceAll(linha, "&nbsp;", " ")
		linha = reAccents.ReplaceAllString(linha, "$1")

		linha = p.Sanitize(linha)
		// linha = strings.ReplaceAll(linha, "&#34;", "")
		// linha = strings.ReplaceAll(linha, "&#39;", "")
		linha = reMarks.ReplaceAllString(linha, " ")
		linha = strings.ReplaceAll(linha, "  ", " ")
		// if len(linha) > 50 {
		// 		for i := 70; i > 0; i-- {
		// 			if linha[i] == ' ' {
		// 				linha = linha[0:i] + "..."
		// 				break
		// 			}
		// 		}
		// 	}
		hinos[i].Busca = strings.ToLower(linha)
	}

	fmt.Println(hinos[0:10])

	dat = []byte(fmt.Sprint(hinos))

	_ = ioutil.WriteFile("hinos.json", dat, 0644)
}