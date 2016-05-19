all:
	docker run --rm \
		-v "$(shell pwd -P)/bin":/go/bin \
		-v "$(shell pwd -P)":/go/src/tagstore \
		golang:1.6.1-alpine \
		go build -o /go/bin/reg tagstore/cmd
