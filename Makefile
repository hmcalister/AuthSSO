EXEC_FILE := AuthSSO

clean:
	rm logs/*
	rm $(EXEC_FILE)
	rm database/database/*

build:
	go build .

dev: build
	./$(EXEC_FILE -debug)

run: build
	./$(EXEC_FILE)