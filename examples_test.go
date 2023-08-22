package jschema_test

import (
	"encoding/json"
	"fmt"

	"github.com/NaturalSelectionLabs/jschema"
)

func ExampleNew() {
	type Node struct {
		// The default tag only accepts json string.
		// So if you want to set a string value "jack",
		// you should use "\"jack\"" instead of "jack" for the field tag
		ID int `json:"id" default:"1"`

		// Use the description tag to set the description of the field
		Children []*Node `json:"children" description:"The children of the node"`
	}

	// Create a schema list instance
	schemas := jschema.New("#/components/schemas")

	// Define a type within the schema
	schemas.Define(Node{})
	schemas.Description(Node{}, "A node in the tree")

	// Marshal the schema list to json string
	out, err := json.MarshalIndent(schemas.JSON(), "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(out))

	// Output:
	// {
	//   "Node": {
	//     "type": "object",
	//     "title": "Node",
	//     "description": "A node in the tree",
	//     "properties": {
	//       "children": {
	//         "type": "array",
	//         "description": "The children of the node",
	//         "items": {
	//           "nullable": true,
	//           "anyOf": [
	//             {
	//               "$ref": "#/components/schemas/Node"
	//             }
	//           ]
	//         }
	//       },
	//       "id": {
	//         "type": "number",
	//         "default": 1
	//       }
	//     },
	//     "required": [
	//       "id",
	//       "children"
	//     ],
	//     "additionalProperties": false
	//   }
	// }
}

func ExampleSchemas() {
	// Create a schema list instance
	schemas := jschema.New("#/components/schemas")

	type A string
	type B int

	type Node struct {
		Name     int         `json:"name"`
		Metadata interface{} `json:"metadata,omitempty"` // omitempty make this field optional
		Version  string      `json:"version"`
		Options  []string    `json:"options"`
	}

	schemas.Define(Node{})
	node := schemas.PeakSchema(Node{})

	// Define default value
	{
		node.Properties["name"].Default = "jack"
	}

	// Make the metadata field accept either A or B
	{
		node.Properties["metadata"] = schemas.AnyOf(A(""), B(0))
	}

	// Define constants
	{
		node.Properties["version"] = schemas.Const("v1")
	}

	// Define enum
	{
		node.Properties["options"].Enum = jschema.ToJValList(1, 2, 3)
	}

	// Marshal the schema list to json string
	out, err := json.MarshalIndent(schemas.JSON(), "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(out))

	// Output:
	// {
	//   "A": {
	//     "type": "string",
	//     "title": "A",
	//     "description": "github.com/NaturalSelectionLabs/jschema_test.A"
	//   },
	//   "B": {
	//     "type": "number",
	//     "title": "B",
	//     "description": "github.com/NaturalSelectionLabs/jschema_test.B"
	//   },
	//   "Node": {
	//     "type": "object",
	//     "title": "Node",
	//     "description": "github.com/NaturalSelectionLabs/jschema_test.Node",
	//     "properties": {
	//       "metadata": {
	//         "anyOf": [
	//           {
	//             "$ref": "#/components/schemas/A"
	//           },
	//           {
	//             "$ref": "#/components/schemas/B"
	//           }
	//         ]
	//       },
	//       "name": {
	//         "type": "number",
	//         "default": "jack"
	//       },
	//       "options": {
	//         "type": "array",
	//         "enum": [
	//           1,
	//           2,
	//           3
	//         ],
	//         "items": {
	//           "type": "string"
	//         }
	//       },
	//       "version": {
	//         "type": "string",
	//         "enum": [
	//           "v1"
	//         ]
	//       }
	//     },
	//     "required": [
	//       "name",
	//       "version",
	//       "options"
	//     ],
	//     "additionalProperties": false
	//   }
	// }
}
