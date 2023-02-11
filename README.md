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
		Age            uint8
		AccountBalance float64
		IsDead         bool
		HomeLoc        location
		WorkLoc        *location
		SocialSecurity *string
	}

	indexJsonStr, err := opensearchutil.GenerateIndexJson(person{})
	if err != nil {
		fmt.Printf("GenerateIndexJson: %v", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", indexJsonStr)
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
   },
   "settings": {
      "number_of_replicas": 2,
      "number_of_shards": 1
   }
}
```