package integr

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/godog"
)

func TestMain(m *testing.M) {
	fmt.Println("Wait 5s for server availability...")
	time.Sleep(5 * time.Second)

	status := godog.RunWithOptions("integr", func(s *godog.Suite) {
		FeatureContext(s)
	}, godog.Options{
		ShowStepDefinitions: false,
		Randomize:           0,
		StopOnFailure:       false,
		Strict:              false,
		NoColors:            false,
		Tags:                "",
		Format:              "progress",
		Concurrency:         0,
		Paths:               []string{"features"},
		Output:              nil,
	})

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}
