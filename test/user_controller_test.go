package test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/controllers"
	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/middleware"
	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/models"
	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var user models.User
var app *fiber.App
var db *gorm.DB

func setupTestDatabase() *gorm.DB {
    config := storage.Config{
        Host:     "localhost",
        Port:     "5432",
        User:     "sid",
        Password: "sidd",
        DbName:   "webapp",
        SSLMode:  "disable",
    }

    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        config.Host, config.Port, config.User, config.Password, config.DbName, config.SSLMode,
    )

    
    var err error

    if db!= nil {
        
        sqlDB, dbErr := db.DB()
        if dbErr == nil {
            
            pingErr := sqlDB.Ping()
            if pingErr == nil {
                log.Println("Database connection is already alive")
                return db 
            } else {
                log.Println("Database ping failed, reopening connection...")
            }
        }
    }

    db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to the database: %v", err)
        return nil
    }

    log.Println("Successfully connected to the database")
    return db
}




func TestCreateUser(t *testing.T) {
    
    app= fiber.New()
    storage.Database = setupTestDatabase()
    // setupTest()
    app.Post("/v1/user", controllers.CreateUser)

    t.Run("Valid User Creation", func(t *testing.T) {
        payload := `{
            "email": "test@example.com",
            "password": "password123",
            "first_name": "John",
            "last_name": "Doe"
        }`
        req := httptest.NewRequest("POST", "/v1/user", strings.NewReader(payload))
        req.Header.Set("Content-Type", "application/json")
        resp, _ := app.Test(req)

        assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

        
        
        err := storage.Database.Where("email = ?", "test@example.com").First(&user).Error
        assert.Nil(t, err)
        assert.Equal(t, "John", user.FirstName)
        assert.Equal(t, "Doe", user.LastName)
    })

   
    t.Run("Invalid Email", func(t *testing.T) {
        payload := `{
            "email": "invalid-email",
            "password": "password123",
            "first_name": "John",
            "last_name": "Doe"
        }`
        req := httptest.NewRequest("POST", "/v1/user", strings.NewReader(payload))
        req.Header.Set("Content-Type", "application/json")
        resp, _ := app.Test(req)

        assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
    })

    
    t.Run("Missing Required Fields", func(t *testing.T) {
        payload := `{
            "email": "test@example.com",
            "password": "password123"
        }`
        req := httptest.NewRequest("POST", "/v1/user", strings.NewReader(payload))
        req.Header.Set("Content-Type", "application/json")
        resp, _ := app.Test(req)

        assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
    })

    
    t.Run("Duplicate Email", func(t *testing.T) {
        payload := `{
            "email": "test@example.com",
            "password": "password123",
            "first_name": "Jane",
            "last_name": "Doe"
        }`
        req := httptest.NewRequest("POST", "/v1/user", strings.NewReader(payload))
        req.Header.Set("Content-Type", "application/json")
        resp, _ := app.Test(req)

        assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
    })
}

func TestGetUser(t *testing.T) {
    app= fiber.New()
    storage.Database = setupTestDatabase()

    // pass,er:=helper.HashPassword("password123")
    // if(er!=nil){
    //     log.Print("error hashing password")
    // }

    user = models.User{
        Email:     "test@example.com",
        Password:  "password123", 
        FirstName: "John",
        LastName:  "Doe",
    }

    app.Get("/v1/user/self", middleware.BasicAuthMiddleware, controllers.GetUser)

    createAuthHeader := func(username, password string) string {
        auth := username + ":" + password
        return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
    }

    t.Run("Successful User Retrieval", func(t *testing.T) {
        req := httptest.NewRequest("GET", "/v1/user/self", nil)
        req.Header.Set("Authorization", createAuthHeader(user.Email, user.Password))
        resp, _ := app.Test(req)

        assert.Equal(t, fiber.StatusOK, resp.StatusCode)

        var userResp models.UserResponse
        json.NewDecoder(resp.Body).Decode(&userResp)
        assert.Equal(t, user.Email, userResp.Email)
        assert.Equal(t, user.FirstName, userResp.FirstName)
        assert.Equal(t, user.LastName, userResp.LastName)
    })

    t.Run("Missing Authorization Header", func(t *testing.T) {
        req := httptest.NewRequest("GET", "/v1/user/self", nil)
        resp, _ := app.Test(req)

        assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
    })

    t.Run("Invalid Authorization Method", func(t *testing.T) {
        req := httptest.NewRequest("GET", "/v1/user/self", nil)
        req.Header.Set("Authorization", "Bearer token")
        resp, _ := app.Test(req)

        assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
    })
    t.Run("Invalid Base64 Encoding", func(t *testing.T) {
        req := httptest.NewRequest("GET", "/v1/user/self", nil)
        req.Header.Set("Authorization", "Basic invalid-base64")
        resp, _ := app.Test(req)

        assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
    })
    t.Run("Invalid Credentials Format", func(t *testing.T) {
        req := httptest.NewRequest("GET", "/v1/user/self", nil)
        req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("invalidformat")))
        resp, _ := app.Test(req)

        assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
    })

    t.Run("Non-existent User", func(t *testing.T) {
        req := httptest.NewRequest("GET", "/v1/user/self", nil)
        req.Header.Set("Authorization", createAuthHeader("fake@example.com", "password"))
        resp, _ := app.Test(req)

        assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
    })

    t.Run("Incorrect Password", func(t *testing.T) {
        req := httptest.NewRequest("GET", "/v1/user/self", nil)
        req.Header.Set("Authorization", createAuthHeader(user.Email, "wrongpassword"))
        resp, _ := app.Test(req)

        assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
    })

    t.Run("Request With Payload", func(t *testing.T) {
        req := httptest.NewRequest("GET", "/v1/user/self", strings.NewReader("payload"))
        req.Header.Set("Authorization", createAuthHeader(user.Email, "password"))
        resp, _ := app.Test(req)

        assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
    })

    t.Run("Request With Query Parameters", func(t *testing.T) {
        req := httptest.NewRequest("GET", "/v1/user/self?param=value", nil)
        req.Header.Set("Authorization", createAuthHeader(user.Email, "password"))
        resp, _ := app.Test(req)

        assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
    })
}

func TestMain(m *testing.M) {
    
    code := m.Run()

    if storage.Database != nil {
        db, _ := storage.Database.DB()
        db.Close()
    }

    os.Exit(code)
}