msgidbind:
	@echo "build install suituphandler"
	@GOOS=windows GOPLATFORM=amd64 go build -o bin/ ./suituphandler
	@go install ./suituphandler
	@go run ./suituphandler/clientgen


asyncdb:
	go install ./dbgenerator/asyncdb.go


handler:
	go install ./suituphandler


install: asyncdb handler


.PHONY: msgidbind asyncdb handler install