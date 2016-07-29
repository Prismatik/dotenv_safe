package dotenv_safe

import (
	"errors"
	"github.com/joho/godotenv"
	"io/ioutil"
	"os"
	"strings"
)

type Config struct {
	Envs     []string
	Examples []string
}

func Load(filenames ...string) {
	envs := envsOrDefault(filenames...)
	LoadMany(Config{
		Envs:     envs,
		Examples: defaultExample,
	})
}

func LoadMany(config Config) {
	err := godotenv.Load(config.Envs...)
	check(err)

	for _, filename := range config.Examples {
		dat, err := ioutil.ReadFile(filename)
		datstr := string(dat)
		check(err)

		lines := strings.Split(datstr, "\n")
		for _, line := range lines {
			if line == "" {
				continue
			}
			frags := strings.Split(line, "=")
			key := frags[0]

			_, isSet := os.LookupEnv(key)

			if !isSet {
				panic(errors.New("Env variable " + key + " is not set"))
			}
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var defaultExample = []string{"example.env"}

func envsOrDefault(s ...string) []string {
	if len(s) == 0 {
		return []string{".env"}
	}
	return s
}

func examplesOrDefault(s ...string) []string {
	if len(s) == 0 {
		return defaultExample
	}
	return s
}
