goose -dir ./migrations postgres "postgres://postgres:password@localhost:5432/oms?sslmode=disable" status

goose -dir ./migrations postgres "postgres://postgres:password@localhost:5432/oms?sslmode=disable" up