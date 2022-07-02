all:
	go build ./cmd/taurine

install:
	go install ./cmd/taurine

test: FORCE
	go run ./test/test.go

gentests:
	go run ./test/test.go --gen-all

clean:
	rm taurine

FORCE: ;