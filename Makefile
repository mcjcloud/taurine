all:
	cd src && go build
	mv ./src/taurine ./taurine

install:
	cd src && go install

runtests:
	cd src && go run ./test/test.go

gentests:
	cd src && go run ./test/test.go --gen-all

clean:
	rm taurine
