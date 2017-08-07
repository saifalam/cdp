CC=cc

FLAGS=-O3

all: life

life: life.go
	#$(CC) $(FLAGS) life.c -o life
	go build -o life life.go

clean:
	rm -f life
