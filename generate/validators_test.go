package generate

import "testing"

// BenchmarkValidateTypes just to check as of now
func BenchmarkValidateTypes(b *testing.B) {

	gen := Generate{
		Models: []Model{
			{Name: "string", Fields: []Field{{Type: "string"}}},
			{Name: "uint", Fields: []Field{{Type: "uint"}}},
			//{"float", []Field{{Type: "float"}}},
			//{"float32", []Field{{Type: "float32"}}},
			//{"float64", []Field{{Type: "float64"}}},
		},
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := gen.ValidateTypes(); err != nil {
			b.Error("Unexpected result:", err)
		}
	}
}
