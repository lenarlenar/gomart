package app

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert/v2"
	mocks "github.com/lenarlenar/gomart/internal/interfaces/mocks"
	"github.com/lenarlenar/gomart/internal/models"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	POST        = "POST"
	registerURL = "/api/user/register"
	loginURL    = "/api/user/login"
)

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authService := mocks.NewMockAuthService(ctrl)
	authStorage := mocks.NewMockAuthStorage(ctrl)
	jwtService := mocks.NewMockJWTService(ctrl)
	app := App{
		AuthStorage: authStorage,
		AuthService: authService,
		JWTService:  jwtService,
	}

	testServer := httptest.NewServer(
		app.SetupRouter(),
	)
	defer testServer.Close()

	fullRegisterURL := testServer.URL + registerURL
	testCases := []struct {
		testName     string
		body         func() io.Reader
		prepare      func()
		expectedCode int
	}{
		{
			testName:     "нет тела запроса",
			expectedCode: http.StatusBadRequest,
		},
		{
			testName: "нет логина",
			body: func() io.Reader {
				userRequest := models.UserRequest{Password: "123"}
				data, _ := json.Marshal(userRequest)
				return bytes.NewBuffer(data)
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			testName: "нет пароля",
			body: func() io.Reader {
				userRequest := models.UserRequest{Login: "user"}
				data, _ := json.Marshal(userRequest)
				return bytes.NewBuffer(data)
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			testName: "пользователь уже зарегистрирован",
			prepare: func() {
				userRequest := models.UserRequest{Login: "user", Password: "123"}
				jwtService.EXPECT().Generate(userRequest.Login).Return("token", nil)
				authService.EXPECT().Register(userRequest).Return(models.ErrDuplicateUser)
			},
			body: func() io.Reader {
				userRequest := models.UserRequest{Login: "user", Password: "123"}
				data, _ := json.Marshal(userRequest)
				return bytes.NewBuffer(data)
			},
			expectedCode: http.StatusConflict,
		},
		{
			testName: "успешная регистрация пользователя",
			prepare: func() {
				userRequest := models.UserRequest{Login: "user", Password: "123"}
				jwtService.EXPECT().Generate(userRequest.Login).Return("token", nil)
				authService.EXPECT().Register(userRequest).Return(nil)
			},
			body: func() io.Reader {
				userRequest := models.UserRequest{Login: "user", Password: "123"}
				data, _ := json.Marshal(userRequest)
				return bytes.NewBuffer(data)
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			var body io.Reader
			if tc.body != nil {
				body = tc.body()
			}

			if tc.prepare != nil {
				tc.prepare()
			}

			req, err := http.NewRequest(POST, fullRegisterURL, body)
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp, err := testServer.Client().Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()
			assert.Equal(t, tc.expectedCode, resp.StatusCode)
		})
	}
}

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authService := mocks.NewMockAuthService(ctrl)
	authStorage := mocks.NewMockAuthStorage(ctrl)
	jwtService := mocks.NewMockJWTService(ctrl)
	app := App{
		AuthStorage: authStorage,
		AuthService: authService,
		JWTService:  jwtService,
	}

	testServer := httptest.NewServer(
		app.SetupRouter(),
	)
	defer testServer.Close()

	fullLoginURL := testServer.URL + loginURL
	testCases := []struct {
		testName     string
		body         func() io.Reader
		prepare      func()
		expectedCode int
	}{
		{
			testName:     "нет тела запроса",
			expectedCode: http.StatusBadRequest,
		},
		{
			testName: "нет логина",
			body: func() io.Reader {
				userRequest := models.UserRequest{Password: "123"}
				data, _ := json.Marshal(userRequest)
				return bytes.NewBuffer(data)
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			testName: "нет пароля",
			body: func() io.Reader {
				userRequest := models.UserRequest{Login: "user"}
				data, _ := json.Marshal(userRequest)
				return bytes.NewBuffer(data)
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			testName: "неверный логин или пароль",
			prepare: func() {
				userRequest := models.UserRequest{Login: "user", Password: "123"}
				authService.EXPECT().Login(userRequest).Return(models.ErrPasswordOrUsernameIsIncorrect)
			},
			body: func() io.Reader {
				userRequest := models.UserRequest{Login: "user", Password: "123"}
				data, _ := json.Marshal(userRequest)
				return bytes.NewBuffer(data)
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			testName: "успешный login",
			prepare: func() {
				userRequest := models.UserRequest{Login: "user", Password: "123"}
				jwtService.EXPECT().Generate(userRequest.Login).Return("token", nil)
				authService.EXPECT().Login(userRequest).Return(nil)
			},
			body: func() io.Reader {
				userRequest := models.UserRequest{Login: "user", Password: "123"}
				data, _ := json.Marshal(userRequest)
				return bytes.NewBuffer(data)
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			var body io.Reader
			if tc.body != nil {
				body = tc.body()
			}

			if tc.prepare != nil {
				tc.prepare()
			}

			req, err := http.NewRequest(POST, fullLoginURL, body)
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp, err := testServer.Client().Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()
			assert.Equal(t, tc.expectedCode, resp.StatusCode)
		})
	}
}
