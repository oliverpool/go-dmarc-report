package report_test

import (
	"fmt"
	"log"
	"os"

	report "github.com/oliverpool/go-dmarc-report"
)

func Example() {
	f, err := os.Open("testdata/test_report.xml.gz")
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()
	agg, err := report.DecodeGzip(f)
	if err != nil {
		log.Panic(err)
	}
	// You can now read the report
	fmt.Println(agg.Err())
	// Output: Some record failed:
	//	* Failure for source IP 10.1.1.1:
	//		* DKIM is not aligned
	//		* SPF is not aligned
	//		* DKIM authentication failed
	//		* SPF authentication failed
}

func ExampleDecodeErr() {
	f, err := os.Open("testdata/test_report.xml")
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()
	err = report.DecodeErr(f)
	fmt.Println(err)
	// Output: Some record failed:
	//	* Failure for source IP 10.1.1.2:
	//		* DKIM is not aligned
	//		* SPF is not aligned
	//		* DKIM authentication failed
	//		* SPF authentication failed
}
