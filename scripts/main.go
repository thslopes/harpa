package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

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
	if len(os.Args) > 1 {
		formatEBD()
		return
	}
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

func formatEBD() {
	// https://escolabiblicadominical.com.br/licao-11-sendo-cautelosos-nas-opressoes/
	dat, err := os.ReadFile("ebd20223/index.html")
	check(err)
	htmlContent := string(dat)
	p := bluemonday.StripTagsPolicy()
	htmlContent = p.Sanitize(htmlContent)
	htmlContent = strings.ReplaceAll(htmlContent, "\t", "")
	for strings.Contains(htmlContent, "  ") {
		htmlContent = strings.ReplaceAll(htmlContent, "  ", " ")
	}
	for strings.Contains(htmlContent, "\n\n") {
		htmlContent = strings.ReplaceAll(htmlContent, "\n\n", "\n")
	}
	i := 0
	linhas := strings.Split(htmlContent, "\n")
	for !strings.HasPrefix(linhas[i], "EBD – Lição ") {
		i++
	}
	out := getTitle(linhas[i], 1)
	for linhas[i] != "TEXTO ÁUREO" {
		i++
	}

	for linhas[i] != "SAIBA TUDO SOBRE A ESCOLA DOMINICAL:" {
		out = append(out, linhas[i])
		i++
	}

	fmt.Println(strings.Join(out, "\n"))
}

func getTitle(s string, n int) []string {
	meses := []string{"", "janeiro", "fevereiro", "março", "abril", "maio", "junho", "julho", "agosto", "setembro", "outubro", "novembro", "dezembro"}
	s = strings.Replace(s, "EBD – ", "", 1)
	x := strings.Split(s, ": ")
	num := x[0]
	title := strings.Split(x[1], " |")[0]
	titleDate := time.Date(2022, 7, 3, 0, 0, 0, 0, time.UTC)
	titleDate = titleDate.AddDate(0, 0, 7*(n-1))
	tdf := fmt.Sprintf(titleDate.Format("2 %s 2006"), meses[titleDate.Month()])
	return []string{num, tdf, title}
}
