GOCMD=go
GORUN=$(GOCMD) run
GOMOD= $(GOCMD) mod
GOTIDY=$(GOMOD) tidy
GOLINT=$(GOCMD) lint

deps :
	${GOTIDY}

postup :
	docker run --name postdev -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres

postdown:
	docker container stop postdev && docker container rm postdev

redsup :
	docker run --name redisdev -p 6379:6379 -d redis

redsdown :
	docker container stop redisdev && docker container rm redisdev

dbup :
	docker exec -it postdev createdb --username=root --owner=root pickablog_dev

dbdown :
	docker exec -it postdev dropdb pickablog_dev

STATE?=up
migrate :
	migrate -database "postgres://root:root@localhost:5432/pickablog_dev?sslmode=disable" -verbose -path db/migration ${STATE}

MODE?=local
run :
	${GORUN} main.go ${MODE}

.PHONY : deps postup postdown dbup dbdown redsup redsdown migrate run