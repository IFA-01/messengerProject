package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"errors"

	"github.com/IFA-01/messenger/internal/repository/queries"

	"github.com/jackc/pgx/v5"
)

func HandleCreateChat(q *queries.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID := r.Context().Value("userID").(int64)

		type parameters struct {
			Username string `json:"username"`
			Name     string `json:"name"`
		}

		params := parameters{}

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error parsing json: %v", err))
			return
		}

		member, err := q.GetUserByUsername(r.Context(), params.Username)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("No user found with nickname: %s", params.Username))
			return
		}

		if member.ID == userID {
			respondWithError(w, http.StatusBadRequest, "cannot create chat with yourself")
			return
		}

		_, err = q.GetUserByID(r.Context(), member.ID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("No user found with id: %d", member.ID))
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
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error checking existing chat: %v", err))
			return
		}

		chat, err := q.CreateChat(r.Context(), queries.CreateChatParams{
			Name:    params.Name,
			IsGroup: false,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error creating chat: %v", err))
			return
		}

		_, err = q.AddChatMember(r.Context(), queries.AddChatMemberParams{
			ChatID: chat.ID,
			UserID: userID,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error adding chat member: %v", err))
			return
		}
		_, err = q.AddChatMember(r.Context(), queries.AddChatMemberParams{
			ChatID: chat.ID,
			UserID: member.ID,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error adding chat member: %v", err))
			return
		}

		respondWithJSON(w, http.StatusCreated, chat)
	}
}

func HandleListChats(q *queries.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int64)

		chats, err := q.ListUsersChats(r.Context(), userID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error listing chats: %v", err))
			return
		}

		if chats == nil {
			chats = []queries.Chat{}
		}
		respondWithJSON(w, http.StatusOK, map[string]interface{}{"chats": chats})
	}
}
