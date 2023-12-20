package file_old

import (
	"testing"
)

func TestFileTypeFromExtension(t *testing.T) {

	type test struct {
		name      string
		extension string
		expected  string
	}

	tcs := []test{
		{
			name:      "image jpg",
			extension: "jpg",
			expected:  "image/jpeg",
		},
		{
			name:      "audio mp3",
			extension: "mp3",
			expected:  "audio/mpeg",
		},
		{
			name:      "audio mp3",
			extension: "tar",
			expected:  "application/x-tar",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			out := fileTypeFromExtension(tc.extension)

			if out.MimeType != tc.expected {
				t.Errorf("FileTypeExtInfo does not match to expected: %s but got: %s", tc.expected, out.MimeType)
			}

		})
	}

}
