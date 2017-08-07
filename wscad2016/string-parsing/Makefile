.PHONY = all clean test mpp

all : spec

clean :
	rm -rf spec *dSYM*

spec : spec.go
	#g++ -std=c++11 -Wall -Wextra -O3 spec.cc -o spec
	go build -o spec spec.go

