package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"kithli-api/firebase"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// UserDataResponse represents the user response excluding password
type UserDataResponse struct {
	ID                 int     `json:"id"`
	Email              string  `json:"email"`
	Phone              *string `json:"phone,omitempty"`
	FirstName          *string `json:"firstName,omitempty"`
	LastName           *string `json:"lastName,omitempty"`
	EmailConfirmed     bool    `json:"emailConfirmed"`
	ActiveSubscription bool    `json:"activeSubscription"`
	ExternalID         string  `json:"external_id"`
	MemberID           *int    `json:"member_id,omitempty"`
	KithID             *int    `json:"kith_id,omitempty"`
	CreatedAt          string  `json:"createdAt"`
	// Member details (if exists)
	MyHeadline     *string `json:"myHeadline,omitempty"`
	AboutMe        *string `json:"aboutMe,omitempty"`
	PostalCode     *string `json:"postalCode,omitempty"`
	StreetAddress  *string `json:"streetAddress,omitempty"`
	AptNumber      *string `json:"aptNumber,omitempty"`
	City           *string `json:"city,omitempty"`
	State          *string `json:"state,omitempty"`
	AdditionalInfo *string `json:"additionalInformation,omitempty"`
}

func ExtractUserIDFromToken(authClient *firebase.FirebaseClient, tokenString string) (string, error) {
	tokenString = strings.TrimSpace(strings.Replace(tokenString, "Bearer ", "", 1)) // Remove 'Bearer ' prefix

	token, err := authClient.FirebaseAuth.VerifyIDToken(context.Background(), tokenString)
	if err != nil {
		return "", err
	}

	uid, ok := token.Claims["user_id"].(string)
	if !ok {
		return "", err
	}
	return uid, nil
}

// GetUserDataHandler retrieves user data by external_id but excludes the password
func GetUserDataHandler(db *sql.DB, authClient *firebase.FirebaseClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		externalID := chi.URLParam(r, "external_id")

		if externalID == "" {
			http.Error(w, `{"error": "external_id is required"}`, http.StatusBadRequest)
			return
		}

		// Extract the token from the request header
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, `{"error": "Missing authentication token"}`, http.StatusUnauthorized)
			return
		}

		// Validate token and extract user ID
		requesterID, err := ExtractUserIDFromToken(authClient, tokenString)
		if err != nil {
			http.Error(w, `{"error": "Invalid token"}`, http.StatusUnauthorized)
			return
		}

		// Ensure the requester is fetching their own data
		if requesterID != externalID {
			http.Error(w, `{"error": "Unauthorized access"}`, http.StatusForbidden)
			return
		}

		ctx := context.Background()
		var userData UserDataResponse

		// Query user data
		err = db.QueryRowContext(ctx, `
			SELECT id, email, phone, first_name, last_name, email_confirmed, active_subscription, external_id, member, kith, created_at, postal_code 
			FROM users 
			WHERE external_id = $1
		`, externalID).Scan(
			&userData.ID, &userData.Email, &userData.Phone, &userData.FirstName, &userData.LastName,
			&userData.EmailConfirmed, &userData.ActiveSubscription, &userData.ExternalID, &userData.MemberID, &userData.KithID, &userData.CreatedAt,  &userData.PostalCode,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, `{"error": "User not found"}`, http.StatusNotFound)
				return
			}
			log.Println("[ERROR] Database error:", err)
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
			return
		}

		// Fetch additional member details if the user has a member ID
		if userData.MemberID != nil {
			err = db.QueryRowContext(ctx, `
				SELECT my_headline, about_me, postal_code, street_address, apt_number, city, state, additional_information
				FROM members WHERE id = $1
			`, *userData.MemberID).Scan(
				&userData.MyHeadline, &userData.AboutMe, &userData.PostalCode, &userData.StreetAddress,
				&userData.AptNumber, &userData.City, &userData.State, &userData.AdditionalInfo,
			)

			if err != nil && err != sql.ErrNoRows {
				log.Println("[ERROR] Failed to fetch member details:", err)
				http.Error(w, `{"error": "Failed to fetch member details"}`, http.StatusInternalServerError)
				return
			}
		}

		// Return user data
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(userData)
		log.Println("[INFO] User data retrieved for external_id:", externalID)
	}
}
