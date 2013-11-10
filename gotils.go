package gotils

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	NonWord = regexp.MustCompile("[^0-9A-Za-z_]")
)

type Slog map[string]interface{}

func (s Slog) String() string {
	var sv string
	parts := make([]string, 0, len(s))

	for k, v := range s {
		switch v.(type) {
		case time.Time: // Format times the way we want them
			sv = v.(time.Time).Format(time.RFC3339Nano)
		default: // Let Go figure out the representation
			sv = fmt.Sprintf("%v", v)
		}
		// If there is a NonWord character then need to quote the value
		if NonWord.MatchString(sv) {
			sv = fmt.Sprintf("%q", sv)
		}
		// Assemble the final part and append it to the array
		parts = append(parts, fmt.Sprintf("%s=%s", k, sv))
	}
	sort.Strings(parts)
	return strings.Join(parts, " ")
}

func (s Slog) FatalError(err error, msg interface{}) {
	s.Error(err, msg)
	os.Exit(1)
}

func (s Slog) Error(err error, msg interface{}) {
	s["at"] = time.Now()
	s["error"] = err
	s["message"] = msg
	fmt.Println(s)
	delete(s, "error")
	delete(s, "message")
}

func LinearSliceContainsString(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

func SliceContainsString(ss []string, s string) bool {
	if sort.StringsAreSorted(ss) {
		idx := sort.SearchStrings(ss, s)
		if idx < len(ss) && ss[idx] == s {
			return true
		}
	} else {
		return LinearSliceContainsString(ss, s)
	}
	return false
}

func Ui64toa(val uint64) string {
	return strconv.FormatUint(val, 10)
}

func Atofloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func PercentFormat(val float64) string {
	return strconv.FormatFloat(val, 'f', 2, 64)
}

func Atouint64(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}

// Checks to see if a path exists or not
func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Returns the value of $env from the OS and if it's empty, returns def
func GetEnvWithDefault(env string, def string) (value string) {
	value = os.Getenv(env)

	if value == "" {
		return def
	}

	return value
}

// Returns the value of $env from the OS and if it's empty, returns def
func GetEnvWithDefaultInt(env string, def int) (int, error) {
	tmp := os.Getenv(env)

	if tmp == "" {
		return def, nil
	}

	return strconv.Atoi(tmp)
}

// Returns the value of $env from the OS and if it's empty, returns def
func GetEnvWithDefaultBool(env string, def bool) (bool, error) {
	tmp := os.Getenv(env)

	if tmp == "" {
		return def, nil
	}

	return strconv.ParseBool(tmp)
}

func GetEnvWithDefaultDuration(env string, def string) (time.Duration, error) {
	tmp := os.Getenv(env)

	if tmp == "" {
		tmp = def
	}

	return time.ParseDuration(tmp)
}

// Returns a slice of sorted strings from the environment or default split on ,
// So "foo,bar" returns ["bar","foo"]
func GetEnvWithDefaultStrings(env string, def string) (value []string) {
	env = GetEnvWithDefault(env, def)
	if env == "" {
		return make([]string, 0)
	}
	value = strings.Split(env, ",")
	if !sort.StringsAreSorted(value) {
		sort.Strings(value)
	}
	return value
}
