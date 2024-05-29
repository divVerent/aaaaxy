// Copyright 2021 Google LLC
//
// Licensed under the Apache Livense, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package flag

import (
	"encoding"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/divVerent/aaaaxy/internal/log"
)

var (
	flagSet           = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	v                 = Int("v", 0, "verbose logging level")                // Must be declared here to prevent cycle.
	batch             = Bool("batch", false, "if set, show no alert boxes") // Must be declared here to prevent cycle.
	loadConfig        = Bool("load_config", true, "enable processing of the configuration file")
	debugPersistFlags = Bool("debug_persist_flags", false, "persist debug_* flags to config (including this one); BEWARE: this can degrade game performance")
)

// SystemDefault performs a GOOS/GOARCH dependent value lookup to be used in flag defaults.
// Map keys shall be */*, GOOS/*, */GOARCH or GOOS/GOARCH.
func SystemDefault[T any](m map[string]T) T {
	k := fmt.Sprintf("%v/%v", runtime.GOOS, runtime.GOARCH)
	if val, found := m[k]; found {
		return val
	}
	k = fmt.Sprintf("%v/*", runtime.GOOS)
	if val, found := m[k]; found {
		return val
	}
	k = fmt.Sprintf("*/%v", runtime.GOARCH)
	if val, found := m[k]; found {
		return val
	}
	return m["*/*"]
}

// Bool creates a bool in our FlagSet.
func Bool(name string, value bool, usage string) *bool {
	return flagSet.Bool(name, value, usage)
}

// Float64 creates a float64 in our FlagSet.
func Float64(name string, value float64, usage string) *float64 {
	return flagSet.Float64(name, value, usage)
}

// Int creates an int in our FlagSet.
func Int(name string, value int, usage string) *int {
	return flagSet.Int(name, value, usage)
}

// String creates a string in our FlagSet.
func String(name string, value string, usage string) *string {
	return flagSet.String(name, value, usage)
}

// Duration creates a Duration in our FlagSet.
func Duration(name string, value time.Duration, usage string) *time.Duration {
	return flagSet.Duration(name, value, usage)
}

// Text creates a flag based on a variable that fulfills TextMarshaler and TextUnmarshaler.
func Text[T any, PT interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
	*T
}](name string, value T, usage string) PT {
	var actual T
	actualPT := PT(&actual)
	valuePT := PT(&value)
	flagSet.TextVar(actualPT, name, valuePT, usage)
	return &actual
}

// Set overrides a flag value. May be used by the menu.
func Set(name string, value interface{}) error {
	switch vT := value.(type) {
	case encoding.TextMarshaler:
		buf, err := vT.MarshalText()
		if err != nil {
			return err
		}
		return flagSet.Set(name, string(buf))
	default:
		return flagSet.Set(name, fmt.Sprint(vT))
	}
}

// Get loads a flag by name.
func Get[T any](name string) T {
	f := flagSet.Lookup(name)
	if f == nil {
		log.Errorf("queried non-existing flag: %v", name)
		var zeroValue T
		return zeroValue
	}
	return f.Value.(flag.Getter).Get().(T)
}

// Config is a JSON serializable type containing the flags.
type Config struct {
	flags map[string]string
}

// MarshalJSON returns the JSON representation of the config.
func (c *Config) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.flags)
}

// UnmarshalJSON loads the config from a JSON object string.
func (c *Config) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &c.flags)
}

// Marshal returns a config object for the currently set flags (both those from the config and command line).
// We only write non-default flag values.
func Marshal() *Config {
	c := &Config{flags: map[string]string{}}
	// Note: VisitAll also sees flags that have been modified by code writing to *flag pointers.
	// Visit would only see flags changed using flag.Set or the command line.
	flagSet.VisitAll(func(f *flag.Flag) {
		// Don't save debug or dump flags.
		if strings.HasPrefix(f.Name, "cheat_") {
			return
		}
		if !*debugPersistFlags && strings.HasPrefix(f.Name, "debug_") {
			return
		}
		if strings.HasPrefix(f.Name, "dump_") {
			return
		}
		if strings.HasPrefix(f.Name, "demo_") {
			return
		}
		if f.Name == "batch" {
			return
		}
		if f.Name == "config_path" {
			return
		}
		if f.Name == "portable" {
			return
		}
		if f.Name == "save_path" {
			return
		}
		if f.Value.String() == f.DefValue {
			return
		}
		c.flags[f.Name] = f.Value.String()
	})
	return c
}

// Cheating returns if any cheats are enabled, and what they are.
func Cheating() (bool, string) {
	cheating := false
	cheats := []string{}
	flagSet.Visit(func(f *flag.Flag) {
		if strings.HasPrefix(f.Name, "cheat_") {
			cheating = true
			cheats = append(cheats, fmt.Sprintf("--%s=%s", f.Name, f.Value.String()))
		}
	})
	return cheating, strings.Join(cheats, " ")
}

// ResetToDefaults returns all flags to their default value.
func ResetToDefaults() {
	flagSet.Visit(func(f *flag.Flag) {
		f.Value.Set(f.DefValue)
	})
}

// ResetFlagToDefault returns a given flag to its default value.
func ResetFlagToDefault(name string) error {
	f := flagSet.Lookup(name)
	if f == nil {
		return fmt.Errorf("resetting non-existing flag: %v", name)
	}
	f.Value.Set(f.DefValue)
	return nil
}

var getConfig func() (*Config, error)

func applyConfig() {
	// Provide verbose level ASAP.
	log.V = v
	log.Batch = batch

	// Skip config loading if so desired.
	// This ability is why flag loading is hard;
	// we need to parse the command line to detect whether we want to load the config,
	// but then we want the command line to have precedence over the config.
	if !*loadConfig {
		log.Infof("config loading was disabled by the command line")
		return
	}
	// Remember which flags have already been set. These will NOT come from the config.
	set := map[string]struct{}{}
	flagSet.Visit(func(f *flag.Flag) {
		set[f.Name] = struct{}{}
	})
	config, err := getConfig()
	if err != nil {
		log.Errorf("could not load config: %v", err)
		return
	}
	if config == nil {
		// Nothing to do.
		return
	}
	for name, value := range config.flags {
		// Don't take from config what's already been overridden.
		if _, found := set[name]; found {
			continue
		}
		err = flagSet.Set(name, value)
		if err != nil {
			log.Errorf("could not apply config value %q=%q: %v", name, value, err)
			continue
		}
	}
}

func showUsage() {
	applyConfig()
	flagSet.PrintDefaults()
}

// Parse parses the command-line flags, then loads the config object using the provided function.
// Should be called initially, before loading config.
func Parse(getSystemDefaults func() (*Config, error)) {
	getConfig = getSystemDefaults
	flagSet.Usage = showUsage
	flagSet.Parse(os.Args[1:])
	applyConfig()
}

// NoConfig can be passed to Parse if the binary wants to do no config file processing.
func NoConfig() (*Config, error) {
	return nil, nil
}

// StringMap is a custom flag type to contain maps from string to T.
func StringMap[T any](name string, value map[string]T, usage string) *map[string]T {
	m := stringMap[T]{m: value}
	flagSet.Var(&m, name, usage)
	return &m.m
}

type stringMap[T any] struct {
	m map[string]T
}

func (m *stringMap[T]) String() string {
	a := make([]string, 0, len(m.m))
	for k := range m.m {
		a = append(a, k)
	}
	sort.Strings(a)
	s := ""
	for _, k := range a {
		if s != "" {
			s += ","
		}
		s += k
		s += "="
		s += fmt.Sprint(m.m[k])
	}
	return s
}

func (m *stringMap[T]) Set(s string) error {
	m.m = map[string]T{}
	if s == "" {
		return nil
	}
	for _, word := range strings.Split(s, ",") {
		kv := strings.SplitN(word, "=", 2)
		switch len(kv) {
		case 2:
			if kv[1] == "" {
				delete(m.m, kv[0])
			} else {
				var v T
				_, err := fmt.Sscanf(kv[1], "%v", &v)
				if err != nil {
					return fmt.Errorf("invalid StringMap flag value, got %q, could no parse contained value %q", s, kv[1])
				}
				m.m[kv[0]] = v
			}
		case 1:
			switch m_m := any(m.m).(type) {
			case map[string]bool:
				if strings.HasPrefix(kv[0], "no") {
					m_m[kv[0][2:]] = false
				} else {
					m_m[kv[0]] = true
				}
			default:
				return fmt.Errorf("missing StringMap flag value, got %q, want items of the form key=value, not %q", s, word)
			}
		default:
			return fmt.Errorf("invalid StringMap flag value, got %q, want items of the form key=value, not %q", s, word)
		}
	}
	return nil
}

func (m *stringMap[T]) Get() interface{} {
	return m.m
}
