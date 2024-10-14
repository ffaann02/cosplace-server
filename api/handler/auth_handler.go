package handler

import (
	"fmt"
	"os"
	"time"

	config "github.com/ffaann02/cosplace-server/internal/config"
	m "github.com/ffaann02/cosplace-server/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var users = []m.User{
	{
		ID:       1,
		Username: "user1",
		Email:    "test@gmail.com",
		Password: "123456",
	},
	{
		ID:       2,
		Username: "user2",
		Email:    "test2@gmail.com",
		Password: "123456",
	},
}

// Register handles user registration
func Register(c *fiber.Ctx) error {
	// Parse the request body
	var registerRequest m.RegisterRequest

	// Get input from query or body
	if err := c.BodyParser(&registerRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}
	fmt.Println("Parsed Request:", registerRequest)

	registerRequest.FirstName = c.Query("firstname")
	registerRequest.LastName = c.Query("lastname")
	registerRequest.PhoneNumber = c.Query("phoneNumber")
	registerRequest.DateOfBirth = c.Query("dateOfBirth")
	registerRequest.Username = c.Query("username")
	registerRequest.Email = c.Query("email")
	registerRequest.Password = c.Query("password")

	if err := c.BodyParser(&registerRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	fmt.Println(registerRequest.FirstName)
	fmt.Println(registerRequest.LastName)
	fmt.Println(registerRequest.PhoneNumber)
	fmt.Println(registerRequest.DateOfBirth)
	fmt.Println(registerRequest.Username)
	fmt.Println(registerRequest.Email)
	fmt.Println(registerRequest.Password)

	// Validate input
	if registerRequest.Username == "" || registerRequest.Email == "" || registerRequest.Password == "" ||
		registerRequest.FirstName == "" || registerRequest.LastName == "" || registerRequest.PhoneNumber == "" || registerRequest.DateOfBirth == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "All fields are required",
		})
	}

	// Check if user already exists
	db := config.MysqlDB()
	var existingUser m.User
	if err := db.Where("email = ?", registerRequest.Email).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Email already registered",
		})
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not hash the password",
		})
	}

	// Create the new user
	newUser := m.User{
		Email:       registerRequest.Email,
		Username:    registerRequest.Username,
		Password:    string(hashedPassword),
		FirstName:   registerRequest.FirstName,
		LastName:    registerRequest.LastName,
		PhoneNumber: registerRequest.PhoneNumber,
		DateOfBirth: registerRequest.DateOfBirth,
	}

	fmt.Println(newUser.Password)

	// // Save the user in the database
	if err := db.Create(&newUser).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create user",
		})
	}

	// Return success response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
	})
}

func Login(c *fiber.Ctx) error {
	var loginRequest m.LoginRequest

	// Get input from query or body
	loginRequest.Email = c.Query("email")
	loginRequest.Username = c.Query("username")
	loginRequest.Password = c.Query("password")

	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	// Access the database connection from the config package
	db := config.MysqlDB()

	var user m.User
	// Query user from the database
	err := db.Where("email = ? OR username = ?", loginRequest.Email, loginRequest.Username).First(&user).Error
	fmt.Println(user)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "No account found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not query the database",
		})
	}

	// Verify password (in production, use hashing)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentials",
		})
	}

	// Create JWT tokens (same as your original code)
	accessClaims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Minute * 15).Unix(),
	}

	refreshClaims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 720).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Secret key not found",
		})
	}

	accessT, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not generate access token",
		})
	}

	refreshT, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not generate refresh token",
		})
	}

	// Set cookies for tokens (same as your original code)
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessT,
		Expires:  time.Now().Add(time.Minute * 15),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshT,
		Expires:  time.Now().Add(time.Hour * 720),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None",
	})

	return c.JSON(fiber.Map{
		"message": "Login success",
	})
}

func Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		MaxAge:   -1,
		Expires:  time.Now().Add(-100 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		MaxAge:   -1,
		Expires:  time.Now().Add(-100 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logout Successfully",
		"user":    users,
	})
}

func Refresh(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Refresh token not provided",
		})
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Secret key not found",
		})
	}

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid refresh token",
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token claims",
		})
	}

	if exp, ok := claims["exp"].(float64); ok {
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Refresh token has expired",
			})
		}
	}

	userId := claims["user_id"].(float64)

	var user m.User
	for _, u := range users {
		if u.ID == int64(userId) {
			user = u
			break
		}
	}

	// if user.ID == 0 {
	// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
	// 		"message": "User not found",
	// 	})
	// }

	accessClaims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Minute * 15).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessT, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not generate access token",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessT,
		Expires:  time.Now().Add(time.Minute * 15),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None",
	})

	return c.JSON(fiber.Map{
		"message": "Token refreshed",
	})
}

func CheckAuth(c *fiber.Ctx) error {
	refresh_token := c.Cookies("refresh_token")
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	fmt.Println(refresh_token)

	if refresh_token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	token, err := jwt.Parse(refresh_token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil // Return the secret key for validation
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["user_id"]
		username := claims["username"]

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message":  "Authenticated",
			"user_id":  userID,
			"username": username,
		})
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"message": "Unauthorized",
	})
}
