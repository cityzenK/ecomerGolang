package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cityzenK/com/types"
	"github.com/gorilla/mux"
)

func TestUserServiceHandlers(t *testing.T) {
	userStore := &mockUserStore{}
	handler := NewHandler(userStore)

	t.Run(
		"should fail if the user payload is invalid",
		func(t *testing.T) {
			payload := types.RegisterUserPayload{
				FirstName: "eliot",
				Email:     "",
				Password:  "asd",
				LastName:  "user123",
			}

			marshaled, _ := json.Marshal(payload)

			req, err := http.NewRequest(
				http.MethodPost,
				"/register",
				bytes.NewBuffer(marshaled),
			)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			router := mux.NewRouter()

			router.HandleFunc("/register", handler.handleRegister)

			router.ServeHTTP(rr, req)

			if rr.Code != http.StatusBadRequest {
				t.Errorf("expected code status %d, got %d", http.StatusBadRequest, rr.Code)
			}
		})
}

type mockUserStore struct{}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	return nil, fmt.Errorf("user not found")
}
func (m *mockUserStore) GetUserByID(id int) (*types.User, error) {
	return nil, nil
}
func (m *mockUserStore) CreateUser(types.User) error {
	return nil
}
