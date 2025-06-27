package handlers

import (
    "database/sql"
    "myapp/internal/models"
    "net/http"
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
)

type UserHandler struct {
    db *sql.DB
}

func NewUserHandler(db *sql.DB) *UserHandler {
    return &UserHandler{db: db}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
    rows, err := h.db.Query("SELECT id, name, email, created_at, updated_at FROM users ORDER BY id")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    var users []models.User
    for rows.Next() {
        var user models.User
        err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        users = append(users, user)
    }

    c.JSON(http.StatusOK, users)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()

    err := h.db.QueryRow(
        "INSERT INTO users (name, email, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id",
        user.Name, user.Email, user.CreatedAt, user.UpdatedAt,
    ).Scan(&user.ID)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetUser(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    var user models.User
    err = h.db.QueryRow(
        "SELECT id, name, email, created_at, updated_at FROM users WHERE id = $1", id,
    ).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)

    if err == sql.ErrNoRows {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user.ID = id
    user.UpdatedAt = time.Now()

    _, err = h.db.Exec(
        "UPDATE users SET name = $1, email = $2, updated_at = $3 WHERE id = $4",
        user.Name, user.Email, user.UpdatedAt, id,
    )

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    _, err = h.db.Exec("DELETE FROM users WHERE id = $1", id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}