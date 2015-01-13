simplexml
=========

simplexml provides a simple API to read, create and manipulate XML documents at run time in pure Go.

[Documentation on GoDoc](https://godoc.org/github.com/kylewolfe/simplexml)

[![Build Status](https://travis-ci.org/kylewolfe/simplexml.svg?branch=master)](https://travis-ci.org/kylewolfe/simplexml)

**This package is currently in alpha and subject to change.**

## Roadmap

- ☑ Basic xml creation, updates
- ☑ Basic xml search
- ☑ New SimpleXMLElement from io.Reader
- ☑ Reuse portions of encoding/xml when possible
- ☑ CDATA support (through API)
- ☐ CDATA support from an io.Reader
- ☑ Pretty XML (formatted with new lines and tabs)
- ☐ Comment support
- ☐ Unmarshal support with used/unused key results (can still use encoding/xml.Unmarshal with SimpleXMLElement.String())
- ☐ xpath search

## Usage

### From scratch

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
```
### From a reader

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
```

```xml
<?xml version="1.0" encoding="UTF-8"?>
<Catalog xmlns:b="api.books.localhost">
	<b:done>true</b:done>
	<b:books id="0">
		<name>Book Title 0</name>
	</b:books>
	<b:books id="1">
		<name>Book Title 1</name>
		<type>Fiction</type>
	</b:books>
</Catalog>
```
