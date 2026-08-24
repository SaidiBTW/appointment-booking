docker run -d \
  --name my-postgres \
  -p 5432:5432 \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=appointment_booking \
  -v $(pwd)/api/init.sql:/docker-entrypoint-initdb.d/init.sql \
  postgres:16