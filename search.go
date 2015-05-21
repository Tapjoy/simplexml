package simplexml

// TagSearch is a slice of *Tag
type TagSearch []*Tag

// ByName searches through the children Tags of each element in TagSearch looking for case sensitive matches of Name and returns a new TagSearch of the results. Namespace is ignored.
func (se TagSearch) ByName(s string) TagSearch {
	var r TagSearch

	for _, v := range se {
		for _, v2 := range v.Tags() {
			if v2.Name == s {
				r = append(r, v2)
			}
		}
	}

	return r
}

// One returns the top result off of a TagSearch
func (se TagSearch) One() *Tag {
	if len(se) > 0 {
		return se[0]
	}
	return nil
}