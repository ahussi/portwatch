// Package pluginhost implements a simple plugin registry for portwatch.
//
// Plugins must implement the Plugin interface:
//
//	type MyPlugin struct{}
//	func (p *MyPlugin) Name() string                    { return "my-plugin" }
//	func (p *MyPlugin) Init(cfg map[string]string) error { return nil }
//	func (p *MyPlugin) Close() error                    { return nil }
//
// Register the plugin with a Host:
//
//	h := pluginhost.New()
//	h.Register(&MyPlugin{}, map[string]string{"key": "value"})
//
// Retrieve it later:
//
//	p, ok := h.Get("my-plugin")
//
// Shut everything down cleanly:
//
//	h.CloseAll()
package pluginhost
