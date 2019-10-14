build:
	@echo "=============building Local API============="
	docker build -f Dockerfile -t main .

start: build
	@echo "=============starting api locally============="
	docker-compose up -d

logs:
	docker-compose logs -f

stop:
	docker-compose down

test:
	go test -v -cover ./...

clean: down
	@echo "=============cleaning up============="
	rm -f api
	docker system prune -f
	docker volume prune -f

start-dev:
	@echo "=============starting api in development mode============="
	compileDaemon -build="go build -o bin/harvard-arts-reverse-proxy ." -command="./bin/harvard-arts-reverse-proxy" -color -graceful-kill