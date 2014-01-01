# Comments starting with #: below are remake GNU Makefile comments. See
# https://github.com/rocky/remake/wiki/Rake-tasks-for-gnu-make

.PHONY: all exports test check clean

#: Same as make go-fish
all: go-fish

#: The non-GNU Readline REPL front-end to the go-interactive evaluator
go-fish: repl_imports.go main.go repl.go
	go build -o go-fish main.go

#: The GNU Readline REPL front-end to the go-interactive evaluator
go-fish-grl: repl_imports.go main_grl.go repl.go
	go build -o go-fish-grl main_grl.go

main.go: repl_imports.go

#: Subsidiary program to import packages into go-fish
make_env: make_env.go
	go build make_env.go

#: Recreate extracted imports
repl_imports.go: make_env
	./make_env > repl_imports.go

#: Check stuff
test: make_env.go
	go test -v

#: Same as test
check: test

clean:
	for file in make_env go-fish go-fish-grl ; do \
		[ -e $$file ] && rm $$file; \
	done

install:
	go install
	[ -x ./go-fish ] && cp ./go-fish $$GOBIN/go-fish
	[ -x ./go-fish-grl ] && cp ./go-fish $$GOBIN/go-fish-grl
