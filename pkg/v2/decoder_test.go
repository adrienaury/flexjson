package v2_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	v2 "github.com/adrienaury/ordered-json/pkg/v2"
)

func TestV2(t *testing.T) {
	data := `
	{
		"clientIdentiferData": {
			"clientID": "89012",
			"webPropertyID": "89012-01",
			"campaignID": "004"
		},
		"journeyKeyData": {
			"campaignCreationPhase": "Teach",
			"purchaseJourneyPhase": "Awareness",
			"topicCategory": "Organizer",
			"topicSubCategory": "Event Team",
			"topicSubSubCategory": "Null"
		},
		"uploadedData": [{
				"dataCategoryLabel": "name of organization",
				"dataStrings": [{
					"string": "UBM"
				}]
			},
			{
				"dataCategoryLabel": "our mission",
				"dataStrings": [{
					"string": "We Serve the Needs of B2B Buyers"
				}]
			},
			{
				"dataCategoryLabel": "our history",
				"dataStrings": [{
					"string": "UBM was fonded in 1948",
					"string2": "It is was purchased by Informa in 2017",
					"string3": "We Serve the Needs of B2B Buyers"
				}]
			},
			{
				"dataCategoryLabel": "our differentiator",
				"dataStrings": [{
					"string": "UBM is the largest event management company in US",
					"string2": "It hosts more than 100 events annually",
					"string3": "Close to a million buyers and sellers are brought together by UBM"
				}]
			}
		]
	}
	`

	decoder := v2.NewDecoder(json.NewDecoder(strings.NewReader(data)))
	result := decoder.Decode()
	fmt.Println(result)
}
