package logging

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	MaxSizeKey     = "MAX_SIZE"
	MaxAgeKey      = "MAX_AGE"
	MaxBackupsKey  = "MAX_BACKUPS"
	LocalTimeKey   = "LOCAL_TIME"
	CompressKey    = "COMPRESS"
)

type lumberjackSink struct {
	*lumberjack.Logger
}

// Sync implements zap.Sink. The remaining methods are implemented
// by the embedded *lumberjack.Logger.
func (lumberjackSink) Sync() error { return nil }

func init() {
	err := zap.RegisterSink("lumberjack", func(u *url.URL) (zap.Sink, error) {
		// There are three cases:
		//
		// opaque: <schema>:<non-slash opaque> e.g. lumberjack:foo.log
		// no authority: <schema>:/<path> e.g. lumberjack:/tmp/foo.log
		// with authority: <schema>://<authority>/<path> e.g. lumberjack://localhost/tmp/foo.log
		//
		// So theoretically, considering only opaque and path is enough. If the URI is opaque, u.Path will not appear.
		//
		// reference:
		// - https://stackoverflow.com/questions/31635991/difference-between-opaque-and-hierarchical-uri
		// - https://www.rfc-editor.org/rfc/rfc3986#page-7
		path := u.Opaque
		if path == "" {
			path = u.Path
			if path == "" {
				panic(fmt.Sprintf("no output path specified: %#+v", u))
			}
		}

		logger := lumberjack.Logger{Filename: path}
		for k, v := range u.Query() {
			var err error
			switch strings.ToUpper(k) {
			case MaxSizeKey:
				logger.MaxSize, err = strconv.Atoi(v[0])
				if err != nil {
					panic(fmt.Sprintf("invalid maxsize: %s", u))
				}
			case MaxAgeKey:
				logger.MaxAge, err = strconv.Atoi(v[0])
				if err != nil {
					panic(fmt.Sprintf("invalid maxage: %s", u))
				}
			case MaxBackupsKey:
				logger.MaxBackups, err = strconv.Atoi(v[0])
				if err != nil {
					panic(fmt.Sprintf("invalid maxbackups: %s", u))
				}
			case LocalTimeKey:
				logger.LocalTime = true
			case CompressKey:
				logger.Compress = true
			}
		}

		return lumberjackSink{
			Logger: &logger,
		}, nil
	})
	if err != nil {
		panic(err)
	}
}