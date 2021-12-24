package main

import (
	"github.com/gomarkdown/markdown"
	"os"
	"strings"
)

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
	// TODO: tidy up error checking
	if err != nil {
		panic(err)
	}

	// TODO: Only do this if directory doesn't exist
	err = os.Mkdir("build", 0755)
	if err != nil {
		panic(err)
	}

	// TODO: Generate an index.html that acts as the main directory

	for _, entry := range contents {
		// TODO: Determine 'doc name' i.e. filename without extension
		// Replace underscores with spaces, e.g. "multiple_sclerosis" -> "multiple sclerosis"
		mdName := entry.Name()
		htmlName := strings.Replace(mdName, ".md", ".html", -1)
		md, err := os.ReadFile(contentDirectory + "/" + mdName)
		if err != nil {
			panic(err)
		}

		body := markdown.ToHTML(md, nil, nil)
		html := header + "<h1>" + mdName + "</h1>\n" + string(body) + footer

		// TODO: Overwrite file if exists
		file, err := os.Create("build/" + htmlName)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		_, err = file.WriteString(html)
		if err != nil {
			panic(err)
		}
	}
}
