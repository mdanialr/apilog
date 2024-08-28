package apilog

// NewConfig return new Config after applying given options.
func NewConfig(opts ...ConfigOpt) *Config {
	var c Config
	// apply options
	for _, opt := range opts {
		opt(&c)
	}
	return &c
}

type (
	// Config required object that holds any necessary data used by each log output implementation
	Config struct {
		nr   NRConfig
		file FileConfig
	}
	// NRConfig specific config for new relic as the log output
	NRConfig struct {
		name    string
		license string
	}
	// FileConfig specific config for file as the log output
	FileConfig struct {
		path string
		size int
		age  int
		num  int
	}
)

// ConfigOpt options for Config. Not recommended to be used directly.
type ConfigOpt func(*Config)

// WithNRAppName set new relic application name.
func WithNRAppName(name string) ConfigOpt {
	return func(c *Config) {
		c.nr.name = name
	}
}

// WithNRLicense set new relic license.
func WithNRLicense(lic string) ConfigOpt {
	return func(c *Config) {
		c.nr.license = lic
	}
}

// WithFilePath set target readable directory + local file which the log data
// will be written.
func WithFilePath(p string) ConfigOpt {
	return func(c *Config) {
		c.file.path = p
	}
}

// WithFileSize set maximum file log size in megabytes before got rotated.
func WithFileSize(size int) ConfigOpt {
	return func(c *Config) {
		c.file.size = size
	}
}

// WithFileAge set maximum number of days to retain old log files based on the
// timestamp encoded in their filename
func WithFileAge(age int) ConfigOpt {
	return func(c *Config) {
		c.file.age = age
	}
}

// WithFileMaxBackup set maximum number of old log files to retain before got
// removed.
func WithFileMaxBackup(max int) ConfigOpt {
	return func(c *Config) {
		c.file.num = max
	}
}
