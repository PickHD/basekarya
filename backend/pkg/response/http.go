package response

import (
	"encoding/base64"
	"encoding/json"
	"math"
	"time"

	"github.com/labstack/echo/v4"
)

type (
	baseResponse struct {
		Message string `json:"message"`
		Data    any    `json:"data"`
		Error   any    `json:"error"`
		Meta    *Meta  `json:"meta,omitempty"`
	}

	Cursor struct {
		ID        uint      `json:"id"`
		SortValue time.Time `json:"sort_value"`
	}

	Meta struct {
		Limit int `json:"limit"`

		Page      int   `json:"page,omitempty"`
		TotalPage int64 `json:"total_page,omitempty"`
		TotalData int64 `json:"total_data,omitempty"`

		HasNext    bool   `json:"has_next,omitempty"`
		NextCursor string `json:"next_cursor,omitempty"`
	}
)

// NewResponses return dynamic JSON responses
func NewResponses[T any](ctx echo.Context, statusCode int, message string, data T, err error, meta *Meta) error {
	var errVal any

	if err != nil {
		if he, ok := err.(*echo.HTTPError); ok {
			errVal = he.Message
		} else {
			errVal = err.Error()
		}
	}

	if statusCode < 400 {
		return ctx.JSON(statusCode, &baseResponse{
			Message: message,
			Data:    data,
			Error:   nil,
			Meta:    meta,
		})

	}

	return ctx.JSON(statusCode, &baseResponse{
		Message: message,
		Data:    data,
		Error:   errVal,
		Meta:    nil,
	})
}

// NewMetaOffset return meta of offset pagination
func NewMetaOffset(page, limit int, totalData int64) *Meta {
	totalPage := int64(math.Ceil(float64(totalData) / float64(limit)))
	return &Meta{
		Page:      page,
		Limit:     limit,
		TotalPage: totalPage,
		TotalData: totalData,
	}
}

// NewMetaCursor return meta of cursor pagination
func NewMetaCursor(limit int, hasNext bool, nextCursorData interface{}) *Meta {
	encoded := ""
	if hasNext && nextCursorData != nil {
		encoded = encodeCursor(nextCursorData)
	}

	return &Meta{
		Limit:      limit,
		HasNext:    hasNext,
		NextCursor: encoded,
	}
}

func encodeCursor(data interface{}) string {
	b, _ := json.Marshal(data)
	return base64.StdEncoding.EncodeToString(b)
}

// DecodeCursor helper to decode cursor string become targeted struct
func DecodeCursor(cursor string, targetStruct interface{}) error {
	if cursor == "" {
		return nil
	}

	b, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, targetStruct)
}
