# OpenSearchUtil

Given an object, makes an OpenSearch index template.

## Example

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
		DOB            time.Time
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
        "type": "basic_date_time"
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
