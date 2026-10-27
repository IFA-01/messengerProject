package handlers

import (
	"encoding/json"
	"net/http"

	"errors"

	"github.com/IFA-01/messenger/internal/middleware"
	"github.com/IFA-01/messenger/internal/repository/queries"

	"github.com/jackc/pgx/v5"
)

func HandleCreateChat(q *queries.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID, ok := middleware.UserIDFromContext(r.Context())
		if !ok {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		type parameters struct {
			Username string `json:"username"`
			Name     string `json:"name"`
		}

		params := parameters{}

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		member, err := q.GetUserByUsername(r.Context(), params.Username)
		if errors.Is(err, pgx.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "user not found")
			return
		}
		if err != nil {
			respondWithInternalError(w, err)
			return
		}

		if member.ID == userID {
			respondWithError(w, http.StatusBadRequest, "cannot create chat with yourself")
			return
		}

		existingChat, err := q.FindDirectChatsBetweenUsers(r.Context(), queries.FindDirectChatsBetweenUsersParams{
			UserID:   userID,
			UserID_2: member.ID,
		})
		if err == nil {
			respondWithJSON(w, http.StatusOK, existingChat)
			return
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			respondWithInternalError(w, err)
			return
		}

		chat, err := q.CreateChat(r.Context(), queries.CreateChatParams{
			Name:    params.Name,
			IsGroup: false,
		})
		if err != nil {
			respondWithInternalError(w, err)
			return
		}

		_, err = q.AddChatMember(r.Context(), queries.AddChatMemberParams{
			ChatID: chat.ID,
			UserID: userID,
		})
		if err != nil {
			respondWithInternalError(w, err)
			return
		}
		_, err = q.AddChatMember(r.Context(), queries.AddChatMemberParams{
			ChatID: chat.ID,
			UserID: member.ID,
		})
		if err != nil {
			respondWithInternalError(w, err)
			return
		}

		respondWithJSON(w, http.StatusCreated, chat)
	}
}

func HandleListChats(q *queries.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.UserIDFromContext(r.Context())
		if !ok {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		chats, err := q.ListUsersChats(r.Context(), userID)
		if err != nil {
			respondWithInternalError(w, err)
			return
		}

		if chats == nil {
			chats = []queries.Chat{}
		}
		respondWithJSON(w, http.StatusOK, map[string]interface{}{"chats": chats})
	}
}
