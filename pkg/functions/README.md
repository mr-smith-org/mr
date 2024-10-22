# Kuma Functions

This document provides an overview of the functions available in this package. These functions are designed to extract and manipulate data from OpenAPI 2.0 files and convert data into YAML format for use in Go templates.

## Table of Contents

- [Parser Functions](#parser-functions)
  - [ToYaml](#toyaml)
- [Group Functions](#group-functions)
  - [GroupByKey](#groupbykey)
- [OpenAPI Functions](#openapi-functions)
  - [GetRefFrom](#getreffrom)
  - [GetPathsByTag](#getpathsbytag)
  - [GetParamsByType](#getparamsbytype)

## Functions

### Parser Functions

#### ToYaml

ToYaml converts a Go data structure into a YAML-formatted slice of strings, where each element represents a line of the resulting YAML. This function is useful for rendering data in a human-readable format for templating purposes.

**Signature:**

```go
func ToYaml(data interface{}) []string
```

**Parameters:**

- `data`: The data structure to be converted into YAML.

**Returns:**

A slice of strings containing the YAML representation of data.

**Example**

```yaml
# input:
#    {
#         "items":  {
#             "subItems": [ "item1",   "item2",  "item3"],
#         },
#   }
data:
   {{ range toYaml . }}{{ . }}{{end}}
#output:
# data:
#   items:
#     subItems:
#         - item1
#         - item2
#         - item3
```

---

### Group Functions

#### GroupByKey

GroupByKey organizes a slice of maps based on a specified key. It returns a map where the keys are unique values found at the specified key in the input items, and the values are slices of items that share that key value. This function is useful for categorizing data in Go applications.

**Signature:**

```go
func GroupByKey(data []interface{}, key string) map[string]interface{
```

**Parameters:**

- `data`: A slice of maps where each map represents an item with key-value pairs.
- `key`: The key used to group items in the data slice.

**Returns:**

A map where each key represents a unique value from the specified key in the data slice, and the corresponding value is a slice of items that have that key value.

**Example**

```yaml
# input:
# data:
#   - {"name": "Alice", "department": "HR"}
#   - {"name": "Bob", "department": "IT"}
#   - {"name": "Charlie", "department": "HR"}
#   - {"name": "Dave", "department": "IT"}
# key: "department"

"{{ groupByKey .data .key }}"
# output:
# {
#   "HR": [
#     {"name": "Alice", "department": "HR"},
#     {"name": "Charlie", "department": "HR"}
#   ],
#   "IT": [
#     {"name": "Bob", "department": "IT"},
#     {"name": "Dave", "department": "IT"}
#   ]
# }
```

---

### OpenAPI Functions

#### GetRefFrom

GetRefFrom extracts the reference identifier from an OpenAPI 2.0 object if it exists. The function expects the reference to be in the format of a JSON pointer within the OpenAPI specification.

**Signature:**

```go
func GetRefFrom(object map[string]interface{}) string
```

**Parameters:**

- `object`: A map representing an OpenAPI object which might contain a $ref field.

**Returns:**

A string containing the reference identifier or an empty string if no valid reference is found.

**Example**

```yaml
# input:
# {
#     "schema": {
#         "$ref": "#/definitions/ApiResponse"
#     }
# }
"{{getRefFrom  .schema}}"
# output:
# "ApiResponse"
```

---

#### GetPathsByTag

GetPathsByTag filters OpenAPI paths by a specified tag. It returns a subset of paths that are associated with the provided tag, useful for generating documentation for specific sections of an API.

**Signature:**

```go
func GetPathsByTag(paths map[string]interface{},
tag string) map[string]interface{}
```

**Parameters:**

- `paths`: A map where keys are path names and values are path item objects.
- `tag`: The tag used to filter paths.

**Returns:**

A map containing the paths that include the specified tag.

**Example**

```yaml
# input:
# paths: {
#   "/pet/{petId}/uploadImage": {
#       "post": {
#         "tags": ["pet"],
#         "summary": "uploads an image",
#       }
#       ...
#   },
#   "/user": {
#       "post": {
#         "tags": ["user"],
#         "summary": "create a user",
#       },
#   },
#   "/user/{userId}": {
#       "put": {
#         "tags": ["user"],
#         "summary": "update user by id",
#       },
#       "get": {
#         "tags": ["user"],
#         "summary": "get user by id",
#       }
#       ...
#  }
#}
"{{ getPathsByTag .paths "user" }}"
# output:
# {
#   "/user": {
#       "post": {
#         "tags": ["user"],
#         "summary": "create a user",
#       },
#   },
#   "/user/{userId}": {
#       "put": {
#         "tags": ["user"],
#         "summary": "update user by id",
#       },
#       "get": {
#         "tags": ["user"],
#         "summary": "get user by id",
#       }
#       ...
#  }
#}
```

---

#### GetParamsByType

GetParamsByType filters parameters based on their `in` field type (e.g., query, header, path, formData). This function helps in extracting parameters of a certain type from an OpenAPI operation.

**Signature:**

```go
func GetParamsByType(params []interface{},
paramType string) []interface{}
```

**Parameters:**

- `params`: A slice of parameter objects.
- `paramType`: The type of parameters to filter by (e.g., "query").

**Returns:**

A slice of parameters that match the specified type.

**Example**

```yaml
# input:
#     {
#       "parameters": [
#             {
#                 "name": "petId",
#                 "in": "path",
#             },
#             {
#                 "name": "body",
#                 "in": "body",
#             },
#             {
#                 "name": "userId",
#                 "in": "path,
#             }
#     ],
# }
"{{ getParamsByType $data.parameters "path" }}"
# output:
# [
#     {
#         "name": "userId"
#     },
#      {
#        "name": "petId"
#      }
#  ]
```

---
