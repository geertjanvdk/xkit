package xsql

import (
	"fmt"
	"os"
	"testing"
)

const defaultDSNMySQL = "root:@tcp(127.0.0.1:3306)/xkit"

var (
	testExitCode int
	testErr      error
	testDSNMySQL = defaultDSNMySQL
)

func testTearDown() {
	if testErr != nil {
		testExitCode = 1
		fmt.Println(testErr)
	}
}

func TestMain(m *testing.M) {
	defer func() { os.Exit(testExitCode) }()
	defer testTearDown()

	if v, have := os.LookupEnv("XKIT_DSN_MYSQL"); have {
		testDSNMySQL = v
	}

	testExitCode = m.Run()
}
