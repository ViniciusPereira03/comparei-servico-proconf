#!/usr/bin/env bash

MONGO_HOST="localhost"
MONGO_PORT="27017"
MONGO_USER="root"
MONGO_PASS="proconfdb"
AUTH_DB="admin"
DB_NAME="proconfdb"

docker exec -it mongo-proconf mongosh -u root -p proconfdb --authenticationDatabase admin
