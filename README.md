# OpenSearchUtil

Utilities for working with OpenSearch.

- **IndexGenerator**: given an object, makes an OpenSearch index template,
- **RequestBodyBuilder**: generates a request body for `POST _bulk`.

## IndexGenerator

```go
package main

import (
	_ "embed"
	"fmt"
	"github.com/varfrog/opensearchutil"
	"os"
)

func main() {
	type location struct {
		FullAddress string
		Confirmed   bool
	}
	type person struct {
		Name           string
		Email          string `opensearch:"type:keyword"`
		DOB            opensearchutil.TimeBasicDateTimeNoMillis
		Age            uint8
		AccountBalance float64
		IsDead         bool
		HomeLoc        location
		WorkLoc        *location
		SocialSecurity *string
	}

	builder := opensearchutil.NewMappingPropertiesBuilder()
	jsonGenerator := opensearchutil.NewIndexGenerator()

	mappingProperties, err := builder.BuildMappingProperties(person{})
	if err != nil {
		fmt.Printf("BuildMappingProperties: %v", err)
		os.Exit(1)
	}

	indexJson, err := jsonGenerator.GenerateIndexJson(mappingProperties)
	if err != nil {
		fmt.Printf("GenerateIndexJson: %v", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", string(indexJson))
}
```

Output:
```json
{
  "mappings": {
    "properties": {
      "account_balance": {
        "type": "float"
      },
      "age": {
        "type": "integer"
      },
      "dob": {
        "format": "basic_date_time_no_millis",
        "type": "date"
      },
      "email": {
        "type": "keyword"
      },
      "home_loc": {
        "properties": {
          "confirmed": {
            "type": "boolean"
          },
          "full_address": {
            "type": "text"
          }
        }
      },
      "is_dead": {
        "type": "boolean"
      },
      "name": {
        "type": "text"
      },
      "social_security": {
        "type": "text"
      },
      "work_loc": {
        "properties": {
          "confirmed": {
            "type": "boolean"
          },
          "full_address": {
            "type": "text"
          }
        }
      }
    }
  }
}
```

The resulting JSON contents is then used in a request to the [Create index API request](https://opensearch.org/docs/1.0/opensearch/rest-api/create-index/). Also specify "settings" and "aliases" that suit your needs.


## RequestBodyBuilder

```go
package main

import (
	"crypto/tls"
	"fmt"
	"github.com/opensearch-project/opensearch-go"
	"github.com/varfrog/opensearchutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	client, err := opensearch.NewClient(opensearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Addresses: []string{"https://localhost:9200"},
		Username:  "admin",
		Password:  "admin",
	})
	if err != nil {
		fmt.Printf("NewClient: %v\n", err)
		os.Exit(1)
	}

	type address struct {
		PostalCode uint32 `json:"postal_code"`
	}
	type person struct {
		ID      string                       `json:"id"      opensearch:"type=keyword"`
		Name    string                       `json:"name"`
		Age     uint8                        `json:"age"`
		DOB     opensearchutil.TimeBasicDate `json:"dob"`
		Address address                      `json:"address"`
	}

	ann := person{
		ID:      "680",
		Name:    "Ann",
		DOB:     opensearchutil.TimeBasicDate(time.Date(1983, 01, 01, 0, 0, 0, 0, time.UTC)),
		Age:     40,
		Address: address{PostalCode: 10000},
	}
	bob := person{
		ID:      "720",
		Name:    "Bob",
		DOB:     opensearchutil.TimeBasicDate(time.Date(1985, 01, 01, 0, 0, 0, 0, time.UTC)),
		Age:     38,
		Address: address{PostalCode: 38000},
	}
	carl := person{
		ID:      "850",
		Name:    "Carl",
		DOB:     opensearchutil.TimeBasicDate(time.Date(1960, 01, 01, 0, 0, 0, 0, time.UTC)),
		Age:     63,
		Address: address{PostalCode: 10000},
	}

	builder := opensearchutil.NewRequestBodyBuilder()

	indexName := "people_1"
	body, err := builder.BuildIndexBody([]opensearchutil.ObjectWrapper{
		{ID: ann.ID, Index: indexName, Object: ann},
		{ID: bob.ID, Index: indexName, Object: bob},
		{ID: carl.ID, Index: indexName, Object: carl},
	})
	if err != nil {
		fmt.Printf("Bulk: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Request body: \n%s\n", body)

	resp, err := client.Bulk(strings.NewReader(body))
	if err != nil {
		fmt.Printf("Bulk: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Response: %v\n", resp.Status())
}
```

Output:
```
Request body: 
{"index":{"_index":"people_1","_id":"680"}}
{"id":"680","name":"Ann","age":40,"dob":"19830101","address":{"postal_code":10000}}
{"index":{"_index":"people_1","_id":"720"}}
{"id":"720","name":"Bob","age":38,"dob":"19850101","address":{"postal_code":38000}}
{"index":{"_index":"people_1","_id":"850"}}
{"id":"850","name":"Carl","age":63,"dob":"19600101","address":{"postal_code":10000}}

Response: 200 OK

```
