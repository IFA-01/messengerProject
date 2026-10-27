package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/IFA-01/messenger/internal/auth"
	"github.com/IFA-01/messenger/internal/models"
	"github.com/IFA-01/messenger/internal/repository/queries"
	"github.com/jackc/pgx/v5"
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
			respondWithError(w, http.StatusBadRequest, "Invalid JSON body")
			return
		}
		user, err := q.GetUserByEmail(r.Context(), params.Email)
		if errors.Is(err, pgx.ErrNoRows) {
			respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}
		if err != nil {
			respondWithInternalError(w, err)
			return
		}

		if err := auth.CheckPassword(params.Password, user.PasswordHash); err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}

		token, err := auth.GenerateToken(user.ID, jwtSecret)
		if err != nil {
			respondWithInternalError(w, err)
			return
		}

		type Loginresponse struct {
			Token string      `json:"token"`
			User  models.User `json:"user"`
		}

		respondWithJSON(w, http.StatusOK, Loginresponse{
			Token: token,
			User:  models.DatabaseUserToUser(user),
		})
	}
}
