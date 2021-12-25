all:
	cd src && go build
	mv ./src/taurine ./taurine

install:
	cd src && go install

runtests:
	cd src && go run ./test/test.go

clean:
	rm taurine
