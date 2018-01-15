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

func (suite *YamlTestSuite) TestGetStrings() {
	// Create tmp yaml object
	fname := suite.writeToTempFile(testYaml)
	config, err := New(fname)
	suite.Nil(err)
	suite.NotNil(config)
	defer suite.removeTempFile(fname) // clean up

	// Should be able to get keys that exists
	s := config.GetStrings([]string{"y", "a.host", "x.b.c.d"})
	suite.Equal([]string{"jack", "localhost", "hello"}, s)

	// Should fail if one of the keys doesnt exist
	s = config.GetStrings([]string{"y", "a.fake.not.real.host", "x.b.c.d"})
	suite.Equal([]string{"jack", "", "hello"}, s)

	// Should fail if one of the keys is wrong type
	s = config.GetStrings([]string{"z", "a.host", "x.b.c.d"})
	suite.Equal([]string{"", "localhost", "hello"}, s)
}

func (suite *YamlTestSuite) TestGetInts() {
	// Create tmp yaml object
	fname := suite.writeToTempFile(testYaml)
	config, err := New(fname)
	suite.Nil(err)
	suite.NotNil(config)
	defer suite.removeTempFile(fname) // clean up

	// Should be able to get keys that exists
	s := config.GetInts([]string{"z", "a.port", "x.b.c.e"})
	suite.Equal([]int{100, 5000, 100}, s)

	// Should fail if one of the keys doesnt exist
	s = config.GetInts([]string{"z", "a.fake.not.real.host", "x.b.c.e"})
	suite.Equal([]int{100, 0, 100}, s)

	// Should fail if one of the keys is wrong type
	s = config.GetInts([]string{"y", "a.port", "x.b.c.d"})
	suite.Equal([]int{0, 5000, 0}, s)
}

func (suite *YamlTestSuite) TestGetBool() {
	// Create tmp yaml object
	fname := suite.writeToTempFile(boolYaml)
	config, err := New(fname)
	suite.Nil(err)
	suite.NotNil(config)
	defer suite.removeTempFile(fname) // clean up

	// Should be able to get a key that exists
	b := config.GetBool("i.j")
	suite.Equal(true, b)

	// Should be able to get a key that exists
	b = config.GetBool("x")
	suite.Equal(false, b)

	// Should be able to get a key that exists
	b = config.GetBool("a.b.c.d")
	suite.Equal(true, b)

	// Should fail if key does not exist
	b = config.GetBool("a.keydoesntexist")
	suite.Equal(false, b)

	// Should fail if key is empty
	b = config.GetBool("")
	suite.Equal(false, b)

	// Should fail if a middle part of path doesnt exist
	b = config.GetBool("x.b.fake.d")
	suite.Equal(false, b)

	// Should fail if we try to get a key that is not a bool
	b = config.GetBool("y")
	suite.Equal(false, b)
}

func (suite *YamlTestSuite) TestGetInt() {
	// Create tmp yaml object
	fname := suite.writeToTempFile(testYaml)
	config, err := New(fname)
	suite.Nil(err)
	suite.NotNil(config)
	defer suite.removeTempFile(fname) // clean up

	// Should be able to get a key that exists
	i := config.GetInt("a.port")
	suite.Equal(5000, i)

	// Should be able to get a key that exists
	i = config.GetInt("z")
	suite.Equal(100, i)

	// Should be able to get a key that exists
	i = config.GetInt("x.b.c.e")
	suite.Equal(100, i)

	// Should fail if key does not exist
	i = config.GetInt("a.keydoesntexist")
	suite.Equal(0, i)

	// Should fail if key is empty
	i = config.GetInt("")
	suite.Equal(0, i)

	// Should fail if a middle part of path doesnt exist
	i = config.GetInt("x.b.fake.d")
	suite.Equal(0, i)

	// Should fail if we try to get a key that is not an int
	i = config.GetInt("a.host")
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
	s := config.GetString("a.host")
	suite.Equal("localhost", s)

	// Should be able to get a key that exists
	s = config.GetString("y")
	suite.Equal("jack", s)

	// Should be able to get a key that exists
	s = config.GetString("x.b.c.d")
	suite.Equal("hello", s)

	// Should fail if key does not exist
	s = config.GetString("a.keydoesntexist")
	suite.Equal("", s)

	// Should fail if key does not exist
	s = config.GetString("")
	suite.Equal("", s)

	// Should fail if a middle part of path doesnt exist
	s = config.GetString("x.b.fake.d")
	suite.Equal("", s)

	// Should fail if we try to access key that is not string
	s = config.GetString("a.port")
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
