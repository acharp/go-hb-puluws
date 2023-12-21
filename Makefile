build::
	GOOS=linux GOARCH=amd64 go build -o bootstrap ./handler/handler.go
	zip -j ./handler/handler.zip bootstrap

clean::
	rm bootstrap
	rm ./handler/handler.zip
