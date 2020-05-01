// SPDX-License-Identifier: MIT

package ast

import (
	"io/ioutil"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/token"
)

func TestGetTagName(t *testing.T) {
	a := assert.New(t)

	p, err := token.NewParser(core.Block{Data: []byte("  <root>xx</root>")})
	a.NotError(err).NotNil(p)
	root, err := getTagName(p)
	a.NotError(err).Equal(root, "root")

	p, err = token.NewParser(core.Block{Data: []byte("<!-- xx -->  <root>xx</root>")})
	a.NotError(err).NotNil(p)
	root, err = getTagName(p)
	a.NotError(err).Equal(root, "root")

	p, err = token.NewParser(core.Block{Data: []byte("<!-- xx -->   <root>xx</root>")})
	a.NotError(err).NotNil(p)
	root, err = getTagName(p)
	a.NotError(err).Equal(root, "root")

	// 无效格式
	p, err = token.NewParser(core.Block{Data: []byte("<!-- xx   <root>xx</root>")})
	a.NotError(err).NotNil(p)
	root, err = getTagName(p)
	a.Error(err).Equal(root, "")

	// 无效格式
	p, err = token.NewParser(core.Block{Data: []byte("</root>")})
	a.NotError(err).NotNil(p)
	root, err = getTagName(p)
	a.Equal(err, ErrNoDocFormat).Equal(root, "")

	// io.EOF
	p, err = token.NewParser(core.Block{Data: []byte("<!-- xx -->")})
	a.NotError(err).NotNil(p)
	root, err = getTagName(p)
	a.Equal(err, ErrNoDocFormat).Equal(root, "")
}

func loadAPIDoc(a *assert.Assertion) *APIDoc {
	data, err := ioutil.ReadFile("./testdata/doc.xml")
	a.NotError(err).NotNil(data)

	doc := &APIDoc{}
	a.NotNil(doc)

	a.NotError(doc.Parse(core.Block{
		Location: core.Location{
			URI:   "doc.xml",
			Range: core.Range{},
		},
		Data: data,
	}))

	return doc
}

func TestAPIDoc(t *testing.T) {
	a := assert.New(t)
	doc := loadAPIDoc(a)

	a.Equal(doc.Version, &VersionAttribute{
		Base: token.Base{
			UsageKey: "usage-apidoc-version",
			Range: core.Range{
				Start: core.Position{Character: 8, Line: 2},
				End:   core.Position{Character: 23, Line: 2},
			},
			XMLName: String{
				Range: core.Range{
					Start: core.Position{Character: 8, Line: 2},
					End:   core.Position{Character: 15, Line: 2},
				},
				Value: "version",
			},
		},

		Value: String{
			Range: core.Range{
				Start: core.Position{Character: 17, Line: 2},
				End:   core.Position{Character: 22, Line: 2},
			},
			Value: "1.1.1",
		},
	})

	a.Equal(len(doc.Tags), 2)
	tag := &Tag{
		Base: token.Base{
			UsageKey: "usage-apidoc-tags",
			Range: core.Range{
				Start: core.Position{Character: 4, Line: 10},
				End:   core.Position{Character: 47, Line: 10},
			},
			XMLName: String{
				Range: core.Range{
					Start: core.Position{Character: 5, Line: 10},
					End:   core.Position{Character: 8, Line: 10},
				},
				Value: "tag",
			},
		},
		Name: &Attribute{
			Base: token.Base{
				UsageKey: "usage-tag-name",
				Range: core.Range{
					Start: core.Position{Character: 9, Line: 10},
					End:   core.Position{Character: 20, Line: 10},
				},
				XMLName: String{
					Range: core.Range{
						Start: core.Position{Character: 9, Line: 10},
						End:   core.Position{Character: 13, Line: 10},
					},
					Value: "name",
				},
			},
			Value: String{
				Range: core.Range{
					Start: core.Position{Character: 15, Line: 10},
					End:   core.Position{Character: 19, Line: 10},
				},
				Value: "tag1",
			},
		},
		Title: &Attribute{
			Base: token.Base{
				UsageKey: "usage-tag-title",
				Range: core.Range{
					Start: core.Position{Character: 21, Line: 10},
					End:   core.Position{Character: 44, Line: 10},
				},
				XMLName: String{
					Range: core.Range{
						Start: core.Position{Character: 21, Line: 10},
						End:   core.Position{Character: 26, Line: 10},
					},
					Value: "title",
				},
			},
			Value: String{
				Range: core.Range{
					Start: core.Position{Character: 28, Line: 10},
					End:   core.Position{Character: 43, Line: 10},
				},
				Value: "tag description",
			},
		},
	}
	a.Equal(doc.Tags[0], tag)

	tag = doc.Tags[1]
	a.Equal(tag.Deprecated.V(), "1.0.1").
		Empty(tag.XMLNameEnd.Value).
		Equal(tag.UsageKey, "usage-apidoc-tags")

	a.Equal(2, len(doc.Servers))
	srv := doc.Servers[0]
	a.Equal(srv.Name.V(), "admin").
		Equal(srv.URL.V(), "https://api.example.com/admin").
		Nil(srv.Description).
		Equal(srv.Summary.V(), "admin api")

	srv = doc.Servers[1]
	a.Equal(srv.Name.V(), "client").
		Equal(srv.URL.V(), "https://api.example.com/client").
		Equal(srv.Deprecated.V(), "1.0.1").
		Equal(srv.Description.V(), "\n        <p>client api</p>\n        ")

	a.NotNil(doc.License).
		Equal(doc.License.Text.V(), "MIT").
		Equal(doc.License.URL.V(), "https://opensource.org/licenses/MIT")

	a.NotNil(doc.Contact).
		Equal(doc.Contact.Name.V(), "test").
		Equal(doc.Contact.URL.V(), "https://example.com").
		Equal(doc.Contact.Email.V(), "test@example.com")

	a.Contains(doc.Description.V(), "<h2>h2</h2>").
		NotContains(doc.Description.V(), "</description>")

	a.True(doc.tagExists("tag1")).
		False(doc.tagExists("not-exists"))

	a.True(doc.serverExists("admin")).
		False(doc.serverExists("not-exists"))

	a.Equal(2, len(doc.Mimetypes)).
		Equal("application/xml", doc.Mimetypes[0].Content.Value)
}

func TestAPIDoc_all(t *testing.T) {
	a := assert.New(t)

	data, err := ioutil.ReadFile("./testdata/all.xml")
	a.NotError(err).NotNil(data)
	doc := &APIDoc{}
	a.NotError(doc.Parse(core.Block{Data: data}))

	a.Equal(doc.Version.V(), "1.1.1").False(doc.Version.XMLName.IsEmpty())

	a.Equal(len(doc.Tags), 2)
	tag := doc.Tags[0]
	a.Equal(tag.Name.V(), "tag1").
		NotEmpty(tag.Title.V())
	tag = doc.Tags[1]
	a.Equal(tag.Deprecated.V(), "1.0.1").
		Equal(tag.Name.V(), "tag2")

	a.Equal(2, len(doc.Servers))
	srv := doc.Servers[0]
	a.Equal(srv.Name.V(), "admin").
		Equal(srv.URL.V(), "https://api.example.com/admin").
		Nil(srv.Description)

	a.True(doc.tagExists("tag1")).
		False(doc.tagExists("not-exists"))

	a.True(doc.serverExists("admin")).
		False(doc.serverExists("not-exists"))

	a.Equal(2, len(doc.Mimetypes)).
		Equal("application/xml", doc.Mimetypes[0].Content.Value)

	// api
	a.Equal(1, len(doc.Apis))
}

func loadAPI(a *assert.Assertion) *API {
	doc := loadAPIDoc(a)

	data, err := ioutil.ReadFile("./testdata/api.xml")
	a.NotError(err).NotNil(data)

	a.NotError(doc.Parse(core.Block{Data: data}))
	return doc.Apis[0]
}

func TestAPI(t *testing.T) {
	a := assert.New(t)
	api := loadAPI(a)

	a.Equal(api.Version.V(), "1.1.0").
		Equal(2, len(api.Tags))

	a.Equal(len(api.Responses), 2)
	resp := api.Responses[0]
	a.Equal(resp.Mimetype.V(), "json").
		Equal(resp.Status.V(), 200).
		Equal(resp.Type.V(), TypeObject).
		Equal(len(resp.Items), 3)
	sex := resp.Items[1]
	a.Equal(sex.Type.V(), TypeString).
		Equal(sex.Default.V(), "male").
		Equal(len(sex.Enums), 2)
	example := resp.Examples[0]
	a.Equal(example.Mimetype.V(), "json").
		NotEmpty(example.Content.Value.Value)

	a.Equal(1, len(api.Requests))
	req := api.Requests[0]
	a.Equal(req.Mimetype.V(), "json").
		Equal(req.Headers[0].Name.V(), "authorization")

	// callback
	cb := api.Callback
	a.Equal(cb.Method.V(), "POST").
		Equal(cb.Requests[0].Type.V(), TypeObject).
		Equal(cb.Requests[0].Mimetype.V(), "json").
		Equal(cb.Responses[0].Status.V(), 200)
}