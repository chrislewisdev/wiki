package main

import (
	"os"
    "strings"
    "github.com/gomarkdown/markdown"
)

func main() {
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

    contentDirectory := "content"
	contents, err := os.ReadDir(contentDirectory)
    if err != nil {
        panic(err)
    }

    err = os.Mkdir("build", 0755)
    if err != nil {
        panic(err)
    }

    for _, entry := range contents {
        mdName := entry.Name()
        htmlName := strings.Replace(mdName, ".md", ".html", -1)
        md, err := os.ReadFile(contentDirectory + "/" + mdName)
        if err != nil {
            panic(err)
        }

        body := markdown.ToHTML(md, nil, nil)
        html := header + "<h1>" + mdName + "</h1>\n" + string(body) + footer

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
