CC=gcc

FLAGS=-O3 -lm

all: histogram

histogram: histogram.go
	#$(CC) histogram.c $(FLAGS)  -o histogram
	go build -o histogram histogram.go

clean:
	rm -f histogram
