package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/IFA-01/messenger/internal/middleware"
	"github.com/IFA-01/messenger/internal/repository/queries"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

func HandleCreateMessage(q *queries.Queries) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.UserIDFromContext(r.Context())
		if !ok {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		type parameters struct {
			ChatID  int64  `json:"chat_id"`
			Content string `json:"content"`
		}
		params := parameters{}

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		if params.Content == "" {
			respondWithError(w, http.StatusBadRequest, "content is required")
			return
		}

		_, err := q.GetChatByID(r.Context(), params.ChatID)
		if errors.Is(err, pgx.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "chat not found")
			return
		}
		if err != nil {
			respondWithInternalError(w, err)
			return
		}

		isChatMember, err := q.IsChatMember(r.Context(), queries.IsChatMemberParams{
			ChatID: params.ChatID,
			UserID: userID,
		})
		if err != nil {
			respondWithInternalError(w, err)
			return
		}
		if !isChatMember {
			respondWithError(w, http.StatusForbidden, "you are not a member of this chat")
			return
		}

		message, err := q.CreateMessage(r.Context(), queries.CreateMessageParams{
			ChatID:   params.ChatID,
			SenderID: userID,
			Content:  params.Content,
		})
		if err != nil {
			respondWithInternalError(w, err)
			return
		}

		respondWithJSON(w, http.StatusCreated, message)
	}
}

func HandleListMessages(q *queries.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.UserIDFromContext(r.Context())
		if !ok {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		chatIDstr := chi.URLParam(r, "chatID")
		chatID, err := strconv.ParseInt(chatIDstr, 10, 64)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid chat id")
			return
		}

		isChatMember, err := q.IsChatMember(r.Context(), queries.IsChatMemberParams{
			ChatID: chatID,
			UserID: userID,
		})
		if err != nil {
			respondWithInternalError(w, err)
			return
		}
		if !isChatMember {
			respondWithError(w, http.StatusForbidden, "you are not a member of this chat")
			return
		}

		messages, err := q.ListChatMessages(r.Context(), chatID)
		if err != nil {
			respondWithInternalError(w, err)
			return
		}

		respondWithJSON(w, http.StatusOK, map[string]interface{}{
			"messages": messages,
		})
	}
}
