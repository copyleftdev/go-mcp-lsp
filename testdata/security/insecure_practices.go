package security

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

// InsecureHashFunction uses MD5 which is cryptographically broken
func InsecureHashFunction(password string) string {
	hasher := md5.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}

// VulnerableToSQLInjection doesn't protect against SQL injection
func VulnerableToSQLInjection(username string) string {
	query := fmt.Sprintf("SELECT * FROM users WHERE username = '%s'", username)
	return query
}

// XSSVulnerability is vulnerable to cross-site scripting
func XSSVulnerability(w http.ResponseWriter, r *http.Request) {
	input := r.URL.Query().Get("input")
	fmt.Fprintf(w, "<div>%s</div>", input) // Direct inclusion of user input
}

// HardcodedCredentials embeds sensitive data directly in code
func HardcodedCredentials() {
	apiKey := "1234567890abcdef"
	username := "admin"
	password := "supersecret"
	
	fmt.Println("Connecting with credentials:", username, password, apiKey)
}
