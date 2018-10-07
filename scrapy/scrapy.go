package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

func main() {
	url := `http://cw.hubwiz.com/card/c/549a704f88dba0136c371703/1/1/1/`

	f, err := os.Create("tutorial.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for {
		fmt.Println(url)
		resp, err := http.Get(url)
		if err != nil {
			break
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			break
		}
		src := string(body)
		//fmt.Println(html)

		article, _ := regexp.Compile(`(?is)<article.*</article>`)
		m := article.FindStringSubmatch(src)
		fmt.Println(m[0])
		f.Write([]byte(m[0]))

		next, _ := regexp.Compile(`(?isU)<a class="btn btn-default next" href="(.*)">`)
		m = next.FindStringSubmatch(src)
		fmt.Println(m[1])
		url = "http://cw.hubwiz.com" + m[1]
	}
}
