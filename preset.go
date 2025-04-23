package logging

import (
	"fmt"
	"net/url"

	"go.uber.org/zap"
)

// DevPreset returns a log configuration for the development environment.
// If the filename parameter is provided, the logs will be output to the specified file and use lumberjack for log rotation.
func DevPreset(filename string) zap.Config {
	cfg := zap.NewDevelopmentConfig()
	if filename != "" {
		cfg.OutputPaths = append(cfg.OutputPaths, fmt.Sprintf("lumberjack:%s", url.PathEscape(filename)))
	}
	return cfg
}

// ProdPreset returns a log configuration for the production environment.
// If the filename parameter is provided, the logs will be output to the specified file and use lumberjack for log rotation.
func ProdPreset(filename string) zap.Config {
	cfg := zap.NewProductionConfig()
	if filename != "" {
		cfg.OutputPaths = append(cfg.OutputPaths, fmt.Sprintf("lumberjack:%s", url.PathEscape(filename)))
	}
	return cfg
}