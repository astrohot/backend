APP_NAME = astrohot/api

### Docker Tasks ###
build: # Build the container
	docker build -t $(APP_NAME) .

run: # Run the container using default port
	docker container run -p 8080:8080 --env-file .env $(APP_NAME)

up: build run # Build and run the container
