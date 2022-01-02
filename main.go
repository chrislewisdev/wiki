package main

import (
	"github.com/gomarkdown/markdown"
	"os"
	"regexp"
	"strings"
)

type document struct {
	name		string
	mdFile		string
	htmlFile	string
}

func toDocName(mdFile string) string {
	return strings.Replace(strings.Split(mdFile, ".")[0], "_", " ", -1)
}

func toHtmlName(mdFile string) string {
	return strings.Replace(mdFile, ".md", ".html", 1)
}

func newDocument(mdFile string) document {
	return document{toDocName(mdFile), mdFile, toHtmlName(mdFile)}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func toSentenceCase(str string) string {
	return strings.ToUpper(str[0:1]) + str[1:len(str)]
}

func getFiles(directory string) []document {
	contents, err := os.ReadDir(directory)
	check(err)

	docs := []document{}
	for _, entry := range contents {
		docs = append(docs, newDocument(entry.Name()))
	}
	return docs
}

func ensureDirectoryExists(directory string) {
	stat, err := os.Stat(directory)
	if err != nil && os.IsNotExist(err) {
		// TODO: Replace magic constant
		err = os.Mkdir(directory, 0755)
		check(err)
	} else if !stat.IsDir() {
		panic("Cannot create '" + directory + "' directory; a file by that name exists")
	}
}

func generateIndex(docs []document) string {
	var index = ""

	for _, doc := range docs {
		link := "[" + doc.name + "](./" + doc.htmlFile + ")"
		index = index + " - " + link + "\n"
	}

	return index
}

func autolink(doc document, md string, docs []document) string {
	for _, otherDoc := range docs {
		if doc.name != otherDoc.name {
			// This doesn't account for words appearing at the end of a sentence
			regex := regexp.MustCompile("(?i)( " + otherDoc.name + " )")
			md = regex.ReplaceAllString(string(md), "[$1](./" + otherDoc.htmlFile + ")")
		}
	}
	return md
}

func renderHtml(md string, docName string) string {
	// TODO: Probably put this into a file and load up as a template string
	header := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Document</title>
</head>
<body>
	<a href="./index.html">index</a>
	`
	footer := "</body></html>"

	body := markdown.ToHTML([]byte(md), nil, nil)
	html := header + "<h1>" + toSentenceCase(docName) + "</h1>\n" + string(body) + footer

	return html
}

func writeFile(filename string, contents string) {
	file, err := os.Create(filename)
	check(err)
	defer file.Close()

	_, err = file.WriteString(contents)
	check(err)
}

func main() {
	contentDirectory := "content"
	buildDirectory := "build"

	docs := getFiles(contentDirectory)

	ensureDirectoryExists(buildDirectory)

	// Generate index
	// TODO: Incorporate about information into index.html, and make it more brief
	writeFile(buildDirectory + "/index.html", renderHtml(generateIndex(docs), "index"))

	// Render out all md -> html files
	for _, doc := range docs {
		mdBytes, err := os.ReadFile(contentDirectory + "/" + doc.mdFile)
		check(err)

		md := string(mdBytes)
		md = autolink(doc, md, docs)

		html := renderHtml(md, doc.name)

		writeFile(buildDirectory + "/" + doc.htmlFile, html)		
	}
}
