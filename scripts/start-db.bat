@echo off
docker-compose up -d
timeout /t 5
echo PostgreSQL container is running 