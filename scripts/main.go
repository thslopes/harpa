package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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
}

func (h hino) String() string {
	// return fmt.Sprintf("%d %s\n", h.Id, h.Linha)
	return fmt.Sprintf("\n{\"id\": %d,\"nome\": \"%s\",\n\"letra\": \"%s\",\n\"linha\": \"%s\"},", h.Id, h.Nome, strings.ReplaceAll(h.Letra, "\"", "\\\""), strings.ReplaceAll(h.Linha, "\"", "\\\""))
}

func main() {
	dat, err := os.ReadFile("../src/app/hinos.ts")
	check(err)

	var hinos []hino
	p := bluemonday.StripTagsPolicy()

	err = json.Unmarshal(dat, &hinos)
	check(err)

	for i := range hinos {
		linha := hinos[i].Letra
		linha = strings.ReplaceAll(linha, "<br />", " ")
		linha = strings.ReplaceAll(linha, "&nbsp;", " ")
		linha = p.Sanitize(linha)
		linha = strings.ReplaceAll(linha, "  ", " ")
		if len(linha) > 50 {
			for i := 70; i > 0; i-- {
				if linha[i] == ' ' {
					linha = linha[0:i] + "..."
					break
				}
			}
		}
		hinos[i].Linha = linha
	}

	dat = []byte(fmt.Sprint(hinos))

	_ = ioutil.WriteFile("hinos.json", dat, 0644)
}
