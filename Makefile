migrate-down:
	- migrate -database postgres://root:@localhost:26257/sso?sslmode=disable -path internal/constant/query/schemas -verbose down
migrate-up:
	- migrate -database postgres://root:@localhost:26257/sso?sslmode=disable -path internal/constant/query/schemas -verbose up