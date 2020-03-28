# go pagination and filtering
This is a package for pagination and filtering all slice and arrays of struct

## Overview
Package gopagination is a package that you can filter your <b>slice</b> and <b>Array</b> of any struct.
This can help you to filter and pagin any data,
for example you can use this package to fliter or pagination datas in web api.
Clinet send us what data needs and this package send response that client expected.

## Examples

### Installation

For install this package you should execute this 

```go
go get -u github.com/kohestanimahdi/go-pagination-filtering
```

### Structs

#### PaginationConfig

The main struct of pagination, this has six property :

```go
type PaginationConfig struct {
	SortColumn        string                // The column name of struct that we want to sort by them
	IsAscending       bool                  // Ascending sort or Descending  ?
	Take              int                   // The count of data that we need
	Skip              int                	// The count of data that we want to skip 
	AndLogic          bool               	// And(&) between ExpressionFilters ?
	ExpressionFilters []ExpressionFilter 	// Filters - If we need filters
}
```

#### ExpressionFilter

The struct of Filters, this has three property :

```go
	PropertyName string      // The column name of struct that we want to filter on this
	Value        interface{} // Value of that column name
	Comparison   int         // Comparison
```
Comparison is one of :
<b> 
 0 for Equal (=) 				// This use for numeric,time and string fields
 1 for LessThan (<)				// This use for numeric and time fields
 2 for LessThanOrEqual (<=)		// This use for numeric and time fields
 3 for GreaterThan (>)			// This use for numeric and time fields
 4 for GreaterThanOrEqual (>=)	// This use for numeric and time fields
 5 for NotEqual (!=)			// This use for numeric,time and string fields
 6 for Contains 				// This use for string fields
 7 for StartsWith 				// This use for string fields
 8 for EndsWith 				// This use for string fields
 </b>
 
### Examples
