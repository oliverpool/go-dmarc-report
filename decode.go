package report

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
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

// DecodeFile decodes dmarc aggregate report based on its name
func DecodeFile(filename string, r io.Reader) (*Aggregate, error) {
	ext := filepath.Ext(filename)
	if ext == ".gz" {
		return DecodeGzip(r)
	}
	if ext != ".zip" {
		// not .gz and not .zip: try to decode it as .xml
		return Decode(r)
	}
	// .zip must be fully read to decode
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	br := bytes.NewReader(buf)
	return DecodeZip(br, br.Size())
}
