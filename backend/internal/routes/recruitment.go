package routes

import (
	"basekarya-backend/pkg/constants"

	"github.com/labstack/echo/v4"
)

func (r *Router) SetupRecruitmentRoutes(e *echo.Group) {
	e.POST("/requisitions", r.container.RecruitmentHandler.CreateRequisition, r.container.AuthMiddleware.GrantPermission(constants.CREATE_REQUISITION))
	e.GET("/requisitions", r.container.RecruitmentHandler.GetRequisitions, r.container.AuthMiddleware.GrantPermission(constants.VIEW_REQUISITION))
	e.GET("/requisitions/:id", r.container.RecruitmentHandler.GetRequisitionDetail, r.container.AuthMiddleware.GrantPermission(constants.VIEW_REQUISITION))
	e.PUT("/requisitions/:id/submit", r.container.RecruitmentHandler.SubmitRequisition, r.container.AuthMiddleware.GrantPermission(constants.CREATE_REQUISITION))
	e.PUT("/requisitions/:id/action", r.container.RecruitmentHandler.RequisitionAction, r.container.AuthMiddleware.GrantPermission(constants.APPROVAL_REQUISITION))
	e.PUT("/requisitions/:id/close", r.container.RecruitmentHandler.CloseRequisition, r.container.AuthMiddleware.GrantPermission(constants.APPROVAL_REQUISITION))
	e.DELETE("/requisitions/:id", r.container.RecruitmentHandler.DeleteRequisition, r.container.AuthMiddleware.GrantPermission(constants.APPROVAL_REQUISITION))

	e.POST("/requisitions/:id/applicants", r.container.RecruitmentHandler.AddApplicant, r.container.AuthMiddleware.GrantPermission(constants.CREATE_APPLICANT))
	e.GET("/requisitions/:id/applicants", r.container.RecruitmentHandler.GetApplicants, r.container.AuthMiddleware.GrantPermission(constants.VIEW_APPLICANT))

	e.GET("/applicants/:id", r.container.RecruitmentHandler.GetApplicantDetail, r.container.AuthMiddleware.GrantPermission(constants.VIEW_APPLICANT))
	e.PUT("/applicants/:id/stage", r.container.RecruitmentHandler.UpdateApplicantStage, r.container.AuthMiddleware.GrantPermission(constants.UPDATE_APPLICANT))
}
