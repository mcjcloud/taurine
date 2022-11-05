all:
	go build

install:
	go install

test: FORCE
	go run ./test/test.go

gentests:
	go run ./test/test.go --gen-all

clean:
	rm taurine *.ll out.o out.s *.out

FORCE: ;