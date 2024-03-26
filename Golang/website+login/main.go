package main

import (
	"html/template"
	"math/rand"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

type EmailConfig struct {
	Mailadress  string `json:"mailadress"`
	Wachtwoord  string `json:"wachtwoord"`
	Smtp_server string `json:"smtp_server"`
	Smtp_poort  int    `json:"smtp_poort"`
}

// Define a template for the HTML login form
var loginTpl = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Login</title>
</head>
<body>
    <h1>Login</h1>
    <form method="post" action="/login">
        <label for="username">Username:</label>
        <input type="text" id="username" name="username" required><br>
        <label for="password">Password:</label>
        <input type="password" id="password" name="password" required><br>
        <button type="submit">Login</button>
    </form>
</body>
</html>
`))

// Define a template for the HTML request application form
var requestAppTpl = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Request Application</title>
</head>
<body>
    <h1>Request Application</h1>
    <form method="post" action="/request">
        <label for="app-email">Email:</label>
        <input type="email" id="app-email" name="email" required>
        <button type="submit">Request Application</button>
    </form>
</body>
</html>
`))

// Map to store session IDs
var sessions = make(map[string]struct{})

func main() {
	http.HandleFunc("/", showLogin)
	http.HandleFunc("/login", login)
	http.HandleFunc("/reset", resetPassword) // Accessed without login
	http.HandleFunc("/request", requestApplicationHandler)
	http.ListenAndServe(":8080", nil)
}

func showLogin(w http.ResponseWriter, r *http.Request) {
	loginTpl.Execute(w, nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	// Get the username and password from the form submission
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Call the PowerShell script to authenticate the user
	cmd := exec.Command("powershell", "-File", "C:/Users/Administrator/Desktop/Authenticate.ps1", "-Username", username, "-Password", password)
	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, "Failed to authenticate user", http.StatusInternalServerError)
		return
	}

	// Check the output of the script
	if strings.TrimSpace(string(output)) == "Authenticated" {
		// If authentication succeeds, redirect to the home page or any other authorized page
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	// If authentication fails, render the login form again with an error message
	loginTpl.Execute(w, "Invalid username or password")
}

func requestApplicationHandler(w http.ResponseWriter, r *http.Request) {
	// Check if user is logged in
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Verify session ID
	if _, ok := sessions[cookie.Value]; ok {
		// Render the request application form
		requestAppTpl.Execute(w, nil)
		return
	}

	// If session ID is not valid or user is not logged in, redirect to login page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Function to generate a random string of given length
func generateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator with the current time

	// Characters allowed in the random string
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%&*-_"

	// Create a byte slice to store the random string
	randomString := make([]byte, length)

	// Generate random characters and append them to the byte slice
	for i := 0; i < length; i++ {
		randomString[i] = charset[rand.Intn(len(charset))]
	}

	// Convert the byte slice to a string and return
	return string(randomString)
}

func resetPassword(w http.ResponseWriter, r *http.Request) {
	// Get the username from the form submission
	//username := r.FormValue("username")

	// Call the PowerShell script to get the email address
	//cmd := exec.Command("powershell", "-File", "C:/Users/Administrator/Desktop/GetEmail.ps1", "-Username", username)
	//email, err := cmd.Output()
	//if err != nil {
	//	http.Error(w, "Failed to retrieve email address", http.StatusInternalServerError)
	//	//return
	//}

	// Send an email to the retrieved email address (replace with your email sending logic)
	//sendEmail(string(email))

	// Respond to the user indicating that the password reset email has been sent
	w.Write([]byte("Password reset email has been sent to the provided email address"))
}
