gometalinter ./... || exit 1
go test ./... -covermode=atomic || exit 1
