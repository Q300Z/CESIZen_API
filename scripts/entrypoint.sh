#!/bin/sh

echo "🔄 Waiting for the database to be ready..."
echo "Ping database : $DB_HOST:$DB_PORT" 
until nc -z -v -w30 "$DB_HOST" "$DB_PORT"
do
  echo "⏳ Waiting for database connection at $DB_HOST:$DB_PORT..."
  sleep 2
done

echo "🚀 Running Prisma migration (db push)..."
go run github.com/steebchen/prisma-client-go db push --schema internal/database/prisma/schema.prisma


if [ $? -ne 0 ]; then
  echo "❌ Migration failed. Exiting..."
  exit 1
fi

echo "✅ Migration completed successfully."

# Check if in development mode for start with air or without air
if [ "$GIN_MODE" = "debug" ]; then
  echo "🚀 Starting application with air..."
  air -c .air.toml
else
  echo "🚀 Starting application in production mode..."
  ./app
fi
