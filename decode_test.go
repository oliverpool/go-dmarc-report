package report

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	a := assert.New(t)

	f, err := os.Open("testdata/test_report.xml")
	a.NoError(err)
	defer f.Close()
	agg, err := Decode(f)
	a.NoError(err)

	a.NotNil(agg.XMLName)
	a.Equal(Metadata{
		OrgName:          "Test Inc.",
		Email:            "postmaster@test",
		ReportID:         "1.id.0",
		ExtraContactInfo: "http://test/help",
		DateRange: DateRange{
			Begin: Time{time.Unix(1524182400, 0)},
			End:   Time{time.Unix(1524268799, 0)},
		},
	}, agg.Metadata)

	pct := 100
	a.Equal(PolicyPublished{
		Domain:     "test.net",
		ADKIM:      "r",
		ASPF:       "r",
		Policy:     "none",
		Percentage: &pct,
	}, agg.PolicyPublished)

	a.Len(agg.Records, 2)
	a.Equal(Row{
		SourceIP: "192.168.1.1",
		Count:    5,
		PolicyEvaluated: PolicyEvaluated{
			Disposition: "none",
			DKIM:        "pass",
			SPF:         "pass",
		},
	}, agg.Records[0].Row)
	a.Equal(Identifiers{
		HeaderFrom: "test.net",
	}, agg.Records[0].Identifiers)
	a.Equal(AuthResults{
		DKIMAuthResult{
			Domain:   "test.net",
			Result:   "pass",
			Selector: "selector",
		},
		SPFAuthResult{
			Domain: "test.net",
			Result: "pass",
			Scope:  "mfrom",
		},
	}, agg.Records[0].AuthResults)
}

func TestDecodeGzip(t *testing.T) {
	a := assert.New(t)

	f, err := os.Open("testdata/test_report.xml.gz")
	a.NoError(err)
	defer f.Close()
	agg, err := DecodeGzip(f)
	a.NoError(err)
	a.Equal("Test Inc.", agg.Metadata.OrgName)
}

func TestDecodeZip(t *testing.T) {
	a := assert.New(t)

	f, err := os.Open("testdata/test_report.zip")
	a.NoError(err)
	stat, err := f.Stat()
	a.NoError(err)
	defer f.Close()
	agg, err := DecodeZip(f, stat.Size())
	a.NoError(err)
	a.Equal("Test Inc.", agg.Metadata.OrgName)
}

func TestRecordErr(t *testing.T) {
	a := assert.New(t)

	f, err := os.Open("testdata/test_report.xml")
	a.NoError(err)
	defer f.Close()
	agg, err := Decode(f)
	a.NoError(err)
	a.NoError(agg.Records[0].Err())
	a.Error(agg.Records[1].Err())
	a.Equal(`Failure for source IP 10.1.1.2:
	* DKIM is not aligned
	* SPF is not aligned
	* DKIM authentication failed
	* SPF authentication failed
`, agg.Records[1].Err().Error())
}

func TestAggregateErr(t *testing.T) {
	a := assert.New(t)

	f, err := os.Open("testdata/test_report.xml.gz")
	a.NoError(err)
	defer f.Close()
	agg, err := DecodeGzip(f)
	a.NoError(err)
	a.Error(agg.Err())
	a.Equal(`Some record failed:
	* Failure for source IP 10.1.1.1:
		* DKIM is not aligned
		* SPF is not aligned
		* DKIM authentication failed
		* SPF authentication failed
`, agg.Err().Error())
}
