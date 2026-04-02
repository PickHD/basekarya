package auth

import (
	"basekarya-backend/pkg/logger"
	"basekarya-backend/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{s}
}

func (h *Handler) Login(ctx echo.Context) error {
	var req LoginRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	if err := ctx.Validate(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	resp, err := h.service.Login(ctx.Request().Context(), req.Username, req.Password)
	if err != nil {
		logger.Errorw("Login failed : ", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Login Success", resp, nil, nil)
}

func (h *Handler) ForgotPassword(ctx echo.Context) error {
	var req SendOrResendOTPRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	if err := ctx.Validate(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	err := h.service.SendOrResendOTP(ctx.Request().Context(), &req)
	if err != nil {
		logger.Errorw("Forgot Password failed : ", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Forgot Password Success", nil, nil, nil)
}

func (h *Handler) ResendOTP(ctx echo.Context) error {
	var req SendOrResendOTPRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	if err := ctx.Validate(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	err := h.service.SendOrResendOTP(ctx.Request().Context(), &req)
	if err != nil {
		logger.Errorw("Resend OTP failed : ", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Resend OTP Success", nil, nil, nil)
}

func (h *Handler) VerifyOTP(ctx echo.Context) error {
	var req VerifyOTPRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	if err := ctx.Validate(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	resp, err := h.service.VerifyOTP(ctx.Request().Context(), &req)
	if err != nil {
		logger.Errorw("Verify OTP failed : ", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Verify OTP Success", resp, nil, nil)
}

func (h *Handler) ResetPassword(ctx echo.Context) error {
	var req ResetPasswordRequest
	if err := ctx.Bind(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	if err := ctx.Validate(&req); err != nil {
		return response.NewResponses[any](ctx, http.StatusBadRequest, "Invalid Request", nil, err, nil)
	}

	err := h.service.ResetPassword(ctx.Request().Context(), &req)
	if err != nil {
		logger.Errorw("Reset Password failed : ", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Reset Password Success", nil, nil, nil)
}
