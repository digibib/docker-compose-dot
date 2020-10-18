package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"log"

	"github.com/awalterschulze/gographviz"
	yaml "gopkg.in/yaml.v2"
)

type config struct {
	Version  string
	Networks map[string]network
	Volumes  map[string]volume
	Services map[string]service
}

type network struct {
	Driver, External string
	DriverOpts       map[string]string "driver_opts"
}

type volume struct {
	Driver, External string
	DriverOpts       map[string]string "driver_opts"
}

type MapOrArrayWrapper []string

type service struct {
	ContainerName                     string "container_name"
	Image                             string
	Networks, Ports, Volumes, Command, Links []string
	VolumesFrom                       []string "volumes_from"
	DependsOn                         []string "depends_on"
	CapAdd                            []string "cap_add"
	Build                             struct{ Context, Dockerfile string }
	Environment                       MapOrArrayWrapper
}

func (w *MapOrArrayWrapper) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var envsArray []string
	var envsMap map[string]string
	if err := unmarshal(&envsMap); err == nil {
		for key, val := range envsMap {
			envsArray = append(envsArray, key + "=" + val)
		}
	}

	if len(envsArray) == 0 {
		if err := unmarshal(&envsArray); err != nil {
			return err
		}
	}
	*w = envsArray
	return nil
}

func nodify(s string) string {
	return strings.Replace(s, "-", "_", -1)
}

func main() {
	var (
		bytes   []byte
		err     error
		graph   *gographviz.Graph
		project string
	)

	if len(os.Args) < 2 {
		log.Fatal("Need input file!")
	}

	bytes, err = ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	// Parse it as YML
	data := &config{}
	err = yaml.Unmarshal(bytes, &data)
	if err != nil {
		log.Fatal(err)
	}

	// Create directed graph
	graph = gographviz.NewGraph()
	graph.SetName(project)
	graph.SetDir(true)

	// Add legend
	graph.AddSubGraph(project, "cluster_legend", map[string]string{"label": "Legend"})
	graph.AddNode("cluster_legend", "legend_service",
		map[string]string{"shape": "plaintext",
			"label": "<<TABLE BORDER='0'>" +
				"<TR><TD BGCOLOR='lightblue'><B>container_name</B></TD></TR>" +
				"<TR><TD BGCOLOR='lightgrey'><FONT POINT-SIZE='9'>ports ext:int</FONT></TD></TR>" +
				"<TR><TD BGCOLOR='orange'><FONT POINT-SIZE='9'>volumes host:container</FONT></TD></TR>" +
				"<TR><TD BGCOLOR='pink'><FONT POINT-SIZE='9'>environment</FONT></TD></TR>" +
				"</TABLE>>",
		})

	/** NETWORK NODES **/
	for name := range data.Networks {
		graph.AddNode(project, nodify(name), map[string]string{
			"label":     fmt.Sprintf("\"Network: %s\"", name),
			"style":     "filled",
			"shape":     "box",
			"fillcolor": "palegreen",
		})
	}

	/** SERVICE NODES **/
	for name, service := range data.Services {
		var attrs = map[string]string{"shape": "plaintext", "label": "<<TABLE BORDER='0'>"}
		attrs["label"] += fmt.Sprintf("<TR><TD BGCOLOR='lightblue'><B>%s</B></TD></TR>", name)

		if service.Ports != nil {
			for _, port := range service.Ports {
				attrs["label"] += fmt.Sprintf("<TR><TD BGCOLOR='lightgrey'><FONT POINT-SIZE='9'>%s</FONT></TD></TR>", port)
			}
		}
		if service.Volumes != nil {
			for _, vol := range service.Volumes {
				attrs["label"] += fmt.Sprintf("<TR><TD BGCOLOR='orange'><FONT POINT-SIZE='9'>%s</FONT></TD></TR>", vol)
			}
		}
		if service.Environment != nil {
			for _, v := range service.Environment {
				attrs["label"] += fmt.Sprintf("<TR><TD BGCOLOR='pink'><FONT POINT-SIZE='9'>%s</FONT></TD></TR>", v)
			}
		}
		attrs["label"] += "</TABLE>>"
		graph.AddNode(project, nodify(name), attrs)
	}
	/** EDGES **/
	for name, service := range data.Services {
		// Links to networks
		if service.Networks != nil {
			for _, linkTo := range service.Networks {
				if strings.Contains(linkTo, ":") {
					linkTo = strings.Split(linkTo, ":")[0]
				}
				graph.AddEdge(nodify(name), nodify(linkTo), true,
					map[string]string{"dir": "none"})
			}
		}
		// volumes_from
		if service.VolumesFrom != nil {
			for _, linkTo := range service.VolumesFrom {
				graph.AddEdge(nodify(name), nodify(linkTo), true,
					map[string]string{"style": "dashed", "label": "volumes_from"})
			}
		}
		// depends_on
		if service.DependsOn != nil {
			for _, linkTo := range service.DependsOn {
				graph.AddEdge(nodify(name), nodify(linkTo), true,
					map[string]string{"style": "dashed", "label": "depends_on"})
			}
		}
		// links
		if service.Links != nil {
			for _, linkTo := range service.Links {
				if strings.Contains(linkTo, ":") {
					linkTo = strings.Split(linkTo, ":")[0]
				}
				graph.AddEdge(nodify(name), nodify(linkTo), true,
					map[string]string{"style": "dashed", "label": "links"})
			}
		}
	}
	fmt.Print(graph)
}
