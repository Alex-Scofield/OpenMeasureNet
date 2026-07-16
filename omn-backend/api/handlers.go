package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type SignupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Code     string `json:"code" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type NodeResponse struct {
	ID           int    `json:"id"`
	DashboardURL string `json:"dashboard_url"`
}

func SignupHandler(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var codeID int
	err := pool.QueryRow(context.Background(),
		"SELECT id FROM invite_codes WHERE code = $1 AND used = FALSE", req.Code).
		Scan(&codeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or used invite code"})
		return
	}

	hash, err := HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	var userID int
	err = pool.QueryRow(context.Background(),
		"INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id",
		req.Email, hash).Scan(&userID)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		return
	}

	_, err = pool.Exec(context.Background(),
		"UPDATE invite_codes SET used = TRUE, used_by = $1 WHERE id = $2", userID, codeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark code as used"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "account created"})
}

func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userID int
	var hash string
	err := pool.QueryRow(context.Background(),
		"SELECT id, password_hash FROM users WHERE email = $1", req.Email).
		Scan(&userID, &hash)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := CheckPassword(hash, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := GenerateToken(userID, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func LogoutHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func GetNodesHandler(c *gin.Context) {
	userID := c.GetInt("user_id")

	rows, err := pool.Query(context.Background(),
		"SELECT id, dashboard_uid FROM nodes WHERE user_id = $1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query nodes"})
		return
	}
	defer rows.Close()

	grafanaURL := os.Getenv("GRAFANA_URL")
	var nodes []NodeResponse

	for rows.Next() {
		var id int
		var uid string
		if err := rows.Scan(&id, &uid); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to scan node"})
			return
		}
		nodes = append(nodes, NodeResponse{
			ID:           id,
			DashboardURL: grafanaURL + "/d/" + uid + "?orgId=1&kiosk",
		})
	}

	if nodes == nil {
		nodes = []NodeResponse{}
	}

	c.JSON(http.StatusOK, nodes)
}

func GetNodeHandler(c *gin.Context) {
	var id int
	if _, err := fmt.Sscan(c.Param("id"), &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid node id"})
		return
	}

	var uid string
	err := pool.QueryRow(context.Background(),
		"SELECT dashboard_uid FROM nodes WHERE id = $1", id).Scan(&uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "node not found"})
		return
	}

	grafanaURL := os.Getenv("GRAFANA_URL")
	c.JSON(http.StatusOK, NodeResponse{
		ID:           id,
		DashboardURL: grafanaURL + "/d/" + uid + "?orgId=1&kiosk",
	})
}

func GetMapHandler(c *gin.Context) {
	grafanaURL := os.Getenv("GRAFANA_URL")
	c.JSON(http.StatusOK, gin.H{
		"dashboard_url": grafanaURL + "/d/all-nodes?orgId=1&kiosk",
	})
}
