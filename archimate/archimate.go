// used to generate PlantUML files and demo of all elements and relationships (see ARCHIMATE.md).
// usage:
// go run archimate.go
// java -jar ~/lib/plantuml.jar archimate-spec.puml
package main

import (
	"log"
	"os"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"
)

const motivation = "#ccccff"
const strategy = "#f5deaa"
const business = "#ffffb5"
const application = "#b5ffff"
const technology = "#c9e7b7"
const physical = "#c9e7b7"
const implementation = "#ffe0e0"

type element struct {
	Name       string
	Definition string
	StereoType string
	Style      string
	Colour     string
}

type elementGroup struct {
	Name     string
	Elements []element
}

type relationship struct {
	Name       string
	Definition string
	Plant      string
}

type relationshipGroup struct {
	Name          string
	Relationships []relationship
}

type archimate struct {
	Elements      []elementGroup
	Relationships []relationshipGroup
}

var archi = archimate{
	Elements: []elementGroup{
		{Name: "Motivation Elements",
			Elements: []element{
				{Name: "Stakeholder", Definition: "Represents the role of an individual, team, or organization (or classes thereof) that represents their interests in the effects of the architecture.",
					StereoType: "<<motivation-stakeholder>>", Style: "archimate", Colour: motivation},
				{Name: "Driver", Definition: "Represents an external or internal condition that motivates an organization to define its goals and implement the changes necessary to achieve them.",
					StereoType: "<<motivation-driver>>", Style: "archimate", Colour: motivation},
				{Name: "Assessment", Definition: "Represents the result of an analysis of the state of affairs of the enterprise with respect to some driver.",
					StereoType: "<<motivation-assessment>>", Style: "archimate", Colour: motivation},
				{Name: "Goal", Definition: "Represents a high-level statement of intent, direction, or desired end state for an organization and its stakeholders.",
					StereoType: "<<motivation-goal>>", Style: "archimate", Colour: motivation},
				{Name: "Outcome", Definition: "Represents an end result.",
					StereoType: "<<motivation-outcome>>", Style: "archimate", Colour: motivation},
				{Name: "Principle", Definition: "Represents a statement of intent defining a general property that applies to any system in a certain context in the architecture.",
					StereoType: "<<motivation-principle>>", Style: "archimate", Colour: motivation},
				{Name: "Requirement", Definition: "Represents a statement of need defining a property that applies to a specific system as described by the architecture.",
					StereoType: "<<motivation-requirement>>", Style: "archimate", Colour: motivation},
				{Name: "Constraint", Definition: "Represents a factor that limits the realization of goals.",
					StereoType: "<<motivation-constraint>>", Style: "archimate", Colour: motivation},
				{Name: "Meaning", Definition: "Represents the knowledge or expertise present in, or the interpretation given to, a concept in a particular context.",
					StereoType: "<<motivation-meaning>>", Style: "archimate", Colour: motivation},
				{Name: "Value", Definition: "Represents the relative worth, utility, or importance of a concept.",
					StereoType: "<<motivation-value>>", Style: "archimate", Colour: motivation},
			},
		},
		{Name: "Strategy Elements",
			Elements: []element{
				{Name: "Resource", Definition: "Represents an asset owned or controlled by an individual or organization.",
					StereoType: "<<strategy-resource>>", Style: "archimate", Colour: strategy},
				{Name: "Capability", Definition: "Represents an ability that an active structure element, such as an organization, person, or system, possesses.",
					StereoType: "<<strategy-capability>>", Style: "archimate", Colour: strategy},
				{Name: "Value stream", Definition: "Represents a sequence of activities that  create an overall result for a customer, stakeholder, or end user.",
					StereoType: "<<strategy-value-stream>>", Style: "archimate", Colour: strategy},
				{Name: "Course of action", Definition: "Represents an approach or plan for configuring some capabilities and resources of the enterprise, undertaken to achieve a goal.",
					StereoType: "<<strategy-course-of-action>>", Style: "archimate", Colour: strategy},
			},
		},
		{Name: "Business Layer Elements",
			Elements: []element{
				{Name: "Business Actor", Definition: "Represents a business entity that is capable of performing behavior.",
					StereoType: "<<business-actor>>", Style: "archimate", Colour: business},
				{Name: "Business role", Definition: "Represents the responsibility for performing specific behavior, to which an actor can be assigned, or the part an actor plays in a particular action or event.",
					StereoType: "<<business-role>>", Style: "archimate", Colour: business},
				{Name: "Business collaboration", Definition: "Represents an aggregate of two or more business internal active structure Elements that work together to perform collective behavior.",
					StereoType: "<<business-collaboration>>", Style: "archimate", Colour: business},
				{Name: "Business interface", Definition: "Represents a point of access where a business service is made available to the environment.",
					StereoType: "<<business-interface>>", Style: "archimate", Colour: business},
				{Name: "Business process", Definition: "Represents a sequence of business behaviors that achieves a specific result such as a defined set of products or business services.",
					StereoType: "<<business-process>>", Style: "archimate", Colour: business},
				{Name: "Business function", Definition: "Represents a collection of business behavior based on a chosen set of criteria (typically required business resources and/or competencies), closely aligned to an organization, but not necessarily explicitly governed by the organization.",
					StereoType: "<<business-function>>", Style: "archimate", Colour: business},
				{Name: "Business interaction", Definition: "Represents a unit of collective business behavior performed by (a collaboration of) two or more business actors, business roles, or business collaborations.",
					StereoType: "<<business-interaction>>", Style: "archimate", Colour: business},
				{Name: "Business event", Definition: "Represents an organizational state change.",
					StereoType: "<<business-event>>", Style: "archimate", Colour: business},
				{Name: "Business service", Definition: "Represents explicitly defined behavior that a business role, business actor, or business collaboration exposes to its environment.",
					StereoType: "<<business-service>>", Style: "archimate", Colour: business},
				{Name: "Business object", Definition: "Represents a concept used within a particular business domain.",
					StereoType: "<<business-object>>", Style: "archimate", Colour: business},
				{Name: "Contract", Definition: "Represents a formal or informal specification of an agreement between a provider and a consumer that specifies the rights and obligations associated with a product and establishes functional and non-functional parameters for interaction.",
					StereoType: "<<business-contract>>", Style: "archimate", Colour: business},
				{Name: "Representation", Definition: "Represents a perceptible form of the information carried by a business object.",
					StereoType: "<<business-representation>>", Style: "archimate", Colour: business},
				{Name: "Product", Definition: "Represents a coherent collection of services and/or passive structure Elements, accompanied by a contract/set of agreements, which is offered as a whole to (internal or external) customers.",
					StereoType: "<<business-product>>", Style: "archimate", Colour: business},
			},
		},
		{Name: "Application Layer Elements",
			Elements: []element{
				{Name: "Application component", Definition: "Represents an encapsulation of application functionality aligned to implementation structure, which is modular and replaceable.",
					StereoType: "<<application-component>>", Style: "archimate", Colour: application},
				{Name: "Application collaboration", Definition: "Represents an aggregate of two or more application internal active structure Elements that work together to perform collective application behavior.",
					StereoType: "<<application-collaboration>>", Style: "archimate", Colour: application},
				{Name: "Application interface", Definition: "Represents a point of access where application services are made available to a user, another application component, or a node.",
					StereoType: "<<application-interface>>", Style: "archimate", Colour: application},
				{Name: "Application function", Definition: "Represents automated behavior that can be performed by an application component.",
					StereoType: "<<application-function>>", Style: "archimate", Colour: application},
				{Name: "Application interaction", Definition: "Represents a unit of collective application behavior performed by (a collaboration of) two or more application components.",
					StereoType: "<<application-interaction>>", Style: "archimate", Colour: application},
				{Name: "Application process", Definition: "Represents a sequence of application behaviors that achieves a specific result.",
					StereoType: "<<application-process>>", Style: "archimate", Colour: application},
				{Name: "Application event", Definition: "Represents an application state change.",
					StereoType: "<<application-event>>", Style: "archimate", Colour: application},
				{Name: "Application service", Definition: "Represents an explicitly defined exposed application behavior.",
					StereoType: "<<application-service>>", Style: "archimate", Colour: application},
				{Name: "Data object", Definition: "Represents data structured for automated processing.",
					StereoType: "<<application-data-object>>", Style: "archimate", Colour: application},
			},
		},
		{Name: "Technology Layer Elements",
			Elements: []element{
				{Name: "Node", Definition: "Represents a computational or physical resource that hosts, manipulates, or interacts with other computational or physical resources.",
					StereoType: "<<technology-node>>", Style: "archimate", Colour: technology},
				{Name: "Device", Definition: "Represents a physical IT resource upon which system software and artifacts may be stored or deployed for execution.",
					StereoType: "<<technology-device>>", Style: "archimate", Colour: technology},
				{Name: "System software", Definition: "Represents software that provides or contributes to an environment for storing, executing, and using software or data deployed within it.",
					StereoType: "<<technology-system-software>>", Style: "archimate", Colour: technology},
				{Name: "Technology collaboration", Definition: "Represents an aggregate of two or more technology internal active structure Elements that work together to perform collective technology behavior.",
					StereoType: "<<technology-collaboration>>", Style: "archimate", Colour: technology},
				{Name: "Technology interface", Definition: "Represents a point of access where technology services offered by a node can be accessed.",
					StereoType: "<<technology-interface>>", Style: "archimate", Colour: technology},
				{Name: "Path", Definition: "Represents a link between two or more nodes, through which these nodes can exchange data, energy, or material.",
					StereoType: "<<technology-path>>", Style: "archimate", Colour: technology},
				{Name: "Communication network", Definition: "Represents a set of structures that connects nodes for transmission, routing, and reception of data.",
					StereoType: "<<technology-network>>", Style: "archimate", Colour: technology},
				{Name: "Technology function", Definition: "Represents a collection of technology behavior that can be performed by a node.",
					StereoType: "<<technology-function>>", Style: "archimate", Colour: technology},
				{Name: "Technology process", Definition: "Represents a sequence of technology behaviors that achieves a specific result.",
					StereoType: "<<technology-process>>", Style: "archimate", Colour: technology},
				{Name: "Technology interaction", Definition: "Represents a unit of collective technology behavior performed by (a collaboration of) two or more nodes.",
					StereoType: "<<technology-interaction>>", Style: "archimate", Colour: technology},
				{Name: "Technology event", Definition: "Represents a technology state change.",
					StereoType: "<<technology-event>>", Style: "archimate", Colour: technology},
				{Name: "Technology service", Definition: "Represents an explicitly defined exposed technology behavior.",
					StereoType: "<<technology-service>>", Style: "archimate", Colour: technology},
				{Name: "Artifact", Definition: "Represents a piece of data that is used or produced in a software development process, or by deployment and operation of an IT system.",
					StereoType: "<<technology-artifact>>", Style: "archimate", Colour: technology},
			},
		},
		{Name: "Physical Elements",
			Elements: []element{
				{Name: "Equipment", Definition: "Represents one or more physical machines, tools, or instruments that can create, use, store, move, or transform materials.",
					StereoType: "<<physical-equipment>>", Style: "archimate", Colour: physical},
				{Name: "Facility", Definition: "Represents a physical structure or environment.",
					StereoType: "<<physical-facility>>", Style: "archimate", Colour: physical},
				{Name: "Distribution network", Definition: "Represents a physical network used to transport materials or energy.",
					StereoType: "<<physical-distribution-network>>", Style: "archimate", Colour: physical},
				{Name: "Material", Definition: "Represents tangible physical matter or energy.",
					StereoType: "<<physical-material>>", Style: "archimate", Colour: physical},
			},
		},
		{Name: "Migration Elements",
			Elements: []element{
				{Name: "Work package", Definition: "Represents a series of actions identified and designed to achieve specific results within specified time and resource constraints.",
					StereoType: "<<implementation-workpackage>>", Style: "archimate", Colour: implementation},
				{Name: "Deliverable", Definition: "Represents a precisely-defined result of a work package.",
					StereoType: "<<implementation-deliverable>>", Style: "archimate", Colour: implementation},
				{Name: "Implementation event", Definition: "Represents a state change related to implementation or migration.",
					StereoType: "<<implementation-event>>", Style: "archimate", Colour: implementation},
				{Name: "Plateau", Definition: "Represents a relatively stable state of the architecture that exists during a limited period of time.",
					StereoType: "<<implementation-plateau>>", Style: "archimate", Colour: "#e0ffe0"},
				{Name: "Gap", Definition: "Represents a statement of difference between two plateaus.",
					StereoType: "<<implementation-gap>>", Style: "archimate", Colour: "#e0ffe0"},
			},
		},
		{Name: "Composite Elements",
			Elements: []element{
				{Name: "Location", Definition: "A location represents a conceptual or physical place or position where concepts are located (e.g., structure Elements) or performed (e.g., behavior Elements).",
					StereoType: "<<location>>", Style: "archimate", Colour: "#fbb875"},
				{Name: "Grouping", Definition: "The grouping element aggregates or composes concepts that belong together based on some common characteristic.",
					StereoType: "<<grouping>>", Style: "Folder", Colour: "#ffffff"},
			},
		},
	},
	Relationships: []relationshipGroup{
		{Name: "Structural Relationships",
			Relationships: []relationship{
				{Name: "Composition", Definition: "Represents that an element consists of one or more other concepts.",
					Plant: "*-$orientation-"},
				{Name: "Aggregation", Definition: "Represents that an element combines one or more other concepts.",
					Plant: "o-$orientation-"},
				{Name: "Assignment", Definition: "Represents the allocation of responsibility, performance of behavior, storage, or execution.",
					Plant: "0-$orientation->"},
				{Name: "Realization", Definition: "Represents that an entity plays a critical role in the creation, achievement, sustenance, or operation of a more abstract entity.",
					Plant: "~$orientation~|>"},
			},
		},
		{Name: "Dependency Relationships",
			Relationships: []relationship{
				{Name: "Serving", Definition: "Represents that an element provides its functionality to another element.",
					Plant: "-$orientation->"},
				{Name: "Access", Definition: "Represents the ability of behavior and active structure elements to observe or act upon passive structure elements.",
					Plant: "~$orientation~"},
				{Name: "Access R", Definition: "Role modifier for PlantUML",
					Plant: "<~$orientation~"},
				{Name: "Access W", Definition: "Role modifier for PlantUML",
					Plant: "~$orientation~>"},
				{Name: "Access RW", Definition: "Role modifier for PlantUML",
					Plant: "<~$orientation~>"},
				{Name: "Influence", Definition: "Represents that an element affects the implementation or achievement of some motivation element.",
					Plant: ".$orientation.>"},
				{Name: "Association", Definition: "Represents an unspecified relationship, or one that is not represented by another ArchiMate relationship.",
					Plant: "-$orientation-"},
			},
		},
		{Name: "Dynamic Relationships",
			Relationships: []relationship{
				{Name: "Triggering", Definition: "Represents a temporal or causal relationship between elements.",
					Plant: "-$orientation->>"},
				{Name: "Flow", Definition: "Represents transfer from one element to another.",
					Plant: ".$orientation.>>"},
			},
		},
		{Name: "Other Relationships",
			Relationships: []relationship{
				{Name: "Specialization", Definition: "Represents that an element is a particular kind of another element.",
					Plant: "-$orientation-|>"},
			},
		},
	},
}

func main() {
	a, err := template.ParseFiles("archimate.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	af, err := os.Create("../Archimate.puml")
	if err != nil {
		log.Fatal(err)
	}
	defer af.Close()

	err = a.Execute(af, archi)
	if err != nil {
		log.Fatal(err)
	}

	s, err := template.ParseFiles("archimate-spec.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	sf, err := os.Create("archimate-spec.puml")
	if err != nil {
		log.Fatal(err)
	}
	defer sf.Close()

	err = s.Execute(sf, archi)
	if err != nil {
		log.Fatal(err)
	}

	m, err := template.ParseFiles("archimate-md.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	mf, err := os.Create("ARCHIMATE.md")
	if err != nil {
		log.Fatal(err)
	}
	defer mf.Close()

	err = m.Execute(mf, archi)
	if err != nil {
		log.Fatal(err)
	}
}

func (e element) Cmd() string {
	s := strings.ReplaceAll(strings.Title(e.Name), " ", "")

	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[n:]
}

func (e element) Label() string {
	return strings.ReplaceAll(e.Name, " ", "\\n")
}

func (r relationship) Cmd() string{
	s := strings.ReplaceAll(r.Name, " ", "")

	if s == "" {
		return ""
	}
	rr, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(rr)) + s[n:]
}