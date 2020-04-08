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
go get github.com/kohestanimahdi/go-pagination-filtering
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
type ExpressionFilter struct {
	PropertyName string      // The column name of struct that we want to filter on this
	Value        interface{} // Value of that column name
	Comparison   int         // Comparison
}
```
Comparison is one of :
<br>

 0 for Equal (=) 				// This use for numeric,time and string fields
 <br>
 1 for LessThan (<)				// This use for numeric and time fields
 <br>
 2 for LessThanOrEqual (<=)		// This use for numeric and time fields
 <br>
 3 for GreaterThan (>)			// This use for numeric and time fields
 <br>
 4 for GreaterThanOrEqual (>=)	// This use for numeric and time fields
 <br>
 5 for NotEqual (!=)			// This use for numeric,time and string fields
 <br>
 6 for Contains 				// This use for string fields
 <br>
 7 for StartsWith 				// This use for string fields
 <br>
 8 for EndsWith 				// This use for string fields

 <br>
### Examples

Imagine we have a struct named `Person` like this :
```go
type Person struct {
	Name 	  string
	Average   float32
	Birthday  time.Time
}
```

And a slice of Persons name 'persons' like :

```go
var persons []Person
	p1 := Person{
		Name:     "Mahdi Kohestani",
		Average:  19.5,
		Birthday: time.Date(1997, time.December, 25, 0, 0, 0, 0, time.UTC),
	}
	p2 := Person{
		Name:     "Mahdi Malverdi",
		Average:  18,
		Birthday: time.Date(1997, time.November, 3, 0, 0, 0, 0, time.UTC),
	}
	p3 := Person{
		Name:     "Amir Sartipi",
		Average:  19,
		Birthday: time.Date(1997, time.September, 28, 0, 0, 0, 0, time.UTC),
	}
	p4 := Person{
		Name:     "Mahdi Rangraz",
		Average:  12,
		Birthday: time.Date(1998, time.February, 10, 0, 0, 0, 0, time.UTC),
	}
	p5 := Person{
		Name:     "Mohammad Saleh",
		Average:  13,
		Birthday: time.Date(1998, time.June, 16, 0, 0, 0, 0, time.UTC),
	}
	persons = append(persons, p1, p2, p3, p4, p5)
```

If we want all datas and sort ascending  by `Name` column

```go
pg := pagination.PaginationConfig{
	SortColumn:        "Name", //Sort by `Name` column
	IsAscending:       true,   // Ascending
	Take:              0,      // Get all datas
	Skip:              0,      // skip zero
	AndLogic:          true,   // And between filters
	ExpressionFilters: nil,	   // We do not need filter
}
newdatas,err :=pg.DoPagination(persons)
if err!= nil{
	panic(err)
}
fmt.Println(newdatas)
```

If we want all datas that sort ascending  by `Name` column and `Average` grather than 15
```go
pg := pagination.PaginationConfig{
	SortColumn:        "Name", //Sort by `Name` column
	IsAscending:       true,   // Ascending
	Take:              0,      // Get all datas
	Skip:              0,      // skip zero
	AndLogic:          true,   // And between filters
	ExpressionFilters: make([]pagination.ExpressionFilter, 0),
}
pg.ExpressionFilters = append(pg.ExpressionFilters, pagination.ExpressionFilter{
	PropertyName: "Average",
	Value:        15.0,
	Comparison:   3,
})
newdatas,err :=pg.DoPagination(persons)
if err!= nil{
	panic(err)
}
fmt.Println(newdatas)
```

If we want datas that sort ascending  by `Name` column and `Average` grather than 15 or `Name` contains 'Mahdi'
```go
pg := pagination.PaginationConfig{
	SortColumn:        "Name", //Sort by `Name` column
	IsAscending:       true,   // Ascending
	Take:              0,      // Get all datas
	Skip:              0,      // skip zero
	AndLogic:          false,   // Or between filters
	ExpressionFilters: make([]pagination.ExpressionFilter, 0),
}
pg.ExpressionFilters = append(pg.ExpressionFilters, pagination.ExpressionFilter{
	PropertyName: "Average",
	Value:        15.0,
	Comparison:   3,
})
pg.ExpressionFilters = append(pg.ExpressionFilters, pagination.ExpressionFilter{
	PropertyName: "Name",
	Value:        "Mahdi",
	Comparison:   6,
})
newdatas,err :=pg.DoPagination(persons)
if err!= nil{
	panic(err)
}
fmt.Println(newdatas)
```

If we want only 2 data  that sort ascending  by `Name` column and `Average` grather than 15 or `Name` contains 'Mahdi'
```go
pg := pagination.PaginationConfig{
	SortColumn:        "Name", //Sort by `Name` column
	IsAscending:       true,   // Ascending
	Take:              2,      // Get 2 data
	Skip:              0,      // skip zero
	AndLogic:          false,  // Or between filters
	ExpressionFilters: make([]pagination.ExpressionFilter, 0),
}
pg.ExpressionFilters = append(pg.ExpressionFilters, pagination.ExpressionFilter{
	PropertyName: "Average",
	Value:        15.0,
	Comparison:   3,
})
pg.ExpressionFilters = append(pg.ExpressionFilters, pagination.ExpressionFilter{
	PropertyName: "Name",
	Value:        "Mahdi",
	Comparison:   6,
})
newdatas,err :=pg.DoPagination(persons)
if err!= nil{
	panic(err)
}
fmt.Println(newdatas)
```

If we want only second and third data (skip 1) that sort ascending  by `Name` column and `Average` grather than 15 or `Name` contains 'Mahdi'
```go
pg := pagination.PaginationConfig{
	SortColumn:        "Name", //Sort by `Name` column
	IsAscending:       true,   // Ascending
	Take:              2,      // Get 2 data
	Skip:              1,      // skip zero
	AndLogic:          false,  // Or between filters
	ExpressionFilters: make([]pagination.ExpressionFilter, 0),
}
pg.ExpressionFilters = append(pg.ExpressionFilters, pagination.ExpressionFilter{
	PropertyName: "Average",
	Value:        15.0,
	Comparison:   3,
})
pg.ExpressionFilters = append(pg.ExpressionFilters, pagination.ExpressionFilter{
	PropertyName: "Name",
	Value:        "Mahdi",
	Comparison:   6,
})
newdatas,err :=pg.DoPagination(persons)
if err!= nil{
	panic(err)
}
fmt.Println(newdatas)
```
