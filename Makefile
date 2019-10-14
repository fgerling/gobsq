all: check ls-rr-caasp

ls-rr-caasp:  cmd/ls-rr-caasp/main.go pkg/obs/
	go build ./cmd/ls-rr-caasp
check: cmd/check/main.go pkg/obs/
	go build ./cmd/check
test: check ls-rr-caasp
	./check
	./ls-rr-caasp
clean:
	rm ./check ./ls-rr-caasp
