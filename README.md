# find-events-berlin
Scrapes the web for finding interesting events in Berlin

## To use
1. Do `go install`
2. Run locally via `go run main.go`

## APIs

### Scrape for events in Berlin
We scrape the web by looking at the events from https://www.visitberlin.de/en/blog/ and summarizing it through ChatGPT

```
# Get events for the current month
http://localhost:8080/events/berlin

# Get events for the next month
http://localhost:8080/events/berlin?month=may
```
