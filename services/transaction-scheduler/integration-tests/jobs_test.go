// +build integration

package integrationtests

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

// JobsTestSuite is a test suite for Transaction API jobs controller
type JobsTestSuite struct {
	suite.Suite
	baseURL string
	env     *IntegrationEnvironment
}

func (s *JobsTestSuite) TestJobs_Validation() {
	s.T().Run("test", func(t *testing.T) {
		assert.Nil(t, nil)
	})
}
