package pluginhost

// NoopPlugin is a do-nothing Plugin useful for testing and as a template for
// real plugin implementations.
type NoopPlugin struct {
	name   string
	Cfg    map[string]string
	Closed bool
}

// NewNoopPlugin returns a NoopPlugin with the given name.
func NewNoopPlugin(name string) *NoopPlugin {
	return &NoopPlugin{name: name}
}

// Name implements Plugin.
func (n *NoopPlugin) Name() string { return n.name }

// Init implements Plugin; it stores cfg for later inspection.
func (n *NoopPlugin) Init(cfg map[string]string) error {
	n.Cfg = cfg
	return nil
}

// Close implements Plugin.
func (n *NoopPlugin) Close() error {
	n.Closed = true
	return nil
}
