package finder

import (
	"os"
	"regexp"
)

func FindEnvBetweenQuotes(str string) map[string]string {
	envMap := make(map[string]string)
	rx := regexp.MustCompile(`(?s)` + regexp.QuoteMeta("${") + `(.*?)` + regexp.QuoteMeta("}"))
	matches := rx.FindAllStringSubmatch(str, -1)
	for _, v := range matches {
		envValue := os.Getenv(v[1])
		if envValue != "" {
			envMap[v[0]] = envValue
		}
	}
	return envMap
}
