build-server:
	docker build ./ -t ghbypass-server:latest

run-server:
	docker run -it -p 8080:8080 ghbypass-server:latest
