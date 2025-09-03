# OpenSearchUtil

Utilities for working with OpenSearch.

- **IndexGenerator**: given an object, makes an OpenSearch index template,
- **MappingPropertiesBuilder and generation of index mappings**: given an object, makes an OpenSearch index mapping,
- **Field types**: go types for struct fields that:
  - when the struct is marshalled into JSON, the fields get marshalled into valid OpenSearch types,
  - when generating an index mapping JSON, the fields get assigned the appropriate OpenSearch type and format.

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/varfrog/opensearchutil)

## Complete Example

This example demonstrates all features of OpenSearchUtil in a single runnable program:

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
	// Define nested structs to demonstrate object mapping
	type location struct {
		FullAddress string
		Confirmed   bool
	}

	// Main struct demonstrating all features
	type person struct {
		// Basic types with default mappings
		Name           string
		Age            uint8
		AccountBalance float64
		IsDead         bool
		Aliases        []string

		// Custom field type overrides
		Email          string `opensearch:"type:keyword"`
		SocialSecurity *string `opensearch:"index_prefixes:min_chars=3;max_chars=5"`

		// Custom date types with proper OpenSearch formats
		DOB opensearchutil.TimeBasicDateTimeNoMillis
		CreatedAt opensearchutil.TimeBasicDateTime
		LastLogin opensearchutil.TimeBasicDate

		// Nested objects (value and pointer)
		HomeLoc location
		WorkLoc *location

		// Array of objects
		Addresses []location

		// Advanced text field features
		Title     string `opensearch:"type:text,analyzer:standard,search_analyzer:english,copy_to:all_text;searchable_text"`
		AllText   string `opensearch:"type:text"`
		SearchableText string `opensearch:"type:text"`
	}

	// Create sample data
	now := time.Now()
	samplePerson := person{
		Name:           "John Doe",
		Age:            30,
		AccountBalance: 1500.75,
		IsDead:         false,
		Aliases:        []string{"Johnny", "JD"},
		Email:          "john@example.com",
		SocialSecurity: opensearchutil.MakePtr("123-45-6789"),
		DOB:            opensearchutil.TimeBasicDateTimeNoMillis(now.AddDate(-30, 0, 0)),
		CreatedAt:      opensearchutil.TimeBasicDateTime(now),
		LastLogin:      opensearchutil.TimeBasicDate(now),
		HomeLoc: location{
			FullAddress: "123 Main St, Anytown, USA",
			Confirmed:   true,
		},
		WorkLoc: &location{
			FullAddress: "456 Business Ave, Anytown, USA",
			Confirmed:   true,
		},
		Addresses: []location{
			{FullAddress: "789 Secondary St, Anytown, USA", Confirmed: false},
		},
		Title:           "Software Engineer",
		AllText:         "",
		SearchableText:  "",
	}

	// 1. Generate mapping properties from struct
	builder := opensearchutil.NewMappingPropertiesBuilder()
	mappingProperties, err := builder.BuildMappingProperties(samplePerson)
	if err != nil {
		fmt.Printf("BuildMappingProperties: %v\n", err)
		os.Exit(1)
	}

	// 2. Generate complete index JSON with settings
	indexGenerator := opensearchutil.NewIndexGenerator()
	indexJson, err := indexGenerator.GenerateIndexJson(
		mappingProperties,
		&opensearchutil.IndexSettings{
			NumberOfShards:   opensearchutil.MakePtr(uint16(2)),
			NumberOfReplicas: opensearchutil.MakePtr(uint16(1)),
		},
		opensearchutil.WithStrictMapping(true),
	)
	if err != nil {
		fmt.Printf("GenerateIndexJson: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== OpenSearch Index JSON ===")
	fmt.Printf("%s\n\n", string(indexJson))

	// 3. Generate just the mappings JSON
	mappingsJson, err := indexGenerator.GenerateMappingsJson(mappingProperties)
	if err != nil {
		fmt.Printf("GenerateMappingsJson: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== Mappings Only ===")
	fmt.Printf("%s\n\n", string(mappingsJson))

	// 4. Demonstrate JSON marshalling of custom date types
	documentJson, err := json.MarshalIndent(&samplePerson, "", "  ")
	if err != nil {
		fmt.Printf("json.MarshalIndent: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== Document JSON (for indexing) ===")
	fmt.Printf("%s\n", string(documentJson))
}
```

Resulting JSON in the "OpenSearch Index JSON" from the code above (note it's also possible to only get the mappings, to update an existing index):
```json
{
   "mappings": {
      "dynamic": "strict",
      "properties": {
         "account_balance": {
            "type": "float"
         },
         "addresses": {
            "properties": {
               "confirmed": {
                  "type": "boolean"
               },
               "full_address": {
                  "type": "text"
               }
            }
         },
         "age": {
            "type": "integer"
         },
         "aliases": {
            "type": "text"
         },
         "all_text": {
            "type": "text"
         },
         "created_at": {
            "format": "basic_date_time",
            "type": "date"
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
         "last_login": {
            "format": "basic_date",
            "type": "date"
         },
         "name": {
            "type": "text"
         },
         "searchable_text": {
            "type": "text"
         },
         "social_security": {
            "index_prefixes": {
               "max_chars": "5",
               "min_chars": "3"
            },
            "type": "text"
         },
         "title": {
            "analyzer": "standard",
            "copy_to": [
               "all_text",
               "searchable_text"
            ],
            "search_analyzer": "english",
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
      "number_of_replicas": 1,
      "number_of_shards": 2
   }
}
```

The resulting JSON can be used directly with the [OpenSearch Create Index API](https://opensearch.org/docs/1.0/opensearch/rest-api/create-index/).