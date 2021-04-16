package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{\"message\":\"pong\"}", w.Body.String())
}

func TestPostTaskRoute(t *testing.T) {
	router := setupRouter()

	t.Run("タスク情報を登録したいが、登録する情報が指定されていないためエラー", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/task", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "{\"error\":\"missing form body\"}", w.Body.String())
	})

	t.Run("タスク情報を登録できる", func(t *testing.T) {
		w := httptest.NewRecorder()

		jsonBody := []byte(`{"content":"hoge","done":false}`)
		req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(jsonBody))
		req.Header.Add("content-type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var registered Task
		err := json.Unmarshal(w.Body.Bytes(), &registered)
		if assert.NoError(t, err) {
			assert.Equal(t, "hoge", registered.Content)
			assert.Equal(t, false, registered.Done)
			assert.NotEmpty(t, registered.ID)
		}

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", fmt.Sprintf("/task/%s", registered.ID), nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var getted Task
		err = json.Unmarshal(w.Body.Bytes(), &getted)
		if assert.NoError(t, err) {
			assert.Equal(t, registered.Content, getted.Content)
			assert.Equal(t, registered.Done, getted.Done)
			assert.Equal(t, registered.ID, getted.ID)
		}
	})
}

func TestPutTaskRoute(t *testing.T) {
	router := setupRouter()

	t.Run("タスク情報を更新したいが、更新する情報が指定されていないためエラー", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/task/a", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "{\"error\":\"missing form body\"}", w.Body.String())
	})

	t.Run("タスク情報を更新できる", func(t *testing.T) {
		w := httptest.NewRecorder()

		jsonBody := []byte(`{"content":"hoge","done":false}`)
		req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(jsonBody))
		req.Header.Add("content-type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var registered Task
		err := json.Unmarshal(w.Body.Bytes(), &registered)
		if assert.NoError(t, err) {
			assert.Equal(t, "hoge", registered.Content)
			assert.Equal(t, false, registered.Done)
			assert.NotEmpty(t, registered.ID)
		}

		w = httptest.NewRecorder()

		jsonBody2 := []byte(`{"content":"hoge2","done":true}`)
		req, _ = http.NewRequest("PUT", fmt.Sprintf("/task/%s", registered.ID), bytes.NewBuffer(jsonBody2))
		req.Header.Add("content-type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var updated Task
		err = json.Unmarshal(w.Body.Bytes(), &updated)
		if assert.NoError(t, err) {
			assert.Equal(t, "hoge2", updated.Content)
			assert.Equal(t, true, updated.Done)
			assert.Equal(t, registered.ID, updated.ID)
		}
	})
}

func TestDeleteTaskRoute(t *testing.T) {
	router := setupRouter()
	t.Run("タスク情報を削除できる", func(t *testing.T) {
		w := httptest.NewRecorder()

		jsonBody := []byte(`{"content":"hoge","done":false}`)
		req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(jsonBody))
		req.Header.Add("content-type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var registered Task
		err := json.Unmarshal(w.Body.Bytes(), &registered)
		if assert.NoError(t, err) {
			assert.Equal(t, "hoge", registered.Content)
			assert.Equal(t, false, registered.Done)
			assert.NotEmpty(t, registered.ID)
		}

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("DELETE", fmt.Sprintf("/task/%s", registered.ID), nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", fmt.Sprintf("/task/%s", registered.ID), nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
