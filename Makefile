prod-example:
	cd docker && docker-compose --env-file .env.example -f docker-compose.yml up --build