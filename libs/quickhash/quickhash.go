package quickhash

import (
	"encoding/binary"
	"github.com/spf13/afero"
	"github.com/twmb/murmur3"
	"io"
)

// Quickhash is heavily inspired on https://github.com/kalafut/imohash
// but with support for aferoFS  https://github.com/spf13/afero

const Size = 16

// SampleThreshold Files smaller than this will be hashed in their entirety.
const SampleThreshold = 128 * 1024
const SampleSize = 16 * 1024

var emptyArray = [Size]byte{}

type QuickHasher struct {
	hasher          murmur3.Hash128
	sampleSize      int
	sampleThreshold int
	bytesAdded      int
}

// NewQuick returns a new QuickHasher using the default sample size and sample threshold values.
func NewQuick() QuickHasher {
	return NewCustomQuick(SampleSize, SampleThreshold)
}

// NewCustomQuick returns a new QuickHasher using the provided sample size
// and sample threshhold values. The entire file will be hashed
// (i.e. no sampling), if sampleSize < 1.
func NewCustomQuick(sampleSize, sampleThreshold int) QuickHasher {
	h := QuickHasher{
		hasher:          murmur3.New128(),
		sampleSize:      sampleSize,
		sampleThreshold: sampleThreshold,
	}

	return h
}

// SumFile hashes a file using the parameters.
func (qh *QuickHasher) SumFile(f afero.File) ([Size]byte, error) {
	fStat, err := f.Stat()
	if err != nil {
		return emptyArray, err
	}
	return qh.hashCore(f, fStat.Size())
}

type SeekReader interface {
	io.Reader
	io.Seeker
}

// hashCore hashes a SectionReader using the ImoHash parameters.
func (qh *QuickHasher) hashCore(f SeekReader, fSize int64) ([Size]byte, error) {
	var result [Size]byte

	qh.hasher.Reset()

	if fSize < int64(qh.sampleThreshold) || qh.sampleSize < 1 {
		if _, err := io.Copy(qh.hasher, f); err != nil {
			return emptyArray, err
		}
	} else {
		buffer := make([]byte, qh.sampleSize)
		_, err := f.Read(buffer)
		if err != nil {
			return emptyArray, err
		}
		// these writess never fail
		qh.hasher.Write(buffer)

		_, err = f.Seek(fSize/2, 0)
		if err != nil {
			return emptyArray, err
		}
		_, err = f.Read(buffer)
		if err != nil {
			return emptyArray, err
		}
		qh.hasher.Write(buffer)

		_, err = f.Seek(int64(-qh.sampleSize), 2)
		if err != nil {
			return emptyArray, err
		}
		_, err = f.Read(buffer)
		if err != nil {
			return emptyArray, err
		}
		qh.hasher.Write(buffer)
	}

	sizeBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(sizeBytes, uint64(fSize))
	qh.hasher.Write(sizeBytes)

	hash := qh.hasher.Sum(nil)

	binary.PutUvarint(hash, uint64(fSize))
	copy(result[:], hash)

	return result, nil
}
