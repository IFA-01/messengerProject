package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IFA-01/messenger/internal/auth"
	"github.com/IFA-01/messenger/internal/models"
	"github.com/IFA-01/messenger/internal/repository/queries"
)

func HandleLogin(q *queries.Queries, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		params := parameters{}

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&params)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("error parsing JSON: %v", err))
			return
		}
		user, err := q.GetUserByEmail(r.Context(), params.Email)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Error getting user by email: %v:", err))
			return
		}

		if err := auth.CheckPassword(params.Password, user.PasswordHash); err != nil {
			respondWithError(w, 400, fmt.Sprintf("Invalid password: %v", err))
			return
		}

		token, err := auth.GenerateToken(user.ID, jwtSecret)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Error generating token: %v", err))
			return
		}

		responseUser := models.DatabaseUserToUser(user)
		responseUser.PasswordHash = ""

		type Loginresponse struct {
			Token string      `json:"token"`
			User  models.User `json:"user"`
		}

		respondWithJSON(w, 200, Loginresponse{
			Token: token,
			User:  responseUser,
		})
	}
}
