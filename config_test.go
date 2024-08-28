package apilog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	t.Run("Applying all available options", func(t *testing.T) {
		cnf := NewConfig(
			WithNRAppName("apilog"),
			WithNRLicense("license"),
			WithFilePath("/var/log/app.log"),
			WithFileSize(100),
			WithFileAge(7),
			WithFileMaxBackup(7),
		)

		// assert all values
		assert.Equal(t, "apilog", cnf.nr.name)
		assert.Equal(t, "license", cnf.nr.license)
		assert.Equal(t, "/var/log/app.log", cnf.file.path)
		assert.Equal(t, 100, cnf.file.size)
		assert.Equal(t, 7, cnf.file.age)
		assert.Equal(t, 7, cnf.file.num)
	})
}
