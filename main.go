package main

import (
	"github.com/gomarkdown/markdown"
	"io/fs"
	"os"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func toSentenceCase(str string) string {
	return strings.ToUpper(str[0:1]) + str[1:len(str)]
}

func getFiles(directory string) []fs.DirEntry {
	contents, err := os.ReadDir(directory)
	check(err)
	return contents
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

func toDocName(mdName string) string {
	return strings.Replace(strings.Split(mdName, ".")[0], "_", " ", -1)
}

func toHtmlName(mdName string) string {
	return strings.Replace(mdName, ".md", ".html", 1)
}

func generateIndex(content []fs.DirEntry) []byte {
	var index = ""

	for _, entry := range content {
		mdName := entry.Name()
		link := "[" + toDocName(mdName) + "](./" + toHtmlName(mdName) + ")"
		index = index + " - " + link + "\n"
	}

	return []byte(index)
}

func renderHtml(md []byte, docName string) string {
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

	body := markdown.ToHTML(md, nil, nil)
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

	// TODO: Use this information to build up an index for inter-linking of articles
	content := getFiles(contentDirectory)

	ensureDirectoryExists(buildDirectory)

	// Generate index
	writeFile(buildDirectory + "/index.html", renderHtml(generateIndex(content), "index"))

	// Render out all md -> html files
	for _, entry := range content {
		mdName := entry.Name()
		docName := toDocName(mdName)
		htmlName := toHtmlName(mdName)

		md, err := os.ReadFile(contentDirectory + "/" + mdName)
		check(err)

		html := renderHtml(md, docName)		

		writeFile(buildDirectory + "/" + htmlName, html)		
	}
}
