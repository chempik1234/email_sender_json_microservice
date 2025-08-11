build:
	docker build ./service -f .\build\Dockerfile.emailer -t email_microservice
