# OpenSearchUtil

Given an object, makes an OpenSearch index template.

## Example

This is from a test case:

```go
var _ = Describe("GenerateIndexJson", func() {
	It("Generates an index JSON string", func() {
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
		str, err := GenerateIndexJson(person{})
		Expect(err).To(BeNil())
		Expect(str).To(Equal(`{
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
}`))
	})
})
```