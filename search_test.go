package simplexml

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTagSearch(t *testing.T) {
	Convey("Given a simple document with multiple children and varying depths", t, func() {
		root := NewTag("root")

		foo := NewTag("foo")
		foo_0 := NewTag("foo_0")
		foo_0.AddAfter(NewTag("foo_0_0"), nil)
		foo_0.AddAfter(NewTag("foo_0_1"), nil)
		foo_0.AddAfter(NewTag("foo_0_1"), nil)
		foo.AddAfter(foo_0, nil)

		foo_1 := NewTag("foo_1")
		foo_1.AddAfter(NewTag("foo_1_0"), nil)
		foo_1.AddAfter(NewTag("foo_1_1"), nil)
		foo.AddAfter(foo_1, nil)

		root.AddAfter(foo, nil)
		root.AddAfter(NewTag("bar"), nil)

		s := TagSearch{root}

		Convey("ByName(\"Foo\") on root should return nil", func() {
			So(s.ByName("Foo"), ShouldBeNil)

			Convey("One() from this result should return nil", func() {
				So(s.ByName("Foo").One(), ShouldBeNil)
			})
		})

		Convey("ByName(\"foo\") on root should return 1 result", func() {
			So(len(s.ByName("foo")), ShouldEqual, 1)
		})

		Convey("ByName(\"foo\").ByName(\"foo_0\") should return 1 result", func() {
			So(len(s.ByName("foo").ByName("foo_0")), ShouldEqual, 1)
		})

		Convey("ByName(\"foo\").ByName(\"foo_0\").ByName(\"foo_0_1\") should return 2 results", func() {
			So(len(s.ByName("foo").ByName("foo_0").ByName("foo_0_1")), ShouldEqual, 2)

			Convey("One() from this result should return a pointer to a Tag", func() {
				So(s.ByName("foo").ByName("foo_0").ByName("foo_0_1").One(), ShouldHaveSameTypeAs, &Tag{})
			})
		})
	})
}

func TestAttributeSearch(t *testing.T) {
	Convey("Given a new AttributeSearch from a slice of attributes", t, func() {
		as := AttributeSearch{
			&Attribute{Prefix: "ns1", Name: "foo", Value: "fooval"},
			&Attribute{Name: "foo2", Value: "foo2val"},
			&Attribute{Name: "foo2", Value: "foo2val2"},
		}

		Convey("ByName(\"Foo\") should return nil", func() {
			So(as.ByName("Foo"), ShouldBeNil)

			Convey("One() from this result should return nil", func() {
				So(as.ByName("Foo").One(), ShouldBeNil)
			})
		})

		Convey("ByName(\"foo\") should return 1 result", func() {
			So(len(as.ByName("foo")), ShouldEqual, 1)
		})

		Convey("ByName(\"foo2\") should return 2 results", func() {
			So(len(as.ByName("foo2")), ShouldEqual, 2)
		})
	})
}
