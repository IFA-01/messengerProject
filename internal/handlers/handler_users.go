package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IFA-01/messenger/internal/auth"
	"github.com/IFA-01/messenger/internal/models"
	"github.com/IFA-01/messenger/internal/repository/queries"
)

func HandlerCreateUser(q *queries.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Username string `json:"username"`
			Nickname string `json:"nickname"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		decoder := json.NewDecoder(r.Body)

		params := parameters{}

		err := decoder.Decode(&params)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("error parsing JSON: %v", err))
			return
		}

		hashedPassword, err := auth.HashPassword(params.Password)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("error hashing password: %v", err))
			return
		}

		user, err := q.CreateUser(r.Context(), queries.CreateUserParams{
			Username:     params.Username,
			Nickname:     params.Nickname,
			Email:        params.Email,
			PasswordHash: string(hashedPassword),
		})

		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("cant create user %v", err))
			return
		}

		response := models.DatabaseUserToUser(user)
		response.PasswordHash = ""

		respondWithJSON(w, 200, response)
	}
}

func HandleGetUser(q *queries.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int64)
		user, err := q.GetUserByID(r.Context(), userID)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Error getting user by ID: %v", err))
			return
		}
		response := models.DatabaseUserToUser(user)
		response.PasswordHash = ""
		respondWithJSON(w, 200, response)
	}
}
