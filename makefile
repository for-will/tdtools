msgidbind:
	@echo "build install suituphandler"
	@GOOS=windows GOPLATFORM=amd64 go build -o bin/ ./suituphandler
	@go install ./suituphandler
	@go run ./suituphandler/clientgen
