package information

import "fmt"

//We use ldflags
var (
	Version     = "No Version Provided"
	GitCommit   = "No GitCommit Provided"
	GitDescribe = "No GitDescribe Provided"
	GitDirty    = "No GitDirty Provided"
	BuildTime   = "No BuildTime Provided"
)

//Print information on project
func Print() string {
	return fmt.Sprintf("%q\t%q\t%q\t%q\t%q", Version, BuildTime, GitCommit, GitDescribe, GitDirty)
}
