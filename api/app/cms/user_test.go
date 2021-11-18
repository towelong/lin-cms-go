package cms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jianfengye/collection"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/towelong/lin-cms-go/internal/config"
	"github.com/towelong/lin-cms-go/internal/domain/dto"
	"github.com/towelong/lin-cms-go/internal/domain/model"
	"github.com/towelong/lin-cms-go/internal/middleware"
	"github.com/towelong/lin-cms-go/internal/pkg/jwt"
	"github.com/towelong/lin-cms-go/internal/service/mock"
	"github.com/towelong/lin-cms-go/pkg/response"
	"github.com/towelong/lin-cms-go/pkg/token"
	"github.com/towelong/lin-cms-go/pkg/validator"
)

type testUserAPI struct {
	userAPI     *UserAPI
	userService *mockservice.MockIUserService
	handler     *gin.Engine
}

func initUserService(t *testing.T) *testUserAPI {
	validator.InitValidator()
	config.LoadConfig()
	ctrl := gomock.NewController(t)
	userService := mockservice.NewMockIUserService(ctrl)
	groupService := mockservice.NewMockIGroupService(ctrl)
	iToken := jwt.NewJWTMaker()
	u := &UserAPI{
		UserService: userService,
		JWT:         iToken,
		Auth: middleware.Auth{
			JWT:          iToken,
			UserService:  userService,
			GroupService: groupService,
		},
	}
	handler := gin.New()
	handler.Use(middleware.ErrorHandler)
	group := handler.Group("/cms")
	u.RegisterServer(group)
	return &testUserAPI{
		userAPI:     u,
		userService: userService,
		handler:     handler,
	}
}

func TestUserAPI_Login(t *testing.T) {
	service := initUserService(t)
	mockUsername := "root"
	// mockPassword := "123456"
	mockUser := model.User{
		Username: mockUsername,
		BaseModel: model.BaseModel{
			ID:         1,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		},
	}

	loginForm := dto.UserLoginDTO{
		Username: "root",
		Password: "123456",
	}

	loginForm2 := dto.UserLoginDTO{
		Username: "root",
		Password: "",
	}

	testcases := []struct {
		name       string
		param      dto.UserLoginDTO
		stub       func(userService *mockservice.MockIUserService)
		wantStatus int
		checkBody  func(t *testing.T, record *httptest.ResponseRecorder)
	}{
		{
			name:  "case1 输入正确的账号密码",
			param: loginForm,
			stub: func(userService *mockservice.MockIUserService) {
				userService.EXPECT().
					VerifyUser(gomock.Any(), gomock.Not("")).
					Times(1).
					Return(mockUser, nil)
			},
			wantStatus: 200,
			checkBody: func(t *testing.T, record *httptest.ResponseRecorder) {
				body := record.Body
				all, bodyErr := ioutil.ReadAll(body)
				require.NoError(t, bodyErr)
				var tokens token.Tokens
				err := json.Unmarshal(all, &tokens)
				require.NoError(t, err)
				require.NotEmpty(t, tokens.AccessToken)
				require.NotEmpty(t, tokens.RefreshToken)
				accessValid, accessErr := service.userAPI.JWT.VerifyAccessToken(tokens.AccessToken)
				refreshValid, refreshErr := service.userAPI.JWT.VerifyRefreshToken(tokens.RefreshToken)
				require.NoError(t, accessErr)
				require.NoError(t, refreshErr)
				require.Equal(t, mockUser.ID, accessValid.Identity)
				require.Equal(t, mockUser.ID, refreshValid.Identity)
			},
		},
		{
			name:  "case2 未输入密码",
			param: loginForm2,
			stub: func(userService *mockservice.MockIUserService) {
				userService.EXPECT().
					VerifyUser(gomock.Any(), gomock.Eq("")).
					Times(0)
			},
			wantStatus: 400,
			checkBody: func(t *testing.T, record *httptest.ResponseRecorder) {
				body := record.Body
				all, bodyErr := ioutil.ReadAll(body)
				require.NoError(t, bodyErr)
				var res response.Response
				err := json.Unmarshal(all, &res)
				require.NoError(t, err)
				require.Equal(t, int64(10030), res.Code)
			},
		},
		{
			name: "case3 未输入账号",
			param: dto.UserLoginDTO{
				Username: "",
				Password: "123123",
			},
			stub: func(userService *mockservice.MockIUserService) {
				userService.EXPECT().
					VerifyUser(gomock.Eq(""), gomock.Any()).
					Times(0)
			},
			wantStatus: 400,
			checkBody: func(t *testing.T, record *httptest.ResponseRecorder) {
				body := record.Body
				all, bodyErr := ioutil.ReadAll(body)
				require.NoError(t, bodyErr)
				var res response.Response
				err := json.Unmarshal(all, &res)
				require.NoError(t, err)
				require.Equal(t, int64(10030), res.Code)
			},
		},
		{
			name: "case4 未输入账号和密码",
			param: dto.UserLoginDTO{
				Username: "",
				Password: "",
			},
			stub: func(userService *mockservice.MockIUserService) {
				userService.EXPECT().
					VerifyUser(gomock.Eq(""), gomock.Eq("")).
					Times(0)
			},
			wantStatus: 400,
			checkBody: func(t *testing.T, record *httptest.ResponseRecorder) {
				body := record.Body
				all, bodyErr := ioutil.ReadAll(body)
				require.NoError(t, bodyErr)
				var res response.Response
				err := json.Unmarshal(all, &res)
				require.NoError(t, err)
				require.Equal(t, int64(10030), res.Code)
			},
		},
		{
			name: "case5 输入错误的账号",
			param: dto.UserLoginDTO{
				Username: "root1",
				Password: "123123",
			},
			stub: func(userService *mockservice.MockIUserService) {
				userService.EXPECT().
					VerifyUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(model.User{}, response.NewResponse(10031))
			},
			wantStatus: 400,
			checkBody: func(t *testing.T, record *httptest.ResponseRecorder) {
				body := record.Body
				all, bodyErr := ioutil.ReadAll(body)
				require.NoError(t, bodyErr)
				var res response.Response
				err := json.Unmarshal(all, &res)
				require.NoError(t, err)
				require.Equal(t, int64(10031), res.Code)
			},
		},
		{
			name: "case6 输入错误的密码",
			param: dto.UserLoginDTO{
				Username: "root",
				Password: "123123",
			},
			stub: func(userService *mockservice.MockIUserService) {
				userService.EXPECT().
					VerifyUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(model.User{}, response.NewResponse(10032))
			},
			wantStatus: 400,
			checkBody: func(t *testing.T, record *httptest.ResponseRecorder) {
				body := record.Body
				all, bodyErr := ioutil.ReadAll(body)
				require.NoError(t, bodyErr)
				var res response.Response
				err := json.Unmarshal(all, &res)
				require.NoError(t, err)
				require.Equal(t, int64(10032), res.Code)
			},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			testcase.stub(service.userService)
			b, err := json.Marshal(&testcase.param)
			require.NoError(t, err)
			uri := "/cms/user/login"
			request := httptest.NewRequest(http.MethodPost, uri, bytes.NewReader(b))
			w := httptest.NewRecorder()
			service.handler.ServeHTTP(w, request)
			require.Equal(t, testcase.wantStatus, w.Code)
			testcase.checkBody(t, w)
		})
	}
}

func TestName(t *testing.T) {
	var userGroups []model.UserGroup
	intColl := collection.NewIntCollection([]int{1, 2, 3})
	intColl.Map(func(item interface{}, key int) interface{} {
		groupId := item.(int)
		return &model.UserGroup{
			UserID:  1,
			GroupID: groupId,
		}
	}).ToObjs(&userGroups)
	fmt.Println(userGroups)
}
