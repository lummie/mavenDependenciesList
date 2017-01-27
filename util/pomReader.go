package util

import (
	"io"
	"encoding/xml"
)

// Module represents the XML Module elements
type Modules struct {
	M []string `xml:"module"`
}

// Dependency represents the XML Dependency elements
type Dependency struct {
	GroupId string `xml:"groupId"`
	ArtifactId string `xml:"artifactId"`
	Version string `xml:"version"`
	Scope string `xml:"scope"`
}

// Dependencies represents the XML Dependencies elements
type Dependencies struct {
	Dependency []Dependency `xml:"dependency"`
}

// Pom represents the XML Root elements
type Pom struct {
	ModelVersion string `xml:"modelVersion"`
	GroupId string `xml:"groupId"`
	ArtifactId string `xml:"artifactId"`
	Modules Modules `xml:"modules"`
	Dependencies Dependencies`xml:"dependencies"`
}


// ReadPom reads the Pom XML and returns a Pom structure
func ReadPom(reader io.Reader) (*Pom, error){
	pom := Pom{}
	if err := xml.NewDecoder(reader).Decode(&pom); err != nil {
		return nil, err
	}
	return &pom, nil
}