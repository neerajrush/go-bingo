.PHONY: all

all: clean game.pb.go

clean:
	@echo "Cleaning files *.go"
	@rm -f *.go

game.pb.go: game.proto
	@echo "Compiling file $< to generate: $@"
	protoc -I../../../.. -I. $< --go_out=plugins=grpc:.
