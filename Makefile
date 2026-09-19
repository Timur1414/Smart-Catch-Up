EXCLUDE_DIRS := grep -v '/mocks' | grep -v '/migrations' | grep -v '/domain' | grep -vE '/pkg$$' | grep -vE 'github.com/Timur1414/Smart-Catch-Up$$'

.PHONY: deps
deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.0
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
	#go install github.com/mailru/easyjson/...@v0.9.2
