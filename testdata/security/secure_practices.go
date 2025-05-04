package security

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"golang.org/x/crypto/argon2"
)

type Credentials struct {
	Username string
	Password []byte
	Salt     []byte
}

func SecureHashPassword(password string) (string, string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", "", fmt.Errorf("failed to generate salt: %w", err)
	}
	
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	
	return hex.EncodeToString(hash), hex.EncodeToString(salt), nil
}

func ProtectedSQLQuery(db *sql.DB, username string) (*sql.Rows, error) {
	query := "SELECT * FROM users WHERE username = ?"
	return db.Query(query, username)
}

func XSSProtection(w http.ResponseWriter, r *http.Request) {
	input := r.URL.Query().Get("input")
	tmpl := template.Must(template.New("response").Parse("<div>{{.}}</div>"))
	if err := tmpl.Execute(w, input); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

func EnvironmentCredentials(ctx context.Context) (string, string, error) {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return "", "", fmt.Errorf("API_KEY environment variable not set")
	}
	
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	if username == "" || password == "" {
		return "", "", fmt.Errorf("database credentials not properly configured")
	}
	
	return username, password, nil
}
