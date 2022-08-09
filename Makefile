migrate-down:
	- migrate -database cockroachdb://root@localhost:26257/defaultdb?sslmode=disable -path internal/constant/query/schemas -verbose down
migrate-up:
	- migrate -database cockroachdb://root@localhost:26257/defaultdb?sslmode=disable -path internal/constant/query/schemas -verbose up
migrate-create:
	- migrate create -ext sql -dir internal/constant/query/schemas -seq create_users_table
go-test:
	- go test -v ./... -p=1 -count=1