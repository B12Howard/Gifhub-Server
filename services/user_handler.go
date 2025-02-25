package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"kithli-api/firebase"
	"log"
	"net/http"
	"time"

	"firebase.google.com/go/auth"
	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
)

type UserQuery struct {
	Uid string `json:"uid"`
}

type UserTypes struct {
	Name            string
	File_size_limit int
	Usage_limit     int
	Id              int
	Created_at      sql.NullTime `json:"created_at"`
	Updated_at      sql.NullTime `json:"updated_at"`
}

type UserRes struct {
	id           int
	uid          string
	user_type_id int
	created_at   sql.NullTime
	ut           UserTypes
}

type UserRoleLimits struct {
	id              int
	max_gif_time    int
	file_size_limit int
	usage_limit     int
}

type UserRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password,omitempty"`
	Token       string `json:"token,omitempty"` // Token from Google/Apple
	Provider    string `json:"provider"`        // "email", "google", "apple"
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	PostalCode  string `json:"postalCode"`
	PhoneNumber string `json:"phoneNumber"`
}

func GetUser(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// var data UserQuery
		// decoder := json.NewDecoder(r.Body)
		// decoder.DisallowUnknownFields()

		// errDecode := decoder.Decode(&data)

		// if errDecode != nil {
		// 	log.Println(errDecode)
		// 	render.JSON(w, r, ("Missing uid"))
		// }

		// retrievedUser := &UserRes{}
		// retrievedUserType := &UserTypes{}

		// row := db.QueryRow(`SELECT Users.id, Users.created_at, UserTypes.name,  UserTypes.file_size_limit,  UserTypes.created_at, uid FROM users Users INNER JOIN user_types UserTypes ON Users.user_type_id=UserTypes.id WHERE uid = $1`, data.Uid)
		// err := row.Scan(&retrievedUser.id, &retrievedUser.created_at, &retrievedUserType.Name, &retrievedUserType.File_size_limit, &retrievedUserType.Created_at, &retrievedUser.uid)

		// if err != nil {
		// 	log.Println(err)
		// 	render.Status(r, http.StatusNotFound)
		// 	render.JSON(w, r, ("No User found with uid " + data.Uid))

		// 	return
		// }

		// payload := map[string]interface{}{
		// 	"id":  &retrievedUser.id,
		// 	"uid": &retrievedUser.uid,
		// 	"ut":  &retrievedUserType,
		// }
		payload := map[string]interface{}{
			"id":  "123d",
			"uid": "123abc",
			"ut":  time.Now().Local().String(),
		}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, payload)
	}
}

func SetUserUsage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi"))
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi"))
}

// func GetUserUsageById()

// func SetUserUsage(w http.ResponseWriter, r *http.Request) {
// 	render.Status(r, http.StatusCreated)
// 	render.JSON(w, r, map[string]string{"stuff": "post"})
// }
// func PutHandler(w http.ResponseWriter, r *http.Request) {
// 	render.Status(r, http.StatusCreated)
// 	render.JSON(w, r, map[string]string{"stuff": "put"})
// }

// func DeleteHandler(w http.ResponseWriter, r *http.Request) {
// 	render.Status(r, http.StatusCreated)
// 	render.JSON(w, r, map[string]string{"stuff": "delete"})
// }

func CreateUserHandler(db *sql.DB, firebaseClient *firebase.FirebaseClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("hello!")
		var req UserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		fmt.Printf("request body %+v", req)
		ctx := context.Background()
		var userID int
		var firebaseUID, email string
		var passwordHash string
		emailConfirmed := false

		switch req.Provider {
		case "email":
			// Create user with email and password

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, "Error hashing password: "+err.Error(), http.StatusInternalServerError)
				return
			}

			params := (&auth.UserToCreate{}).
				Email(req.Email).
				Password(req.Password)
			passwordHash = string(hashedPassword)
			user, err := firebaseClient.FirebaseAuth.CreateUser(ctx, params)
			// if FirebaseAuth == nil {
			// 	log.Println("Error: FirebaseAuth is nil")
			// 	http.Error(w, "FirebaseAuth is not initialized", http.StatusInternalServerError)
			// 	return
			// }
			log.Printf("after firebase auth")

			if err != nil {
				http.Error(w, "Error creating user: "+err.Error(), http.StatusInternalServerError)
				return
			}

			link, err := firebaseClient.FirebaseAuth.EmailVerificationLink(ctx, email)
			if err != nil {
				log.Println("[ERROR] Failed to generate verification link:", err)
				http.Error(w, "Failed to send verification email", http.StatusInternalServerError)
				return
			}

			// TODO: Send `link` via your email provider (e.g., SendGrid, SMTP)
			log.Printf("[INFO] Verification email link generated: %s", link)

			email = req.Email
			firebaseUID = user.UID

		case "google", "apple":
			// Verify token and extract user info
			token, err := firebaseClient.FirebaseAuth.VerifyIDToken(ctx, req.Token)
			if err != nil {
				http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
				return
			}

			emailClaim, emailExists := token.Claims["email"]
			if !emailExists {
				http.Error(w, "Email not found in token", http.StatusBadRequest)
				return
			}

			email = emailClaim.(string)
			emailConfirmed = true // Email is verified in OAuth flows
			firebaseUID = token.UID

		default:
			http.Error(w, "Unsupported provider", http.StatusBadRequest)
			return
		}

		dbErr := db.QueryRowContext(ctx, `
		INSERT INTO users (email, password_hash, email_confirmed, active_subscription, external_id, created_at, first_name, last_name, postal_code, phone)
		VALUES ($1, $2, $3, $4, $5, NOW(), $6, $7, $8, $9)
		ON CONFLICT (email) DO UPDATE SET external_id = EXCLUDED.external_id
		RETURNING id`, email, passwordHash, emailConfirmed, false, firebaseUID, req.FirstName, req.LastName, req.PostalCode, req.PhoneNumber,
		).Scan(&userID)

		fmt.Printf("%d", userID)

		if dbErr != nil {
			http.Error(w, "Error creating user: "+dbErr.Error(), http.StatusInternalServerError)
			return
		}

		response := json.NewEncoder(w).Encode(map[string]string{
			"uid":   firebaseUID,
			"email": email,
		})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		log.Println("[INFO] User registered successfully:", response)
	}
}
