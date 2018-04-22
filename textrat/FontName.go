package textrat

import (
	"fmt"
)

type FontName struct {
	Family, Name string
	FullName     string
}

func (s FontName) String() string {
	return fmt.Sprintf("VectorFont(%s, family : %s, name : %s)", s.FullName, s.Family, s.Name)
}
