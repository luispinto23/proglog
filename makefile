compile:
	protoc api/v1/*proto --go_out=. --go_opts=paths=source_relative --proto_path .

test:
	go test -race ./...