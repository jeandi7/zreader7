package zinterpreter

import (
	"fmt"
	"testing"
)

func TestReadSchema(t *testing.T) {
	tests := []struct {
		input       string
		expectError bool
	}{
		{input: `definition monsujet { }`, expectError: false},
		{input: `definition monsujet { relation marelation1: monsujet2 }`, expectError: false},
		{input: `definition monsujet { relation marelation1: monsujet2 | monsujet3 }`, expectError: false},
		{input: `definition monsujet { relation marelation1: monsujet2 | monsujet3  relation marelation2: monsujet2 }`, expectError: false},
		{input: `definition monsujet { } definition monsujet2 { } definition maressource { relation marelation: monsujet | monsujet2   }`, expectError: false},
		{input: `definition monsujet { } definition monsujet2 { } definition maressource { relation marelation: monsujet | monsujet2  relation marelation2: monsujet | monsujet2 | monsujet3  }`, expectError: false},
		{input: `definition monsujet {`, expectError: true}, // syntax error
		{input: `definition user { } definition group { relation member2: user }  definition organization { relation member: group#member2 }`, expectError: false},
	}

	for _, tt := range tests {
		lexer := NewLexer(tt.input)
		lexer.NextToken()

		_, err := lexer.ReadZSchema()

		if tt.expectError && err == nil {
			t.Errorf("expected an error but got none for input: %s", tt.input)
		}
		if !tt.expectError && err != nil {
			t.Errorf("did not expect an error but got one for input: %s, error: %v", tt.input, err)
		}

	}

}

func TestCreateIDs(t *testing.T) {
	tests := []struct {
		input       string
		expectError bool
	}{
		{input: `definition user { } definition group { relation member2: user }  definition organization { relation member: group#member2 }`, expectError: false},
		{input: `definition user { } definition document { relation reader: user | group#member } definition group {  relation member: user }`, expectError: false},
	}

	for _, tt := range tests {
		lexer := NewLexer(tt.input)
		lexer.NextToken()
		z, err := lexer.ReadZSchema()

		if tt.expectError && err == nil {
			t.Errorf("expected an error but got none for input: %s", tt.input)
		}
		if !tt.expectError && err != nil {
			t.Errorf("did not expect an error but got one for input: %s, error: %v", tt.input, err)
		}

		mydraw := PlantUMLArchimateSchema{Zdefs: z}
		mydraw.createIDforZdef()
		fmt.Printf("essai:")
	}

}

func TestDoubleRelationInDefinition(t *testing.T) {
	tests := []struct {
		input       string
		expectError bool
	}{
		// {input: `definition user { } definition document { relation reader: user relation reader: user } `, expectError: true},
		// {input: `definition user { } definition document { relation reader: user | group#member relation member: user } definition group { relation member: user } `, expectError: true},
		{input: `definition user { } definition document { relation reader: user | group#member | group#member relation member: user } definition group { relation member: user } `, expectError: true},
	}

	for _, tt := range tests {
		lexer := NewLexer(tt.input)
		lexer.NextToken()
		z, err := lexer.ReadZSchema()

		mydraw := PlantUMLArchimateSchema{Zdefs: z}
		mydraw.createIDforZdef()

		if tt.expectError && err != nil {
			t.Errorf("expected an error : %s %s", tt.input, err)
		}
	}
}

func TestWildcardRelationInDefinition(t *testing.T) {
	tests := []struct {
		input       string
		expectError bool
	}{

		{input: `definition user { } definition document { relation reader: user: *} `, expectError: false},
	}

	for _, tt := range tests {
		lexer := NewLexer(tt.input)
		lexer.NextToken()
		z, err := lexer.ReadZSchema()

		mydraw := PlantUMLArchimateSchema{Zdefs: z}
		mydraw.createIDforZdef()

		if tt.expectError && err != nil {
			t.Errorf("expected an error : %s %s", tt.input, err)
		}
	}
}

func TestPermissions1(t *testing.T) {
	tests := []struct {
		input       string
		expectError bool
	}{
		{input: `definition user { }  definition document { relation writer: user }`, expectError: false},
		{input: `definition user { }  definition document { relation writer: user  relation reader: user permission edit = writer }`, expectError: false},
		{input: `definition user { }  definition document { permission view = reader + writer }`, expectError: false},
		{input: `definition user { }  definition document { permission view = reader->view + writer }`, expectError: false},
		{input: `definition user { }  definition document { permission view = (reader) }`, expectError: false},
		{input: `definition user { }  definition document { permission view = ((reader)->writer) }`, expectError: false},
		{input: `definition user { }  definition document { relation writer: user  relation reader: user permission edit = writer permission view = reader + edit }`, expectError: false},
		{input: `definition user { }  definition document { permission view = reader.any(writer) }`, expectError: false},
	}

	for _, tt := range tests {
		lexer := NewLexer(tt.input)
		lexer.NextToken()
		_, err := lexer.ReadZSchema()

		if tt.expectError && err == nil {
			t.Errorf("expected an error but got none for input: %s", tt.input)
		}
		if !tt.expectError && err != nil {
			t.Errorf("did not expect an error but got one for input: %s, error: %v", tt.input, err)
		}

	}

}

func TestPermissions2(t *testing.T) {
	tests := []struct {
		input       string
		expectError bool
	}{
		{input: `definition user { }  definition document { relation writer: user  relation reader: user permission edit = writer permission view = reader + edit }`, expectError: false},
	}

	for _, tt := range tests {
		lexer := NewLexer(tt.input)
		lexer.NextToken()
		z, err := lexer.ReadZSchema()

		mydraw := PlantUMLArchimateSchema{Zdefs: z}
		mydraw.createIDforZdef()

		if tt.expectError && err == nil {
			t.Errorf("expected an error but got none for input: %s", tt.input)
		}
		if !tt.expectError && err != nil {
			t.Errorf("did not expect an error but got one for input: %s, error: %v", tt.input, err)
		}

	}

}

func TestPermissions3(t *testing.T) {
	tests := []struct {
		input       string
		expectError bool
	}{
		{input: `definition user { }  definition document { relation writer: user  relation reader: user permission writer = writer permission view = reader + edit }`, expectError: false},
	}

	for _, tt := range tests {
		lexer := NewLexer(tt.input)
		lexer.NextToken()
		z, err := lexer.ReadZSchema()

		mydraw := PlantUMLArchimateSchema{Zdefs: z}
		mydraw.createIDforZdef()

		if tt.expectError && err == nil {
			t.Errorf("expected an error but got none for input: %s", tt.input)
		}
		if !tt.expectError && err != nil {
			t.Errorf("did not expect an error but got one for input: %s, error: %v", tt.input, err)
		}

	}

}

func TestGeneration1(t *testing.T) {
	tests := []struct {
		input       string
		expectError bool
	}{

		{input: `definition user { } definition group { relation viewer : user } definition document { relation group2: group permission view = group2->viewer } `, expectError: false},
	}

	for _, tt := range tests {
		lexer := NewLexer(tt.input)
		lexer.NextToken()
		z, err := lexer.ReadZSchema()

		mydraw := PlantUMLArchimateSchema{Zdefs: z}
		mydraw.Generate("testG1")

		if tt.expectError && err != nil {
			t.Errorf("expected an error : %s %s", tt.input, err)
		}
	}
}
