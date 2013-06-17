GOPATH=$(PWD)
export GOPATH
tests testpkgs : APP_ENV=test
export APP_ENV

.PHONY: all tests testpkgs bin/shotgun bin/{%% .Name %%} assets

all: tests testpkgs bin/shotgun bin/{%% .Name %%} assets

bin/{%% .Name %%}:
	echo $(GOPATH)
	go build -o bin/{%% .Name %%} {%% .Name %%}

bin/shotgun:
	echo $(GOPATH)
	go build -o bin/shotgun github.com/DanielHeath/shotgun-go

tests: testpkgs
	cd src && find {%% .Name %%} -iname *_test.go | xargs -n 1 dirname | sort | uniq | xargs -n 1 go test

testpkgs:
	cd src && find {%% .Name %%} -iname *_test.go | xargs -n 1 dirname | sort | uniq | xargs go test -i

clean:
	rm -rf pkg/* public/* bin/* log/*

assets:
	bundle exec rake assets
