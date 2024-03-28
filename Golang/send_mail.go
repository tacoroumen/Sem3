package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/smtp"
	"os"
)

type EmailConfig struct {
	Mailadress  string `json:"mailadress"`
	Wachtwoord  string `json:"wachtwoord"`
	Smtp_server string `json:"smtp_server"`
	Smtp_poort  int    `json:"smtp_poort"`
}

func main() {
	// Parse command-line flags
	var fullname, username, password, email string
	flag.StringVar(&fullname, "fullname", "", "User's full name")
	flag.StringVar(&username, "username", "", "User's username")
	flag.StringVar(&password, "password", "", "User's password")
	flag.StringVar(&email, "email", "", "User's email address")
	flag.Parse()

	if fullname == "" || username == "" || password == "" || email == "" {
		fmt.Println("Usage: send_email -fullname <fullname> -username <username> -password <password> -email <email>")
		return
	}

	// Open and read the email configuration from a JSON file
	file, err := os.Open("C:/Users/Administrator/Desktop/Website/mail.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := EmailConfig{}
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println(err)
		return
	}

	auth := smtp.PlainAuth("", config.Mailadress, config.Wachtwoord, config.Smtp_server)

	to := []string{email}
	subject := "Subject: " + "Hello " + fullname + " welcome to tacoroumen" + "\r\n"
	body := "Hello " + fullname + " Your account has been created successfully." + "\r\n\r\n Your login credentials are \r\n Username: " + username + "\r\n Password: " + password + "\r\n"
	msg := []byte(subject +
		"\r\n" +
		body)

	err = smtp.SendMail(fmt.Sprintf("%s:%d", config.Smtp_server, config.Smtp_poort), auth, config.Mailadress, to, msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Email sent successfully!")
}
