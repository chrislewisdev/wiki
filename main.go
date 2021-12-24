package main

import (
	"fmt"
	"os"
)

func main() {
	contents, err := os.ReadDir("content")
    if err != nil {
        panic(err)
    }

    for _, entry := range contents {
        fmt.Println(entry.Name(), entry.IsDir())
    }
}
