package biz

import (
	"coauth/internal/entities"
	"coauth/internal/model"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 模拟 AuthorizeRepo 接口
type MockAuthorizeRepo struct {
	FindClientByIDFunc func(ctx context.Context, clientID string) (*entities.App, error)
	FindAuthCodeFunc   func(ctx context.Context, code string) (*entities.AuthCode, error)
	SaveAuthCodeFunc   func(ctx context.Context, authCode *entities.AuthCode) error
	SaveTokenFunc      func(ctx context.Context, token *entities.Token) error
	SaveClientFunc     func(ctx context.Context, client *entities.App) error
}

func (m *MockAuthorizeRepo) FindClientByID(ctx context.Context, clientID string) (*entities.App, error) {
	return m.FindClientByIDFunc(ctx, clientID)
}

func (m *MockAuthorizeRepo) FindAuthCode(ctx context.Context, code string) (*entities.AuthCode, error) {
	return m.FindAuthCodeFunc(ctx, code)
}

func (m *MockAuthorizeRepo) SaveAuthCode(ctx context.Context, authCode *entities.AuthCode) error {
	return m.SaveAuthCodeFunc(ctx, authCode)
}

func (m *MockAuthorizeRepo) SaveToken(ctx context.Context, token *entities.Token) error {
	return m.SaveTokenFunc(ctx, token)
}

func (m *MockAuthorizeRepo) FindToken(ctx context.Context, accessToken string) (*entities.Token, error) {
	return nil, nil
}

func (m *MockAuthorizeRepo) SaveClient(ctx context.Context, client *entities.App) error {
	return m.SaveClientFunc(ctx, client)
}

func TestAuthorizeUseCase_Authorize(t *testing.T) {
	Convey("Test AuthorizeUseCase Authorize method", t, func() {
		// 创建模拟的 AuthorizeRepo
		mockRepo := &MockAuthorizeRepo{}
		log, _ := zap.NewDevelopment()
		uc := NewAuthorizeUseCase(mockRepo, log)

		// 准备测试数据
		clientID := uuid.New().String()
		req := &model.AuthorizeRequest{
			ClientID:    clientID,
			RedirectURI: "",
			State:       "test-state",
		}
		client := &entities.App{
			ID:          clientID,
			Secret:      "test-secret",
			RedirectUrl: "http://example.com/callback",
		}

		Convey("When client exists and save auth code succeeds", func() {
			// 设置模拟方法的行为
			mockRepo.FindClientByIDFunc = func(ctx context.Context, clientID string) (*entities.App, error) {
				return client, nil
			}
			mockRepo.SaveAuthCodeFunc = func(ctx context.Context, authCode *entities.AuthCode) error {
				return nil
			}

			// 调用 Authorize 方法
			reply, err := uc.Authorize(context.Background(), req)

			// 验证结果
			So(err, ShouldBeNil)
			So(reply.Code, ShouldNotBeEmpty)
			So(reply.State, ShouldEqual, req.State)
			So(reply.RedirectURI, ShouldNotBeEmpty)
		})

		Convey("When client does not exist", func() {
			// 设置模拟方法的行为
			mockRepo.FindClientByIDFunc = func(ctx context.Context, clientID string) (*entities.App, error) {
				return nil, status.Errorf(codes.NotFound, "client not found")
			}

			// 调用 Authorize 方法
			reply, err := uc.Authorize(context.Background(), req)

			// 验证结果
			So(reply, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "client not found")
		})

		Convey("When save auth code fails", func() {
			// 设置模拟方法的行为
			mockRepo.FindClientByIDFunc = func(ctx context.Context, clientID string) (*entities.App, error) {
				return client, nil
			}
			mockRepo.SaveAuthCodeFunc = func(ctx context.Context, authCode *entities.AuthCode) error {
				return status.Errorf(codes.Internal, "failed to save auth code")
			}

			// 调用 Authorize 方法
			reply, err := uc.Authorize(context.Background(), req)

			// 验证结果
			So(reply, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "failed to save auth code")
		})
	})
}

func TestAuthorizeUseCase_Token(t *testing.T) {
	Convey("Test AuthorizeUseCase Token method", t, func() {
		// 创建模拟的 AuthorizeRepo
		mockRepo := &MockAuthorizeRepo{}
		log, _ := zap.NewDevelopment()
		uc := NewAuthorizeUseCase(mockRepo, log)

		// 准备测试数据
		clientID := uuid.New().String()
		clientSecret := uuid.New().String()
		code := uuid.New().String()
		_ = uuid.New().String()
		authCode := &entities.AuthCode{
			Code:        code,
			ClientID:    clientID,
			UserID:      "test-user",
			RedirectURL: "http://example.com/callback",
			ExpiresAt:   time.Now().Add(time.Minute * 10),
		}
		client := &entities.App{
			ID:     clientID,
			Secret: clientSecret,
		}

		Convey("When all conditions are met", func() {
			// 设置模拟方法的行为
			mockRepo.FindClientByIDFunc = func(ctx context.Context, clientID string) (*entities.App, error) {
				return client, nil
			}
			mockRepo.FindAuthCodeFunc = func(ctx context.Context, code string) (*entities.AuthCode, error) {
				return authCode, nil
			}
			mockRepo.SaveTokenFunc = func(ctx context.Context, token *entities.Token) error {
				return nil
			}

			req := &model.TokenRequest{
				GrantType:    "authorization_code",
				Code:         code,
				ClientID:     clientID,
				ClientSecret: clientSecret,
				DataType:     "json",
			}

			// 调用 Token 方法
			reply, err := uc.Token(context.Background(), req)

			// 验证结果
			So(err, ShouldBeNil)
			So(reply.AccessToken, ShouldNotBeEmpty)
			So(reply.TokenType, ShouldEqual, "Bearer")
			So(reply.ExpiresIn, ShouldEqual, int64(3600))
			So(reply.RefreshToken, ShouldNotBeEmpty)
		})

		Convey("When invalid grant type", func() {
			req := &model.TokenRequest{
				GrantType:    "invalid_grant_type",
				Code:         code,
				ClientID:     clientID,
				ClientSecret: clientSecret,
				DataType:     "json",
			}

			// 调用 Token 方法
			reply, err := uc.Token(context.Background(), req)

			// 验证结果
			So(reply, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "invalid grant type")
		})

		Convey("When invalid client id", func() {
			// 设置模拟方法的行为
			mockRepo.FindClientByIDFunc = func(ctx context.Context, clientID string) (*entities.App, error) {
				return nil, status.Errorf(codes.Unauthenticated, "invalid client id")
			}

			req := &model.TokenRequest{
				GrantType:    "authorization_code",
				Code:         code,
				ClientID:     clientID,
				ClientSecret: clientSecret,
				DataType:     "json",
			}

			// 调用 Token 方法
			reply, err := uc.Token(context.Background(), req)

			// 验证结果
			So(reply, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "invalid client id")
		})

		Convey("When invalid client credentials", func() {
			// 设置模拟方法的行为
			mockRepo.FindClientByIDFunc = func(ctx context.Context, clientID string) (*entities.App, error) {
				return &entities.App{
					ID:     clientID,
					Secret: "wrong_secret",
				}, nil
			}

			req := &model.TokenRequest{
				GrantType:    "authorization_code",
				Code:         code,
				ClientID:     clientID,
				ClientSecret: clientSecret,
				DataType:     "json",
			}

			// 调用 Token 方法
			reply, err := uc.Token(context.Background(), req)

			// 验证结果
			So(reply, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "invalid client credentials")
		})

		Convey("When invalid authorization code", func() {
			// 设置模拟方法的行为
			mockRepo.FindClientByIDFunc = func(ctx context.Context, clientID string) (*entities.App, error) {
				return client, nil
			}
			mockRepo.FindAuthCodeFunc = func(ctx context.Context, code string) (*entities.AuthCode, error) {
				return nil, status.Errorf(codes.Unauthenticated, "invalid authorization code")
			}

			req := &model.TokenRequest{
				GrantType:    "authorization_code",
				Code:         code,
				ClientID:     clientID,
				ClientSecret: clientSecret,
				DataType:     "json",
			}

			// 调用 Token 方法
			reply, err := uc.Token(context.Background(), req)

			// 验证结果
			So(reply, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "invalid authorization code")
		})

		Convey("When authorization code expired", func() {
			// 设置模拟方法的行为
			mockRepo.FindClientByIDFunc = func(ctx context.Context, clientID string) (*entities.App, error) {
				return client, nil
			}
			expiredAuthCode := *authCode
			expiredAuthCode.ExpiresAt = time.Now().Add(-time.Minute * 10)
			mockRepo.FindAuthCodeFunc = func(ctx context.Context, code string) (*entities.AuthCode, error) {
				return &expiredAuthCode, nil
			}

			req := &model.TokenRequest{
				GrantType:    "authorization_code",
				Code:         code,
				ClientID:     clientID,
				ClientSecret: clientSecret,
				DataType:     "json",
			}

			// 调用 Token 方法
			reply, err := uc.Token(context.Background(), req)

			// 验证结果
			So(reply, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "authorization code expired")
		})

		Convey("When failed to save token", func() {
			// 设置模拟方法的行为
			mockRepo.FindClientByIDFunc = func(ctx context.Context, clientID string) (*entities.App, error) {
				return client, nil
			}
			mockRepo.FindAuthCodeFunc = func(ctx context.Context, code string) (*entities.AuthCode, error) {
				return authCode, nil
			}
			mockRepo.SaveTokenFunc = func(ctx context.Context, token *entities.Token) error {
				return status.Errorf(codes.Internal, "failed to save token")
			}

			req := &model.TokenRequest{
				GrantType:    "authorization_code",
				Code:         code,
				ClientID:     clientID,
				ClientSecret: clientSecret,
				DataType:     "json",
			}

			// 调用 Token 方法
			reply, err := uc.Token(context.Background(), req)

			// 验证结果
			So(reply, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "failed to save token")
		})
	})
}

func TestAuthorizeUseCase_AddClient(t *testing.T) {
	Convey("Test AuthorizeUseCase AddClient method", t, func() {
		// 创建模拟的 AuthorizeRepo
		mockRepo := &MockAuthorizeRepo{}
		log, _ := zap.NewDevelopment()
		uc := NewAuthorizeUseCase(mockRepo, log)

		// 准备测试数据
		req := &model.AddClientRequest{
			ClientName:    "TestClient",
			CallbackURL:   "http://example.com/callback",
			HomePage:      "http://example.com",
			ClientSummary: "This is a test client",
			Scopes:        []string{"scope1", "scope2"},
			UserID:        "testUserID",
		}

		Convey("When save client succeeds", func() {
			// 设置模拟方法的行为
			mockRepo.SaveClientFunc = func(ctx context.Context, client *entities.App) error {
				return nil
			}

			// 调用 AddClient 方法
			reply, err := uc.AddClient(context.Background(), req)

			// 验证结果
			So(err, ShouldBeNil)
			So(reply.ClientID, ShouldNotBeEmpty)
		})

		Convey("When save client fails", func() {
			// 设置模拟方法的行为
			mockRepo.SaveClientFunc = func(ctx context.Context, client *entities.App) error {
				return status.Errorf(codes.Internal, "failed to save client")
			}

			// 调用 AddClient 方法
			reply, err := uc.AddClient(context.Background(), req)

			// 验证结果
			So(reply, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "failed to save client")
		})
	})
}
