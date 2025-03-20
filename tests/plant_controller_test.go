package tests

import (
	_ "api/controllers"
	_ "api/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/stretchr/testify/assert"
)

type Respons struct {
	Err  bool        `json:"err"`
	Data interface{} `json:"data"`
}

func init() {
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}

func TestPostPlant(t *testing.T) {
	r, _ := http.NewRequest("POST", "/plant", bytes.NewBuffer([]byte(`{"name": "Rose", "type": "Flower"}`)))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", "Bearer test_token")

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp Respons
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Nil(t, err)
	assert.False(t, resp.Err)
}

func TestPostPlantInvalidJSON(t *testing.T) {
	r, _ := http.NewRequest("POST", "/plants", bytes.NewBuffer([]byte(`{"name": "Rose", "type":`))) // Неверный JSON
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", "Bearer test_token")

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp Respons
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Nil(t, err)
	assert.True(t, resp.Err)
	assert.Equal(t, "Invalid JSON format", resp.Data)
}
