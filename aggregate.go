package report

import (
	"encoding/xml"
	"fmt"
	"strings"
)

// Aggregate represents root of dmarc report struct
type Aggregate struct {
	XMLName         xml.Name        `xml:"feedback"`
	Metadata        Metadata        `xml:"report_metadata"`
	PolicyPublished PolicyPublished `xml:"policy_published"`
	Records         []Record        `xml:"record"`
}

type errorSlice struct {
	title  string
	errors []error
}

func (es errorSlice) ErrorOrNil() error {
	if len(es.errors) == 0 {
		return nil
	}
	return es
}

func (es errorSlice) Error() string {
	s := es.title
	if s == "" {
		s = "error"
	}
	s += "\n"
	for _, e := range es.errors {
		idented := strings.ReplaceAll(strings.TrimSpace(e.Error()), "\n", "\n\t")
		s += "\t* " + idented + "\n"
	}
	return s
}

// Err returns a non-nil error if any of the record as an Err
func (agg Aggregate) Err() error {
	serr := errorSlice{
		title: "Some record failed:",
	}
	for _, r := range agg.Records {
		err := r.Err()
		if err == nil {
			continue
		}
		serr.errors = append(serr.errors, err)
	}
	return serr.ErrorOrNil()
}

// Metadata represents feedback>report_metadata section
type Metadata struct {
	OrgName          string    `xml:"org_name"`
	Email            string    `xml:"email"`
	ExtraContactInfo string    `xml:"extra_contact_info"`
	ReportID         string    `xml:"report_id"`
	DateRange        DateRange `xml:"date_range"`
}

// PolicyPublished represents feedback>policy_published section
type PolicyPublished struct {
	Domain     string `xml:"domain"`
	ADKIM      string `xml:"adkim"`
	ASPF       string `xml:"aspf"`
	Policy     string `xml:"p"`
	SPolicy    string `xml:"sp"`
	Percentage *int   `xml:"pct"`
}

// Record represents feedback>record section
type Record struct {
	Row         Row         `xml:"row"`
	Identifiers Identifiers `xml:"identifiers"`
	AuthResults AuthResults `xml:"auth_results"`
}

// Err returns a non-nil error if any of the DMARC policy failed (or is missing)
func (r Record) Err() error {
	serr := errorSlice{
		title: "Failure for source IP " + r.Row.SourceIP + ":",
	}
	if !r.FinalDispositionSuccess() {
		serr.errors = append(serr.errors, fmt.Errorf("DMARC disposition failed"))
	}
	if !r.DKIMAligned() {
		serr.errors = append(serr.errors, fmt.Errorf("DKIM is not aligned"))
	}
	if !r.SPFAligned() {
		serr.errors = append(serr.errors, fmt.Errorf("SPF is not aligned"))
	}
	if !r.DKIMSuccess() {
		serr.errors = append(serr.errors, fmt.Errorf("DKIM authentication failed"))
	}
	if !r.SPFSuccess() {
		serr.errors = append(serr.errors, fmt.Errorf("SPF authentication failed"))
	}
	return serr.ErrorOrNil()
}

// FinalDispositionSuccess is the result of the domain’s policy combined with the DKIM and SPF aligned policy results.
func (r Record) FinalDispositionSuccess() bool {
	return r.Row.PolicyEvaluated.Disposition == "none"
}

// DKIMAligned returns true if the DKIM is aligned:
// Domain in the RFC5322.From header matches the domain in the "d=" tag in the DKIM signature.
func (r Record) DKIMAligned() bool {
	return r.Row.PolicyEvaluated.DKIM == "pass"
}

// SPFAligned returns true if the SPF is aligned:
// Domain in the RFC5322.From header matches the domain in the RFC5321.MailFrom field
func (r Record) SPFAligned() bool {
	return r.Row.PolicyEvaluated.SPF == "pass"
}

// DKIMSuccess returns true if the DKIM authenticated was successful.
func (r Record) DKIMSuccess() bool {
	return r.AuthResults.DKIM.Result == "pass"
}

// SPFSuccess returns true if the SPF authenticated was successful.
func (r Record) SPFSuccess() bool {
	return r.AuthResults.SPF.Result == "pass"
}

// Row represents feedback>record>row section
type Row struct {
	SourceIP        string          `xml:"source_ip"`
	Count           int             `xml:"count"`
	PolicyEvaluated PolicyEvaluated `xml:"policy_evaluated"`
}

// PolicyEvaluated represents feedback>record>row>policy_evaluated section
type PolicyEvaluated struct {
	Disposition string `xml:"disposition"`
	DKIM        string `xml:"dkim"`
	SPF         string `xml:"spf"`
}

// Identifiers represents feedback>record>identifiers section
type Identifiers struct {
	HeaderFrom string `xml:"header_from"`
}

// AuthResults represents feedback>record>auth_results section
type AuthResults struct {
	DKIM DKIMAuthResult `xml:"dkim"`
	SPF  SPFAuthResult  `xml:"spf"`
}

// DKIMAuthResult represents feedback>record>auth_results>dkim sections
type DKIMAuthResult struct {
	Domain   string `xml:"domain"`
	Result   string `xml:"result"`
	Selector string `xml:"selector"`
}

// SPFAuthResult represents feedback>record>auth_results>spf section
type SPFAuthResult struct {
	Domain string `xml:"domain"`
	Result string `xml:"result"`
	Scope  string `xml:"scope"`
}
