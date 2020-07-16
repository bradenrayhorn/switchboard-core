test_report:
	go test -v -coverprofile=./reports/coverage.txt -covermode=atomic -coverpkg=./... ./...
	go tool cover -html=./reports/coverage.txt

