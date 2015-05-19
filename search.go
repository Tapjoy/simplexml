package simplexml

// Search is a slice of *Tag
type Search []*Tag

// ByName returns a Search of Tags that have a case sesnsitive match on Tag name alone, ignoring namespace.
func (se Search) ByName(s string) Search {
	var r Search

	for _, v := range se {
		if v.Name == s {
			r = append(r, v)
		}
	}

	return r
}

// One returns the top result off of a Search
func (se Search) One() *Tag {
	if len(se) > 0 {
		return se[0]
	}
	return nil
}
