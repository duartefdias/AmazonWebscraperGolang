# Amazon WebScraper (Go)

A simple amazon webscraper written in Go. Scrapes amazon's products based on a search terms and generates a csv with the results.

Usage:
```
go run webscraper.go <search term> <pages to scrape>
```

Example:
```
go run webscraper.go skateboard 5
```

CSV Output fields:
```
price
name
url
```