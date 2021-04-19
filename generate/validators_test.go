package generate

import "testing"

var (
	gen = Generate{
		Models: []Model{
			{Name: "string", Fields: []Field{{Type: "string"}}},
			{Name: "uint", Fields: []Field{{Type: "uint"}}},
			{Name: "Address", Fields: []Field{{Type: "string", Name: "AddressLine1"}}},
			{Name: "Contact", Fields: []Field{{Type: "Address", Name: "[]AddressLine"}}},
			{Name: "Person", Fields: []Field{{Type: "Contact", Name: "ContactInfo"}}},
		},
	}
)

// BenchmarkValidateTypes just to check as of now
func BenchmarkValidateTypes(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := gen.ValidateTypes(); err != nil {
			b.Error("Unexpected result:", err)
		}
	}
}

func BenchmarkSetFieldCatebgory(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := gen.SetFieldCategory(); err != nil {
			b.Error("Unexpected result:", err)
		}
	}

}
