package simplexml

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"

	"encoding/xml"
	"strings"

	"fmt"
	"strconv"
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
			sxml.Value = "bar"

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

			Convey("Given a manual set of Value", func() {
				sxml.Value = "bar"

				Convey("String() should panic", func() {
					So(func() { sxml.String() }, ShouldPanic)
				})
			})

			Convey("Given a value of 'baz' for the new child element", func() {
				sxml.Children[0].Value = "baz"

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

func TestSearch(t *testing.T) {
	Convey("Given a new Element 'Foo' and a Search of Foo", t, func() {
		foo := New(xml.Name{Local: "Foo"})
		searchFoo := foo.Search()

		Convey("MatchParentName for Foo should return 1 result", func() {
			So(len(searchFoo.MatchParentName(xml.Name{Local: "Foo"})), ShouldEqual, 1)
		})

		Convey("MatchChildName for Foo should return 0 results", func() {
			So(len(searchFoo.MatchChildName(xml.Name{Local: "Foo"})), ShouldEqual, 0)
		})

		Convey("Given an attribute of type=foo", func() {
			foo.AddAttribute(xml.Attr{Name: xml.Name{Local: "type"}, Value: "foo"})

			Convey("MatchParentAttr should return 1 result", func() {
				So(len(searchFoo.MatchParentAttr(xml.Attr{Name: xml.Name{Local: "type"}, Value: "foo"})), ShouldEqual, 1)
			})
		})

		Convey("Given a child element of 'Bar'", func() {
			bar := foo.AddChild(xml.Name{Local: "Bar"})

			Convey("MatchParentName(Foo).MatchChildName(Bar) from Foo should return 1 result", func() {
				So(len(searchFoo.MatchParentName(xml.Name{Local: "Foo"}).MatchChildName(xml.Name{Local: "Bar"})), ShouldEqual, 1)
			})

			Convey("Given an attribute of type=bar", func() {
				bar.AddAttribute(xml.Attr{Name: xml.Name{Local: "type"}, Value: "bar"})

				Convey("MatchParentAttr should return 0 results", func() {
					So(len(searchFoo.MatchParentAttr(xml.Attr{Name: xml.Name{Local: "type"}, Value: "bar"})), ShouldEqual, 0)
				})

				Convey("MatchChildAttr should return 1 result", func() {
					So(len(searchFoo.MatchChildAttr(xml.Attr{Name: xml.Name{Local: "type"}, Value: "bar"})), ShouldEqual, 1)
				})
			})

			Convey("Given two children of 'Baz'", func() {
				baz1 := bar.AddChild(xml.Name{Local: "Baz"})
				baz2 := bar.AddChild(xml.Name{Local: "Baz"})

				Convey("Given an attribute of type=baz", func() {
					baz1.AddAttribute(xml.Attr{Name: xml.Name{Local: "type"}, Value: "baz"})
					baz2.AddAttribute(xml.Attr{Name: xml.Name{Local: "type"}, Value: "baz"})

					Convey("MatchParentAttr should return 0 results", func() {
						So(len(searchFoo.MatchParentAttr(xml.Attr{Name: xml.Name{Local: "type"}, Value: "baz"})), ShouldEqual, 0)
					})

					Convey("MatchChildAttr should return 0 results", func() {
						So(len(searchFoo.MatchChildAttr(xml.Attr{Name: xml.Name{Local: "type"}, Value: "baz"})), ShouldEqual, 0)
					})
				})

				Convey("MatchChildNameDeep from Foo should return 2 results", func() {
					So(len(searchFoo.MatchChildNameDeep(xml.Name{Local: "Baz"})), ShouldEqual, 2)
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
	//	</b:books>
	//	<b:books id="2">
	//		<name>Book Title 2</name>
	//	</b:books>
	//</Catalog>
}

func ExampleMatch() {
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

	// get books with id of 1
	fmt.Println(len(sxml.Search().MatchParentName(
		xml.Name{
			Local: "Catalog",
		},
	).MatchChildName(
		xml.Name{
			Local: "books",
			Space: "api.books.localhost",
		},
	).MatchParentAttr(
		xml.Attr{
			Name:  xml.Name{Local: "id"},
			Value: "2",
		})))

	// or more quickly

	fmt.Println(len(sxml.Search().MatchChildAttrDeep(
		xml.Attr{
			Name:  xml.Name{Local: "id"},
			Value: "2",
		})))

	//Output:
	//1
	//1
}
