simplexml [![Build Status](https://travis-ci.org/kylewolfe/simplexml.svg?branch=master)](https://travis-ci.org/kylewolfe/simplexml) [![Coverage Status](https://coveralls.io/repos/kylewolfe/simplexml/badge.svg)](https://coveralls.io/r/kylewolfe/simplexml) [![GoDoc](http://godoc.org/github.com/kylewolfe/simplexml?status.svg)](http://godoc.org/github.com/kylewolfe/simplexml) 
=========

**This package is currently in alpha and subject to change.**

simplexml provides a simple API to read, create and manipulate XML documents at run time in pure Go.

## Roadmap

- ☐ CDATA support from an io.Reader
- ☐ Comment support
- ☐ Unmarshal support with used/unused key results (can still use encoding/xml.Unmarshal with Element.String())
- ☐ xpath search
- ☑ Basic xml creation, updates
- ☑ Basic xml search
- ☑ New Element from io.Reader
- ☑ Reuse portions of encoding/xml when possible
- ☑ CDATA support (through API)
- ☑ Pretty XML (formatted with new lines and tabs)

## Usage

### From Scratch

```go
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
```

#### Output
```xml
<?xml version="1.0" encoding="UTF-8"?>
<Catalog xmlns:b="api.books.localhost">
	<b:books id="0">
		<name>Book Title 0</name>
	</b:books>
	<b:books id="1">
		<name>Book Title 1</name>
	</b:books>
	<b:books id="2">
		<name>Book Title 2</name>
	</b:books>
</Catalog>
}
```
### From a Reader

```go
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
```

#### Output
```xml
<?xml version="1.0" encoding="UTF-8"?>
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
</Catalog>
```
### Searching

```go
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
```