package deps

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/jackspirou/govend/app/helpers"
)

//
// Note: The deps package is built very much like a config package.

// Public vars in the deps package.
var (
	Dir  string
	List []string
)

// DepsConfig represents the deps.json file.
type DepsConfig struct {
	Dir     string   `json:"dir"`
	DepList []string `json:"deps"`
}

// Init reads the deps.json file and sets defaults.
func init() {

	// Check that a deps.json file exists.
	if _, err := os.Stat("./deps.json"); os.IsNotExist(err) {
		log.Fatalln("Unable to find deps.json.")
	}

	// Read the deps.json file into a slice of bytes.
	bytes, err := ioutil.ReadFile("./deps.json")
	helpers.Check(err)

	// Build the deps config object.
	d := &DepsConfig{}
	err = json.Unmarshal(bytes, d)
	helpers.Check(err)

	// Assign the deps dir to Dir.
	Dir = d.Dir

	// If no deps dir was provided, set the default.
	if len(d.Dir) < 1 {
		Dir = "./vendor"
	}

	// Assign the dependency list to List.
	List = d.DepList
}
