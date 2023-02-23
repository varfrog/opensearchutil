# OpenSearchUtil

Utilities for working with OpenSearch.

- **IndexGenerator**: given an object, makes an OpenSearch index template,
- **Field types**: go types for struct fields that:
  - when the struct is marshalled into JSON, the fields get marshalled into valid OpenSearch types,
  - when generating an index mapping JSON, the fields get assigned the appropriate OpenSearch type and format.

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


## Field types

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/varfrog/opensearchutil"
	"os"
	"time"
)

func main() {
	type dates struct {
		DateA opensearchutil.TimeBasicDate             `json:"date_a"`
		DateB opensearchutil.TimeBasicDateTimeNoMillis `json:"date_b"`
		DateC opensearchutil.TimeBasicDateTime         `json:"date_c"`
	}

	t := time.Now()
	d := dates{
		DateA: opensearchutil.TimeBasicDate(t),
		DateB: opensearchutil.TimeBasicDateTimeNoMillis(t),
		DateC: opensearchutil.TimeBasicDateTime(t),
	}

	// When generating an index mapping JSON, the fields get assigned the approprate OpenSearch type and format
	//
	mappingProperties, err := opensearchutil.NewMappingPropertiesBuilder().BuildMappingProperties(d)
	if err != nil {
		fmt.Printf("BuildMappingProperties: %v", err)
		os.Exit(1)
	}
	jsonBytes, err := opensearchutil.NewIndexGenerator().GenerateIndexJson(mappingProperties)
	if err != nil {
		fmt.Printf("GenerateIndexJson: %v", err)
		os.Exit(1)
	}
	fmt.Printf("Mapping JSON:\n%s\n", string(jsonBytes))

	// When marshalling into JSON, the fields marshall into the approprate formats:
	//
	jsonBytes, err = json.MarshalIndent(&d, "", "  ")
	if err != nil {
		fmt.Printf("json.MarshalIndent: %v", err)
		os.Exit(1)
	}
	fmt.Printf("\nDocument body:\n%s\n", string(jsonBytes))
}
```

Output:
```
Mapping JSON:
{
   "mappings": {
      "properties": {
         "date_a": {
            "format": "basic_date",
            "type": "date"
         },
         "date_b": {
            "format": "basic_date_time_no_millis",
            "type": "date"
         },
         "date_c": {
            "format": "basic_date_time",
            "type": "date"
         }
      }
   }
}

Document body:
{
  "date_a": "20230223",
  "date_b": "20230223T224633+02:00",
  "date_c": "20230223T224633.808+02:00"
}
```