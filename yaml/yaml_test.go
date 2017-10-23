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
y: jack
z: 100
a: 
  host: localhost
  port: 5000
x: 
  b:
    c:
      d: hello
      e: 100
`
	boolYaml = `
--- 
x: false
y: 100
i: 
  j: true
a: 
  b:
    c:
      d: true
      e: false
`
	invalidYaml = `
--- 
a: 
host: localhost
  port: 5000
b: 
  host: localhost
  port: 3001
c: 
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

func (suite *YamlTestSuite) TestGetBool() {
	// Create tmp yaml object
	fname := suite.writeToTempFile(boolYaml)
	config, err := New(fname)
	suite.Nil(err)
	suite.NotNil(config)
	defer suite.removeTempFile(fname) // clean up

	// Should be able to get a key that exists
	b, err := config.GetBool("i.j")
	suite.Nil(err)
	suite.Equal(true, b)

	// Should be able to get a key that exists
	b, err = config.GetBool("x")
	suite.Nil(err)
	suite.Equal(false, b)

	// Should be able to get a key that exists
	b, err = config.GetBool("a.b.c.d")
	suite.Nil(err)
	suite.Equal(true, b)

	// Should fail if key does not exist
	_, err = config.GetBool("a.keydoesntexist")
	suite.NotNil(err)

	// Should fail if key is empty
	_, err = config.GetBool("")
	suite.NotNil(err)

	// Should fail if a middle part of path doesnt exist
	_, err = config.GetBool("x.b.fake.d")
	suite.NotNil(err)

	// Should fail if we try to get a key that is not a bool
	_, err = config.GetBool("y")
	suite.NotNil(err)
}

func (suite *YamlTestSuite) TestGetInt() {
	// Create tmp yaml object
	fname := suite.writeToTempFile(testYaml)
	config, err := New(fname)
	suite.Nil(err)
	suite.NotNil(config)
	defer suite.removeTempFile(fname) // clean up

	// Should be able to get a key that exists
	i, err := config.GetInt("a.port")
	suite.Nil(err)
	suite.Equal(5000, i)

	// Should be able to get a key that exists
	i, err = config.GetInt("z")
	suite.Nil(err)
	suite.Equal(100, i)

	// Should be able to get a key that exists
	i, err = config.GetInt("x.b.c.e")
	suite.Nil(err)
	suite.Equal(100, i)

	// Should fail if key does not exist
	i, err = config.GetInt("a.keydoesntexist")
	suite.NotNil(err)
	suite.Equal(0, i)

	// Should fail if key is empty
	i, err = config.GetInt("")
	suite.NotNil(err)
	suite.Equal(0, i)

	// Should fail if a middle part of path doesnt exist
	i, err = config.GetInt("x.b.fake.d")
	suite.NotNil(err)
	suite.Equal(0, i)

	// Should fail if we try to get a key that is not an int
	i, err = config.GetInt("a.host")
	suite.NotNil(err)
	suite.Equal(0, i)
}

func (suite *YamlTestSuite) TestGetString() {
	// Create tmp yaml object
	fname := suite.writeToTempFile(testYaml)
	config, err := New(fname)
	suite.Nil(err)
	suite.NotNil(config)
	defer suite.removeTempFile(fname) // clean up

	// Should be able to get a key that exists
	s, err := config.GetString("a.host")
	suite.Nil(err)
	suite.Equal("localhost", s)

	// Should be able to get a key that exists
	s, err = config.GetString("y")
	suite.Nil(err)
	suite.Equal("jack", s)

	// Should be able to get a key that exists
	s, err = config.GetString("x.b.c.d")
	suite.Nil(err)
	suite.Equal("hello", s)

	// Should fail if key does not exist
	s, err = config.GetString("a.keydoesntexist")
	suite.NotNil(err)
	suite.Equal("", s)

	// Should fail if key does not exist
	s, err = config.GetString("")
	suite.NotNil(err)
	suite.Equal("", s)

	// Should fail if a middle part of path doesnt exist
	s, err = config.GetString("x.b.fake.d")
	suite.NotNil(err)
	suite.Equal("", s)

	// Should fail if we try to access key that is not string
	s, err = config.GetString("a.port")
	suite.NotNil(err)
	suite.Equal("", s)
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
