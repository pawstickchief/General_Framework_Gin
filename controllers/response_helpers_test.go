package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func TestResCodeMsg(t *testing.T) {
	if got := CodeInvalidParam.Msg(); got != codeMsgMap[CodeInvalidParam] {
		t.Fatalf("expected %s, got %s", codeMsgMap[CodeInvalidParam], got)
	}

	var unknownCode ResCode = 99999
	if got := unknownCode.Msg(); got != codeMsgMap[CodeServerBusy] {
		t.Fatalf("expected fallback message %s, got %s", codeMsgMap[CodeServerBusy], got)
	}
}

func TestResponseError(t *testing.T) {
	c, w := setupTestContext()
	ResponseError(c, CodeInvalidParam)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var rd ResponseData
	if err := json.Unmarshal(w.Body.Bytes(), &rd); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if rd.Code != CodeInvalidParam || rd.Msg != codeMsgMap[CodeInvalidParam] || rd.Data != nil {
		t.Fatalf("unexpected response payload: %#v", rd)
	}
}

func TestResponseSuccess(t *testing.T) {
	c, w := setupTestContext()
	payload := map[string]string{"hello": "world"}
	ResponseSuccess(c, payload)

	var rd ResponseData
	if err := json.Unmarshal(w.Body.Bytes(), &rd); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if rd.Code != CodeSuccess {
		t.Fatalf("expected code %d, got %d", CodeSuccess, rd.Code)
	}

	if rd.Msg != codeMsgMap[CodeSuccess] {
		t.Fatalf("expected success message, got %v", rd.Msg)
	}

	if rd.Data == nil {
		t.Fatalf("expected data to be present")
	}
}

func TestResponseErrorWithMsg(t *testing.T) {
	c, w := setupTestContext()
	ResponseErrorWithMsg(c, CodeInvalidParam, "custom error")

	var rd ResponseData
	if err := json.Unmarshal(w.Body.Bytes(), &rd); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if rd.Msg != "custom error" {
		t.Fatalf("expected custom message, got %v", rd.Msg)
	}
}

func TestResponseSystemDataSuccess(t *testing.T) {
	c, w := setupTestContext()
	data := []int{1, 2, 3}
	ResponseSystemDataSuccess(c, data)

	var rd ResponseData
	if err := json.Unmarshal(w.Body.Bytes(), &rd); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if rd.Code != CodeSuccess || rd.Msg != codeMsgMap[CodeSuccess] {
		t.Fatalf("unexpected response payload: %#v", rd)
	}

	if rd.Data == nil {
		t.Fatalf("expected data to be present")
	}
}
