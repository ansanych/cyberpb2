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
			--python_out ./python/cyberpb2 \
			--grpc_python_out ./python/cyberpb2 \
			./proto/cyber.proto
	python3 -m grpc_tools.protoc -I ./proto \
			--python_out ./python/cyberpb2 \
			--grpc_python_out ./python/cyberpb2 \
			./proto/cybermetrica.proto
	python3 -m grpc_tools.protoc -I ./proto \
			--python_out ./python/cyberpb2 \
			--grpc_python_out ./python/cyberpb2 \
			./proto/cyberfuel.proto
	# Fix imports in generated _grpc files for package-relative access
	sed -i 's/^import cyber_pb2 as/from . import cyber_pb2 as/' ./python/cyberpb2/*_pb2_grpc.py