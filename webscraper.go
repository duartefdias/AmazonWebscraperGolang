package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type AmazonProduct struct {
	Name  string
	Price string
	Url   string
}

func main() {
	var products []AmazonProduct

	// If no arguments are passed, print usage and exit
	if len(os.Args) < 3 {
		fmt.Println("Expected usage: go run webscraper.go <search term> <pages to scrape>")
		fmt.Println("Try: go run webscraper.go basketball 3")
		os.Exit(1)
	}

	// Define target URL
	baseUrl := "https://www.amazon.com/s?k=" + os.Args[1] + "&page="

	// Define amount of pages to scrape from command line
	pagesToScrape, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Error converting string to int: ", err)
		os.Exit(1)
	}

	// Instantiate default collector
	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36"

	// Called before performing an HTTP request with Visit()
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
	})

	// Called if an error occurs during the request
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong: ", err)
	})

	// Called right after OnResponse() if the received content is HTML
	c.OnHTML("div.puisg-col-inner", func(e *colly.HTMLElement) {

		// Ignore invalid elements
		h2Text := e.ChildText("h2")
		if h2Text == "" {
			fmt.Println("Element does not contain a h2 child element, skipping...")
			return
		}
		if strings.Contains(h2Text, "Sponsored") {
			fmt.Println("Element is an ad, skipping...")
			return
		}

		price := e.ChildText("span.a-price-whole")
		if price == "" {
			fmt.Println("Element does not contain a price child element, skipping...")
			return
		}

		fmt.Println("Valid product found!")

		var product AmazonProduct

		product.Name = h2Text
		product.Url = e.Request.AbsoluteURL(e.ChildAttr("h2 a.a-link-normal", "href"))
		priceWhole := e.ChildText("span.a-price-whole")
		priceFraction := e.ChildText("span.a-price-fraction")
		product.Price = priceWhole + priceFraction

		products = append(products, product)
	})

	// Called after all OnHTML() callback executions
	c.OnScraped(func(r *colly.Response) {
		fmt.Println(r.Request.URL, " scraped!")
	})

	// Iterate over pages
	for i := 1; i <= pagesToScrape; i++ {
		url := baseUrl + fmt.Sprintf("%d", i)
		err := c.Visit(url)
		if err != nil {
			log.Fatalln("Failed to make request: ", err)
		}
	}

	// Open csv file for writing
	file, err := os.Create("products.csv")
	if err != nil {
		log.Fatalln("Failed to create output CSV file", err)
	}
	defer file.Close()

	// Initialize csv writer
	writer := csv.NewWriter(file)

	// Write CSV headers
	headers := []string{
		"price",
		"name",
		"url",
	}
	writer.Write(headers)

	// Iterate over products and write to CSV
	for _, product := range products {
		data := []string{
			product.Price,
			product.Name,
			product.Url,
		}
		writer.Write(data)
	}
}
