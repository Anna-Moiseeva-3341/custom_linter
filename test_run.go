package mylinterproject

import "log/slog"

func TestRun() {
	password := 123

	// running few basic checks
	slog.Info("Bad start")
	slog.Warn("русский язык")
	slog.Debug("!!!error")
	slog.Info("🔥")
	slog.Error("my pass is: ", "password_value", password)
}
