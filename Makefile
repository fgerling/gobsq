all: lsrepos lsrr

lsrr:  cmd/lsrr/main.go pkg/obs/
	go build ./cmd/lsrr

lsrepos: cmd/lsrepos/main.go pkg/obs/
	go build ./cmd/lsrepos

test: lsrepos lsrr
	./lsrepos
	./lsrr

clean:
	rm ./lsrepos ./lsrr

install_lsrepos: lsrepos
	install lsrepos $(prefix)/bin/lsrepos

install_lsrr: lsrr
	install lsrr $(prefix)/bin/lsrr

install: install_lsrepos install_lsrr
