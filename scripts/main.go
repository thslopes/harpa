package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/microcosm-cc/bluemonday"

	"log"

	"golang.org/x/net/html"

	"github.com/bmaupin/go-epub"
	"github.com/bmaupin/go-htmlutil"
)

const (
	effectiveGoCoverImg      = "assets/covers/capa.png"
	effectiveGoFilename      = "lbap20223.epub"
	effectiveGoSectionTag    = "h1"
	effectiveGoTitle         = "Lições Bĩblicas Professor 2022-3"
	effectiveGoTitleFilename = "title.xhtml"
	revistaUrl               = "assets/full.html"
	epubCSSFile              = "assets/ebub.css"
	preFontFile              = "assets/fonts/SourceCodePro-Regular.ttf"
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

	if len(os.Args) > 1 && os.Args[1] == "epub" {
		err := buildEffectiveGo()
		if err != nil {
			log.Printf("Error building Effective Go: %s", err)
		}
		return
	}

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
	dat, err := os.ReadFile("ebd20223/full.md")
	check(err)
	htmlContent := string(dat)
	linhas := strings.Split(htmlContent, "\n")
	i := 0
	lessonLines := []string{}
	for _, l := range linhas {
		if strings.HasPrefix(l, "# LI") {
			lessonFile(i, lessonLines)
			i += 1
			lessonLines = []string{}
		}
		lessonLines = append(lessonLines, l)
	}
	lessonFile(i, lessonLines)
}

func lessonFile(i int, lines []string) {
	dat := []byte(strings.Join(lines, "\n"))
	fileName := fmt.Sprintf("ebd20223/%02d.md", i)
	if i == 0 {
		fileName = "ebd20223/README.md"
	}
	_ = ioutil.WriteFile(fileName, dat, 0644)
}

type epubSection struct {
	title    string
	filename string
	nodes    []html.Node
}

func buildEffectiveGo() error {
	r, err := os.Open(revistaUrl)
	if err != nil {
		return err
	}
	// defer func() {
	// if err := resp.Body.Close(); err != nil {
	// panic(err)
	// }
	// }()

	doc, err := html.Parse(r)
	if err != nil {
		log.Fatalf("Parse error: %s", err)
	}

	pageNode := htmlutil.GetFirstHtmlNode(doc, "div", "", "")

	sections := []epubSection{}
	sectionFilename := effectiveGoTitleFilename
	// Don't add a title so it doesn't get added to the TOC
	section := &epubSection{
		filename: effectiveGoTitleFilename,
	}
	internalLinks := make(map[string]string)

	initFound := false

	// Iterate through each child node
	for node := pageNode.FirstChild; node != nil; node = node.NextSibling {
		if !initFound && getId(node) != "lições-bíblicas" {
			continue
		}
		initFound = true
		// If we find the section tag
		if node.Type == html.ElementNode && node.Data == effectiveGoSectionTag {
			// Add the previous section to the slice of sections
			sections = append(sections, *section)

			sectionTitle := node.FirstChild.Data
			sectionFilename = titleToFilename(sectionTitle)

			// Start a new section
			section = &epubSection{
				filename: sectionFilename,
				title:    sectionTitle,
			}
		}

		section.nodes = append(section.nodes, *node)

		// Map internal links to their section filename
		for _, idNode := range htmlutil.GetAllHtmlNodes(node, "", "id", "") {
			for _, attr := range idNode.Attr {
				if attr.Key == "id" {
					internalLinks[attr.Val] = fmt.Sprintf("%s#%s", sectionFilename, attr.Val)
				}
			}
		}
	}

	// Make sure the last section gets added
	sections = append(sections, *section)

	e := epub.NewEpub(effectiveGoTitle)
	// effectiveGoCoverImgPath, err := filepath.Abs(effectiveGoCoverImg)
	effectiveGoCoverImgPath, err := e.AddImage(effectiveGoCoverImg, "cover.png")
	if err != nil {
		return err
	}
	e.SetCover(effectiveGoCoverImgPath, "")

	epubCSSPath, err := e.AddCSS(epubCSSFile, "")
	if err != nil {
		return err
	}

	_, err = e.AddFont(preFontFile, "")
	if err != nil {
		return err
	}

	// Iterate through each section and add it to the EPUB
	for _, section := range sections {
		sectionContent := ""
		for _, sectionNode := range section.nodes {
			// Fix internal links so they work after splitting page into sections
			for _, linkNode := range htmlutil.GetAllHtmlNodes(&sectionNode, "a", "", "") {
				for i, attr := range linkNode.Attr {
					if attr.Key == "href" && strings.HasPrefix(attr.Val, "#") {
						linkNode.Attr[i].Val = internalLinks[attr.Val[1:]]
					}
				}
			}

			nodeContent, err := htmlutil.HtmlNodeToString(&sectionNode)
			if err != nil {
				return err
			}
			sectionContent += nodeContent
		}

		_, err := e.AddSection(sectionContent, section.title, section.filename, epubCSSPath)
		if err != nil {
			return err
		}
	}

	err = e.Write(effectiveGoFilename)
	if err != nil {
		return err
	}

	return nil
}

func getId(n *html.Node) string {
	for _, attr := range n.Attr {
		if attr.Key == "id" {
			return attr.Val
		}
	}
	return ""
}

func titleToFilename(title string) string {
	title = strings.ToLower(title)
	title = strings.Replace(title, " ", "-", -1)

	return fmt.Sprintf("%s.xhtml", title)
}
