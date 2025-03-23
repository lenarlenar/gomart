package app

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"github.com/golang-jwt/jwt/v5"
	mocks "github.com/lenarlenar/gomart/internal/interfaces/mocks"
	"github.com/lenarlenar/gomart/internal/models"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	GET       = "GET"
	ordersURL = "/api/user/orders"
)

func TestCreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authService := mocks.NewMockAuthService(ctrl)
	authStorage := mocks.NewMockAuthStorage(ctrl)
	jwtService := mocks.NewMockJWTService(ctrl)
	orderService := mocks.NewMockOrdersService(ctrl)
	accrualService := mocks.NewMockAccrualService(ctrl)
	app := App{
		AuthStorage:    authStorage,
		AuthService:    authService,
		JWTService:     jwtService,
		OrdersService:  orderService,
		AccrualService: accrualService,
	}

	testServer := httptest.NewServer(
		app.SetupRouter(),
	)
	defer testServer.Close()

	fullURL := testServer.URL + ordersURL
	testCases := []struct {
		testName     string
		targetURL    string
		body         func() io.Reader
		prepare      func()
		expectedCode int
	}{
		{
			testName: "должен создать заказ",
			prepare: func() {
				jwtToken := jwt.NewWithClaims(
					jwt.SigningMethodHS256,
					jwt.MapClaims{
						"sub": "login",
					})
				user := models.User{ID: "user-id", Login: "user", HashPass: "hash"}
				jwtService.EXPECT().Validate("token").Return(jwtToken, nil)
				orderService.EXPECT().Check("order-id").Return(true)
				orderService.EXPECT().CreateOrder("order-id", "user-id").Return(nil)
				authStorage.EXPECT().GetUser("login").Return(&user, nil)
				accrualService.EXPECT().Calculate("order-id")
			},
			body: func() io.Reader {
				return bytes.NewBuffer([]byte("order-id"))
			},
			expectedCode: http.StatusAccepted,
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

			req, err := http.NewRequest(POST, fullURL, body)
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer token")

			resp, err := testServer.Client().Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedCode, resp.StatusCode)
		})
	}
}

func TestGetOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authService := mocks.NewMockAuthService(ctrl)
	authStorage := mocks.NewMockAuthStorage(ctrl)
	jwtService := mocks.NewMockJWTService(ctrl)
	orderService := mocks.NewMockOrdersService(ctrl)
	accrualService := mocks.NewMockAccrualService(ctrl)
	app := App{
		AuthStorage:    authStorage,
		AuthService:    authService,
		JWTService:     jwtService,
		OrdersService:  orderService,
		AccrualService: accrualService,
	}

	testServer := httptest.NewServer(
		app.SetupRouter(),
	)
	defer testServer.Close()

	fullURL := testServer.URL + ordersURL
	testCases := []struct {
		testName     string
		method       string
		prepare      func()
		expectedCode int
		expectedRes  string
	}{
		{
			testName: "вернуть список заказов",
			prepare: func() {
				jwtToken := jwt.NewWithClaims(
					jwt.SigningMethodHS256,
					jwt.MapClaims{
						"sub": "user",
					})

				// Определяем пользователя
				user := models.User{ID: "user-id", Login: "user", HashPass: "hash"}
				authStorage.EXPECT().GetUser("user").Return(&user, nil)
				jwtService.EXPECT().Validate("token").Return(jwtToken, nil)
				orderService.EXPECT().GetOrders("user-id").Return([]models.Order{
					{
						ID:         "order-id",
						Status:     "StatusNew",
						Accrual:    nil,
						UploadedAt: time.Date(2009, 11, 17, 0, 0, 0, 0, time.UTC),
					},
				}, nil)
			},
			expectedCode: http.StatusOK,
			expectedRes:  `[{"number":"order-id","status":"StatusNew","uploaded_at":"2009-11-17T00:00:00Z"}]`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			if tc.prepare != nil {
				tc.prepare()
			}

			req, err := http.NewRequest(GET, fullURL, nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer token")

			resp, err := testServer.Client().Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			responseBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedCode, resp.StatusCode)
			assert.Equal(t, tc.expectedRes, string(responseBody))
		})
	}
}
