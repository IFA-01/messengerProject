package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/IFA-01/messenger/internal/repository/queries"
	"github.com/go-chi/chi/v5"
)

func HandleCreateMessage(q *queries.Queries) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int64)

		type parameters struct {
			ChatID  int64  `json:"chat_id"`
			Content string `json:"content"`
		}
		params := parameters{}

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error parsing json: %v", err))
			return
		}

		if params.Content == "" {
			respondWithError(w, http.StatusBadRequest, "content is required")
			return
		}

		_, err := q.GetChatByID(r.Context(), params.ChatID)
		if err != nil {
			respondWithError(w, http.StatusNotFound, fmt.Sprintf("chat not found with id: %d", params.ChatID))
			return
		}

		isChatMember, err := q.IsChatMember(r.Context(), queries.IsChatMemberParams{
			ChatID: params.ChatID,
			UserID: userID,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error checking chat member: %v", err))
			return
		}
		if !isChatMember {
			respondWithError(w, http.StatusBadRequest, "you are not a member of this chat")
			return
		}

		message, err := q.CreateMessage(r.Context(), queries.CreateMessageParams{
			ChatID:   params.ChatID,
			SenderID: userID,
			Content:  params.Content,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error creating message: %v", err))
			return
		}

		respondWithJSON(w, http.StatusCreated, message)
	}
}

func HandleListMessages(q *queries.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int64)

		chatIDstr := chi.URLParam(r, "chatID")
		chatID, err := strconv.ParseInt(chatIDstr, 10, 64)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("invalid chat id: %s", chatIDstr))
			return
		}

		isChatMember, err := q.IsChatMember(r.Context(), queries.IsChatMemberParams{
			ChatID: chatID,
			UserID: userID,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error checking chat member: %v", err))
			return
		}
		if !isChatMember {
			respondWithError(w, http.StatusForbidden, "you are not a member of this chat")
			return
		}

		messages, err := q.ListChatMessages(r.Context(), chatID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error listing messages: %v", err))
			return
		}

		respondWithJSON(w, http.StatusOK, map[string]interface{}{
			"messages": messages,
		})
	}
}
