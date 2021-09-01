.PHONY: build clean 

build: clean
	export GO111MODULE=on
	env CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bin/register lambdas/register/register.go
	env CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bin/login lambdas/login/login.go
	env CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bin/post_transaction lambdas/post_transaction/post_transaction.go
	env CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bin/get_transactions lambdas/get_transactions/get_transactions.go
	env CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bin/send_notification lambdas/send_notification/send_notification.go
	env CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bin/event_processor lambdas/event_processor/event_processor.go

clean:
	rm -rf ./bin
