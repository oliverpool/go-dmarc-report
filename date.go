package report

import (
	"encoding/xml"
	"time"
)

// DateRange represents feedback>report_metadata>date_range section
type DateRange struct {
	Begin Time `xml:"begin" json:"begin"`
	End   Time `xml:"end" json:"end"`
}

// Time is the custom time for DateRange.Begin and DateRange.End values
type Time struct {
	time.Time
}

// UnmarshalXML unmarshals unix timestamp to time.Time
func (t *Time) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v int64
	err := d.DecodeElement(&v, &start)
	t.Time = time.Unix(v, 0)
	return err
}
