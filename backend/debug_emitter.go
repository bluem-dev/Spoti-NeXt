package backend

import "sync"

// DebugEmitter allows the backend to send log lines to the frontend debug panel.
// app.go sets the emit function once at startup; backend packages call DebugLog.
// If no emitter is set, messages are silently dropped (safe for tests / CLI use).

var (
	debugEmitMu sync.RWMutex
	debugEmitFn func(string)
)

// SetDebugEmitter registers the function that forwards messages to the frontend.
// Called once from app.go after the Wails context is available.
func SetDebugEmitter(fn func(string)) {
	debugEmitMu.Lock()
	defer debugEmitMu.Unlock()
	debugEmitFn = fn
}

// DebugLog sends a message to the frontend debug panel (if an emitter is set).
func DebugLog(msg string) {
	debugEmitMu.RLock()
	fn := debugEmitFn
	debugEmitMu.RUnlock()
	if fn != nil {
		fn(msg)
	}
}
