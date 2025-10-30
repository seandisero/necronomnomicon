if [ -f .env ]; then
	source .env
fi 
cd sql/schema
goose turso "$DB_URL?authToken=$DB_TOKEN" down
