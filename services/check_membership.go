package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

// UserMembershipCheckRequest represents the request body
type UserMembershipCheckRequest struct {
	UID string `json:"uid"`
}

// UserMembershipCheckResponse represents the API response
type UserMembershipCheckResponse struct {
	HasMember bool `json:"has_member"`
}

// CheckUserMembershipHandler checks if a user has a `member` foreign key
func CheckUserMembershipHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UserMembershipCheckRequest
		// Decode the JSON request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "Invalid request payload"}`, http.StatusBadRequest)
			return
		}
		if req.UID == "" {
			http.Error(w, `{"error": "UID is required"}`, http.StatusBadRequest)
			return
		}

		ctx := context.Background()
		var memberID sql.NullString

		// Query to check if user has a `member` foreign key
		err := db.QueryRowContext(ctx, `SELECT member FROM users WHERE external_id = $1`, req.UID).Scan(&memberID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, `{"error": "User not found"}`, http.StatusNotFound)
				return
			}
			log.Println("[ERROR] Database error:", err)
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
			return
		}
		// Response indicating whether the user has a `member` foreign key
		response := UserMembershipCheckResponse{HasMember: memberID.Valid}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}