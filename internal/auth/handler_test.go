package auth

import (
	"net/http"
	"net/http/httptest"
	"simple_gin_server/configs"
	"simple_gin_server/internal/moks"

	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterHandler_ValidCredentials(t *testing.T) {
	// Настройка теста
	gin.SetMode(gin.TestMode)
	mockService := new(moks.MockAuthService)
	conf := configs.LoadConfig()
	handler := NewAuthHandler(mockService, conf)

	// Ожидаемый вызов
	mockService.On("Register", mock.Anything, "test@example.com", "va123*tro").Return(nil)

	// Создание тестового контекста
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Устанавливаем validatedData в контекст
	c.Set("validatedData", &RegisterRequest{
		Email:    "test@example.com",
		Password: "va123*tro",
	})

	// Вызов хендлера
	handler.RegisterHandler(c)

	// Проверки
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"message": "user registered"}`, w.Body.String())
	mockService.AssertExpectations(t)
}
