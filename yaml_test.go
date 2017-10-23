package yaml

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

const (
	// Ugly but needed for hardcoding yaml which is sensitive to spacing
	testYaml = `
--- 
facebook: 
  host: localhost
  port: 5000
hacker-news: 
  host: localhost
  port: 3001
reddit: 
  host: localhost
  port: 3001
`
	invalidYaml = `
--- 
facebook: 
host: localhost
  port: 5000
hacker-news: 
  host: localhost
  port: 3001
reddit: 
  host: localhost
  port: 3001
`
)

type YamlTestSuite struct {
	suite.Suite
}

func (suite *YamlTestSuite) SetupSuite() {
	// Disable logging while testing
	log.SetOutput(ioutil.Discard)
}

func (suite *YamlTestSuite) removeTempFile(fname string) {
	suite.Nil(os.Remove(fname))
}

// Helper function for writing to a temporary file for testing
func (suite *YamlTestSuite) writeToTempFile(contents string) string {
	tmpfile, err := ioutil.TempFile("/tmp", "yaml")
	suite.Nil(err)

	_, err = tmpfile.Write([]byte(contents))
	suite.Nil(err)

	err = tmpfile.Close()
	suite.Nil(err)

	return tmpfile.Name()
}

func (suite *YamlTestSuite) TestGetString() {
	// Create tmp yaml object
	fname := suite.writeToTempFile(testYaml)
	config, err := New(fname)
	suite.Nil(err)
	suite.NotNil(config)
	defer suite.removeTempFile(fname) // clean up

	// Should be able to get a key that exists
	s, err := config.GetString("facebook.host")
	suite.Nil(err)
	suite.Equal(s, "localhost")

	// Should fail if key does not exist
	s, err = config.GetString("facebook.keydoesntexist")
	suite.NotNil(err)
	suite.Equal(s, "")
}

func (suite *YamlTestSuite) TestNew() {
	// Trying to read a yaml file that doesnt exist should fail
	y, err := New("/var/logs/test/this/never/will/exist/at/least/i/hope/not")
	suite.Nil(y)
	suite.NotNil(err)

	// Trying to read a yaml file with a garbage path should fail
	y, err = New("what even is this")
	suite.Nil(y)
	suite.NotNil(err)

	// Opening an existing file with invalid yaml should fail
	fname := suite.writeToTempFile(invalidYaml)
	y, err = New(fname)
	suite.NotNil(err)
	suite.Nil(y)
	suite.removeTempFile(fname)

	// Opening a valid yaml file should work
	fname = suite.writeToTempFile(testYaml)
	y, err = New(fname)
	suite.Nil(err)
	suite.NotNil(y)
	suite.removeTempFile(fname)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(YamlTestSuite))
}
