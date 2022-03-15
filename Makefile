all: lsrepos lsrr lsmu

lsmu:  cmd/lsmu/main.go
	go build ./cmd/lsmu

lsrr:  cmd/lsrr/main.go
	go build ./cmd/lsrr

lsrepos: cmd/lsrepos/main.go
	go build ./cmd/lsrepos

test: lsrepos lsrr
	./lsrepos
	./lsrr
	./lsmu

clean:
	rm ./lsrepos ./lsrr ./lsmu

go_install_lsmu: lsmu
	go install ./cmd/lsmu

go_install_lsrepos: lsrepos
	go install ./cmd/lsrepos

go_install_lsrr: lsrr
	go install ./cmd/lsrr

install_lsmu: lsmu
	install lsmu $(prefix)/bin/lsmu

install_lsrepos: lsrepos
	install lsrepos $(prefix)/bin/lsrepos

install_lsrr: lsrr
	install lsrr $(prefix)/bin/lsrr

go_install: go_install_lsrepos go_install_lsrr go_install_lsmu

install: install_lsrepos install_lsrr install_lsmu

