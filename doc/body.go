// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

import (
	"bytes"
	"strings"

	"github.com/caixw/apidoc/doc/lexer"
	"github.com/caixw/apidoc/doc/schema"
)

// Request 表示用户请求所表示的数据。
type Request = Body

// Response 表示一次请求或是返回的数据。
type Response struct {
	Body
	Status string `yaml:"status" json:"status"`
}

// Body 表示请求和返回的共有内容
type Body struct {
	Mimetype string         `yaml:"mimetype,omitempty" json:"mimetype,omitempty"`
	Headers  []*Header      `yaml:"headers,omitempty" json:"headers,omitempty"`
	Type     *schema.Schema `yaml:"type" json:"type"`
	Examples []*Example     `yaml:"examples,omitempty" json:"examples,omitempty"`
}

// Header 报头
type Header struct {
	Name     string `yaml:"name" json:"name"`                             // 参数名称
	Summary  string `yaml:"summary" json:"summary"`                       // 参数介绍
	Optional bool   `yaml:"optional,omitempty" json:"optional,omitempty"` // 是否可以为空
}

// Example 示例
type Example struct {
	Mimetype string `yaml:"mimetype" json:"mimetype"`
	Summary  string `yaml:"summary,omitempty" json:"summary,omitempty"`
	Value    string `yaml:"value" json:"value"` // 示例内容
}

// 解析示例代码，格式如下：
//  @apiExample application/json
//  {
//      "id": 1,
//      "name": "name",
//  }
func (body *Body) parseExample(tag *lexer.Tag) error {
	lines := tag.Lines(2)
	if len(lines) != 2 {
		return tag.ErrInvalidFormat()
	}

	words := lexer.SplitWords(lines[0], 2)

	if body.Examples == nil {
		body.Examples = make([]*Example, 0, 3)
	}

	example := &Example{
		Mimetype: string(words[0]),
		Value:    string(lines[1]),
	}
	if len(words) == 2 { // 如果存在简介
		example.Summary = string(words[1])
	}

	body.Examples = append(body.Examples, example)

	return nil
}

func (body *Body) parseHeader(tag *lexer.Tag) error {
	data := tag.Words(3)
	if len(data) != 3 {
		return tag.ErrInvalidFormat()
	}

	if body.Headers == nil {
		body.Headers = make([]*Header, 0, 3)
	}

	body.Headers = append(body.Headers, &Header{
		Name:     string(data[0]),
		Summary:  string(data[2]),
		Optional: isOptional(data[1]),
	})

	return nil
}

var requiredBytes = []byte("required")

func isOptional(data []byte) bool {
	return !bytes.Equal(data, requiredBytes)
}

func newResponse(l *lexer.Lexer, tag *lexer.Tag) (*Response, error) {
	data := tag.Words(3)
	if len(data) != 3 {
		return nil, tag.ErrInvalidFormat()
	}

	s := &schema.Schema{}
	if err := s.Build(tag, nil, data[1], nil, data[2]); err != nil {
		return nil, err
	}
	resp := &Response{
		Body: Body{
			Mimetype: string(data[1]),
		},
	}

	for tag, eof := l.Tag(); !eof; tag, eof = l.Tag() {
		switch strings.ToLower(tag.Name) {
		case "@apiexample":
			if err := resp.parseExample(tag); err != nil {
				return nil, err
			}
		case "@apiheader":
			if err := resp.parseHeader(tag); err != nil {
				return nil, err
			}
		case "@apiparam":
			data := tag.Words(4)
			if len(data) != 4 {
				return nil, tag.ErrInvalidFormat()
			}

			if err := s.Build(tag, data[0], data[1], data[2], data[3]); err != nil {
				return nil, err
			}
		default:
			l.Backup(tag)
			return resp, nil
		}
	}

	return resp, nil
}