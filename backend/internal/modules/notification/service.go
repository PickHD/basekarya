package notification

import (
	"basekarya-backend/internal/infrastructure"
	"basekarya-backend/pkg/utils"
	"context"
	"encoding/json"
)

type Service interface {
	SendNotification(ctx context.Context, userID uint,
		Type string,
		Title string,
		Message string, relatedID uint) error
	GetList(ctx context.Context, userID uint) ([]NotificationListResponse, error)
	MarkAsRead(ctx context.Context, id uint, userID uint) error
	DeleteReadOlderThan(days int) error
	BlastNotification(ctx context.Context, userIDs []uint,
		Type string,
		Title string,
		Message string, relatedID uint) error
}

type service struct {
	wsHub *infrastructure.Hub
	repo  Repository
}

func NewService(wsHub *infrastructure.Hub, repo Repository) Service {
	return &service{wsHub, repo}
}

func (s *service) SendNotification(ctx context.Context, userID uint,
	notifType string,
	title string,
	message string,
	relatedID uint) error {

	notification := Notification{
		CompanyID: utils.GetCompanyIDFromCtx(ctx),
		UserID:    userID,
		Type:      notifType,
		Title:     title,
		Message:   message,
		RelatedID: relatedID,
		IsRead:    false,
	}

	err := s.repo.Create(ctx, &notification)
	if err != nil {
		return err
	}

	payload := NotificationRequest{
		ID:        notification.ID,
		UserID:    userID,
		Type:      notifType,
		Title:     title,
		Message:   message,
		RelatedID: relatedID,
	}

	data, err := json.Marshal(&payload)
	if err != nil {
		return err
	}

	s.wsHub.SendToUser(payload.UserID, data)

	return nil
}

func (s *service) GetList(ctx context.Context, userID uint) ([]NotificationListResponse, error) {
	data, err := s.repo.FindAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return []NotificationListResponse{}, nil
	}

	var responses []NotificationListResponse
	for _, n := range data {
		responses = append(responses, NotificationListResponse{
			ID:        n.ID,
			UserID:    n.UserID,
			Type:      n.Type,
			Title:     n.Title,
			Message:   n.Message,
			RelatedID: n.RelatedID,
			IsRead:    n.IsRead,
			CreatedAt: n.CreatedAt,
		})
	}

	return responses, nil
}

func (s *service) MarkAsRead(ctx context.Context, id uint, userID uint) error {
	_, err := s.repo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return err
	}

	err = s.repo.MarkAsRead(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) DeleteReadOlderThan(days int) error {
	return s.repo.DeleteReadOlderThan(context.Background(), days)
}

func (s *service) BlastNotification(ctx context.Context, userIDs []uint,
	notifType string,
	title string,
	message string,
	relatedID uint) error {

	if len(userIDs) == 0 {
		return nil
	}

	var notifications []*Notification
	for _, userID := range userIDs {
		notifications = append(notifications, &Notification{
			CompanyID:  utils.GetCompanyIDFromCtx(ctx),
			UserID:    userID,
			Type:      notifType,
			Title:     title,
			Message:   message,
			RelatedID: relatedID,
			IsRead:    false,
		})
	}

	err := s.repo.CreateBatch(ctx, notifications)
	if err != nil {
		return err
	}

	// Construct messages for pipelined WebSocket broadcast
	var wsMessages []infrastructure.Message
	for _, notification := range notifications {
		payload := NotificationRequest{
			ID:        notification.ID,
			UserID:    notification.UserID,
			Type:      notifType,
			Title:     title,
			Message:   message,
			RelatedID: relatedID,
		}

		data, err := json.Marshal(&payload)
		if err == nil {
			wsMessages = append(wsMessages, infrastructure.Message{
				TargetUserID: payload.UserID,
				Data:         data,
			})
		}
	}

	s.wsHub.BroadcastPipelined(wsMessages)

	return nil
}
