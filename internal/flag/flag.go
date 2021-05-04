// Copyright 2021 Google LLC
//
// Licensed under the Apache License, SaveGameVersion 2.0 (the "License");
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
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

var (
	flagSet = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	loadConfig = Bool("load_config", true, "enable processing of the configuration file")
)

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

// Set overrides a flag value. May be used by the menu.
func Set(name string, value interface{}) error {
	return flagSet.Set(name, fmt.Sprint(value))
}

// Get loads a flag by name.
func Get(name string) interface{} {
	f := flagSet.Lookup(name)
	if f == nil {
		log.Printf("Queried non-existing flag: %v", name)
		return ""
	}
	return f.Value.(flag.Getter).Get()
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
	flagSet.Visit(func(f *flag.Flag) {
		// Don't save debug or dump flags.
		if strings.HasPrefix(f.Name, "cheat_") {
			return
		}
		if strings.HasPrefix(f.Name, "debug_") {
			return
		}
		if strings.HasPrefix(f.Name, "dump_") {
			return
		}
		if f.Value.String() == f.DefValue {
			return
		}
		c.flags[f.Name] = f.Value.String()
	})
	return c
}

// Cheating returns if any cheats are enabled.
func Cheating() bool {
	cheating := false
	flagSet.Visit(func(f *flag.Flag) {
		if strings.HasPrefix(f.Name, "cheat_") {
			cheating = true
		}
	})
	return cheating
}

// ResetToDefaults returns all flags to their default value.
func ResetToDefaults() {
	flagSet.Visit(func(f *flag.Flag) {
		f.Value.Set(f.DefValue)
	})
}

var defaultUsage func()
var getConfig func() (*Config, error)

func applyConfig() {
	// Skip config loading if so desired.
	// This ability is why flag loading is hard;
	// we need to parse the command line to detect whether we want to load the config,
	// but then we want the command line to have precedence over the config.
	if !*loadConfig {
		log.Printf("config loading was disabled by the command line")
		return
	}
	// Remember which flags have already been set. These will NOT come from the config.
	set := map[string]struct{}{}
	flagSet.Visit(func(f *flag.Flag) {
		set[f.Name] = struct{}{}
	})
	config, err := getConfig()
	if err != nil {
		log.Printf("could not load config: %v", err)
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
			log.Printf("could not apply config value %q=%q: %v", name, value, err)
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
func Parse(getDefaults func() (*Config, error)) {
	getConfig = getDefaults
	flagSet.Usage = showUsage
	flagSet.Parse(os.Args[1:])
	applyConfig()
}

// NoConfig can be passed to Parse if the binary wants to do no config file processing.
func NoConfig() (*Config, error) {
	return nil, nil
}

// StringMap is a custom flag type to contain maps from string to string.
func StringMap(name string, value map[string]string, usage string) *map[string]string {
	m := stringMap{m: value}
	flagSet.Var(&m, name, usage)
	return &m.m
}

type stringMap struct {
	m map[string]string
}

func (m *stringMap) String() string {
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
		s += m.m[k]
	}
	return s
}

func (m *stringMap) Set(s string) error {
	m.m = map[string]string{}
	if s == "" {
		return nil
	}
	for _, word := range strings.Split(s, ",") {
		kv := strings.SplitN(word, "=", 2)
		switch len(kv) {
		case 2:
			m.m[kv[0]] = kv[1]
		default:
			return fmt.Errorf("invalid StringMap flag value, got %q, want something of the form key1=value1,key2=value2,...", s)
		}
	}
	return nil
}

func (m *stringMap) Get() interface{} {
	return m.m
}
