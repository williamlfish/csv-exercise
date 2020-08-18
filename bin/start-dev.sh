#!/usr/bin/env sh
set -e
G='\033[0;32m'
NC='\033[0m' # No Color

colorP() { 
    printf "\n\n$G $1 $NC \n\n" 
}
colorP "checking for env"
if test -f .env; then
    colorp "sourcing env"
    set -a; source .env
fi
colorP "building binary"

go build ./cmd/dirwatcher/ 


colorP "a little clean up"
docker-compose down

colorP "docker compose up!!"
docker-compose up -d --remove-orphans

colorP "let that db warm up"
sleep 5s

colorP "migrations baby!!"
./bin/local-migrate.sh

./dirwatcher $1 $2 $2 $4

