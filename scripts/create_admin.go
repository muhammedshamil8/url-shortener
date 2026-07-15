package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/muhammedshamil8/url-shortener/internal/auth"
	"github.com/muhammedshamil8/url-shortener/internal/config"
	"github.com/muhammedshamil8/url-shortener/internal/database"
	"github.com/muhammedshamil8/url-shortener/internal/repository"
	"golang.org/x/term"
)

func main() {
	// Flags
	flagUsername := flag.String("username", "", "Username for the admin user")
	flagEmail := flag.String("email", "", "Email for the admin user")
	flagPassword := flag.String("password", "", "Password for the admin user")
	flag.Parse()

	// Load env
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found, using existing environment variables")
	}

	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		fmt.Printf("Error: Invalid configuration: %v\n", err)
		os.Exit(1)
	}

	// Connect to database
	db, err := database.InitDB(cfg.DB)
	if err != nil {
		fmt.Printf("Error: Failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	repo := repository.New(db)

	reader := bufio.NewReader(os.Stdin)

	username := *flagUsername
	if username == "" {
		fmt.Print("Enter admin username: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading username: %v\n", err)
			os.Exit(1)
		}
		username = strings.TrimSpace(input)
	}

	if len(username) < 3 || len(username) > 50 {
		fmt.Println("Error: Username must be between 3 and 50 characters")
		os.Exit(1)
	}

	email := *flagEmail
	if email == "" {
		fmt.Print("Enter admin email: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading email: %v\n", err)
			os.Exit(1)
		}
		email = strings.TrimSpace(input)
	}

	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		fmt.Println("Error: Invalid email format")
		os.Exit(1)
	}

	password := *flagPassword
	if password == "" {
		fmt.Print("Enter admin password: ")
		// Try to read password without echo
		fd := int(os.Stdin.Fd())
		if term.IsTerminal(fd) {
			bytePassword, err := term.ReadPassword(fd)
			if err != nil {
				fmt.Printf("\nError reading password: %v\n", err)
				os.Exit(1)
			}
			fmt.Println() // print newline after password input
			password = string(bytePassword)
		} else {
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("Error reading password: %v\n", err)
				os.Exit(1)
			}
			password = strings.TrimSpace(input)
		}
	}

	if len(password) < 8 || len(password) > 72 {
		fmt.Println("Error: Password must be between 8 and 72 characters")
		os.Exit(1)
	}

	// Check if user already exists
	existingUser, _ := repo.GetUserByEmail(email)
	if existingUser != nil {
		fmt.Printf("Error: A user with email %s already exists\n", email)
		os.Exit(1)
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		fmt.Printf("Error: Failed to hash password: %v\n", err)
		os.Exit(1)
	}

	// Create admin user
	id, err := repo.CreateUserWithRole(username, email, hashedPassword, "admin")
	if err != nil {
		fmt.Printf("Error: Failed to create admin user: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully created admin user '%s' (ID: %d, Email: %s) with role 'admin'\n", username, id, email)
}
