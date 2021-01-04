package implement

import (
	"fmt"
	"regexp"
)

// Implement is a struct to carry methods for lang specific implementations
type Implement struct{}

// New is to create new instance of Implement type
func New() *Implement {
	return &Implement{}
}

// IsValidIdentifier is to check whether the field is a valid identifier or not
func (i *Implement) IsValidIdentifier(fielden string) bool {

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
}

// GetFuncReturnType is to give the return type of the function
func (i *Implement) GetFuncReturnType(x interface{}) string {
	return fmt.Sprintf("%T", x)
}
