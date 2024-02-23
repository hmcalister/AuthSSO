EXEC_FILE := AuthSSO

clean:
	rm logs/*
	rm $(EXEC_FILE)
	rm database/database/*

sqlcGenerate:
	cd database; sqlc generate

build:
	go build .

all: sqlcGenerate build

dev: build
	./$(EXEC_FILE -debug)

run: build
	./$(EXEC_FILE)