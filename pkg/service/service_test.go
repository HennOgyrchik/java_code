package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"java_code/pkg/db/psql"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestService_Wallet(t *testing.T) {
	ctx := context.Background()

	testDB := psql.New("postgres://postgres:123@192.168.31.197:5432/postgres?connect_timeout=5&sslmode=disable", 50*time.Second)
	testDB.Start()
	app := New(ctx, &testDB)

	r := gin.Default()
	r.POST("/api/v1/wallet", app.Wallet)

	id, _ := uuid.FromString("123e4567-e89b-12d3-a456-426655440000")

	type data struct {
		WalletId      uuid.UUID
		OperationType string
		Amount        float64
	}
	tests := []struct {
		data     data
		expected int
	}{

		{data{WalletId: id, OperationType: "deposit", Amount: 100}, http.StatusOK},
		{data{id, "fqwfqawsd", 100}, http.StatusBadRequest},
		{data{id, "WITHDRAW", 50616891685460}, http.StatusBadRequest},
		{data{id, "witHdraW", 100}, http.StatusOK},
	}

	for _, tt := range tests {
		jsonData, _ := json.Marshal(tt.data)

		w := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)
		assert.Equal(t, tt.expected, w.Code)

	}

}

func TestService_Wallets(t *testing.T) {
	ctx := context.Background()

	testDB := psql.New("postgres://postgres:123@192.168.31.197:5432/postgres?connect_timeout=5&sslmode=disable", 50*time.Second)
	testDB.Start()
	app := New(ctx, &testDB)

	r := gin.Default()
	r.GET("/api/v1/wallets/:uuid", app.Wallets)

	tests := []struct {
		url      string
		expected int
	}{
		{"123e4567-e89b-12d3-a456-426655440000", 200},
		{"gwsgdafb wgfwfg", 400},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/wallets/"+test.url), nil)

		r.ServeHTTP(w, req)
		assert.Equal(t, test.expected, w.Code)
	}

}

func TestStress(t *testing.T) {
	ctx := context.Background()

	testDB := psql.New("postgres://postgres:123@127.0.0.1:5432/postgres?connect_timeout=50&sslmode=disable", 5*time.Second)
	testDB.Start()

	defer testDB.Stop()

	app := New(ctx, &testDB)

	r := gin.Default()
	r.POST("/api/v1/wallet", app.Wallet)

	id, _ := uuid.FromString("123e4567-e89b-12d3-a456-426655440000")

	type data struct {
		WalletId      uuid.UUID
		OperationType string
		Amount        float64
	}

	jsonData, _ := json.Marshal(data{WalletId: id, OperationType: "deposit", Amount: 1})

	n := 1

	var wg sync.WaitGroup
	wg.Add(n)

	startTime := time.Now()

	for i := 0; i < n; i++ {

		go func() {

			w := httptest.NewRecorder()

			req, _ := http.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
			wg.Done()
		}()

	}

	wg.Wait()

	endTime := time.Now()
	fmt.Println(endTime.Sub(startTime))
}
