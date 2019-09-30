package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var parts = make(map[string]map[string][][]string)

func PrintLine(from string, end string) {
	front := strings.Split(from, ".")
	rear := strings.Split(end, ".")

	if len(front) != 4 || len(rear) != 4 {
		return
	}

	top := front[0]
	second := front[1] + "-" + rear[1]

	if _, exists := parts[top]; !exists {
		parts[top] = make(map[string][][]string)
	}
	item := []string{front[2], rear[2]}
	parts[top][second] = append(parts[top][second], item)
}

func getIpPool(url string) {
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	if res != nil {
		defer res.Body.Close()
	}
	fmt.Printf("status code: %d %s\n", res.StatusCode, res.Status)
	doc, oerr := goquery.NewDocumentFromReader(res.Body)
	if oerr != nil {
		panic(oerr)
	}
	doc.Find("table tbody tr").Each(func(i int, s *goquery.Selection) {
		if s.Find("td").Length() < 2 {
			return
		}
		front := s.Find("td").Eq(0).Text()
		rear := s.Find("td").Eq(1).Text()

		front = strings.TrimSpace(front)
		rear = strings.TrimSpace(rear)

		PrintLine(front, rear)
	})
}

func saveToFile(file string, content string) {
	err := ioutil.WriteFile(file, []byte(content), os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func main() {
	// step1: 爬取大陆ip段
	url := "http://ip.bczs.net/country/CN"
	output := "ips.json"
	getIpPool(url)

	jsonstr, err := json.Marshal(parts)
	if err != nil {
		fmt.Println("json.Marshal failed:", err)
		return
	}

	// step2: 保存结果到json文件中
	saveToFile(output, string(jsonstr))

	fmt.Println("finished, ip saved to file: " + output)
}
