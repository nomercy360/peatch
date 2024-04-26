docker exec -it postgres psql -U postgres -d postgres -c 'DROP DATABASE IF EXISTS peatch WITH (FORCE)'
docker exec -it postgres psql -U postgres -d postgres -c 'CREATE DATABASE peatch;'
migrate -source file://${PWD}/scripts/migrations -database postgres://postgres:mysecretpassword@localhost:5432/peatch?sslmode=disable up
cd tests && npx ts-mocha test.ts --timeout=10000
cd ../..
