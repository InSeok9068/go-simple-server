package config

import "fmt"

func AppDatabaseURL(serviceName string) string {
	if IsDevEnv() {
		return fmt.Sprintf("file:./projects/%s/data/data.db?_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=foreign_keys(ON)", serviceName)
	} else {
		return fmt.Sprintf("file:/srv/%s/data/data.db?_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=foreign_keys(ON)", serviceName)
	}
}

func LogDatabaseURL() string {
	if IsDevEnv() {
		return "file:./shared/log/auxiliary.db?_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)"
	} else {
		return "file:/srv/log/auxiliary.db?_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)"
	}
}
