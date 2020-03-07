package linesection

import (
	"bytes"
	"gotest.tools/assert"
	"io"
	"testing"
)

var testString = `0123456789
01234567890123456789
0123456789012345678901234567890123456789
01234567890123456789012345678901234567890123456789012345678901234567890123456789
0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789`

func TestNewLineSectionReader(t *testing.T) {
	var bslice [4000]byte
	var err error
	var n int
	var lr [20]*LineSectionReader

	f := bytes.NewReader([]byte(testString))

	for i := range lr {
		lr[i] = NewLineSectionReader(f, int64(i*10), 10)
		var b bytes.Buffer
		io.Copy(&b, lr[i])

		t.Logf("%d %v %#v", i, lr[i], b.String())
	}

	curr := lr[0]
	assert.Equal(t, curr.lineStart, int64(0))
	assert.Equal(t, curr.lineEnd, int64(11))
	assert.Equal(t, curr.Pos(), int64(0))
	assert.Equal(t, curr.Size(), int64(11))
	n, err = curr.ReadAt(bslice[:], curr.Size()-1)
	assert.Equal(t, err, io.EOF)
	assert.Equal(t, n, 1)
	assert.Equal(t, bslice[0], byte('\n'))

	curr = lr[1]
	assert.Equal(t, curr.lineStart, int64(12))
	assert.Equal(t, curr.lineEnd, int64(12+20))
	assert.Equal(t, curr.Pos(), int64(12))
	assert.Equal(t, curr.Size(), int64(20))
	n, err = curr.ReadAt(bslice[:1], 0)
	assert.NilError(t, err)
	assert.Equal(t, n, 1)
	assert.Assert(t, bslice[0] != byte('\n'))
	n, err = curr.ReadAt(bslice[:], curr.Size()-1)
	assert.Equal(t, err, io.EOF)
	assert.Equal(t, n, 1)
	assert.Equal(t, bslice[0], byte('\n'))

	// This one is an empty section because the newline start is pushed past the end
	curr = lr[2]
	assert.Equal(t, curr.lineStart, int64(33))
	assert.Equal(t, curr.lineEnd, int64(33))
	assert.Equal(t, curr.Pos(), int64(33))
	assert.Equal(t, curr.Size(), int64(0))
	n, err = curr.ReadAt(bslice[:1], 0)
	assert.Equal(t, err, io.EOF)
	assert.Equal(t, n, 0)

	curr = lr[7]
	assert.Equal(t, curr.lineStart, int64(74))
	assert.Equal(t, curr.lineEnd, int64(74+80))
	assert.Equal(t, curr.Pos(), int64(74))
	assert.Equal(t, curr.Size(), int64(80))
	n, err = curr.ReadAt(bslice[:1], 0)
	assert.NilError(t, err)
	assert.Equal(t, n, 1)
	assert.Assert(t, bslice[0] != byte('\n'))
	n, err = curr.ReadAt(bslice[:], curr.Size()-1)
	assert.Equal(t, err, io.EOF)
	assert.Equal(t, n, 1)
	assert.Equal(t, bslice[0], byte('\n'))
}
