package jwt_secret

import (
	"strings"
	"testing"
	"time"
)

const (
	testSecretKey = "test_secret_key_for_jwt_12345"
	testUserID    = uint(1)
	testAdminID   = uint(2)
)

// setupTest 測試前置設置
func setupTest(t *testing.T) {
	t.Helper()
	SetJwtSecret(testSecretKey)
}

func TestSetJwtSecret(t *testing.T) {
	tests := []struct {
		name   string
		secret string
	}{
		{"正常密鑰", "my_secret_key"},
		{"短密鑰", "123"},
		{"長密鑰", strings.Repeat("a", 100)},
		{"空密鑰", ""}, // 空密鑰也能生成，但安全性低
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetJwtSecret(tt.secret)
			// 驗證是否能夠正常生成 Token（間接驗證設置成功）
			token, err := GenerateToken(LoginUser, 1)
			if err != nil {
				t.Errorf("Failed to generate token with secret '%s': %v", tt.name, err)
			}
			if token == "" {
				t.Error("Expected non-empty token")
			}
			t.Logf("Secret: %s, Token length: %d", tt.name, len(token))
		})
	}
}

func TestGenerateToken(t *testing.T) {
	setupTest(t)

	tests := []struct {
		name      string
		role      LoginRole
		loginID   uint
		wantErr   bool
		errSubstr string
	}{
		{
			name:    "用戶角色_有效ID",
			role:    LoginUser,
			loginID: testUserID,
			wantErr: false,
		},
		{
			name:    "管理員角色_有效ID",
			role:    LoginAdmin,
			loginID: testAdminID,
			wantErr: false,
		},
		{
			name:    "用戶角色_ID為0",
			role:    LoginUser,
			loginID: 0,
			wantErr: false, // ID 為 0 也應該允許
		},
		{
			name:      "無效角色",
			role:      "invalid_role",
			loginID:   1,
			wantErr:   true,
			errSubstr: "LoginRole not allow",
		},
		{
			name:      "空角色",
			role:      "",
			loginID:   1,
			wantErr:   true,
			errSubstr: "LoginRole not allow",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.role, tt.loginID)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got nil")
					return
				}
				if tt.errSubstr != "" && !strings.Contains(err.Error(), tt.errSubstr) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errSubstr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if token == "" {
				t.Error("Expected non-empty token")
				return
			}

			// 驗證 Token 格式（JWT 應該有 3 個部分，用 . 分隔）
			parts := strings.Split(token, ".")
			if len(parts) != 3 {
				t.Errorf("Expected JWT to have 3 parts, got %d", len(parts))
			}

			t.Logf("Generated token length: %d", len(token))
		})
	}
}

func TestParseToken(t *testing.T) {
	setupTest(t)

	t.Run("解析用戶Token", func(t *testing.T) {
		// 生成用戶 Token
		token, err := GenerateToken(LoginUser, testUserID)
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}

		// 解析 Token
		claims, err := ParseToken(token)
		if err != nil {
			t.Fatalf("Failed to parse token: %v", err)
		}

		// 驗證 Claims
		if claims.UserID != testUserID {
			t.Errorf("Expected UserID=%d, got %d", testUserID, claims.UserID)
		}
		if claims.AdminID != 0 {
			t.Errorf("Expected AdminID=0, got %d", claims.AdminID)
		}

		// 驗證 MapClaims
		if iss, ok := claims.MapClaims["iss"].(string); !ok || iss != "gin-blog" {
			t.Errorf("Expected issuer='gin-blog', got %v", claims.MapClaims["iss"])
		}
	})

	t.Run("解析管理員Token", func(t *testing.T) {
		// 生成管理員 Token
		token, err := GenerateToken(LoginAdmin, testAdminID)
		if err != nil {
			t.Fatalf("Failed to generate admin token: %v", err)
		}

		// 解析 Token
		claims, err := ParseToken(token)
		if err != nil {
			t.Fatalf("Failed to parse admin token: %v", err)
		}

		// 驗證 Claims
		if claims.AdminID != testAdminID {
			t.Errorf("Expected AdminID=%d, got %d", testAdminID, claims.AdminID)
		}
		if claims.UserID != 0 {
			t.Errorf("Expected UserID=0, got %d", claims.UserID)
		}
	})

	t.Run("立即生成和解析", func(t *testing.T) {
		// 測試 Token 立即可用
		token, _ := GenerateToken(LoginUser, 100)
		claims, err := ParseToken(token)
		if err != nil {
			t.Errorf("Token should be valid immediately after generation: %v", err)
		}
		if claims.UserID != 100 {
			t.Errorf("Expected UserID=100, got %d", claims.UserID)
		}
	})
}

func TestParseTokenErrors(t *testing.T) {
	setupTest(t)

	tests := []struct {
		name        string
		token       string
		description string
	}{
		{
			name:        "空Token",
			token:       "",
			description: "empty token",
		},
		{
			name:        "格式錯誤_只有一部分",
			token:       "invalid_jwt_token",
			description: "malformed token",
		},
		{
			name:        "格式錯誤_只有兩部分",
			token:       "header.payload",
			description: "malformed token",
		},
		{
			name:        "無效簽名",
			token:       "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJhZG1pbl9pZCI6MCwiTWFwQ2xhaW1zIjp7ImV4cCI6MTcyMTIwNzI1MywiaXNzIjoiZ2luLWJsb2cifX0.invalid_signature",
			description: "invalid signature",
		},
		{
			name:        "不同密鑰簽名",
			token:       "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjowLCJhZG1pbl9pZCI6MSwiTWFwQ2xhaW1zIjp7ImV4cCI6MTcyMTIwNzI1MywiaXNzIjoiZ2luLWJsb2cifX0.y4Ku16plzvIUUPoCnF08xSG9JAOFgijv83ZNerxjjjo",
			description: "signed with different secret",
		},
		{
			name:        "無效Base64",
			token:       "not.valid.base64!!!",
			description: "invalid base64",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ParseToken(tt.token)
			if err == nil {
				t.Errorf("Expected error for %s, but got nil", tt.description)
				return
			}
			if claims != nil {
				t.Errorf("Expected nil claims, got %+v", claims)
			}
			t.Logf("Got expected error: %v", err)
		})
	}
}

func TestTokenExpiration(t *testing.T) {
	setupTest(t)

	// 註：此測試無法真正測試過期，因為實際的 Token 過期時間是 1 小時
	// 這裡只是驗證 exp claim 是否正確設置
	t.Run("驗證過期時間設置", func(t *testing.T) {
		beforeGen := time.Now()
		token, err := GenerateToken(LoginUser, testUserID)
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}
		afterGen := time.Now()

		claims, err := ParseToken(token)
		if err != nil {
			t.Fatalf("Failed to parse token: %v", err)
		}

		// 驗證過期時間是否存在
		exp, ok := claims.MapClaims["exp"].(float64)
		if !ok {
			t.Fatal("Expected exp claim to exist")
		}

		expTime := time.Unix(int64(exp), 0)

		// 過期時間應該在生成時間的 1 小時後（允許一些誤差）
		expectedExpMin := beforeGen.Add(59 * time.Minute)
		expectedExpMax := afterGen.Add(61 * time.Minute)

		if expTime.Before(expectedExpMin) || expTime.After(expectedExpMax) {
			t.Errorf("Expected expiration between %v and %v, got %v",
				expectedExpMin, expectedExpMax, expTime)
		}

		t.Logf("Token expiration: %v (in %v)", expTime, time.Until(expTime))
	})
}

func TestTokenRoleSegregation(t *testing.T) {
	setupTest(t)

	t.Run("用戶Token不包含AdminID", func(t *testing.T) {
		token, _ := GenerateToken(LoginUser, 123)
		claims, _ := ParseToken(token)

		if claims.UserID != 123 {
			t.Errorf("Expected UserID=123, got %d", claims.UserID)
		}
		if claims.AdminID != 0 {
			t.Errorf("Expected AdminID=0 for user token, got %d", claims.AdminID)
		}
	})

	t.Run("管理員Token不包含UserID", func(t *testing.T) {
		token, _ := GenerateToken(LoginAdmin, 456)
		claims, _ := ParseToken(token)

		if claims.AdminID != 456 {
			t.Errorf("Expected AdminID=456, got %d", claims.AdminID)
		}
		if claims.UserID != 0 {
			t.Errorf("Expected UserID=0 for admin token, got %d", claims.UserID)
		}
	})
}

// BenchmarkGenerateToken 測試 Token 生成性能
func BenchmarkGenerateToken(b *testing.B) {
	SetJwtSecret(testSecretKey)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = GenerateToken(LoginUser, uint(i%1000))
	}
}

// BenchmarkParseToken 測試 Token 解析性能
func BenchmarkParseToken(b *testing.B) {
	SetJwtSecret(testSecretKey)
	token, _ := GenerateToken(LoginUser, testUserID)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = ParseToken(token)
	}
}

// BenchmarkGenerateAndParse 測試完整流程性能
func BenchmarkGenerateAndParse(b *testing.B) {
	SetJwtSecret(testSecretKey)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		token, _ := GenerateToken(LoginUser, uint(i%1000))
		_, _ = ParseToken(token)
	}
}
