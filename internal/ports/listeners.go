package ports

// SystemListenerLookup resolves listeners through the current platform's
// native process inspection mechanism.
type SystemListenerLookup struct{}

// Listener reports whether port is occupied and every owner PID visible to the
// current process. A port without a listener is not an error.
func (SystemListenerLookup) Listener(port int) (Listener, error) {
	return systemListener(port)
}
