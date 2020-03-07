package linesection

import (
	"bufio"
	"io"
)

// LineSectionReader is a SectionReader that seeks its start and end to the
// next newline.
type LineSectionReader struct {
	*io.SectionReader
	start, end         int64
	lineStart, lineEnd int64
}

func (m *LineSectionReader) Pos() int64 {
	return m.lineStart
}

func defaultLineSectionReader(r io.ReaderAt, off int64, n int64) *LineSectionReader {
	return &LineSectionReader{
		SectionReader: io.NewSectionReader(r, off, n),
		start:         off,
		end:           off + n,
		lineStart:     off,
		lineEnd:       off + n,
	}
}

func NewLineSectionReader(r io.ReaderAt, off int64, n int64) *LineSectionReader {
	var err error
	var start int64 = off
	var end int64 = off + n

	sr := io.NewSectionReader(r, off, n+4096)
	// Find the newline after "off"
	// Except that "off 0" is a special case of file start.
	if off > 0 {
		br := bufio.NewReader(sr)
		line, err := br.ReadSlice('\n')
		if err != nil && err != io.EOF {
			return defaultLineSectionReader(r, off, n)
		}
		start = off + int64(len(line)) + 1
	}

	if start > end {
		return defaultLineSectionReader(r, start, 0)
	}

	// Find the newline after "off + n"
	if _, err = sr.Seek(n, io.SeekStart); err != nil {
		return defaultLineSectionReader(r, start, end-start)
	}
	br := bufio.NewReader(sr)
	line, err := br.ReadSlice('\n')
	if err != nil && err != io.EOF {
		return defaultLineSectionReader(r, start, end-start)
	}
	end = off + n + int64(len(line))

	// Create the new SectionReader and the rest of the fields.
	nr := &LineSectionReader{
		SectionReader: io.NewSectionReader(r, start, end-start),
		start:         off,
		end:           off + n,
		lineStart:     start,
		lineEnd:       end,
	}
	return nr
}
