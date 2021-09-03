mock:
	mockgen -package mockservice -destination internal/service/mock/group.go github.com/towelong/lin-cms-go/internal/service IGroupService
test:
	go test ./... -coverprofile=cover.txt
	go tool cover -html=cover.txt