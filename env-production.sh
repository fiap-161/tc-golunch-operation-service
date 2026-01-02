#!/bin/bash

# Production Service Environment Variables
export PRODUCTION_SERVICE_PORT=8083
export PRODUCTION_SERVICE_DB_HOST=localhost
export PRODUCTION_SERVICE_DB_PORT=5434
export PRODUCTION_SERVICE_DB_USER=golunch_prod
export PRODUCTION_SERVICE_DB_PASSWORD=golunch_prod123
export OPERATION_SERVICE_DB_NAME=golunch_operation

# Database URL para GORM
export DATABASE_URL="host=localhost user=golunch_oper password=golunch_oper123 dbname=golunch_operation port=5434 sslmode=disable TimeZone=America/Sao_Paulo"

# Order Service URL for HTTP communication
export CORE_SERVICE_URL=http://localhost:8081

# Payment Service URL for HTTP communication  
export PAYMENT_SERVICE_URL=http://localhost:8082

echo "Production Service environment variables set:"
echo "PORT: $PRODUCTION_SERVICE_PORT"
echo "DB: $PRODUCTION_SERVICE_DB_HOST:$PRODUCTION_SERVICE_DB_PORT/$PRODUCTION_SERVICE_DB_NAME"
echo "Core Service: $CORE_SERVICE_URL"
echo "Payment Service: $PAYMENT_SERVICE_URL"