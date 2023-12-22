run: gen
	go run main.go

tidy:
	go mod tidy

gen:
	templ generate

dev: gen
	air