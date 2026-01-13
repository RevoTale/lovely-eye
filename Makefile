prod-example:
	cd docker && docker-compose --env-file .env.example -f docker-compose.yml up --build

test:
	@echo "Running dashboard linting and type checking..."
	cd ./dashboard && bun run lint && bun run tsc
	@echo "\nRunning Go unit tests..."
	cd ./server && DASHBOARD_PATH=../dashboard/dist go test -v ./... --count=2