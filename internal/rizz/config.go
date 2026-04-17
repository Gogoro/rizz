package rizz

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config holds user-overridable settings loaded from ~/.config/rizz/config.toml.
type Config struct {
	Theme    string            `toml:"theme"`
	Keybinds map[string]string `toml:"keybinds"`
}

// defaultActionKeys maps each rebindable action to its built-in default key.
// Users can add alternate keys for these actions via the config file.
var defaultActionKeys = map[string]string{
	"quit":          "q",
	"help":          "?",
	"filter":        "/",
	"command":       ":",
	"commit-msgs":   "m",
	"view-toggle":   "v",
	"view-all":      "a",
	"reset":         "r",
	"next-unviewed": "U",
	"next-file":     "n",
	"prev-file":     "p",
	"enter-diff":    "enter",
	"leave-diff":    "esc",
}

// ConfigPath returns the canonical location of the config file.
func ConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "rizz", "config.toml")
}

// LoadConfig reads the config file if present. Returns a zero config when
// no file exists. Parse errors are returned to the caller.
func LoadConfig() (*Config, error) {
	cfg := &Config{Keybinds: map[string]string{}}

	path := ConfigPath()
	if path == "" {
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return cfg, err
	}
	if _, err := toml.Decode(string(data), cfg); err != nil {
		return cfg, err
	}
	if cfg.Keybinds == nil {
		cfg.Keybinds = map[string]string{}
	}
	return cfg, nil
}

// KeyRemap converts the config's action->customKey map into a runtime
// customKey->defaultKey lookup. Unknown actions are skipped silently.
func (c *Config) KeyRemap() map[string]string {
	out := map[string]string{}
	for action, custom := range c.Keybinds {
		if custom == "" {
			continue
		}
		if def, ok := defaultActionKeys[action]; ok && custom != def {
			out[custom] = def
		}
	}
	return out
}
