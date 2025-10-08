#!/bin/sh
set -e

# Wait for database to be ready
echo "Waiting for database to be ready..."
while ! nc -z db 5432; do
  sleep 1
  echo "Retrying database connection..."
done
echo "Database is ready!"

# Run database migrations
echo "Running database migrations..."
migrate -path /app/migrations -database "$DATABASE_URL" up

# Check if superadmin exists
echo "Checking for existing superadmin..."
ADMIN_CHECK=$(psql $DATABASE_URL -t -c "SELECT EXISTS (SELECT 1 FROM users WHERE roles @> '{SuperAdmin}' LIMIT 1);")
ADMIN_EXISTS=$(echo $ADMIN_CHECK | tr -d '[:space:]')
echo "Admin check result: $ADMIN_EXISTS"

# Display current users for debugging
echo "Current users in database:"
psql $DATABASE_URL -c "SELECT id, email, roles FROM users;"

# Seed superadmin if none exists
if [ "$ADMIN_EXISTS" = "f" ]; then
    echo "Creating superadmin account..."
    /app/seed --table users
    echo "Superadmin account created successfully"
else
    echo "Superadmin already exists, skipping creation"
fi

# Ensure uploads directory permissions
if [ ! -d "/app/uploads" ]; then
    mkdir -p /app/uploads
fi

# Start the application
echo "Starting the application..."
exec "$@" # This ensures proper signal forwarding