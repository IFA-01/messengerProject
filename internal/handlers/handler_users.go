package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/IFA-01/messenger/internal/auth"
	"github.com/IFA-01/messenger/internal/middleware"
	"github.com/IFA-01/messenger/internal/models"
	"github.com/IFA-01/messenger/internal/repository/queries"
	"github.com/jackc/pgx/v5"
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
			respondWithError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		if params.Username == "" || params.Email == "" || params.Password == "" {
			respondWithError(w, http.StatusBadRequest, "username, email and password are required")
			return
		}
		hashedPassword, err := auth.HashPassword(params.Password)
		if err != nil {
			respondWithInternalError(w, err)
			return
		}
		user, err := q.CreateUser(r.Context(), queries.CreateUserParams{
			Username:     params.Username,
			Nickname:     params.Nickname,
			Email:        params.Email,
			PasswordHash: hashedPassword,
		})
		if isUniqueViolation(err) {
			respondWithError(w, http.StatusConflict, "username or email already taken")
			return
		}
		if err != nil {
			respondWithInternalError(w, err)
			return
		}
		respondWithJSON(w, http.StatusCreated, models.DatabaseUserToUser(user))
	}
}

func HandleGetUser(q *queries.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.UserIDFromContext(r.Context())
		if !ok {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		user, err := q.GetUserByID(r.Context(), userID)
		if errors.Is(err, pgx.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "user not found")
			return
		}
		if err != nil {
			respondWithInternalError(w, err)
			return
		}
		respondWithJSON(w, http.StatusOK, models.DatabaseUserToUser(user))
	}
}
