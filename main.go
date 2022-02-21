package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func writeFile(data, filename string) {
	file, error := os.Create(filename)
	defer file.Close()
	check(error)

	file.WriteString(data)
}

func main() {
	url := "https://techcrunch.com/"

	response, error := http.Get(url)
	defer response.Body.Close()
	check(error)

	if response.StatusCode > http.StatusBadRequest {
		fmt.Println("Status Code: ", response.StatusCode)
	}

	doc, error := goquery.NewDocumentFromReader(response.Body)
	check(error)

	file, error := os.Create("posts.csv")
	check(error)
	writer := csv.NewWriter(file)

	doc.Find("div.river").
		Find("div.post-block").
		Each(func(index int, item *goquery.Selection) {
			h2 := item.Find("h2")
			title := strings.TrimSpace(h2.Text())
			url, _ := h2.Find("a").Attr("href")

			excerpt := strings.TrimSpace(item.Find("div.post-block__content").Text())
			fmt.Println(title, url)

			posts := []string{title, url, excerpt}

			writer.Write(posts)
		})

	writer.Flush()
}
