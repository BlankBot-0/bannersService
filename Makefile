PACKAGE=banners/cmd/banners
CURDIR=$(shell pwd)

build:
	go build -o ${CURDIR}/bin/app ${PACKAGE}

test:
	go test

run-all:
	docker-compose up --force-recreate --build