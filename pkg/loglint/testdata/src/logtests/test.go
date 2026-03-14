package logtests

import (
	"log/slog"

	"go.uber.org/zap"
)

func TestSlog() {
	password := "1234"

	// check first letter is lowercase
	slog.Info("Starting server") // want "log messages must start with lowercase letter"
	slog.Warn("Added info")      // want "log messages must start with lowercase letter"
	slog.Debug("Warning")        // want "log messages must start with lowercase letter"

	// check english language
	slog.Info("комментарий на русском") // want "log messages must contain only english"
	slog.Warn("english и русский")      // want "log messages must contain only english"
	slog.Debug("美")                     // want "log messages must contain only english"
	slog.Info("qué es esto")            // want "log messages must contain only english"

	// check no emoji
	slog.Info("server started 🚀") // want "log messages must not contain emoji"
	slog.Warn("🔥⚡️")              // want "log messages must not contain emoji"
	slog.Debug("error occured 💀") // want "log messages must not contain emoji"

	// check special symbols
	slog.Info("connection failed!!")              // want "log messages must not contain special symbols"
	slog.Warn("warning: something went wrong...") // want "log messages must not contain special symbols"
	slog.Debug("server stopped;")                 // want "log messages must not contain special symbols"

	// check sensitive data
	slog.Info("user password:" + password) // want "log messages must not contain sensitive data"
	slog.Info(password)                    // want "log messages must not contain sensitive data"
	slog.Debug("api_key=" + password)      // want "log messages must not contain sensitive data"
	slog.Info("token: " + password)        // want "log messages must not contain sensitive data"

	// check false positives
	slog.Info("user logged", slog.String("Pass", "111")) // OK
	slog.Info("pass is authentificated")                 // OK
	slog.Info("server ip is 127.0.0.1")                  // OK
	slog.Info("")                                        // OK
}

func ZapTest() {
	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()
	pass := "secret123"

	// check first letter is lowercase
	logger.Info("Starting server") // want "log messages must start with lowercase letter"
	sugar.Info("Added info")       // want "log messages must start with lowercase letter"
	sugar.Infof("Warning")         // want "log messages must start with lowercase letter"

	// check english language
	logger.Info("комментарий на русском") // want "log messages must contain only english"
	sugar.Warn("english и русский")       // want "log messages must contain only english"
	logger.Debug("美")                     // want "log messages must contain only english"
	sugar.Info("qué es esto")             // want "log messages must contain only english"

	// check no emoji
	logger.Info("server started 🚀") // want "log messages must not contain emoji"
	sugar.Warn("⚡️")                // want "log messages must not contain emoji"
	sugar.Debug("error occured 💀")  // want "log messages must not contain emoji"

	// check special symbols
	logger.Info("connection failed!!")             // want "log messages must not contain special symbols"
	sugar.Warn("warning: something went wrong...") // want "log messages must not contain special symbols"
	logger.Debug("server stopped;")                // want "log messages must not contain special symbols"

	// check sensitive data
	logger.Info("user password:" + pass) // want "log messages must not contain sensitive data"
	sugar.Info(pass)                     // want "log messages must not contain sensitive data"
	logger.Debug("api_key=" + pass)      // want "log messages must not contain sensitive data"
	sugar.Info("token: " + pass)         // want "log messages must not contain sensitive data"

	// check false positives
	sugar.Info("user logged", zap.String("Pass", "111")) // OK
	logger.Info("pass is authentificated")               // OK
	sugar.Info("server ip is 127.0.0.1")                 // OK
	logger.Info("")                                      // OK
}
