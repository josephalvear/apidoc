// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/xml"
	"testing"

	"github.com/issue9/assert"
)

var (
	_ xml.Unmarshaler = &Tag{}
	_ xml.Unmarshaler = &Server{}
)

func TestTag_UnmarshalXML(t *testing.T) {
	a := assert.New(t)

	obj := &Tag{
		Name:        "tag1",
		Deprecated:  "1.1.1",
		Description: "<a>test</a>",
	}
	str := `<Tag name="tag1" deprecated="1.1.1"><![CDATA[<a>test</a>]]></Tag>`

	data, err := xml.Marshal(obj)
	a.NotError(err).Equal(string(data), str)

	obj1 := &Tag{}
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	a.Equal(obj1, obj)

	// 正常，带 CDATA
	obj1 = &Tag{}
	str = `<Tag name="tag1"><![CDATA[text]]></Tag>`
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	a.Equal(obj1.Description, "text")

	// 正常，不带 CDATA
	obj1 = &Tag{}
	str = `<Tag name="tag1">text</Tag>`
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	a.Equal(obj1.Description, "text")

	// 少 name
	str = `<Tag>test</Tag>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// 少 description
	str = `<Tag name="tag1"></Tag>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// 语法错误
	str = `<Tag name="tag1" deprecated="x.0.1">desc</Tag>`
	a.Error(xml.Unmarshal([]byte(str), obj1))
}

func TestServer_UnmarshalXML(t *testing.T) {
	a := assert.New(t)

	obj := &Server{
		Name:        "srv1",
		URL:         "https://api.example.com/srv1",
		Deprecated:  "1.1.1",
		Description: "<a>test</a>",
	}
	str := `<Server name="srv1" url="https://api.example.com/srv1" deprecated="1.1.1"><![CDATA[<a>test</a>]]></Server>`

	data, err := xml.Marshal(obj)
	a.NotError(err).Equal(string(data), str)

	obj1 := &Server{}
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	a.Equal(obj1, obj)

	// 正常，带 CDATA
	obj1 = &Server{}
	str = `<Server name="tag1" url="https://example.com"><![CDATA[text]]></Server>`
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	a.Equal(obj1.Description, "text")

	// 少 name
	str = `<Server>test</Server>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// 少 description
	str = `<Tag name="tag1"></Tag>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// 少 url
	str = `<Tag name="tag1">test</Tag>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// 语法错误
	str = `<Tag name="tag1" deprecated="x.0.1">desc</Tag>`
	a.Error(xml.Unmarshal([]byte(str), obj1))
}