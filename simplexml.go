// Package simplexml provides a simple API to read, write, edit and search XML documents at run time in pure Go.
package simplexml

import (
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"strings"
)

// Element is the base struct for reading, writing and manipulating XML documents.
type Element struct {
	// Declaration is a string that is prepended to the element when using String().
	Declaration string

	// Name is a xml.Name for the current Element.
	Name xml.Name

	// Attributes is a slice of xml.Attr.
	Attributes []xml.Attr

	// PrettyXML determines whether to use new lines and tabs when true in String().
	PrettyXML bool

	// Children is a slice of Element which represent it's inner contents. There should not be Children and a Value at
	// the same time or String() will panic.
	Children []*Element

	// Value is the string value of the Element. There should not be Children and a Value at
	// the same time String() will panic.
	Value string

	// CDATA wraps the value with CDATA tags in when true String(). CDATA is ignored if the current element has Children.
	CDATA bool

	// Parents is useful to determine how far nested within a Element the element is. Needed for PrettyXML.
	Parents []xml.Name

	// NSPrefixes is a map[string]string with the key being the namespace and the value being it's prefix used for converting full
	// namspaces to it's available prefix during String().
	NSPrefixes NSPrefixes
}

// SetValue sets Element.Value but panics if the element already has children.
func (e *Element) SetValue(s string) *Element {
	if len(e.Children) > 0 {
		panic("tried setting value on an element with Children")
	}

	e.Value = s
	return e
}

// AddAttribute appends an attribute to the given Element. AddAttribute returns it's self for function chaining.
func (e *Element) AddAttribute(attr xml.Attr) *Element {
	// append the attribute
	e.Attributes = append(e.Attributes, attr)

	return e
}

// AddNamespace is a wrapper to AddAttribute in addition to adding the prefix to NSPrefixes for conversion in String()
func (e *Element) AddNamespace(prefix string, namespace string) *Element {
	// init NSPrefixes if needed
	if e.NSPrefixes == nil {
		e.NSPrefixes = make(NSPrefixes)
	}

	// convert empty prefix to something parsable. we do this so that a tag of 'xmlns=foo.local' does not add a prefix to sub elements
	// TODO: this seems hacky, better fix?
	if prefix == "" {
		prefix = "!NIL!"
	}

	e.AddAttribute(xml.Attr{Name: xml.Name{Local: prefix, Space: "xmlns"}, Value: namespace})
	e.NSPrefixes[namespace] = prefix
	return e
}

// AddChild adds a Element to its Children, copying information about the parent down to the child. AddChild panics if the element already has a value.
func (e *Element) AddChild(name xml.Name) *Element {
	if e.Value != "" {
		panic("tried adding child on an element with non empty Value")
	}

	n := &Element{Name: name, PrettyXML: e.PrettyXML, Parents: append(e.Parents, e.Name), NSPrefixes: e.NSPrefixes}
	e.Children = append(e.Children, n)
	return n
}

// RemoveChild takes a pointer to a Element and removes it from the current Element's Children recursively. RemoveChild
// will return an error if the memory address was not found.
func (e *Element) RemoveChild(sxml *Element) error {
	found := false
	for k, v := range e.Children {
		if v == sxml {
			e.Children = append(e.Children[:k], e.Children[k+1:]...)
			found = true
		} else {
			if v.RemoveChild(sxml) == nil {
				found = true
			}
		}
	}
	if !found {
		return errors.New("address not found")
	}
	return nil
}

// AllChildren returns a single slice of Element of it's children at any depth
func (e *Element) AllChildren() []*Element {
	res := e.Children

	for _, v := range e.Children {
		res = append(res, v.AllChildren()...)
	}

	return res
}

// SetPrettyXML sets the PrettyXML indicator for the Element and all of it's children recursively. SetPrettyXML returns
// its self for function chaining.
func (e *Element) SetPrettyXML(b bool) *Element {
	e.PrettyXML = b
	for _, v := range e.AllChildren() {
		v.PrettyXML = b
	}
	return e
}

// String prepares the xml document like xml.Marshal.
func (e Element) String() string {
	var s string
	var prefix string
	var suffix string

	// panic on bad Element
	if e.Value != "" && len(e.Children) > 0 {
		panic("have both a non empty Value and Children")
	}

	if e.PrettyXML {
		for i := 0; i < len(e.Parents); i++ {
			prefix = prefix + "\t"
		}
		suffix = "\n"
	}

	// prepare the tag to be used
	startTag := e.Name.Local
	if e.Name.Space != "" && e.NSPrefixes.GetPrefix(e.Name.Space) != "" && e.NSPrefixes.GetPrefix(e.Name.Space) != "!NIL!" {
		startTag = fmt.Sprintf("%v:%v", e.NSPrefixes.GetPrefix(e.Name.Space), e.Name.Local)
	}
	endTag := startTag

	// add attributes to tag
	if len(e.Attributes) > 0 {
		for _, v := range e.Attributes {
			if v.Name.Space != "" {
				startTag = fmt.Sprintf("%s %s:%s=\"%s\"", startTag, v.Name.Space, v.Name.Local, v.Value)
			} else {
				startTag = fmt.Sprintf("%s %s=\"%s\"", startTag, v.Name.Local, v.Value)
			}
		}
	}

	if len(e.Children) > 0 {
		s = s + suffix
		for _, v := range e.Children {
			s = s + v.String()
		}
		s = fmt.Sprintf("%s<%s>%s%s</%s>%s", prefix, startTag, s, prefix, endTag, suffix)
	} else {
		s = html.EscapeString(e.Value)
		if e.CDATA {
			s = fmt.Sprintf("<![CDATA[%s]]>", s)
		}

		s = fmt.Sprintf("%s<%s>%s</%s>%s", prefix, startTag, s, endTag, suffix)
	}

	// add declaration
	s = e.Declaration + s

	return s
}

// Search returns a new Search from the current Element. This is useful for function chaining.
func (e *Element) Search() Search {
	return Search{e}
}

// XPath returns the xpath representation of the current element from it's root element.
func (e Element) XPath() string {
	str := []string{}

	// add parents to path
	for _, v := range e.Parents {
		if v.Space != "" {
			str = append(str, fmt.Sprintf("%s:%s", v.Space, v.Local))
		} else {
			str = append(str, v.Local)
		}
	}

	// add self to path
	if e.Name.Space != "" {
		str = append(str, fmt.Sprintf("%s:%s", e.Name.Space, e.Name.Local))
	} else {
		str = append(str, e.Name.Local)
	}

	return strings.Join(str, "/")
}

// NSPrefixes is a map[string]string with the key being the namespace and the value being it's prefix used for converting full
// namspaces to it's available prefix during String().
type NSPrefixes map[string]string

// GetPrefix returns the prefix used for a namespace
func (p NSPrefixes) GetPrefix(ns string) string {
	return p[ns]
}

// New returns a blank Element with the given element name. Declaration is set to
// '<?xml version=\"1.0\" encoding=\"UTF-8\"?>' by New()
func New(name xml.Name) *Element {
	return &Element{Name: name, Declaration: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n", NSPrefixes: make(NSPrefixes)}
}

// NewFromReader creates a new Element from an existing XML document
func NewFromReader(r io.Reader) (*Element, error) {
	// TODO: CDATA from reader

	var root *Element
	var declaration string
	tree := []*Element{}

	// start up a new xml decoder
	d := xml.NewDecoder(r)
	var start xml.StartElement

	// iterate through the tokens
	for {
		tok, _ := d.Token()
		if tok == nil {
			break
		}

		// switch on token type
		switch t := tok.(type) {
		case xml.StartElement:
			start = t.Copy()
			var n *Element

			if len(tree) == 0 {
				// create out first element and set it as root
				n = New(start.Name)
				root = n
			} else {
				// create a new child from last element in tree
				n = tree[len(tree)-1].AddChild(start.Name)
			}

			// set attributes and new namespaces
			n.Attributes = start.Attr
			for _, v := range n.Attributes {
				if strings.ToLower(v.Name.Space) == "xmlns" { // normal namespace, add it to NSPrefixes
					n.NSPrefixes[v.Value] = v.Name.Local
				} else if strings.ToLower(v.Name.Local) == "xmlns" { // namespace without a prefix, add it to NSPrefixes so we can cleaqr the prefix later
					// TODO: this is hacky, find a better way
					n.NSPrefixes[v.Value] = "!NIL!"
				}
			}

			// add element to tree
			tree = append(tree, n)
		case xml.EndElement:
			// done with the element, drop it from tree and reset start token
			tree = tree[:len(tree)-1]
			start = xml.StartElement{}
		case xml.CharData:
			// assign the value to last element in tree if not whitespace
			if start.Name.Local != "" {
				tree[len(tree)-1].Value = strings.TrimSpace(string(t))
			}
		case xml.ProcInst:
			declaration = fmt.Sprintf("<?%s %s?>\n", t.Target, string(t.Inst))
		default:
			// eat line
		}
	}

	// we should be back down to the root element
	if len(tree) != 0 {
		// TODO: position of failure
		return nil, errors.New("malformed document")
	}

	// set the declaration
	root.Declaration = declaration

	return root, nil
}

// Search is a simplexml.Element that has search capabilities
type Search []*Element

// MatchName returns a new Search where the supplied xml.Name matches the Search.Children()
func (sxmls Search) MatchName(name xml.Name) Search {
	r := Search{}

	// search the top level elements
	for _, v := range sxmls {
		if v.Name.Local == name.Local && v.Name.Space == name.Space {
			r = append(r, v)
		}
	}

	return r
}

// MatchNameDeep returns a new Search where the supplied xml.Name matches the Search.AllChildren()
func (sxmls Search) MatchNameDeep(name xml.Name) Search {
	r := Search{}

	for _, v := range sxmls {
		// search the top level elements
		if v.Name.Local == name.Local && v.Name.Space == name.Space {
			r = append(r, v)
		}

		// search all its children
		for _, v2 := range v.AllChildren() {
			if v2.Name.Local == name.Local && v2.Name.Space == name.Space {
				r = append(r, v2)
			}

		}
	}

	return r
}

// MatchAttr returns a new Search where the supplied xml.Attr matches the Search.Children()
func (sxmls Search) MatchAttr(attr xml.Attr) Search {
	r := Search{}

	// search the top level elements
	for _, v := range sxmls {
		for _, a := range v.Attributes {
			if a.Value == attr.Value && a.Name.Local == attr.Name.Local && a.Name.Space == attr.Name.Space {
				r = append(r, v)
				break
			}
		}
	}

	return r
}

// MatchAttrDeep returns a new Search where the supplied xml.Attr matches the Search.AllChildren()
func (sxmls Search) MatchAttrDeep(attr xml.Attr) Search {
	r := Search{}

	// search the top level elements
	for _, v := range sxmls {
		for _, a := range v.Attributes {
			if a.Value == attr.Value && a.Name.Local == attr.Name.Local && a.Name.Space == attr.Name.Space {
				r = append(r, v)
				break
			}
		}

		// search all its children
		for _, v2 := range v.AllChildren() {
			for _, a := range v2.Attributes {
				if a.Value == attr.Value && a.Name.Local == attr.Name.Local && a.Name.Space == attr.Name.Space {
					r = append(r, v2)
					break
				}
			}
		}
	}

	return r
}

// One returns the top result off of a Search
func (sxmls Search) One() *Element {
	if len(sxmls) > 0 {
		return sxmls[0]
	}
	return nil
}
