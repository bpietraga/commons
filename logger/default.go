package logger

import (
	"fmt"

	"strings"

	"github.com/kr/pretty"
	logsrusapi "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

var currentLogger Logger

func init() {
	currentLogger = NewLogrusLogger(logsrusapi.StandardLogger())
}

func BindFlags(flags *pflag.FlagSet) {
	flags.CountP("loglevel", "v", "Increase logging level")
	flags.Bool("json-logs", false, "Print logs in json format to stderr")
}

func ParseFlags(flags *pflag.FlagSet) {
	level, _ := flags.GetCount("loglevel")
	currentLogger.SetLogLevel(level)
}

func Warnf(format string, args ...interface{}) {
	currentLogger.Warnf(format, args...)
}

func Infof(format string, args ...interface{}) {
	currentLogger.Infof(format, args...)
}

//Secretf is like Tracef, but attempts to strip any secrets from the text
func Secretf(format string, args ...interface{}) {
	currentLogger.Tracef(stripSecrets(fmt.Sprintf(format, args...)))
}

//Prettyf is like Tracef, but pretty prints the entire struct
func Prettyf(msg string, obj interface{}) {
	pretty.Print(obj)
	// currentLogger.Tracef(msg, pretty.Sprint(obj))
}

func Errorf(format string, args ...interface{}) {
	currentLogger.Errorf(format, args...)
}

func Debugf(format string, args ...interface{}) {
	currentLogger.Debugf(format, args...)
}

func Tracef(format string, args ...interface{}) {
	currentLogger.Tracef(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	currentLogger.Fatalf(format, args...)
}

func IsTraceEnabled() bool {
	return currentLogger.IsTraceEnabled()
}

func IsDebugEnabled() bool {
	return currentLogger.IsDebugEnabled()
}

func WithValues(keysAndValues ...interface{}) Logger {
	return currentLogger.WithValues(keysAndValues)
}

func StandardLogger() Logger {
	return currentLogger
}

// stripSecrets takes a YAML or INI formatted text and removes any potentially secret data
// as denoted by keys containing "pass" or "secret" or exact matches for "key"
// the last character of the secret is kept to aid in troubleshooting
func stripSecrets(text string) string {
	out := ""
	for _, line := range strings.Split(text, "\n") {

		var k, v, sep string
		if strings.Contains(line, ":") {
			parts := strings.Split(line, ":")
			k = parts[0]
			if len(parts) > 1 {
				v = parts[1]
			}
			sep = ":"
		} else if strings.Contains(line, "=") {
			parts := strings.Split(line, "=")
			k = parts[0]
			if len(parts) > 1 {
				v = parts[1]
			}
			sep = "="
		} else {
			v = line
		}

		if strings.Contains(k, "pass") || strings.Contains(k, "secret") || strings.Contains(k, "_key") || strings.TrimSpace(k) == "key" || strings.TrimSpace(k) == "token" {
			if len(v) == 0 {
				out += k + sep + "\n"
			} else {
				out += k + sep + "****" + v[len(v)-1:] + "\n"
			}
		} else {
			out += k + sep + v + "\n"
		}
	}
	return out

}
