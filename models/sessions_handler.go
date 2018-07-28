package models

import (
	"errors"
	"regexp"
	"strings"
	"github.com/twinj/uuid"
	"github.com/ievgen-ma/groups-chat/app"
)

type sessionsHandler struct {
	PlatformRE string
}

func NewSessionsHandler() *sessionsHandler {
	return &sessionsHandler{
		PlatformRE: "(web|ios|android|live)+",
	}
}

func (h *sessionsHandler) Create(userID, deviceID, platform, model string, build int, name, remoteAddr string) (*Session, error) {
	s := &Session{
		ID:          uuid.NewV4().String(),
		UserID:      userID,
		Type:        typeByPlatform(platform),
		DeviceID:    deviceID,
		Platform:    platform,
		Model:       model,
		Build:       build,
		Name:        name,
		AccessToken: uuid.NewV4().String(),
		Online:      false,
		CreatedAt:   app.Timestamp(),
		UpdatedAt:   app.Timestamp(),
	}

	if s.DeviceID == "" {
		return nil, errors.New("device_id required")
	}

	if matched, _ := regexp.MatchString(h.PlatformRE, s.Platform); !matched {
		return nil, errors.New("platform invalid")
	}

	if s.Model == "" {
		return nil, errors.New("model invalid")
	}

	if s.Build == 0 {
		return nil, errors.New("build invalid")
	}

	s.IPAddress = strings.Split(remoteAddr, ":")[0]

	if err := s.Create(); err != nil {
		return nil, err
	}

	return s, nil

}

func typeByPlatform(platform string) string {
	switch platform {
	case "ios":
		return "mobile"
	case "android":
		return "mobile"
	default:
		return "web"
	}
}
