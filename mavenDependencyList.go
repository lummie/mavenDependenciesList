package main

import (
	"os"
	"fmt"
	"path/filepath"
	"mavenDependencyList/util"
	"runtime"
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




// channel to receive poms on
var poms chan *util.Pom = make(chan *util.Pom)

// main
func main() {

	runtime.GOMAXPROCS(runtime.GOMAXPROCS(0))

	// array to collect results
	results := make([]*util.Pom,0)

	// channel for notification when scanner is complete
	quit := make(chan bool)

	// start scanner
	go PomScanner("/home/I049472/connected-goods", quit)

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


	for _, dep := range results {
		fmt.Printf("%+v\n", dep.ArtifactId)
	}

	// close the poms channel
	close(poms)
}