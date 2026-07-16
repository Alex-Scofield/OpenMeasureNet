package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

type Observation struct {
	QuantityID uint8
	Value      float32
	Timestamp  float64
	Longitude  float32
	Latitude   float32
}

func main() {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_DB"),
	)

	var err error
	pool, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	if err = pool.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	r := gin.Default()
	r.POST("/upload", handleUpload)

	log.Println("Bulk service listening on :8081")
	if err := r.Run(":8081"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func handleUpload(c *gin.Context) {
	nodeIDStr := c.PostForm("node_id")
	nodeID, err := strconv.Atoi(nodeIDStr)
	if err != nil || nodeID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid node_id"})
		return
	}

	auth := c.GetHeader("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
		return
	}
	token := strings.TrimPrefix(auth, "Bearer ")

	var storedPassword string
	err = pool.QueryRow(context.Background(),
		"SELECT password FROM nodes WHERE id = $1", nodeID).
		Scan(&storedPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid node"})
		return
	}
	if token != storedPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing file"})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read file"})
		return
	}

	if len(data)%21 != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("file size %d is not a multiple of 21 bytes", len(data)),
		})
		return
	}

	if len(data) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "uploaded", "count": 0})
		return
	}

	observations := parseObservations(data)
	count := bulkInsert(observations, nodeID)

	c.JSON(http.StatusOK, gin.H{"message": "uploaded", "count": count})
}

func parseObservations(data []byte) []Observation {
	n := len(data) / 21
	observations := make([]Observation, n)
	for i := 0; i < n; i++ {
		offset := i * 21
		observations[i] = Observation{
			QuantityID: data[offset],
			Value:      math.Float32frombits(binary.BigEndian.Uint32(data[offset+1 : offset+5])),
			Timestamp:  math.Float64frombits(binary.BigEndian.Uint64(data[offset+5 : offset+13])),
			Longitude:  math.Float32frombits(binary.BigEndian.Uint32(data[offset+13 : offset+17])),
			Latitude:   math.Float32frombits(binary.BigEndian.Uint32(data[offset+17 : offset+21])),
		}
	}
	return observations
}

const batchSize = 500

func bulkInsert(observations []Observation, nodeID int) int {
	total := 0
	for i := 0; i < len(observations); i += batchSize {
		end := i + batchSize
		if end > len(observations) {
			end = len(observations)
		}
		batch := observations[i:end]

		args := make([]interface{}, 0, len(batch)*6)
		rows := make([]string, 0, len(batch))
		for j, obs := range batch {
			idx := j * 6
			args = append(args, obs.Timestamp, nodeID, int(obs.QuantityID), obs.Value, obs.Longitude, obs.Latitude)
			rows = append(rows, fmt.Sprintf("(to_timestamp($%d),$%d,$%d,$%d,ST_SetSRID(ST_MakePoint($%d,$%d),4326))",
				idx+1, idx+2, idx+3, idx+4, idx+5, idx+6))
		}

		query := "INSERT INTO observations (time,node_id,quantity,value,location) VALUES " +
			strings.Join(rows, ",") +
			" ON CONFLICT ON CONSTRAINT observations_pkey DO NOTHING"

		tag, err := pool.Exec(context.Background(), query, args...)
		if err != nil {
			log.Printf("Batch insert error: %v", err)
			continue
		}
		total += int(tag.RowsAffected())
	}
	return total
}
