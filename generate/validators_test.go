package generate

import (
	"readyGo/box"
	"readyGo/lang/implement"
	"readyGo/scalar"
	"testing"
)

var ()

// BenchmarkValidateTypes just to check as of now
/*func BenchmarkValidateTypes(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := gen.ValidateTypes(); err != nil {
			b.Error("Unexpected result:", err)
		}
	}
}*/
func LoadGen() *Generate {
	gen := &Generate{
		Models: []Model{
			{Name: "Sample", Fields: []Field{{Type: "string", Name: "name"}}},
			{Name: "Demo", Fields: []Field{{Type: "uint", Name: "age"}}},
			{Name: "Address", Fields: []Field{{Type: "string", Name: "AddressLine"}}},
			{Name: "Contact", Fields: []Field{{Type: "[]Address", Name: "Addresses"}}},
			{Name: "Person", Fields: []Field{{Type: "Contact", Name: "ContactInfo"}}},
		},
	}
	ops := &box.Box{}
	scalar, _ := scalar.New(ops, "configs/scalars.json")
	gen.Scalars = scalar
	imlementer := implement.New()
	gen.Implementer = imlementer
	return gen
}
func BenchmarkSetFieldCatebgory(b *testing.B) {
	gen := LoadGen()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := gen.SetFieldCategory(); err != nil {
			b.Error("Unexpected result:", err)
		}
	}
}

func BenchmarkChangeIden(b *testing.B) {
	gen := LoadGen()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := gen.ChangeIden(); err != nil {
			b.Error("Unexpected result:", err)
		}
	}
}
