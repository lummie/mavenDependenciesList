package util

import (
	"io"
	"encoding/xml"
	"strings"
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

type Property struct {
	XMLName xml.Name `xml:""`
	Value string `xml:",chardata"`
}

type Properties struct {
	Properties []Property `xml:",any"`
}


// Pom represents the XML Root elements
type Pom struct {
	ModelVersion string `xml:"modelVersion"`
	GroupId string `xml:"groupId"`
	ArtifactId string `xml:"artifactId"`
	Modules Modules `xml:"modules"`
	Dependencies Dependencies`xml:"dependencies"`
	Properties Properties `xml:"properties"`
}

// ReadPom reads the Pom XML and returns a Pom structure
func ReadPom(reader io.Reader) (*Pom, error){
	pom := Pom{}
	if err := xml.NewDecoder(reader).Decode(&pom); err != nil {
		return nil, err
	}
	return &pom, nil
}

func (p *Pom) GetProperty(propname string) (bool, string) {
	propname = strings.TrimLeft(propname,"${")
	propname = strings.TrimRight(propname,"}")
	for _, prop := range p.Properties.Properties {
		if prop.XMLName.Local == propname {
			return true, prop.Value
		}
	}
	return false,""
}