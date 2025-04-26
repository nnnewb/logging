package logging

import (
	"net/url"
	"path"
	"path/filepath"

	"go.uber.org/zap"
)

func constructLumberjackURI(filename string) string {
	if filename == "" {
		// normally unreachable.
		panic("file name can not be empty")
	}

	// unify path separator
	filename = filepath.ToSlash(filename)

	// trim redundant slashes.
	// e.g. //var/log/foo.log -> /var/log/foo.log
	filename = path.Clean(filename)

	// construct URI
	uri := url.URL{Scheme: "lumberjack", Opaque: filename}
	return uri.String()
}

// DevPreset returns a log configuration for the development environment.
// If the filename parameter is provided, the logs will be output to the specified file and use lumberjack for log rotation.
func DevPreset(filename string) zap.Config {
	cfg := zap.NewDevelopmentConfig()
	if filename != "" {
		cfg.OutputPaths = append(cfg.OutputPaths, constructLumberjackURI(filename))
	}
	return cfg
}

// ProdPreset returns a log configuration for the production environment.
// If the filename parameter is provided, the logs will be output to the specified file and use lumberjack for log rotation.
func ProdPreset(filename string) zap.Config {
	cfg := zap.NewProductionConfig()
	if filename != "" {
		cfg.OutputPaths = append(cfg.OutputPaths, constructLumberjackURI(filename))
	}
	return cfg
}
