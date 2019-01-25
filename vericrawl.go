package main

import (
    "os"
    "fmt"
    "net/http"
    "io/ioutil"
    "strings"
    "time"
    . "github.com/logrusorgru/aurora"
)

func main() {
    args := os.Args

    if len(args) < 2 {
      panic("No url passed in")
    }

    mainUrl := args[1]

    fmt.Printf("Reading %s/sitemap.txt",mainUrl)
    fmt.Println("")

    resp, err := http.Get(mainUrl + "/sitemap.txt")

    if err != nil {
        fmt.Println("Error!")
        return
    }

    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)

    if err != nil {
        fmt.Println("Body Error!")
        return
    }
    bodyStr := string(body)

    crawlCache := strings.Split(bodyStr, "\n")

    duration := time.Duration(3) * time.Second

    for _, url := range crawlCache {
        plainUrl := strings.Replace(url, "https", "http", 1)
        resp, err := http.Head(plainUrl)
        if err != nil {
            panic(err)
        }
        if resp.StatusCode == 200 {
            fmt.Println(Green(resp.StatusCode), plainUrl)
        }

        if resp.StatusCode < 399 && resp.StatusCode > 200 {
            fmt.Println(Brown(resp.StatusCode), plainUrl)
        }

        if resp.StatusCode > 399 && resp.StatusCode < 500 {
            fmt.Println(Magenta(resp.StatusCode), plainUrl)
        }

        if resp.StatusCode > 499 {
            fmt.Println(Red(resp.StatusCode), plainUrl)
        }
        time.Sleep(duration)
    }
}
