package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

// MemberRequest represents the incoming request
type MemberRequest struct {
	UID                   string  `json:"uid"`
	MyHeadline            string  `json:"myHeadline"`
	AboutMe               *string `json:"aboutMe,omitempty"`
	PostalCode            *string `json:"postalCode,omitempty"`
	StreetAddress         *string `json:"streetAddress,omitempty"`
	AptNumber             *string `json:"aptNumber,omitempty"`
	City                  *string `json:"city,omitempty"`
	State                 *string `json:"state,omitempty"`
	AdditionalInformation *string `json:"additionalInformation,omitempty"`
}

// CreateMemberHandler inserts a new member and updates the user with the member ID
func CreateMemberHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req MemberRequest

		// Decode JSON request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "Invalid request payload"}`, http.StatusBadRequest)
			return
		}

		if req.UID == "" || req.MyHeadline == "" {
			http.Error(w, `{"error": "UID and myHeadline are required"}`, http.StatusBadRequest)
			return
		}

		ctx := context.Background()
		var memberID int

		// Insert the new member into the database
		err := db.QueryRowContext(ctx, `
			INSERT INTO members (my_headline, about_me, postal_code, street_address, apt_number, city, state, additional_information)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id
		`, req.MyHeadline, req.AboutMe, req.PostalCode, req.StreetAddress, req.AptNumber, req.City, req.State, req.AdditionalInformation).Scan(&memberID)

		if err != nil {
			log.Println("[ERROR] Failed to create member:", err)
			http.Error(w, `{"error": "Failed to create member"}`, http.StatusInternalServerError)
			return
		}

		// Update the user with the new member ID
		_, err = db.ExecContext(ctx, `
			UPDATE users SET member = $1 WHERE external_id = $2
		`, memberID, req.UID)

		if err != nil {
			log.Println("[ERROR] Failed to update user with member ID:", err)
			http.Error(w, `{"error": "Failed to update user"}`, http.StatusInternalServerError)
			return
		}

		// Respond with success
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"member_id": memberID,
			"message":   "Member created successfully",
		})
		log.Println("[INFO] Member created successfully for user:", req.UID)
	}
}
