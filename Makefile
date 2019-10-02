app_name = astrohot/api

### docker tasks ###
build: # build the container
	docker build -t $(app_name) .

run: # run the container using default port
	docker container run -p 8080:8080 --env-file .env $(app_name)

up: build run # build and run the container
