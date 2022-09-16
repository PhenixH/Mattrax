package settings

import (
	"sync"

	"github.com/mattrax/Mattrax/internal/db"
)

type Settings struct {
}

// Service allow safely retrieving and setting of server settings
type Service struct {
	settings     Settings
	settingsLock sync.RWMutex
}

// Get safely returns the servers settings
func (s *Service) Get() Settings {
	s.settingsLock.RLock()
	var settings = s.settings
	s.settingsLock.RUnlock()
	return settings
}

// New initialises a new settings service
func New(q *db.Queries) (*Service, error) {
	// settings, err := q.Settings(context.Background())
	// if err != nil {
	// 	return nil, err
	// }

	var settings = Settings{
		// TenantName: "Mattrax"
	}

	return &Service{
		settings: settings,
	}, nil
}
