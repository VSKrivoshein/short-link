start-initial:
	docker-compose --env-file=.env --profile migrator up --build

start:
	docker-compose --env-file=.env up

test:
	docker-compose -f docker-compose.test.yaml --env-file=test.env --profile migrator up --build

down:
	docker-compose down --volumes && docker-compose down