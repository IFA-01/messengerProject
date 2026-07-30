package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IFA-01/messenger/internal/repository/queries"
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
