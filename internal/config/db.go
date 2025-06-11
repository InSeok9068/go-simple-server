package config

import "fmt"

func AppDatabaseURL(serviceName string) string {
	if IsDevEnv() {
		return fmt.Sprintf("file:./projects/%s/data/data.db", serviceName)
	} else {
		return fmt.Sprintf("file:/srv/%s/data/data.db", serviceName)
	}
}

func LogDatabaseURL() string {
	if IsDevEnv() {
		return "file:./shared/log/auxiliary.db"
	} else {
		return "file:/srv/log/auxiliary.db"
	}
}
