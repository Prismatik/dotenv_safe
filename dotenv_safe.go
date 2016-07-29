package dotenv_safe

import (
	"errors"
	"github.com/joho/godotenv"
	"io/ioutil"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Load(filenames ...string) {
	err := godotenv.Load(filenames...)
	check(err)

	dat, err := ioutil.ReadFile("./example.env")
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
