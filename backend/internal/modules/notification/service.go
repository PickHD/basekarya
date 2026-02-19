package notification

import (
	"context"
	"encoding/json"
	"hris-backend/internal/infrastructure"
)

type Service interface {
	SendNotification(userID uint,
		Type string,
		Title string,
		Message string, relatedID uint) error
	GetList(ctx context.Context, userID uint) ([]NotificationListResponse, error)
	MarkAsRead(ctx context.Context, id uint) error
	DeleteReadOlderThan(days int) error
}

type service struct {
	wsHub *infrastructure.Hub
	repo  Repository
}

func NewService(wsHub *infrastructure.Hub, repo Repository) Service {
	return &service{wsHub, repo}
}

func (s *service) SendNotification(userID uint,
	notifType string,
	title string,
	message string,
	relatedID uint) error {

	notification := Notification{
		UserID:    userID,
		Type:      notifType,
		Title:     title,
		Message:   message,
		RelatedID: relatedID,
		IsRead:    false,
	}

	err := s.repo.Create(context.Background(), &notification)
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
	data, err := s.repo.FindAllByUserID(userID)
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

func (s *service) MarkAsRead(ctx context.Context, id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	err = s.repo.MarkAsRead(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) DeleteReadOlderThan(days int) error {
	return s.repo.DeleteReadOlderThan(days)
}
