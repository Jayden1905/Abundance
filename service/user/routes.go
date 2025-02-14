package user

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"

	"github.com/jayden1905/abundance/cmd/pkg/database"
	"github.com/jayden1905/abundance/config"
	"github.com/jayden1905/abundance/service/auth"
	"github.com/jayden1905/abundance/service/email"
	"github.com/jayden1905/abundance/types"
	"github.com/jayden1905/abundance/utils"
)

type Handler struct {
	store  types.UserStore
	mailer email.Mailer
}

func NewHandler(store types.UserStore, mailer email.Mailer) *Handler {
	return &Handler{store: store, mailer: mailer}
}

// RegisterRoutes for Fiber
func (h *Handler) RegisterRoutes(router fiber.Router) {
	rateLimiterEmailVerification := auth.CreateRateLimiter(1, 5*time.Minute, "We have sent you a verification email. Please check your inbox and spam folder.")

	router.Post("/user/auth/login", auth.BlockIfAuthenticated(h.handleLogin))
	router.Post("/user/auth/logout", h.handleLogout)
	router.Post("/user/register", h.handleRegister)
	router.Patch("/user/super-user", h.handleCreateSuperUser)
	router.Get("/users", auth.WithJWTAuth(h.handleGetUsersPaginated, h.store))
	router.Get("/user/:id", auth.WithJWTAuth(h.handleGetUserByID, h.store))
	router.Delete("/user/:id", auth.WithJWTAuth(h.handleDeleteUser, h.store))
	router.Get("/user/auth/status", h.handleIsAuthenticated)
	router.Get("/user/verify/email", h.handleVerifyAccount)
	router.Post("/user/verify/email/resend", rateLimiterEmailVerification, h.handleResendVerificationEmail)
}

// Handler for registering a new user
func (h *Handler) handleRegister(c *fiber.Ctx) error {
	// Parse JSON payload
	var payload types.RegisterUserPayload

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	// Validate the payload
	invalidFields, validationErr := utils.ValidatePayload(payload)
	if validationErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":          "Invalid payload",
			"invalid_fields": invalidFields,
		})
	}

	// Check if the user already exists
	_, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("User with email %s already exists", payload.Email)})
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error hashing password"})
	}

	// Create a new user with unverified status
	err = h.store.CreateUser(c.Context(), &database.User{
		Username:       payload.Username,
		Email:          payload.Email,
		PasswordHash:   hashedPassword,
		RoleID:         utils.ConvertRoleStringToRoleID(payload.Role),
		SubscriptionID: utils.ConvertSubscriptionStringToSubscriptionID(payload.Subscription),
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Generate a verification token
	token, err := auth.GenerateVerificationToken(payload.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Send email asynchronously
	go func() {
		err = h.mailer.SendVerificationEmail(payload.Email, token)
		if err != nil {
			fmt.Printf("Error sending verification email: %v\n", err)
		}
	}()

	// Return success response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"email":   payload.Email,
		"status":  "verification email sent",
	})
}

// Handler for resending verification email
func (h *Handler) handleResendVerificationEmail(c *fiber.Ctx) error {
	var payload types.ResendVerificationEmailPayload

	// Parse JSON payload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	// Validate the payload
	invalidFields, validationErr := utils.ValidatePayload(payload)
	if validationErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":          "Invalid payload",
			"invalid_fields": invalidFields,
		})
	}

	// Check if the user already exists and not verified
	user, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("User with email %s does not exist", payload.Email)})
	}

	if user.IsVerified {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("User with email %s is already verified", payload.Email)})
	}

	// Generate a verification token
	token, err := auth.GenerateVerificationToken(payload.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Send email asynchronously
	go func() {
		err = h.mailer.SendVerificationEmail(payload.Email, token)
		if err != nil {
			fmt.Printf("Error sending verification email: %v\n", err)
		}
	}()

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Verification email sent",
	})
}

// Handler for verifying a user
func (h *Handler) handleVerifyAccount(c *fiber.Ctx) error {
	tokenString := c.Query("token")
	if tokenString == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Token is missing"})
	}

	// Validate the verification token and return email
	email, err := auth.ValidateVerificationToken(tokenString)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Error validating verification token: %v", err)})
	}

	// Get the user by email
	user, err := h.store.GetUserByEmail(email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Error getting user by email: %v", err)})
	}

	if user.IsVerified {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User is already verified"})
	}

	// Update the user verification status
	if err := h.store.UpdateUserVerification(c.Context(), user.ID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Error updating user verification status: %v", err)})
	}

	return c.Redirect(config.Envs.PublicHost+"/", fiber.StatusSeeOther)
}

// Hanlder for login
func (h *Handler) handleLogin(c *fiber.Ctx) error {
	// Parse JSON payload
	var payload types.LoginUserPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	// Validate the payload
	invalidFields, err := utils.ValidatePayload(payload)
	if err != nil {
		// Return the invalid fields if validation fails
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":          "Invalid payload",
			"invalid_fields": invalidFields,
		})
	}

	// Check if the user exists by email
	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "You don't have an account. Please register"})
	}

	if !auth.ComparePasswords(u.PasswordHash, []byte(payload.Password)) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email or password is incorrect"})
	}

	if !u.IsVerified {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Please verify your email"})
	}

	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, int(u.ID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		HTTPOnly: true,                     // Disallow JS access to the cookie
		Secure:   config.Envs.ISProduction, // Set to true in production (HTTPS)
		SameSite: "Lax",                    // Prevent CSRF attacks
		Path:     "/",                      // Valid for the entire site
		MaxAge:   int(config.Envs.JWTExpirationInSeconds),
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token, "expires_in": fmt.Sprintf("%d", config.Envs.JWTExpirationInSeconds)})
}

// Handler for logout
func (h *Handler) handleLogout(c *fiber.Ctx) error {
	// Clear the token cookie by setting an expired cookie
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,                     // Disallow JS access to the cookie
		Secure:   config.Envs.ISProduction, // Set to true in production (HTTPS)
		SameSite: "Lax",                    // Prevent CSRF attacks
		Path:     "/",                      // Valid for the entire site
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Logged out successfully"})
}

// Handler for creating a super user
func (h *Handler) handleCreateSuperUser(c *fiber.Ctx) error {
	var payload types.RegisterUserPayload

	// Parse JSON payload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	// Validate the payload
	invalidFields, err := utils.ValidatePayload(payload)
	if err != nil {
		// Return the invalid fields if validation fails
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":          "Invalid payload",
			"invalid_fields": invalidFields,
		})
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error hashing password"})
	}

	// Check if the user already exists
	u, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		// compare password
		if !auth.ComparePasswords(u.PasswordHash, []byte(payload.Password)) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email or password is incorrect"})
		}

		// check if user is already a super user
		role, err := h.store.GetUserRoleByID(u.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Error getting user role by id: %v", err)})
		}
		if role == "super_user" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User is already a super user"})
		}

		// update user to super user
		h.store.UpdateUserToSuperUser(c.Context(), u.ID)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "User updated to super user successfully",
		})
	}

	// Create a new super user
	err = h.store.CreateSuperUser(c.Context(), &types.User{
		Username:     payload.Username,
		Email:        payload.Email,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Super user created successfully"})
}

// Handler for deleting a user
func (h *Handler) handleDeleteUser(c *fiber.Ctx) error {
	// get user id from context
	userID := auth.GetUserIDFromContext(c)

	// check if user is a super user
	superUser, err := utils.IsSuperUser(userID, h.store)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Error getting user role by id: %v", err)})
	}
	if !superUser {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
	}

	paramsID := c.Params("id")
	// convert id to int
	intID, err := strconv.Atoi(paramsID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid id: %v", err)})
	}
	// convert id to int32
	id := int32(intID)

	// Check if the user exists in the database
	_, err = h.store.GetUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("%v", err)})
	}

	// Check if the user is super_user and is trying to delete themselves
	if id == userID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "You cannot delete yourself"})
	}

	// delete user
	if err := h.store.DeleteUserByID(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Error deleting user by id: %v", err)})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User deleted successfully"})
}

// Handler for getting all users by page
func (h *Handler) handleGetUsersPaginated(c *fiber.Ctx) error {
	userID := auth.GetUserIDFromContext(c)

	// Checking if user have access
	superUser, err := utils.IsSuperUser(userID, h.store)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Error getting user role by id: %v", err)})
	}
	if !superUser {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
	}

	const (
		defaultPageSize = 10
		maxPageSize     = 100
	)

	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")

	page := 1
	pageSize := defaultPageSize

	// Parse page if provided
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Parse pageSize if provided
	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= maxPageSize {
			pageSize = ps
		}
	}

	users, err := h.store.GetUsersPaginated(int32(page), int32(pageSize))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("Error getting users by page and page size")
	}

	return c.Status(fiber.StatusOK).JSON(users)
}

// Handler for getting a user by id
func (h *Handler) handleGetUserByID(c *fiber.Ctx) error {
	userID := auth.GetUserIDFromContext(c)

	superUser, err := utils.IsSuperUser(userID, h.store)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Error getting user role by id: %v", err)})
	}
	if !superUser {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
	}

	stringID := c.Params("id")

	// convert id to int
	intID, err := strconv.Atoi(stringID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid id: %v", err)})
	}

	// convert id to int32
	id := int32(intID)

	u, err := h.store.GetUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Error getting user by id: %v", err)})
	}

	return c.Status(fiber.StatusOK).JSON(u)
}

// Handler for checking if a user is authenticated
func (h *Handler) handleIsAuthenticated(c *fiber.Ctx) error {
	tokenString := c.Cookies("token")

	if tokenString == "" {
		// get token from Authorization header
		tokenString = c.Get("Authorization")
	}

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token is missing"})
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.Envs.JWTSecret), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Error parsing token"})
	}

	if !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token is invalid"})
	}

	// get user id from token
	claims := token.Claims.(jwt.MapClaims)
	str := claims["userID"].(string)

	userIDInt, err := strconv.Atoi(str)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error getting user id from token"})
	}

	userID := int32(userIDInt)

	// get if user exists
	user, err := h.store.GetUserByID(int32(userID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Error getting user by id: %v", err)})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": user,
	})
}

// Handler for updating user password
func (h *Handler) handleUpdateUserPassword(c *fiber.Ctx) error {
	userID := auth.GetUserIDFromContext(c)

	// Parse JSON payload
	var payload types.UpdateUserPasswordPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	// Validate the payload
	invalidFields, err := utils.ValidatePayload(payload)
	if err != nil {
		// Return the invalid fields if validation fails
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":          "Invalid payload",
			"invalid_fields": invalidFields,
		})
	}

	// Check if the user exists
	u, err := h.store.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User not found"})
	}

	// Check if the old password is correct
	if !auth.ComparePasswords(u.PasswordHash, []byte(payload.OldPassword)) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Old password is incorrect"})
	}

	// Hash the new password
	hashedPassword, err := auth.HashPassword(payload.NewPassword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error hashing password"})
	}

	// Update the user password
	if err := h.store.UpdateUserPassword(c.Context(), userID, hashedPassword); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Error updating user password: %v", err)})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Password updated successfully"})
}
