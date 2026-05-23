package utils

import (
	"testing"

	"github.com/labstack/echo/v4"
)

type validStruct struct {
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
}

type minStruct struct {
	Name string `validate:"min=5"`
}

type eqfieldStruct struct {
	Password        string `validate:"required"`
	ConfirmPassword string `validate:"eqfield=Password"`
}

func TestNewValidator(t *testing.T) {
	cv := NewValidator()
	if cv == nil {
		t.Fatal("expected non-nil validator")
	}
	if cv.Validator == nil {
		t.Fatal("expected non-nil underlying validator")
	}
}

func TestValidate_ValidStruct(t *testing.T) {
	cv := NewValidator()
	s := validStruct{Name: "John", Email: "john@example.com"}

	if err := cv.Validate(s); err != nil {
		t.Errorf("expected nil error for valid struct, got %v", err)
	}
}

func TestValidate_InvalidStruct(t *testing.T) {
	cv := NewValidator()
	s := validStruct{}

	err := cv.Validate(s)
	if err == nil {
		t.Fatal("expected error for invalid struct, got nil")
	}
}

func TestValidate_MinTag(t *testing.T) {
	cv := NewValidator()
	s := minStruct{Name: "ab"}

	err := cv.Validate(s)
	if err == nil {
		t.Fatal("expected error for min validation")
	}

	he, ok := err.(*echo.HTTPError)
	if !ok {
		t.Fatalf("expected *echo.HTTPError, got %T", err)
	}
	_ = he.Message
}

func TestValidate_EqfieldTag(t *testing.T) {
	cv := NewValidator()
	s := eqfieldStruct{Password: "abc123", ConfirmPassword: "different"}

	err := cv.Validate(s)
	if err == nil {
		t.Fatal("expected error for eqfield validation")
	}
}

func TestValidate_Eqfield_Success(t *testing.T) {
	cv := NewValidator()
	s := eqfieldStruct{Password: "abc123", ConfirmPassword: "abc123"}

	if err := cv.Validate(s); err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestValidate_MinTag_ErrorMessage(t *testing.T) {
	cv := NewValidator()
	s := minStruct{Name: "ab"}

	err := cv.Validate(s)
	if err == nil {
		t.Fatal("expected error")
	}

	he, ok := err.(*echo.HTTPError)
	if !ok {
		t.Fatalf("expected *echo.HTTPError, got %T", err)
	}

	if he.Code != 400 {
		t.Errorf("expected status 400, got %d", he.Code)
	}

	msg := he.Message.(map[string]interface{})
	errs := msg["errors"].([]ErrorResponse)
	if len(errs) == 0 {
		t.Fatal("expected at least one validation error")
	}
	if errs[0].Field != "Name" {
		t.Errorf("expected field Name, got %s", errs[0].Field)
	}
	expectedMsg := "Name must be at least 5 characters"
	if errs[0].Message != expectedMsg {
		t.Errorf("expected '%s', got '%s'", expectedMsg, errs[0].Message)
	}
}

func TestValidate_EqfieldTag_ErrorMessage(t *testing.T) {
	cv := NewValidator()
	s := eqfieldStruct{Password: "abc123", ConfirmPassword: "different"}

	err := cv.Validate(s)
	if err == nil {
		t.Fatal("expected error")
	}

	he, ok := err.(*echo.HTTPError)
	if !ok {
		t.Fatalf("expected *echo.HTTPError, got %T", err)
	}

	msg := he.Message.(map[string]interface{})
	errs := msg["errors"].([]ErrorResponse)
	if len(errs) == 0 {
		t.Fatal("expected at least one validation error")
	}
	if errs[0].Field != "ConfirmPassword" {
		t.Errorf("expected field ConfirmPassword, got %s", errs[0].Field)
	}
	expectedMsg := "ConfirmPassword must be equal to Password"
	if errs[0].Message != expectedMsg {
		t.Errorf("expected '%s', got '%s'", expectedMsg, errs[0].Message)
	}
}
