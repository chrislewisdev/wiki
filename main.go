package main

import (
	"github.com/gomarkdown/markdown"
	"bytes"
	"text/template"
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

func contains(list []string, element string) bool {
	for _, item := range list {
		if item == element {
			return true
		}
	}

	return false
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

func copyIfExists(filename string, copyName string) {
	stat, err := os.Stat(filename)
	
	if err == nil && stat.Mode().IsRegular() {
		bytes, err := os.ReadFile(filename)
		check(err)
		writeFile(copyName, string(bytes))
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
	// Consider making this info a metadata of the pages rather than hardcoded
	blocklist := []string{"links", "about", "now"}

	for _, otherDoc := range docs {
		if doc.name != otherDoc.name && !contains(blocklist, otherDoc.name) {
			regex := regexp.MustCompile("(?i)(\\W)(" + otherDoc.name + ")(\\W)")
			md = regex.ReplaceAllString(string(md), "$1[$2](./" + otherDoc.htmlFile + ")$3")
		}
	}
	return md
}

func renderHtml(layout *template.Template, md string, docName string) string {
	body := markdown.ToHTML([]byte(md), nil, nil)

	templateData := struct {
		Title string
		Content string
	}{
		Title: toSentenceCase(docName),
		Content: string(body),
	}

	buffer := &bytes.Buffer{}
	err := layout.Execute(buffer, templateData)
	check(err)

	return buffer.String()
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

	ensureDirectoryExists(buildDirectory)

	layout := template.Must(template.ParseFiles("design/layout.html"))
	copyIfExists("design/style.css", "build/style.css")

	docs := getFiles(contentDirectory)

	// Generate index
	// TODO: Incorporate about information into index.html, and make it more brief
	writeFile(buildDirectory + "/index.html", renderHtml(layout, generateIndex(docs), "index"))

	// Render out all md -> html files
	for _, doc := range docs {
		mdBytes, err := os.ReadFile(contentDirectory + "/" + doc.mdFile)
		check(err)

		// TODO: Insert some kind of "last modified" date via Git history?
		md := string(mdBytes)
		md = autolink(doc, md, docs)

		html := renderHtml(layout, md, doc.name)

		writeFile(buildDirectory + "/" + doc.htmlFile, html)		
	}
}
