package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IFA-01/messenger/internal/models"
	"github.com/IFA-01/messenger/internal/repository/queries"
	"golang.org/x/crypto/bcrypt"
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

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
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
