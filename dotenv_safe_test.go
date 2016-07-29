package dotenv_safe

import (
	"io/ioutil"
	"os"
	"testing"
)

func checkFail(target string, t *testing.T) {
	if r := recover().(error); r != nil {
		if r.Error() != target {
			t.Error("Failed with the wrong message", r.Error())
		}
	} else {
		t.Error("Didn't fail")
	}
}

func unset(v string) (msg string) {
	return "Env variable " + v + " is not set"
}

func cleanup() {
	os.Clearenv()
	os.Remove("example.env")
	os.Remove(".env")
	os.Remove("foo.env")
	os.Remove("bar.env")
}

func TestFailureBlankExample(t *testing.T) {
	// With a blank .env and an example.env containing WUT, it should fail because WUT is not set
	defer checkFail(unset("WUT"), t)
	defer cleanup()
	ioutil.WriteFile(".env", []byte(""), 0655)
	ioutil.WriteFile("example.env", []byte("WUT"), 0655)
	Load()
}

func TestSuccessAllBlank(t *testing.T) {
	// With a blank .env and a blank example.env it should all just pass
	defer cleanup()
	ioutil.WriteFile(".env", []byte(""), 0655)
	ioutil.WriteFile("example.env", []byte(""), 0655)
	Load()
}

func TestSuccessSingleMatching(t *testing.T) {
	// With a single matched var it should pass
	defer cleanup()
	ioutil.WriteFile(".env", []byte("FOO=bar"), 0655)
	ioutil.WriteFile("example.env", []byte("FOO"), 0655)
	Load()
}

func TestFailureSingleBlankDeclaration(t *testing.T) {
	// With a single required var that is declared but unset in the .env, it should fail
	defer checkFail(unset("FOO"), t)
	defer cleanup()
	ioutil.WriteFile(".env", []byte("FOO"), 0655)
	ioutil.WriteFile("example.env", []byte("FOO"), 0655)
	Load()
}

func TestSuccessMultipleMatching(t *testing.T) {
	// With multiple matched vars it should pass
	defer cleanup()
	ioutil.WriteFile(".env", []byte("FOO=bar\nBAZ=quux"), 0655)
	ioutil.WriteFile("example.env", []byte("FOO\nBAZ"), 0655)
	Load()
}

func TestFailureMultiple(t *testing.T) {
	// With multiple failed vars it should panic on the first one
	defer checkFail(unset("FOO"), t)
	defer cleanup()
	ioutil.WriteFile(".env", []byte(""), 0655)
	ioutil.WriteFile("example.env", []byte("FOO\nBAZ"), 0655)
	Load()
}

func TestMissingExampleFile(t *testing.T) {
	// With a missing example.env file it should fail
	defer checkFail("open ./example.env: no such file or directory", t)
	defer cleanup()
	ioutil.WriteFile(".env", []byte(""), 0655)
	Load()
}

func TestMissingEnvFile(t *testing.T) {
	// With a missing .env file it should fail
	defer checkFail("open .env: no such file or directory", t)
	defer cleanup()
	ioutil.WriteFile("example.env", []byte(""), 0655)
	Load()
}

func TestCustomEnvFile(t *testing.T) {
	// With a specific location for env file it should pass
	defer cleanup()
	ioutil.WriteFile("example.env", []byte(""), 0655)
	ioutil.WriteFile("foo.env", []byte(""), 0655)
	Load("foo.env")
}

func TestCustomEnvFiles(t *testing.T) {
	// With multiple specific locations for env files it should pass
	defer cleanup()
	ioutil.WriteFile("example.env", []byte(""), 0655)
	ioutil.WriteFile("foo.env", []byte(""), 0655)
	ioutil.WriteFile("bar.env", []byte(""), 0655)
	Load("foo.env", "bar.env")
}

func TestMissingCustom(t *testing.T) {
	// With specific locations for env files where one is missing it should fail
	defer checkFail("open bar.env: no such file or directory", t)
	defer cleanup()
	ioutil.WriteFile("example.env", []byte(""), 0655)
	ioutil.WriteFile("foo.env", []byte(""), 0655)
	Load("foo.env", "bar.env")
}
