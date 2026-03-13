package a

import (
	"log/slog"

	"go.uber.org/zap"
)

func check() {
	var logger *zap.Logger
	password := "12345"
	token := "abc"

	// rule 1: starts with small
	slog.Info("Starting server")   // want "log message should start with a lowercase letter"
	logger.Error("Failed connect") // want "log message should start with a lowercase letter"

	// rule 2: english only
	slog.Info("привет мир")   // want "log message contains non-english characters"
	logger.Info("привет мир") // want "log message contains non-english characters"

	// rule 3: no emojis or special symbols
	slog.Warn("done!!!")     // want "log message contains invalid characters: emoji or special symbols"
	slog.Debug("rocket 🚀")   // want "log message contains invalid characters: emoji or special symbols"
	logger.Warn("done!!!")   // want "log message contains invalid characters: emoji or special symbols"
	logger.Debug("rocket 🚀") // want "log message contains invalid characters: emoji or special symbols"

	// rule 4: sensitive data
	slog.Info("user auth", "pass", password)              // want "potential sensitive data leak in log: variable 'password'"
	slog.Debug("token is here", "t", token)               // want "potential sensitive data leak in log: variable 'token'"
	logger.Info("user auth", "pass", password)            // want "potential sensitive data leak in log: variable 'password'"
	logger.Debug("token is here", zap.String("t", token)) // want "potential sensitive data leak in log: variable 'token'"

	// right ones
	slog.Info("server started")
	slog.Debug("request processed", "id", 123)
	logger.Info("server started")
	logger.Debug("request processed", "id", 123)
}
