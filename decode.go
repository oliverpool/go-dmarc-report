package report

import (
	"archive/zip"
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"io"
	"path/filepath"
)

// Decode a dmarc aggregate report
func Decode(r io.Reader) (*Aggregate, error) {
	agg := Aggregate{}
	return &agg, xml.NewDecoder(r).Decode(&agg)
}

// DecodeGzip decodes a gzipped dmarc aggregate report
func DecodeGzip(r io.Reader) (*Aggregate, error) {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not create gzip reader: %w", err)
	}
	defer gr.Close()

	return Decode(gr)
}

// DecodeZip decodes a zipped dmarc aggregate report
func DecodeZip(r io.ReaderAt, size int64) (*Aggregate, error) {
	zr, err := zip.NewReader(r, size)
	if err != nil {
		return nil, fmt.Errorf("could not create zip reader: %w", err)
	}

	for _, file := range zr.File {
		if file.FileInfo().IsDir() {
			continue
		}
		ext := filepath.Ext(file.Name)
		if ext != ".xml" {
			continue
		}

		fr, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("could not open zipped file %s: %w", file.Name, err)
		}
		defer fr.Close()

		return Decode(fr)
	}

	return nil, fmt.Errorf("no suitable .xml file found in zip")
}

// DecodeErr decodes and check errors in a dmarc aggregate report
func DecodeErr(r io.Reader) error {
	agg, err := Decode(r)
	if err != nil {
		return err
	}
	return agg.Err()
}

// DecodeGzipErr decodes and check errors in a gzipped dmarc aggregate report
func DecodeGzipErr(r io.Reader) error {
	agg, err := DecodeGzip(r)
	if err != nil {
		return err
	}
	return agg.Err()
}

// DecodeZipErr decodes and check errors in a zipped dmarc aggregate report
func DecodeZipErr(r io.ReaderAt, size int64) error {
	agg, err := DecodeZip(r, size)
	if err != nil {
		return err
	}
	return agg.Err()
}
