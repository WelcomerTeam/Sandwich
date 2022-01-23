package internal

// The MIT License (MIT)

// Copyright (c) 2015-present Rapptz

// Permission is hereby granted, free of charge, to any person obtaining a
// copy of this software and associated documentation files (the 'Software'),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND, EXPRESS
// OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

import "strings"

var quotes = map[string]string{
	`"`: `"`,
	"‘": "’",
	"‚": "‛",
	"“": "”",
	"„": "‟",
	"⹂": "⹂",
	"「": "」",
	"『": "』",
	"〝": "〞",
	"﹁": "﹂",
	"﹃": "﹄",
	"＂": "＂",
	"｢": "｣",
	"«": "»",
	"‹": "›",
	"《": "》",
	"〈": "〉",
}

var (
	allQuotes    = []string{}
	allQuotesMap = map[string]bool{}
)

type StringView struct {
	index    int
	buffer   string
	end      int
	previous int
}

func NewStringView(buffer string) (s *StringView) {
	s = &StringView{
		index:    0,
		buffer:   buffer,
		end:      len(buffer),
		previous: 0,
	}

	return s
}

func (s *StringView) Current() (current string, ok bool) {
	if s.EOF() {
		return "", false
	}

	return string(s.buffer[s.index]), true
}

func (s *StringView) EOF() bool {
	return s.index >= s.end
}

func (s *StringView) Undo() {
	s.index = s.previous
}

func (s *StringView) SkipWS() bool {
	pos := 0
	for !s.EOF() {
		if (s.index + pos) >= len(s.buffer) {
			break
		}

		current := string(s.buffer[s.index+pos])
		if strings.TrimSpace(current) != "" {
			break
		}

		pos++
	}

	s.previous = s.index
	s.index += pos

	return s.previous != s.index
}

func (s *StringView) SkipString(str string) bool {
	strlen := len(str)
	if s.buffer[s.index:s.index+strlen] == str {
		s.previous = s.index
		s.index += strlen

		return true
	}

	return false
}

func (s *StringView) ReadRest() (result string) {
	result = s.buffer[s.index:]
	s.previous = s.index
	s.index = s.end

	return
}

func (s *StringView) Read(n int) (result string) {
	result = s.buffer[s.index : s.index+n]
	s.previous = s.index
	s.index += n

	return
}

func (s *StringView) Get() (result string, ok bool) {
	if s.index+1 < len(s.buffer) {
		result = string(s.buffer[s.index+1])
		ok = true
	}

	s.previous = s.index
	s.index++

	return
}

func (s *StringView) GetWord() (result string) {
	pos := 0

	for !s.EOF() {
		if s.index+pos < len(s.buffer) {
			if strings.TrimSpace(string(s.buffer[s.index+pos])) == "" {
				break
			}
			pos++
		} else {
			break
		}
	}

	s.previous = s.index
	result = s.buffer[s.index : s.index+pos]
	s.index += pos

	return
}

func (s *StringView) GetQuotedWord() (result string, ok bool, err error) {
	current, ok := s.Current()
	if !ok {
		return result, false, nil
	}

	var results []string

	var escapedQuotes []string

	closeQuotes, isQuoted := quotes[current]
	if isQuoted {
		results = []string{}
		escapedQuotes = []string{current, closeQuotes}
	} else {
		results = []string{current}
		escapedQuotes = allQuotes
	}

	for !s.EOF() {
		current, ok = s.Get()
		if !ok {
			if isQuoted {
				return result, true, ErrExpectedClosingQuoteError
			}

			result = strings.Join(results, "")

			return result, true, nil
		}

		//	currently we accept strings in the format of "hello world"
		//	to embed a quote inside the string you must escape it: "a \"world\""
		if current == "\\" {
			nextChar, ok := s.Get()
			if !ok {
				if isQuoted {
					return result, true, ErrExpectedClosingQuoteError
				}

				result = strings.Join(results, "")

				return result, true, nil
			}

			ok = contains(escapedQuotes, nextChar)
			if ok {
				results = append(results, nextChar)
			} else {
				// different escape character, ignore it
				s.Undo()
				results = append(results, current)
			}

			continue
		}

		_, ok = allQuotesMap[current]
		if !isQuoted && ok {
			return result, true, ErrUnexpectedQuoteError
		}

		if isQuoted && current == closeQuotes {
			nextChar, ok := s.Get()
			validEOF := !ok || strings.TrimSpace(nextChar) == ""

			if !validEOF {
				return result, true, ErrInvalidEndOfQuotedStringError
			}

			result = strings.Join(results, "")
		}

		if strings.TrimSpace(current) == "" && !isQuoted {
			result = strings.Join(results, "")

			return result, true, nil
		}

		results = append(results, current)
	}

	return result, true, nil
}

func init() {
	for _, v := range quotes {
		allQuotes = append(allQuotes, v)
		allQuotesMap[v] = true
	}

	for k := range quotes {
		allQuotes = append(allQuotes, k)
		allQuotesMap[k] = true
	}
}
