package missing_header_fields

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/scroll-tech/go-ethereum/common"
)

type missingHeader struct {
	headerNum  uint64
	difficulty uint64
	stateRoot  common.Hash
	extraData  []byte
}

type Reader struct {
	file           *os.File
	reader         *bufio.Reader
	sortedVanities map[int][32]byte
	lastReadHeader *missingHeader
}

func NewReader(filePath string) (*Reader, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	r := &Reader{
		file:   f,
		reader: bufio.NewReader(f),
	}

	if err = r.initialize(); err != nil {
		if err = f.Close(); err != nil {
			return nil, fmt.Errorf("failed to close file after initialization error: %w", err)
		}
		return nil, fmt.Errorf("failed to initialize reader: %w", err)
	}

	return r, nil
}

func (r *Reader) initialize() error {
	// reset the reader and last read header
	if _, err := r.file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek to start: %w", err)
	}
	r.reader = bufio.NewReader(r.file)
	r.lastReadHeader = nil

	// read the count of unique vanities
	vanityCount, err := r.reader.ReadByte()
	if err != nil {
		return err
	}

	// read the unique vanities
	r.sortedVanities = make(map[int][32]byte)
	for i := uint8(0); i < vanityCount; i++ {
		var vanity [32]byte
		if _, err = r.reader.Read(vanity[:]); err != nil {
			return err
		}
		r.sortedVanities[int(i)] = vanity
	}

	return nil
}

func (r *Reader) Read(headerNum uint64) (difficulty uint64, stateRoot common.Hash, extraData []byte, err error) {
	if r.lastReadHeader != nil && headerNum < r.lastReadHeader.headerNum {
		if err = r.initialize(); err != nil {
			return 0, common.Hash{}, nil, fmt.Errorf("failed to reinitialize reader due to requested header number being lower than last read header: %w", err)
		}
	}

	if r.lastReadHeader == nil {
		if _, _, err = r.ReadNext(); err != nil {
			return 0, common.Hash{}, nil, err
		}
	}

	if headerNum > r.lastReadHeader.headerNum {
		// skip the headers until the requested header number
		for i := r.lastReadHeader.headerNum; i < headerNum; i++ {
			if _, _, err = r.ReadNext(); err != nil {
				return 0, common.Hash{}, nil, err
			}
		}
	}

	if headerNum == r.lastReadHeader.headerNum {
		return r.lastReadHeader.difficulty, r.lastReadHeader.stateRoot, r.lastReadHeader.extraData, nil
	}

	return 0, common.Hash{}, nil, fmt.Errorf("error reading header number %d: last read header number is %d", headerNum, r.lastReadHeader.headerNum)
}

func (r *Reader) ReadNext() (difficulty uint64, extraData []byte, err error) {
	// read the bitmask
	bitmaskByte, err := r.reader.ReadByte()
	if err != nil {
		return 0, nil, fmt.Errorf("failed to read bitmask: %v", err)
	}

	bits := newBitMaskFromBytes(bitmaskByte)

	// read the vanity index
	vanityIndex, err := r.reader.ReadByte()
	if err != nil {
		return 0, nil, fmt.Errorf("failed to read vanity index: %v", err)
	}

	stateRoot := make([]byte, common.HashLength)
	if _, err := io.ReadFull(r.reader, stateRoot); err != nil {
		return 0, nil, fmt.Errorf("failed to read state root: %v", err)
	}

	seal := make([]byte, bits.sealLen())
	if _, err = io.ReadFull(r.reader, seal); err != nil {
		return 0, nil, fmt.Errorf("failed to read seal: %v", err)
	}

	// construct the extraData field
	vanity := r.sortedVanities[int(vanityIndex)]
	var b bytes.Buffer
	b.Write(vanity[:])
	b.Write(seal)

	// we don't have the header number, so we'll just increment the last read header number
	// we assume that the headers are written in order, starting from 0
	if r.lastReadHeader == nil {
		r.lastReadHeader = &missingHeader{
			headerNum:  0,
			difficulty: uint64(bits.difficulty()),
			stateRoot:  common.BytesToHash(stateRoot),
			extraData:  b.Bytes(),
		}
	} else {
		r.lastReadHeader.headerNum++
		r.lastReadHeader.difficulty = uint64(bits.difficulty())
		r.lastReadHeader.stateRoot = common.BytesToHash(stateRoot)
		r.lastReadHeader.extraData = b.Bytes()
	}

	return difficulty, b.Bytes(), nil
}

func (r *Reader) Close() error {
	return r.file.Close()
}

// bitMask is a bitmask that encodes the following information:
//
//	bit 6: 0 if difficulty is 2, 1 if difficulty is 1
//	bit 7: 0 if seal length is 65, 1 if seal length is 85
type bitMask struct {
	b uint8
}

func newBitMaskFromBytes(b uint8) bitMask {
	return bitMask{b}
}

func (b bitMask) difficulty() int {
	val := (b.b >> 6) & 0x01
	if val == 0 {
		return 2
	} else {
		return 1
	}
}

func (b bitMask) sealLen() int {
	val := (b.b >> 7) & 0x01
	if val == 0 {
		return 65
	} else {
		return 85
	}
}
