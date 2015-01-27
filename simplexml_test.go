package simplexml

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"

	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

const xmlstring = `<?xml version="1.0" encoding="utf-8"?>
<Job xmlns="foo.local">
	<type>Bulk</type>
	<id>abc123</id>
	<result>
		<id>123abc</id>
		<name>foo</name>
	</result>
	<result>
		<id>124abc</id>
		<name>bar</name>
	</result>
</Job>
`

const badxmlstring = `<?xml version="1.0" encoding="utf-8"?>
<Job xmlns="foo.local">
	<type>Bulk</type>
	<id>abc123</id>
	<result>
		<id>123abc</id>
		<name>foo</name>
	</result>
	<result>
		<id>124abc</id>
		<name>bar</name>
	<!--</result>-->
</Job>
`

const soapstring = `<?xml version="1.0" encoding="utf-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:urn="urn:foo.bar">
	<soapenv:Body>
		<urn:Job>
			<type>Bulk</type>
			<id>abc123</id>
			<result>
				<id>123abc</id>
				<name>foo</name>
			</result>
			<result>
				<id>124abc</id>
				<name>bar</name>
			</result>
		</urn:Job>
	</soapenv:Body>
</soapenv:Envelope>
`

func TestValues(t *testing.T) {
	Convey("Given a Element of 'Foo' and space of 'foo.local' prefix 'n1'", t, func() {
		sxml := &Element{Name: xml.Name{Local: "Foo", Space: "foo.local"}}
		sxml.AddNamespace("n1", "foo.local")
		sxml.Declaration = ""

		Convey("Given a string value of 'bar'", func() {
			sxml.SetValue("bar")

			Convey("AddChild() should panic", func() {
				So(func() { sxml.AddChild(xml.Name{Local: "baz"}) }, ShouldPanic)
			})

			Convey("String() of the parent should equal '<n1:Foo xmlns:n1=\"foo.local\">bar</n1:Foo>'", func() {
				So(sxml.String(), ShouldEqual, "<n1:Foo xmlns:n1=\"foo.local\">bar</n1:Foo>")
			})

			Convey("Given Foo CDATA is true", func() {
				sxml.CDATA = true

				Convey("String() of the parent should equal '<n1:Foo xmlns:n1=\"foo.local\"><![CDATA[bar]]></n1:Foo>'", func() {
					So(sxml.String(), ShouldEqual, "<n1:Foo xmlns:n1=\"foo.local\"><![CDATA[bar]]></n1:Foo>")
				})
			})

			Convey("Foo.XPath() should equal 'foo.local:Foo'", func() {
				So(sxml.XPath(), ShouldEqual, "foo.local:Foo")
			})
		})
	})
}

func TestChildren(t *testing.T) {
	Convey("Given a Element of 'Foo' and space of 'foo.local'", t, func() {
		sxml := New(xml.Name{Local: "Foo", Space: "foo.local"})
		sxml.AddNamespace("n1", "foo.local")
		sxml.Declaration = ""

		Convey("Given a child element of 'Bar'", func() {
			sxml.AddChild(xml.Name{Local: "Bar"})

			Convey("SetValue() should panic on the parent", func() {
				So(func() { sxml.SetValue("baz") }, ShouldPanic)
			})

			Convey("Given a manual set of Value", func() {
				sxml.Value = "bar"

				Convey("String() should panic", func() {
					So(func() { sxml.String() }, ShouldPanic)
				})
			})

			Convey("Given a value of 'baz' for the new child element", func() {
				sxml.Children[0].SetValue("baz")

				Convey("String() of the parent should equal '<n1:Foo xmlns:n1=\"foo.local\"><Bar>baz</Bar></n1:Foo>'", func() {
					So(sxml.String(), ShouldEqual, "<n1:Foo xmlns:n1=\"foo.local\"><Bar>baz</Bar></n1:Foo>")
				})

				Convey("Given PrettyXML is true for all elements", func() {
					sxml.PrettyXML = true
					for _, v := range sxml.AllChildren() {
						v.PrettyXML = true
					}

					Convey("String() of the parent should equal '<n1:Foo xmlns:n1=\"foo.local\">\n\t<Bar>baz</Bar>\n</n1:Foo>\n'", func() {
						So(sxml.String(), ShouldEqual, "<n1:Foo xmlns:n1=\"foo.local\">\n\t<Bar>baz</Bar>\n</n1:Foo>\n")
					})
				})

			})

			Convey("Bar.Parents should be a length of 1", func() {
				So(len(sxml.Children[0].Parents), ShouldEqual, 1)
			})

			Convey("Bar.XPath() should equal 'foo.local:Foo/Bar'", func() {
				So(sxml.Children[0].XPath(), ShouldEqual, "foo.local:Foo/Bar")
			})
		})
	})
}

func TestAttr(t *testing.T) {
	Convey("Given an empty Element of 'Foo'", t, func() {
		sxml := New(xml.Name{Local: "Foo"})
		sxml.Declaration = ""

		Convey("Given an attribute of type and value of 'bar'", func() {
			sxml.Attributes = append(sxml.Attributes, xml.Attr{Name: xml.Name{Local: "type"}, Value: "bar"})

			Convey("String() should equal '<Foo type=\"bar\"></Foo>'", func() {
				So(sxml.String(), ShouldEqual, "<Foo type=\"bar\"></Foo>")
			})

			Convey("Given a second attribute of type2, value of 'baz' and space of 'ns1", func() {
				sxml.Attributes = append(sxml.Attributes, xml.Attr{Name: xml.Name{Local: "type2", Space: "ns1"}, Value: "baz"})

				Convey("String() should equal '<Foo type=\"bar\" ns1:type2=\"baz\"></Foo>'", func() {
					So(sxml.String(), ShouldEqual, "<Foo type=\"bar\" ns1:type2=\"baz\"></Foo>")
				})
			})
		})
	})
}

func TestNewFromReaderXML(t *testing.T) {
	Convey("Given an xml string in NewFromReader()", t, func() {
		sxml, err := NewFromReader(strings.NewReader(xmlstring))

		Convey("err should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("String() should equal the original document", func() {
			sxml.PrettyXML = true
			for _, v := range sxml.AllChildren() {
				v.PrettyXML = true
			}

			So(sxml.String(), ShouldEqual, xmlstring)
		})
	})
}

func TestNewFromReaderSoap(t *testing.T) {
	Convey("Given a soap string in NewFromReader()", t, func() {
		sxml, err := NewFromReader(strings.NewReader(soapstring))

		Convey("err should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("String() should equal the original document", func() {
			sxml.PrettyXML = true
			for _, v := range sxml.AllChildren() {
				v.PrettyXML = true
			}

			So(sxml.String(), ShouldEqual, soapstring)
		})
	})
}

func TestNewFromReaderBad(t *testing.T) {
	Convey("Given a bad string in NewFromReader()", t, func() {
		sxml, err := NewFromReader(strings.NewReader(badxmlstring))

		Convey("err.Error() should equal 'malformed document'", func() {
			So(err.Error(), ShouldEqual, "malformed document")
		})

		Convey("The Element should be nil'", func() {
			So(sxml, ShouldBeNil)
		})
	})
}

func TestMatchName(t *testing.T) {
	Convey("Given a new element of Foo", t, func() {
		sxml := New(xml.Name{Local: "Foo"})
		s := Search{sxml}

		Convey("MatchName for Foo should return 1 result", func() {
			res := s.MatchName(xml.Name{Local: "Foo"})

			So(len(res), ShouldEqual, 1)
		})

		Convey("Given two children of Bar to Foo", func() {
			sxml.AddChild(xml.Name{Local: "Bar"})
			sxml.AddChild(xml.Name{Local: "Bar"})

			Convey("MatchName for Bar should return 0 results", func() {
				So(len(s.MatchName(xml.Name{Local: "Bar"})), ShouldEqual, 0)
			})

			Convey("MatchNameDeep for Bar should return 2 results", func() {
				So(len(s.MatchNameDeep(xml.Name{Local: "Bar"})), ShouldEqual, 2)
			})

			Convey("Given another child of Bar to one of the Bar elements", func() {
				sxml.Children[0].AddChild(xml.Name{Local: "Bar"})

				Convey("MatchNameDeep for Bar should return 3 results", func() {
					So(len(s.MatchNameDeep(xml.Name{Local: "Bar"})), ShouldEqual, 3)
				})
			})
		})

		Convey("MatchName for Foo2 should return 0 results", func() {
			res := s.MatchName(xml.Name{Local: "Foo2"})

			So(len(res), ShouldEqual, 0)
		})
	})
}

func TestMatchAttr(t *testing.T) {
	Convey("Given a new element of Foo", t, func() {
		sxml := New(xml.Name{Local: "Foo"})
		s := Search{sxml}

		Convey("And given a bar attribute equal to baz", func() {
			sxml.Attributes = append(sxml.Attributes, xml.Attr{Name: xml.Name{Local: "bar"}, Value: "baz"})

			Convey("MatchAttr() should return 1 result", func() {
				So(len(s.MatchAttr(xml.Attr{Name: xml.Name{Local: "bar"}, Value: "baz"})), ShouldEqual, 1)
			})

			Convey("MatchAttr() for a different value should return 0 results", func() {
				So(len(s.MatchAttr(xml.Attr{Name: xml.Name{Local: "bar"}, Value: "notbaz"})), ShouldEqual, 0)
			})

			Convey("Given a child element of Bar to Foo with the same attribute", func() {
				sxml.AddChild(xml.Name{Local: "Bar"})
				sxml.Children[0].Attributes = append(sxml.Children[0].Attributes, xml.Attr{Name: xml.Name{Local: "bar"}, Value: "baz"})

				Convey("MatchAttr() should return 1 result", func() {
					So(len(s.MatchAttr(xml.Attr{Name: xml.Name{Local: "bar"}, Value: "baz"})), ShouldEqual, 1)
				})

				Convey("MatchAttrDeep() should return 2 results", func() {
					So(len(s.MatchAttrDeep(xml.Attr{Name: xml.Name{Local: "bar"}, Value: "baz"})), ShouldEqual, 2)
				})
			})
		})
	})
}

func ExampleNew() {
	// create a simplexml element
	catalog := New(xml.Name{Local: "Catalog"}).AddNamespace("b", "api.books.localhost")
	catalog.SetPrettyXML(true)

	// add books
	for i := 0; i < 3; i++ {
		book := catalog.AddChild(xml.Name{Local: "books", Space: "api.books.localhost"}).AddAttribute(xml.Attr{
			Name:  xml.Name{Local: "id"},
			Value: strconv.Itoa(i),
		})
		book.AddChild(xml.Name{Local: "name"}).Value = fmt.Sprintf("Book Title %v", strconv.Itoa(i))
	}

	fmt.Println(catalog.String())
	//Output:
	//<?xml version="1.0" encoding="UTF-8"?>
	//<Catalog xmlns:b="api.books.localhost">
	//	<b:books id="0">
	//		<name>Book Title 0</name>
	//	</b:books>
	//	<b:books id="1">
	//		<name>Book Title 1</name>
	//	</b:books>
	//	<b:books id="2">
	//		<name>Book Title 2</name>
	//	</b:books>
	//</Catalog>
}

func ExampleNewFromReader() {
	s := `<?xml version="1.0" encoding="UTF-8"?>
<Catalog xmlns:b="api.books.localhost">
	<b:done>true</b:done>
	<b:books id="0">
		<name>Book Title 0</name>
	</b:books>
	<b:books id="1">
		<name>Book Title 1</name>
	</b:books>
	<b:books id="2">
		<name>Book Title 2</name>
	</b:books>
</Catalog>`

	// create simplexml from reader
	sxml, err := NewFromReader(strings.NewReader(s))
	if err != nil {
		panic(err)
	}

	// find book elements
	books := sxml.Search().MatchNameDeep(xml.Name{Space: "api.books.localhost", Local: "books"})

	// match id from books, add type
	books.MatchAttr(
		xml.Attr{
			Name:  xml.Name{Local: "id"},
			Value: "1",
		}).One().AddChild(xml.Name{Local: "type"}).Value = "Fiction"

	// match id from books, remove the element
	sxml.RemoveChild(books.MatchAttr(
		xml.Attr{
			Name:  xml.Name{Local: "id"},
			Value: "2",
		}).One())

	// set pretty flag on all elements
	sxml.SetPrettyXML(true)

	fmt.Println(sxml.String())
	//Output:
	//<?xml version="1.0" encoding="UTF-8"?>
	//<Catalog xmlns:b="api.books.localhost">
	//	<b:done>true</b:done>
	//	<b:books id="0">
	//		<name>Book Title 0</name>
	//	</b:books>
	//	<b:books id="1">
	//		<name>Book Title 1</name>
	//		<type>Fiction</type>
	//	</b:books>
	//</Catalog>
}
