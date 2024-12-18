package handler

import (
	"fmt"
	"os"
	"time"

	"github.com/ffaann02/cosplace-server/api/helper"
	config "github.com/ffaann02/cosplace-server/internal/config"
	m "github.com/ffaann02/cosplace-server/internal/model"
	v "github.com/ffaann02/cosplace-server/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Register handles user registration
func Register(c *fiber.Ctx) error {
	// Parse the request body
	registerRequest := new(m.RegisterRequest)

	// Get input from query or body
	if err := c.BodyParser(&registerRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}
	fmt.Println("Parsed Request:", registerRequest)

	registerRequest.Username = c.Query("username")
	// Init display name with username
	registerRequest.Email = c.Query("email")
	registerRequest.Password = c.Query("password")
	registerRequest.FirstName = c.Query("firstname")
	registerRequest.LastName = c.Query("lastname")
	registerRequest.DateOfBirth = c.Query("date_of_birth")
	registerRequest.PhoneNumber = c.Query("phone_number")
	registerRequest.Gender = c.Query("gender")

	if err := c.BodyParser(&registerRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid input",
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
	valid, missingField, err := v.ValidateStruct(*registerRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if !valid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Missing field: %s", missingField),
		})
	}

	// Check if user already exists
	db := config.MysqlDB()
	tx := db.Begin()
	var existingUser m.User
	if err := tx.Where("email = ?", registerRequest.Email).First(&existingUser).Error; err == nil {
		tx.Rollback()
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "อีเมลนี้ถูกใช้งานแล้ว",
		})
	}

	if err := tx.Where("username = ?", registerRequest.Username).First(&existingUser).Error; err == nil {
		tx.Rollback()
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "ชื่อผู้ใช้งานนี้ถูกใช้งานแล้ว",
		})
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "เกิดข้อผิดพลาดกับรหัสผ่าน",
		})
	}

	userID, err := helper.GenerateNewUserID(tx)
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "เกิดข้อผิดพลาดในการสร้างบัญชีผู้ใช้งาน",
		})
	}
	// Create the new user
	newUser := m.User{
		UserID:      userID,
		Email:       registerRequest.Email,
		Username:    registerRequest.Username,
		Password:    string(hashedPassword),
		FirstName:   registerRequest.FirstName,
		LastName:    registerRequest.LastName,
		PhoneNumber: registerRequest.PhoneNumber,
		DateOfBirth: registerRequest.DateOfBirth,
		Gender:      registerRequest.Gender,
	}

	// // Save the user in the database
	if err := db.Create(&newUser).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "เกิดข้อผิดพลาดในการสร้างบัญชีผู้ใช้งาน",
		})
	}

	profileID, err := helper.GenerateNewProfileID(tx)
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "เกิดข้อผิดพลาดในการสร้างโปรไฟล์ผู้ใช้งาน",
		})
	}

	newProfile := m.Profile{
		ProfileID:   profileID,
		UserID:      userID,
		DisplayName: registerRequest.Username,
	}

	if err := db.Create(&newProfile).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "เกิดข้อผิดพลาดในการสร้างโปรไฟล์ผู้ใช้งาน",
		})
	}

	tx.Commit()

	// Return success response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "สมัครบัญชีผู้ใช้งานเสร็จสิ้น",
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

	fmt.Println(loginRequest)

	// Access the database connection from the config package
	db := config.MysqlDB()

	var user m.User
	// Query user from the database
	err := db.Where("email = ? OR username = ?", loginRequest.Email, loginRequest.Username).First(&user).Error
	fmt.Println(user)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "ไม่พบบัญชีผู้ใช้งาน โปรดลงทะเบียนหรือตรวจสอบข้อมูลใหม่อีกครั้ง",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "เกิดข้อผิดพลาดบางอย่าง โปรดลองอีกครั้ง",
		})
	}

	// Verify password (in production, use hashing)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "รหัสผ่านไม่ถูกต้อง กรุณาลองใหม่อีกครั้ง",
		})
	}

	var role string
	var seller m.Seller
	if err := db.Where("user_id = ?", user.UserID).First(&seller).Error; err != nil {
		role = "customer"
		seller.SellerID = ""
	} else {
		role = "seller"
	}

	// Create JWT tokens (same as your original code)
	accessClaims := jwt.MapClaims{
		"user_id":   user.UserID,
		"username":  user.Username,
		"role":      role,
		"seller_id": seller.SellerID,
		"exp":       time.Now().Add(time.Minute * 15).Unix(),
	}

	refreshClaims := jwt.MapClaims{
		"user_id":   user.UserID,
		"username":  user.Username,
		"role":      role,
		"seller_id": seller.SellerID,
		"exp":       time.Now().Add(time.Hour * 720).Unix(),
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":       "Login success",
		"user_id":       user.UserID,
		"username":      user.Username,
		"role":          role,
		"seller_id":     seller.SellerID,
		"access_token":  accessT,
		"refresh_token": refreshT,
	})
}

func Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/", // Try to fix logout on production
		Domain:   "",  // Try to fix logout on production
		MaxAge:   -1,
		Expires:  time.Now().Add(-100 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None", // Try to fix logout on production
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/", // Try to fix logout on production
		Domain:   "",  // Try to fix logout on production
		MaxAge:   -1,
		Expires:  time.Now().Add(-100 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None", // Try to fix logout on production
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logout Successfully",
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

	var user m.User
	var role string
	var seller m.Seller
	db := config.MysqlDB()
	if err := db.Where("user_id = ?", user.UserID).First(&seller).Error; err != nil {
		role = "customer"
		seller.SellerID = ""
	} else {
		role = "seller"
	}

	accessClaims := jwt.MapClaims{
		"user_id":   user.UserID,
		"username":  user.Username,
		"role":      role,
		"seller_id": seller.SellerID,
		"exp":       time.Now().Add(time.Minute * 15).Unix(),
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
		role := claims["role"]
		sellerID := claims["seller_id"]

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message":   "Authenticated",
			"user_id":   userID,
			"username":  username,
			"role":      role,
			"seller_id": sellerID,
		})
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"message": "Unauthorized",
	})
}

func ChangePassword(c *fiber.Ctx) error {
	var passwordRequest struct {
		UserID             string `json:"user_id"` // Assuming user_id is set in context
		OldPassword        string `json:"old_password"`
		NewPassword        string `json:"new_password"`
		ConfirmNewPassword string `json:"confirm_password"`
	}

	if err := c.BodyParser(&passwordRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
			"error":   "เกิดข้อผิดพลาดในการรับข้อมูล",
		})
	}

	if passwordRequest.NewPassword != passwordRequest.ConfirmNewPassword {
		fmt.Println(passwordRequest.NewPassword)
		fmt.Println(passwordRequest.ConfirmNewPassword)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "New passwords do not match",
			"error":   "รหัสผ่านใหม่ไม่ตรงกัน",
		})
	}

	db := config.MysqlDB()
	userID := passwordRequest.UserID

	var user m.User
	if err := db.Where("user_id = ?", userID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "User not found",
			"error":   "ไม่พบผู้ใช้งาน",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwordRequest.OldPassword)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Old password does not match",
			"error":   "รหัสผ่านเดิมไม่ถูกต้อง",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordRequest.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error encrypting new password",
			"error":   "เกิดข้อผิดพลาดในการเข้ารหัสรหัสผ่านใหม่",
		})
	}

	user.Password = string(hashedPassword)
	if err := db.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error saving new password",
			"error":   "เกิดข้อผิดพลาดในการบันทึกรหัสผ่านใหม่",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "Password changed successfully",
		"message_th": "เปลี่ยนรหัสผ่านสำเร็จ",
	})
}
