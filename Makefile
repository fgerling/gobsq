check: cmd/check/main.go pkg/
	go build ./cmd/check
test: check
	./check
clean:
	rm ./check
