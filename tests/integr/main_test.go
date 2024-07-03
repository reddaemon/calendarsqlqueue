package integr

import (
	"fmt"
	"github.com/cucumber/godog"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	fmt.Println("Wait 5s for server availability...")
	time.Sleep(5 * time.Second)

	suite := godog.TestSuite{
		ScenarioInitializer: FeatureContext,
		Options: &godog.Options{
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
		},
	}

	if suite.Run() != 0 {
		fmt.Println("Failed")
		os.Exit(1)
	}
}
