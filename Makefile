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

dbup :
	docker exec -it postdev createdb --username=root --owner=root pickablog_dev

dbdown :
	docker exec -it postdev dropdb pickablog_dev

STATE?=up
migrate :
	migrate -database "postgres://root:root@localhost:5432/pickablog_dev?sslmode=disable" -verbose -path db/migration ${STATE}

run :
	${GORUN} main.go

.PHONY : deps postup postdown dbup dbdown migrate run