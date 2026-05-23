package response

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
)

type testResp struct {
	Foo string `json:"foo"`
}

func TestNewResponses_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	data := testResp{Foo: "bar"}
	meta := NewMetaOffset(1, 10, 100)

	err := NewResponses(ctx, 200, "ok", data, nil, meta)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var resp baseResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if resp.Message != "ok" {
		t.Errorf("expected message 'ok', got %s", resp.Message)
	}
	if resp.Error != nil {
		t.Errorf("expected nil error, got %v", resp.Error)
	}
	if resp.Meta == nil {
		t.Error("expected meta to be set")
	}
}

func TestNewResponses_ErrorStatus(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := errors.New("something went wrong")
	meta := NewMetaOffset(1, 10, 100)

	if apiErr := NewResponses[testResp](ctx, 500, "fail", testResp{}, err, meta); apiErr != nil {
		t.Fatalf("unexpected error: %v", apiErr)
	}

	var resp baseResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if resp.Error == nil {
		t.Error("expected error to be set")
	}
	if resp.Error != "internal server error" {
		t.Errorf("expected 'internal server error', got %v", resp.Error)
	}
	if resp.Meta != nil {
		t.Errorf("expected nil meta for error status, got %v", resp.Meta)
	}
}

func TestNewResponses_HTTPError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	httpErr := echo.NewHTTPError(404, "not found")

	if apiErr := NewResponses[testResp](ctx, 404, "fail", testResp{}, httpErr, nil); apiErr != nil {
		t.Fatalf("unexpected error: %v", apiErr)
	}

	var resp baseResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if resp.Error != "not found" {
		t.Errorf("expected 'not found', got %v", resp.Error)
	}
}

func TestNewMetaOffset(t *testing.T) {
	meta := NewMetaOffset(1, 10, 95)

	if meta.Page != 1 {
		t.Errorf("expected Page 1, got %d", meta.Page)
	}
	if meta.Limit != 10 {
		t.Errorf("expected Limit 10, got %d", meta.Limit)
	}
	if meta.TotalPage != 10 {
		t.Errorf("expected TotalPage 10, got %d", meta.TotalPage)
	}
	if meta.TotalData != 95 {
		t.Errorf("expected TotalData 95, got %d", meta.TotalData)
	}
}

func TestNewMetaOffset_EvenDivision(t *testing.T) {
	meta := NewMetaOffset(1, 10, 100)

	if meta.TotalPage != 10 {
		t.Errorf("expected TotalPage 10, got %d", meta.TotalPage)
	}
}

func TestNewMetaOffset_ZeroData(t *testing.T) {
	meta := NewMetaOffset(1, 10, 0)

	if meta.TotalPage != 0 {
		t.Errorf("expected TotalPage 0, got %d", meta.TotalPage)
	}
	if meta.TotalData != 0 {
		t.Errorf("expected TotalData 0, got %d", meta.TotalData)
	}
}

func TestNewMetaCursor_HasNextFalse(t *testing.T) {
	meta := NewMetaCursor(10, false, nil)

	if meta.Limit != 10 {
		t.Errorf("expected Limit 10, got %d", meta.Limit)
	}
	if meta.HasNext != false {
		t.Errorf("expected HasNext false, got %v", meta.HasNext)
	}
	if meta.NextCursor != "" {
		t.Errorf("expected empty NextCursor, got %s", meta.NextCursor)
	}
}

func TestNewMetaCursor_HasNextTrue(t *testing.T) {
	cursorData := Cursor{ID: 42, SortValue: time.Now()}

	meta := NewMetaCursor(10, true, cursorData)

	if meta.HasNext != true {
		t.Errorf("expected HasNext true, got %v", meta.HasNext)
	}
	if meta.NextCursor == "" {
		t.Error("expected non-empty NextCursor")
	}

	decoded, err := base64.StdEncoding.DecodeString(meta.NextCursor)
	if err != nil {
		t.Fatalf("failed to decode cursor: %v", err)
	}
	var result Cursor
	if err := json.Unmarshal(decoded, &result); err != nil {
		t.Fatalf("failed to unmarshal cursor: %v", err)
	}
	if result.ID != 42 {
		t.Errorf("expected cursor ID 42, got %d", result.ID)
	}
}

func TestDecodeCursor_Valid(t *testing.T) {
	original := Cursor{ID: 99, SortValue: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)}
	b, _ := json.Marshal(original)
	encoded := base64.StdEncoding.EncodeToString(b)

	var result Cursor
	if err := DecodeCursor(encoded, &result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != 99 {
		t.Errorf("expected ID 99, got %d", result.ID)
	}
}

func TestDecodeCursor_Empty(t *testing.T) {
	var result Cursor
	if err := DecodeCursor("", &result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != 0 {
		t.Errorf("expected zero Cursor, got %+v", result)
	}
}
