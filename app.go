package env

import (
	"os"
	"path/filepath"
	"time"
)

// ================================================================
//
// ================================================================
type App struct {
	AppTitle   string
	AppHost    string
	AppPath    string
	AppPort    string
	GinMode    string
	Location   *time.Location
	TrustProxy string
}

// ================================================================
//
// ================================================================
func NewApp() (*App, error) {
	loc, err := time.LoadLocation(os.Getenv("TIMEZONE"))
	if err != nil {
		return nil, err
	}

	return &App{
		AppTitle:   os.Getenv("APP_TITLE"),
		AppHost:    os.Getenv("APP_HOST"),
		AppPath:    filepath.Join("/", os.Getenv("APP_PATH")),
		AppPort:    os.Getenv("APP_PORT"),
		GinMode:    os.Getenv("GIN_MODE"),
		Location:   loc,
		TrustProxy: os.Getenv("TRUST_PROXY"),
	}, nil
}

// ================================================================
//
// ================================================================
func (e App) GetAppRootURL() string {
	return "https://" + filepath.Join(e.AppHost, e.AppPath)
}
