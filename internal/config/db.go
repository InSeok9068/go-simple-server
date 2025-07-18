package config

import "fmt"

func AppDatabaseURL(serviceName string) string {
	var url string
	pragmas := `
		&_pragma=journal_mode(WAL)
		&_pragma=synchronous(NORMAL)
		&_pragma=busy_timeout(5000)
		&_pragma=foreign_keys(ON)`
	if IsDevEnv() {
		url = fmt.Sprintf(`file:./projects/%s/data/data.db?%s`, serviceName, pragmas)
	} else {
		url = fmt.Sprintf(`file:/srv/%s/data/data.db?%s`, serviceName, pragmas)
	}
	return url
}

func LogDatabaseURL() string {
	var url string
	pragmas := `
		&_pragma=journal_mode(WAL)
		&_pragma=synchronous(NORMAL)
		&_pragma=busy_timeout(5000)`
	if IsDevEnv() {
		url = fmt.Sprintf(`file:./shared/log/auxiliary.db?%s`, pragmas)
	} else {
		url = fmt.Sprintf(`file:/srv/log/auxiliary.db?%s`, pragmas)
	}
	return url
}
