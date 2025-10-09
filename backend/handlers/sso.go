package handlers

import (
    "crypto/rand"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"

    "farm-management-backend/config"
    "farm-management-backend/database"
    "farm-management-backend/models"

    "github.com/gofiber/fiber/v2"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
)

func googleOAuthConfig() *oauth2.Config {
    return &oauth2.Config{
        ClientID:     config.GetEnv("GOOGLE_CLIENT_ID", ""),
        ClientSecret: config.GetEnv("GOOGLE_CLIENT_SECRET", ""),
        RedirectURL:  config.GetEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/v1/auth/google/callback"),
        Scopes: []string{
            "openid",
            "email",
            "profile",
        },
        Endpoint: google.Endpoint,
    }
}

func randomState() (string, error) {
    b := make([]byte, 32)
    if _, err := rand.Read(b); err != nil {
        return "", err
    }
    return base64.RawURLEncoding.EncodeToString(b), nil
}

// GoogleLoginRedirect starts the OAuth2 flow
func GoogleLoginRedirect() fiber.Handler {
    return func(c *fiber.Ctx) error {
        conf := googleOAuthConfig()
        if conf.ClientID == "" || conf.ClientSecret == "" {
            return c.Status(500).JSON(fiber.Map{"error": "Google SSO not configured"})
        }

        state, err := randomState()
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "failed to generate state"})
        }
        // store state in cookie for 5 mins
        c.Cookie(&fiber.Cookie{
            Name:     "oauth_state",
            Value:    state,
            HTTPOnly: true,
            Secure:   false,
            SameSite: "Lax",
            Expires:  time.Now().Add(5 * time.Minute),
            Path:     "/",
        })

        url := conf.AuthCodeURL(state, oauth2.AccessTypeOffline)
        return c.Redirect(url, http.StatusTemporaryRedirect)
    }
}

type googleUserInfo struct {
    ID            string `json:"id"`
    Email         string `json:"email"`
    VerifiedEmail bool   `json:"verified_email"`
    Name          string `json:"name"`
    GivenName     string `json:"given_name"`
    FamilyName    string `json:"family_name"`
    Picture       string `json:"picture"`
}

// GoogleLoginCallback handles the OAuth2 callback
func GoogleLoginCallback(db *database.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        conf := googleOAuthConfig()

        // Validate state
        stateCookie := c.Cookies("oauth_state")
        state := c.Query("state")
        if state == "" || stateCookie == "" || state != stateCookie {
            return c.Status(401).JSON(fiber.Map{"error": "invalid oauth state"})
        }

        code := c.Query("code")
        if code == "" {
            return c.Status(400).JSON(fiber.Map{"error": "missing code"})
        }

        token, err := conf.Exchange(c.Context(), code)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{"error": fmt.Sprintf("token exchange failed: %v", err)})
        }

        client := conf.Client(c.Context(), token)
        resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("failed to fetch userinfo: %v", err)})
        }
        defer resp.Body.Close()
        body, _ := io.ReadAll(resp.Body)
        var gu googleUserInfo
        if err := json.Unmarshal(body, &gu); err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "failed to parse userinfo"})
        }

        // Find or create user
        var user models.User
        if err := db.Where("email = ?", gu.Email).First(&user).Error; err != nil {
            user = models.User{FullName: gu.Name, Email: gu.Email, PasswordHash: ""}
            if err := db.Create(&user).Error; err != nil {
                return c.Status(500).JSON(fiber.Map{"error": "failed to create user"})
            }
        }

        // Issue JWT
        jwtToken, err := generateJWT(user.ID.String())
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "failed to generate token"})
        }

        // Redirect back to frontend carrying token
        frontend := config.GetEnv("FRONTEND_URL", "http://localhost:5173")
        redirectURL := fmt.Sprintf("%s/login?token=%s", frontend, jwtToken)
        return c.Redirect(redirectURL, http.StatusTemporaryRedirect)
    }
}


