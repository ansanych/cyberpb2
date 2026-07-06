export PATH:=$(PATH):$(shell go env GOPATH)/bin

gen: gen-go gen-py

gen-go:
	protoc -I ./proto \
			--go_out ./proto \
			--go_opt paths=source_relative \
			--go-grpc_out ./proto \
			--go-grpc_opt paths=source_relative \
			./proto/cyber.proto
	protoc -I ./proto \
			--go_out ./proto \
			--go_opt paths=source_relative \
			--go-grpc_out ./proto \
			--go-grpc_opt paths=source_relative \
			./proto/cybermetrica.proto
	protoc -I ./proto \
			--go_out ./proto \
			--go_opt paths=source_relative \
			--go-grpc_out ./proto \
			--go-grpc_opt paths=source_relative \
			./proto/cyberfuel.proto

gen-py:
	python3 -m grpc_tools.protoc -I ./proto \
			--python_out ./proto \
			--grpc_python_out ./proto \
			./proto/cyber.proto
	python3 -m grpc_tools.protoc -I ./proto \
			--python_out ./proto \
			--grpc_python_out ./proto \
			./proto/cybermetrica.proto
	python3 -m grpc_tools.protoc -I ./proto \
			--python_out ./proto \
			--grpc_python_out ./proto \
			./proto/cyberfuel.proto