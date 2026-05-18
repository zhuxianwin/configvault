package audit

// NoopLogger is an audit Logger that discards all entries.
// It is useful in tests or when auditing is disabled.
type NoopLogger struct{}

// Log discards the entry and returns nil.
func (n *NoopLogger) Log(_ Entry) error { return nil }

// Close is a no-op.
func (n *NoopLogger) Close() error { return nil }

// Auditor is the interface satisfied by both Logger and NoopLogger.
type Auditor interface {
	Log(Entry) error
	Close() error
}
