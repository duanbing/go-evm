all: 
	#GDB enabled
	go build -gcflags "-N -l" -ldflags -s main.go
.PHONY : all
