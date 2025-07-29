package friendly

import (
	"errors"
	"github.com/payme50rmb/jigsaw/pkg/logger"
	"github.com/spf13/viper"
)

type Logger struct {
}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Init() error {
	var cfg struct {
		Level  string `json:"level"`
		Format string `json:"format"`
		Output string `json:"output"`
		Path   string `json:"path"`
	}
	if err := viper.UnmarshalKey("logger", &cfg); err != nil {
		return err
	}
	if cfg.Output == "file" {
		if cfg.Path == "" {
			return errors.New("logger.path is required when logger.output is file")
		}
	}
	if cfg.Format == "" {
		cfg.Format = "json"
	}
	if cfg.Level == "" {
		cfg.Level = "info"
	}
	logger.SetLevel(cfg.Level)
	logger.SetFormat(cfg.Format)
	logger.SetOutput(cfg.Output)
	if cfg.Output == "file" {
		logger.SetFilePath(cfg.Path)
	}
	return nil
}
