# Amazon Product WebScraper (Go)

A simple amazon webscraper written in Go. Scrapes amazon's products based on a search term and generates a csv with the results.

Usage:
```
go run webscraper.go <search term> <pages to scrape>
```

Example:
```
go run webscraper.go skateboard 5
```

CSV output fields:
```
price
name
url
```