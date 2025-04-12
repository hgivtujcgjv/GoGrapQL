package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"
const sessionKey = 1

type Session struct {
	UserID   int
	Username string
}

// все серверные методы строго пост , второе что происходит это auth требуют только два метода
// добавить удалить из корзины и посмотреть содержимое
// точнее миделварь парсит пользователя и создает для него контекст
// контекст позволяет сделать запрос уникальным и иидентифицировать пользователя
// вернуть поле токен в любом формате , и потом парсить Authorization для проверки авторизации
// писать сразу резолвер типа интерфейс

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Token" {
				http.Error(w, `{"errors":{"body":["Invalid Authorization format"]}}`, http.StatusUnauthorized)
				return
			}
			tokenData, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				http.Error(w, `{"errors":{"body":["Invalid token"]}}`, http.StatusUnauthorized)
				return
			}
			userData := strings.Split(string(tokenData), ":")
			if len(userData) != 2 {
				http.Error(w, `{"errors":{"body":["Invalid token structure"]}}`, http.StatusUnauthorized)
				return
			}
			user := User{
				Username: userData[0],
				Email:    userData[1],
			}
			ctx := context.WithValue(r.Context(), sessionKey, &user)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateToken(user *User) string {
	tokenData := fmt.Sprintf("%s:%s", user.Username, user.Email)
	encoded := base64.StdEncoding.EncodeToString([]byte(tokenData))
	return encoded
}

type TokenResponse struct {
	Token string `json:"token"`
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var us User
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, `{"errors":{"body":["Error while parsing request for register"]}}`, http.StatusInternalServerError)
		return
	}
	fmt.Println("Received request body:", string(body))
	if err := json.Unmarshal(body, &us); err != nil {
		http.Error(w, `{"errors":{"body":["Invalid JSON format"]}}`, http.StatusBadRequest)
		return
	}
	if us.Email == "" || us.Password == "" || us.Username == "" {
		http.Error(w, `{"errors":{"body":["Mismatched input fields"]}}`, http.StatusBadRequest)
		return
	}
	token := CreateToken(&us)
	w.WriteHeader(http.StatusCreated)
	writeJSONResponse(w, TokenResponse{Token: "Bearer " + token})
}

func writeJSONResponse(w http.ResponseWriter, smth any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(smth); err != nil {
		http.Error(w, `{"errors":{"body":["Failed to encode response"]}}`, http.StatusInternalServerError)
	}
}

func GetApp() http.Handler {
	resolver, err := loadData("testdata.json")
	if err != nil {
		log.Fatalf("Ошибка загрузки данных: %v", err)
	}
	cfg := Config{
		Resolvers: resolver,
	}
	srv := handler.New(NewExecutableSchema(cfg))
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.Use(extension.FixedComplexityLimit(1000))
	//srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})
	mux := http.NewServeMux()
	mux.HandleFunc("/register", RegisterUser)
	authProtected := AuthMiddleware(srv)
	mux.Handle("/query", authProtected)
	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	return mux
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, GetApp())) //заместо GetApp nil при тестировании
}
