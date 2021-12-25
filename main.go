package main

import (
	"github.com/gomarkdown/markdown"
	"os"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
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
    `
	footer := "</body></html>"

	// TODO: Use this information to build up an index for inter-linking of articles
	contentDirectory := "content"
	contents, err := os.ReadDir(contentDirectory)
	check(err)

	buildDirectory := "build"
	stat, err := os.Stat(buildDirectory)
	if err != nil && os.IsNotExist(err) {
		err = os.Mkdir(buildDirectory, 0755)
		check(err)
	} else if !stat.IsDir() {
		panic("Cannot create 'build' directory; a file by that name exists")
	}

	// TODO: Generate an index.html that acts as the main directory

	for _, entry := range contents {
		// TODO: Determine 'doc name' i.e. filename without extension
		// Replace underscores with spaces, e.g. "multiple_sclerosis" -> "multiple sclerosis"
		mdName := entry.Name()
		htmlName := strings.Replace(mdName, ".md", ".html", -1)
		md, err := os.ReadFile(contentDirectory + "/" + mdName)
		check(err)

		body := markdown.ToHTML(md, nil, nil)
		html := header + "<h1>" + mdName + "</h1>\n" + string(body) + footer

		file, err := os.Create("build/" + htmlName)
		check(err)
		defer file.Close()

		_, err = file.WriteString(html)
		check(err)
	}
}
