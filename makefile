i:
	go mod tidy
	go mod download

build:
	go build -o statistic-service cmd/main.go
run:
	./statistic-service
buildAndRun:
	go build -o statistic-service cmd/main.go
	./statistic-service
