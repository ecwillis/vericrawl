package main

import (
    "os"
    "fmt"
    "net/http"
    "io/ioutil"
    "strings"
    "time"
    "regexp"
    . "github.com/logrusorgru/aurora"
)

func main() {
    args := os.Args

    if len(args) < 2 {
      panic("No url passed in")
    }

    mainUrl := args[1]

    fmt.Printf("Reading %s",mainUrl)
    fmt.Println("")


    // Client based w/UserAgent
    client := &http.Client{}

    req, err := http.NewRequest("GET", mainUrl, nil)
    if err != nil {
      panic(err)
    }
    req.Header.Set("User-Agent", "F+T VeriCrawl bot/1.0")

    resp, err := client.Do(req)

    if err != nil {
        panic(err)
    }

    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }

    bodyStr := string(body)
    crawlCache := strings.Split(bodyStr, "\n")

    duration := time.Duration(300) * time.Millisecond
    var badHtmlTest = regexp.MustCompile(`&lt;.*&gt;`)
    var okHtmlTest = regexp.MustCompile(`[a-z-]*=".*&lt;.*&gt;.*"`)
    for _, url := range crawlCache {
        plainUrl := strings.Replace(url, "https", "http", 1)
        subReq, err := http.NewRequest("GET", plainUrl, nil)
        subReq.Header.Set("User-Agent", "F+T VeriCrawl bot/1.0")

        subResp, err := client.Do(subReq)

        if err != nil {
            panic(err)
        }

        defer subResp.Body.Close()

        subBody, err := ioutil.ReadAll(subResp.Body)

        badData := badHtmlTest.FindIndex(subBody)
        okData := okHtmlTest.FindIndex(subBody)

        htmlValid := Green("true ")
        if len(badData) > 0 && len(okData) != len(badData) {
          htmlValid = Red("false")
        }

        if err != nil {
          panic(err)
        }

        retCode := Green(subResp.StatusCode)

        if subResp.StatusCode < 399 && subResp.StatusCode > 200 {
            retCode = Brown(subResp.StatusCode)
        }

        if subResp.StatusCode > 399 && subResp.StatusCode < 500 {
            retCode = Magenta(subResp.StatusCode)
        }

        if subResp.StatusCode > 499 {
            retCode = Red(subResp.StatusCode)
        }

        fmt.Println(retCode, htmlValid, url)
        time.Sleep(duration)
    }
}
