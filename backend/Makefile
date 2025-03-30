db:
	docker compose up -d;
	PGPASSWORD=mycal psql -h localhost -p 5434 -U mycal -d mycal -f ./internal/database/schema.sql


