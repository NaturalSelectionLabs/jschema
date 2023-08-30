package jschema_test

import (
	"encoding/json"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/NaturalSelectionLabs/jschema"
	"github.com/NaturalSelectionLabs/jschema/lib/test"
	"github.com/ysmood/got"
)

func TestTypeName(t *testing.T) {
	g := got.T(t)

	g.Eq(reflect.TypeOf(1).PkgPath(), "")
}

func TestNil(t *testing.T) {
	g := got.T(t)

	type A struct {
		A *A
	}

	c := jschema.New("")

	out := c.Define(A{})

	g.Eq(g.JSON(g.ToJSONString(out)), map[string]interface{}{
		"$ref": `#/$defs/A`, /* len=42 */
	})
}

func TestCommonSchema(t *testing.T) {
	g := got.T(t)

	c := jschema.New("")

	type Node2 struct {
		Map map[string]int
		Any interface{}
	}

	type Node1 struct {
		Str     string
		Num     int    `json:"num,omitempty"`
		Bool    bool   `json:"bool"`
		Ignore  string `json:"-"`
		Slice   []Node1
		Arr     [2]int
		Obj     *Node2
		Enum    test.Enum
		private int //nolint: unused
	}

	c.Define(Node1{})
	c.Define(Node2{})

	g.Eq(g.JSON(g.ToJSONString(c.Define(Node1{}))), map[string]interface{}{
		"$ref": "#/$defs/Node1",
	})

	g.Eq(g.JSON(c.String()), map[string]interface{} /* len=3 */ {
		"Enum": map[string]interface{} /* len=4 */ {
			"description": `github.com/NaturalSelectionLabs/jschema/lib/test.Enum`, /* len=53 */
			"enum": []interface{} /* len=3 cap=4 */ {
				"one",
				"two",
				"three",
			},
			"title": "Enum",
			"type":  "string",
		},
		"Node1": map[string]interface{} /* len=6 */ {
			`additionalProperties` /* len=20 */ : false,
			"description":                        `github.com/NaturalSelectionLabs/jschema_test.Node1`, /* len=50 */
			"properties": map[string]interface{} /* len=7 */ {
				"Arr": map[string]interface{} /* len=4 */ {
					"items": map[string]interface{}{
						"type": "number",
					},
					"maxItems": 2.0,
					"minItems": 2.0,
					"type":     "array",
				},
				"Enum": map[string]interface{}{
					"$ref": "#/$defs/Enum",
				},
				"Obj": map[string]interface{}{
					"anyOf": []interface{} /* len=2 cap=2 */ {
						map[string]interface{}{
							"$ref": "#/$defs/Node2",
						},
						map[string]interface{}{
							"type": "null",
						},
					},
				},
				"Slice": map[string]interface{} /* len=2 */ {
					"items": map[string]interface{}{
						"$ref": "#/$defs/Node1",
					},
					"type": "array",
				},
				"Str": map[string]interface{}{
					"type": "string",
				},
				"bool": map[string]interface{}{
					"type": "boolean",
				},
				"num": map[string]interface{}{
					"type": "number",
				},
			},
			"required": []interface{} /* len=6 cap=8 */ {
				"Str",
				"bool",
				"Slice",
				"Arr",
				"Obj",
				"Enum",
			},
			"title": "Node1",
			"type":  "object",
		},
		"Node2": map[string]interface{} /* len=6 */ {
			`additionalProperties` /* len=20 */ : false,
			"description":                        `github.com/NaturalSelectionLabs/jschema_test.Node2`, /* len=50 */
			"properties": map[string]interface{} /* len=2 */ {
				"Any": map[string]interface{}{
					"type": "object",
				},
				"Map": map[string]interface{} /* len=2 */ {
					`patternProperties` /* len=17 */ : map[string]interface{}{
						"": map[string]interface{}{
							"type": "number",
						},
					},
					"type": "object",
				},
			},
			"required": []interface{} /* len=2 cap=2 */ {
				"Map",
				"Any",
			},
			"title": "Node2",
			"type":  "object",
		},
	})
}

func TestHandler(t *testing.T) {
	g := got.T(t)

	c := jschema.New("")

	type A struct {
		Str string
	}
	type B struct {
		A A
	}

	c.AddHandler(A{}, func() *jschema.Schema {
		return &jschema.Schema{
			Type: "number",
		}
	})

	c.Define(B{})

	g.Eq(g.JSON(c.String()), map[string]interface{} /* len=2 */ {
		"A": map[string]interface{}{
			"type": "number",
		},
		"B": map[string]interface{} /* len=6 */ {
			`additionalProperties` /* len=20 */ : false,
			"description":                        `github.com/NaturalSelectionLabs/jschema_test.B`, /* len=57 */
			"properties": map[string]interface{}{
				"A": map[string]interface{}{
					"type": "number",
				},
			},
			"required": []interface{} /* len=1 cap=1 */ {
				"A",
			},
			"title": "B",
			"type":  "object",
		},
	})
}

func TestTime(t *testing.T) {
	g := got.T(t)

	c := jschema.New("")
	c.AddTimeHandler()
	c.Define(time.Now())

	g.Eq(g.JSON(c.String()), map[string]interface{}{
		`Time` /* len=37 */ : map[string]interface{} /* len=3 */ {
			"description": "time.Time",
			"title":       "Time",
			"type":        "string",
		},
	})
}

func TestBigInt(t *testing.T) {
	g := got.T(t)

	c := jschema.New("")
	c.AddBigIntHandler()
	c.Define(big.Int{})

	g.Eq(g.JSON(c.String()), map[string]interface{}{
		`Int` /* len=36 */ : map[string]interface{} /* len=3 */ {
			"description": "math/big.Int",
			"title":       "Int",
			"type":        "number",
		},
	})
}

func TestNameConflict(t *testing.T) {
	g := got.T(t)

	c := jschema.New("")

	type Time struct {
		Name string
	}

	c.Define(time.Time{})
	c.Define(Time{})

	g.Eq(g.JSON(c.String()), map[string]interface{} /* len=2 */ {
		"Time": map[string]interface{} /* len=4 */ {
			`additionalProperties` /* len=20 */ : false,
			"description":                        "time.Time",
			"title":                              "Time",
			"type":                               "object",
		},
		"Time1": map[string]interface{} /* len=6 */ {
			`additionalProperties` /* len=20 */ : false,
			"description":                        `github.com/NaturalSelectionLabs/jschema_test.Time`, /* len=60 */
			"properties": map[string]interface{}{
				"Name": map[string]interface{}{
					"type": "string",
				},
			},
			"required": []interface{} /* len=1 cap=1 */ {
				"Name",
			},
			"title": "Time",
			"type":  "object",
		},
	})
}

func TestRawMessage(t *testing.T) {
	g := got.T(t)

	c := jschema.New("")
	c.AddJSONRawMessageHandler()

	type A struct {
		A json.RawMessage
	}

	c.Define(A{})

	g.Eq(g.JSON(c.String()), map[string]interface{} /* len=2 */ {
		"A": map[string]interface{} /* len=6 */ {
			`additionalProperties` /* len=20 */ : false,
			"description":                        `github.com/NaturalSelectionLabs/jschema_test.A`, /* len=57 */
			"properties": map[string]interface{}{
				"A": map[string]interface{} /* len=2 */ {
					"description": `encoding/json.RawMessage`, /* len=24 */
					"title":       "RawMessage",
				},
			},
			"required": []interface{} /* len=1 cap=1 */ {
				"A",
			},
			"title": "A",
			"type":  "object",
		},
		"RawMessage": map[string]interface{} /* len=2 */ {
			"description": `encoding/json.RawMessage`, /* len=24 */
			"title":       "RawMessage",
		},
	})
}

func TestRef(t *testing.T) {
	g := got.T(t)

	c := jschema.New("")

	type A struct{}

	type B struct{ A A }

	c.Define(B{})

	g.Eq(c.PeakSchema(A{}).Title, "A")
}

func TestEmbeddedStruct(t *testing.T) {
	g := got.T(t)

	type A struct{ Val int }

	type B struct {
		A
	}

	c := jschema.New("")

	c.Define(B{})

	g.Eq(g.JSON(c.String()), map[string]interface{} /* len=2 */ {
		"A": map[string]interface{} /* len=6 */ {
			`additionalProperties` /* len=20 */ : false,
			"description":                        `github.com/NaturalSelectionLabs/jschema_test.A`, /* len=46 */
			"properties": map[string]interface{}{
				"Val": map[string]interface{}{
					"type": "number",
				},
			},
			"required": []interface{} /* len=1 cap=1 */ {
				"Val",
			},
			"title": "A",
			"type":  "object",
		},
		"B": map[string]interface{} /* len=6 */ {
			`additionalProperties` /* len=20 */ : false,
			"description":                        `github.com/NaturalSelectionLabs/jschema_test.B`, /* len=46 */
			"properties": map[string]interface{}{
				"Val": map[string]interface{}{
					"type": "number",
				},
			},
			"required": []interface{} /* len=1 cap=1 */ {
				"A",
			},
			"title": "B",
			"type":  "object",
		},
	})
}
