package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
)

var (
	PASSWORD = os.Getenv("TODO_PASSWORD")
	jwtKey   = []byte("very-secret-key")
)

// Credentials структура для передачи пароля пользователя при входе.
//
// swagger:model Credentials
type Credentials struct {
	Password string `json:"password"`
}

type Claims struct {
	PasswordHash string `json:"password-hash"`
	jwt.RegisteredClaims
}

func checkPassword(password string) string {
	if password == "" {
		password = "123456"
	}
	return password
}

// signin получает пароль пользователя, выдает JWT токен при успехе.
//
// @Summary      Вход в систему (логин)
// @Description  Получить JWT токен по паролю. Пароль можно задать через переменную окружения TODO_PASSWORD (по умолчанию "123456").
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      Credentials  true  "Пароль пользователя"
// @Success      200          {object}  map[string]string    "JWT токен"
// @Failure      401          {object}  map[string]interface{} "Ошибка: неверный пароль"
// @Failure      500          {object}  map[string]interface{} "Внутренняя ошибка"
// @Router       /api/signin [post]
func signin(w http.ResponseWriter, r *http.Request) {
	var pass Credentials
	json.NewDecoder(r.Body).Decode(&pass)
	err := validator.New().Struct(&pass)
	PASSWORD = checkPassword(PASSWORD)
	if err != nil || pass.Password != PASSWORD {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]any{"error": "Неверный пароль"})
		return
	}

	hashBytes := sha256.Sum256([]byte(pass.Password))
	hashString := fmt.Sprintf("%x", hashBytes)
	claims := &Claims{
		PasswordHash: hashString,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{"token": tokenString})
}

func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		PASSWORD = checkPassword(PASSWORD)
		if len(PASSWORD) > 0 {
			var tokenString string // JWT-токен из куки

			// получаем куку
			cookie, err := r.Cookie("token")
			if err == nil {
				tokenString = cookie.Value
			}

			var valid bool
			token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return jwtKey, nil
			})
			if err != nil {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}

			hashBytes := sha256.Sum256([]byte(PASSWORD))
			hashString := fmt.Sprintf("%x", hashBytes)
			if hashString == token.Claims.(*Claims).PasswordHash {
				valid = true
			}

			if !valid {
				// возвращаем ошибку авторизации 401
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
