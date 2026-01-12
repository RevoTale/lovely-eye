prod-example:
	cd docker && docker-compose --env-file .env.example -f docker-compose.yml up --build
test:
	cd ./server && make test && cd ../dashboard && bun run test