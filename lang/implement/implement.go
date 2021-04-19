package implement

import (
	"fmt"
)

// Implement is a struct to carry methods for lang specific implementations
type Implement struct{}

// New is to create new instance of Implement type
func New() *Implement {
	return &Implement{}
}

// IsValidIdentifier is to check whether the field is a valid identifier or not
/*func (i *Implement) IsValidIdentifier(fielden string) bool {

	// Should not start with the number or any special chars other than _
	// should not contain secial chars other than _
	// should have atleast one char and can have n number of digits

	// according to https://www.geeksforgeeks.org/check-whether-the-given-string-is-a-valid-identifier/
	// It must start with either underscore(_) or any of the characters from the ranges [‘a’, ‘z’] and [‘A’, ‘Z’].
	// There must not be any white space in the string.
	//  And, all the subsequent characters after the first character must not consist of any special characters like $, #, % etc.

	// This tests whether a pattern matches a string.
	match, err := regexp.MatchString(`^[^\d\W]\w*$`, fielden)
	if err != nil {
		return false
	}
	return match
}*/
// IsValidIdentifier is to check whether the field is a valid identifier or not
// It must start with either underscore(_) or any of the characters from the ranges [‘a’, ‘z’] and [‘A’, ‘Z’].
// There must not be any white space in the string.
// And, all the subsequent characters after the first character must not consist of any special characters like $, #, % etc.
func (i *Implement) IsValidIdentifier(iden string) bool {
	if !((iden[0] >= 65 && iden[0] <= 90) || (iden[0] >= 97 && iden[0] <= 122) || string(iden[0]) == "_") {
		return false
	}
	for i := 1; i < len(iden); i++ {
		if !((iden[i] >= 65 && iden[i] <= 90) || (iden[i] >= 97 && iden[i] <= 122) || (iden[i] >= 48 && iden[i] <= 57) || string(iden[i]) == "_") {
			return false
		}
	}
	return true
}

// GetFuncReturnType is to give the return type of the function
func (i *Implement) GetFuncReturnType(x interface{}) string {
	return fmt.Sprintf("%T", x)
}
