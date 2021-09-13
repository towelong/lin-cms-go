mock:
	mockgen -package mockservice -destination internal/service/mock/group.go github.com/towelong/lin-cms-go/internal/service IGroupService
	mockgen -package mockservice -destination internal/service/mock/user.go github.com/towelong/lin-cms-go/internal/service IUserService
test:
	go test ./... -coverprofile=cover.txt
	go tool cover -html=cover.txt