package mgcplog

// LogLevel : type for log level
type LogLevel string

// LogFormat :
type LogFormat int

const (
	// Info :
	Info LogLevel = "info"
	// Warn :
	Warn LogLevel = "warn"
	// Error :
	Error LogLevel = "error"
	// Fatal :
	Fatal LogLevel = "fatal"
	// Panic :
	Panic LogLevel = "panic"
)

const (
	// JSON :
	JSON LogFormat = 0
	// Text :
	Text LogFormat = 1
)

// LogConfiguration : logger configuration
type LogConfiguration struct {
	LogFile     string
	ServiceName string
}

// LogEntity : log entry
type LogEntity struct {
	Timestamp    string   `json:"timestamp" binding:"required"`
	Level        LogLevel `json:"level" binding:"required"`
	Service      string   `json:"service" binding:"required"`
	SessionID    string   `json:"session_id" binding:"required"`
	FunctionName string   `json:"function_name" binding:"required"`
	FileName     string   `json:"file_name" binding:"required"`
	Message      string   `json:"message" binding:"required"`
}
