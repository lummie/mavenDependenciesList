package main

import (
	"os"
	"fmt"
	"path/filepath"
	"mavenDependencyList/util"
	"runtime"
	"strings"
	"sort"
	"flag"
)


// processPom opens the pom file, reads the contents and sends the pom to the poms channel
func processPom(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	pom, err := util.ReadPom(file)
	poms <- pom
}

// visit is the visitor function for selecting the pom.xml files from the scanned files and directories
func visit(path string, f os.FileInfo, err error) error {
	if !f.IsDir() {
		if f.Name() == "pom.xml" {
			go processPom(path)
		}
	}
	return err
}


// PomScanner scans the specified directory and sub-directories for pom files.
func PomScanner(rootDir string, quit chan bool) error{
	if err := filepath.Walk(rootDir, visit); err != nil {
		panic(err)
	}
	quit <- true
	return nil
}


// sort


// channel to receive poms on
var poms chan *util.Pom = make(chan *util.Pom)


type pomInfo []*util.Pom


// sort for artifacts
type ByArtifact pomInfo
func (s ByArtifact) Len() int {
	return len(s)
}
func (s ByArtifact) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByArtifact) Less(i, j int) bool {
	return strings.Compare(s[i].ArtifactId,s[j].ArtifactId) == -1
}

//sort for dependencies
type ByDependency []util.Dependency
func (s ByDependency) Len() int {
	return len(s)
}
func (s ByDependency) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByDependency) Less(i, j int) bool {
	return strings.Compare(s[i].GroupId,s[j].GroupId) == -1
}



// main
func main() {
	runtime.GOMAXPROCS(runtime.GOMAXPROCS(0))

	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Println("USAGE: mavenDependencyList [searchPath] ")
		os.Exit(1)
	}

	// array to collect results
	results := make([]*util.Pom,0)

	// channel for notification when scanner is complete
	quit := make(chan bool)

	// start scanner
	target := flag.Arg(0)
	go PomScanner(target, quit)

	for {
		var quitting bool = false

		select {
		case quitting = <-quit:
			break
		case pom := <- poms:
			results = append(results,pom)
		}

		if quitting {
			break
		}

	}

	sort.Sort(ByArtifact(results))

	for _, artifact := range results {
		if len(artifact.Dependencies.Dependency) > 0 {
			sort.Sort(ByDependency(artifact.Dependencies.Dependency))
			for _, dep := range artifact.Dependencies.Dependency {
				if dep.Scope == ""{
					version := dep.Version
					if strings.HasPrefix(version,"${") {
						if found, versionProp := artifact.GetProperty(version); found == true {
							version = versionProp
						}
					}
					fmt.Printf("%s, %s, %s, %s\n", artifact.ArtifactId, dep.GroupId, dep.ArtifactId, version)
				}
			}
		}
	}

	// close the poms channel
	close(poms)
}