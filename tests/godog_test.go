package tests

import (
	"fmt"
	"time"

	"github.com/DATA-DOG/godog"
)

func main() {

}

func thereAreGodogs(arg1 int) error {
	time.Sleep(2 * time.Second)
	fmt.Println("test1")
	return nil
}

func iEat(arg1 int) error {
	time.Sleep(10 * time.Second)
	fmt.Println("test2")
	return godog.ErrPending
}

func thereShouldBeRemaining(arg1 int) error {
	return godog.ErrPending
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^there are (\d+) godogs$`, thereAreGodogs)
	s.Step(`^I eat (\d+)$`, iEat)
	s.Step(`^there should be (\d+) remaining$`, thereShouldBeRemaining)
}
