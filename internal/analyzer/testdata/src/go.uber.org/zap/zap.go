package zap

type Logger struct{}

func (l *Logger) Info(msg string, fields ...any)  {}
func (l *Logger) Debug(msg string, fields ...any) {}
func (l *Logger) Warn(msg string, fields ...any)  {}
func (l *Logger) Error(msg string, fields ...any) {}

type SugaredLogger struct{}

func (s *SugaredLogger) Info(args ...any)                   {}
func (s *SugaredLogger) Infof(template string, args ...any) {}

func String(key, val string) any    { return nil }
func Int(key string, val int) any   { return nil }
func Bool(key string, val bool) any { return nil }
