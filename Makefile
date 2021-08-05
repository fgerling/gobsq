all: lsrepos lsrr

lsrr:  cmd/lsrr/main.go
	go build ./cmd/lsrr

lsrepos: cmd/lsrepos/main.go
	go build ./cmd/lsrepos

test: lsrepos lsrr
	./lsrepos
	./lsrr

clean:
	rm ./lsrepos ./lsrr

go_install_lsrepos: lsrepos
	go install ./cmd/lsrepos

go_install_lsrr: lsrr
	go install ./cmd/lsrr

install_lsrepos: lsrepos
	install lsrepos $(prefix)/bin/lsrepos

install_lsrr: lsrr
	install lsrr $(prefix)/bin/lsrr

go_install: go_install_lsrepos go_install_lsrr

install: install_lsrepos install_lsrr

