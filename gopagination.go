package gopagination

import (
	"errors"
	"reflect"
	"strings"
	"time"
)

// Equal = 0,
// LessThan = 1,
// LessThanOrEqual = 2,
// GreaterThan = 3,
// GreaterThanOrEqual = 4,
// NotEqual = 5,
// Contains = 6,
// StartsWith = 7,
// EndsWith = 8

//ExpressionFilter is struct for add filtering in slice or array
type ExpressionFilter struct {
	PropertyName string      `json:"propertyname" form:"propertyname" query:"propertyname"`
	Value        interface{} `json:"value" form:"value" query:"value"`
	Comparison   int         `json:"comparison" form:"comparison" query:"comparison"`
}

//PaginationConfig is struct for add pagination in slice or array
type PaginationConfig struct {
	SortColumn        string             `json:"sortcolumn" form:"sortcolumn" query:"sortcolumn"`
	IsAscending       bool               `json:"isascending" form:"isascending" query:"isascending"`
	Take              int                `json:"take" form:"take" query:"take"`
	Skip              int                `json:"skip" form:"skip" query:"skip"`
	AndLogic          bool               `json:"andlogic" form:"andlogic" query:"andlogic"`
	ExpressionFilters []ExpressionFilter `json:"expressionfilters" form:"expressionfilters" query:"expressionfilters"`
}

type filters struct {
	PropertyName string
	Value        interface{}
	Comparison   int
	Type         reflect.Type
}

//DoPagination is function of PaginationConfig struct for start pagination in slice or array - this return an slice and error
func (dc PaginationConfig) DoPagination(datas interface{}) (interface{}, error) {
	v := reflect.ValueOf(datas)
	if !(v.Kind() == reflect.Slice || v.Kind() == reflect.Array) {
		return nil, errors.New("Only Array or Slice")
	}

	if v.Len() <= 0 {
		return nil, errors.New("Empty data")
	}

	var FilterTypes []filters
	for _, filter := range dc.ExpressionFilters {
		fieldtype, err := getexiststype(datas, filter.PropertyName)
		if err == nil {

			filter := filters{
				Comparison:   filter.Comparison,
				PropertyName: filter.PropertyName,
				Type:         fieldtype,
				Value:        filter.Value,
			}

			FilterTypes = append(FilterTypes, filter)
		}
	}
	if dc.SortColumn != "" {
		sortdata(datas, dc.SortColumn, dc.IsAscending)

	}

	if len(FilterTypes) > 0 {
		if dc.AndLogic {
			filterdatas, err := dofilter(datas, FilterTypes, true)
			if err != nil {
				return nil, err
			}
			datas = filterdatas
		} else {
			filterdatas, err := dofilter(datas, FilterTypes, false)
			if err != nil {
				return nil, err
			}
			datas = filterdatas
		}
	}
	if dc.Take != 0 {
		datas = takecountofdatas(datas, dc.Skip, dc.Take)
	}

	return datas, nil
}

func dofilter(datas interface{}, FilterTypes []filters, IsAnd bool) (interface{}, error) {

	datavalues := reflect.ValueOf(datas)
	var newDatas []interface{}
	for i := 0; i < datavalues.Len(); i++ {
		res, err := calculatefilter(datavalues.Index(i), FilterTypes, IsAnd)
		if err != nil {
			return nil, err
		}
		if res {
			newDatas = append(newDatas, datavalues.Index(i).Interface())
		}
	}
	return newDatas, nil
}

func getexiststype(data interface{}, name string) (reflect.Type, error) {

	s := reflect.TypeOf(data)
	val := s.Elem()
	var Type reflect.Type

	for i := 0; i < val.NumField(); i++ {

		if strings.ToLower(val.Field(i).Name) == strings.ToLower(name) {
			Type = val.Field(i).Type
			return Type, nil
		}
	}
	return Type, errors.New("No field found")
}

func getdataoftype(data reflect.Value, name string) reflect.Value {

	Value := data.FieldByNameFunc(func(c string) bool {
		return strings.ToLower(c) == strings.ToLower(name)
	})
	return Value
}

func calculatefilter(data interface{}, FilterTypes []filters, IsAnd bool) (bool, error) {

	for _, filter := range FilterTypes {
		if val, ok := data.(reflect.Value); ok {
			value := getdataoftype(val, filter.PropertyName)
			if strings.Contains(strings.ToLower(filter.Type.Name()), "int") {
				if filter.Comparison == 6 || filter.Comparison == 7 || filter.Comparison == 8 {
					return false, errors.New("Invalid Comparison")
				}
				datavalue := float64(value.Int())
				filtervalue := reflect.ValueOf(filter.Value).Interface().(float64)
				switch filter.Comparison {
				case 0:
					if datavalue == filtervalue {
						if !IsAnd {
							return true, nil
						}
					} else {
						if IsAnd {
							return false, nil
						}
					}
				case 1:
					if datavalue < filtervalue {
						if !IsAnd {
							return true, nil
						}
					} else {
						if IsAnd {
							return false, nil
						}
					}
				case 2:
					if datavalue <= filtervalue {
						if !IsAnd {
							return true, nil
						}
					} else {
						if IsAnd {
							return false, nil
						}
					}
				case 3:
					if datavalue > filtervalue {
						if !IsAnd {
							return true, nil
						}
					} else {
						if IsAnd {
							return true, nil
						}
					}
				case 4:
					if datavalue >= filtervalue {
						if !IsAnd {
							return true, nil
						}
					} else {
						if IsAnd {
							return false, nil
						}
					}
				case 5:
					if datavalue != filtervalue {
						if !IsAnd {
							return true, nil
						}
					} else {
						if IsAnd {
							return false, nil
						}
					}
				}
			}
			if strings.Contains(strings.ToLower(filter.Type.Name()), "string") {
				if filter.Comparison == 1 || filter.Comparison == 2 || filter.Comparison == 3 || filter.Comparison == 4 {
					return false, errors.New("Invalid Comparison")
				}
				datavalue := value.String()
				filtervalue := reflect.ValueOf(filter.Value).String()
				switch filter.Comparison {
				case 0:
					if datavalue == filtervalue {
						if !IsAnd {
							return true, nil
						}
					} else {
						if IsAnd {
							return false, nil
						}
					}
				case 5:
					if datavalue != filtervalue {
						if !IsAnd {
							return true, nil
						}
					} else {
						if IsAnd {
							return false, nil
						}
					}
				case 6:
					if strings.Contains(datavalue, filtervalue) {
						if !IsAnd {
							return true, nil
						}
					} else {
						if IsAnd {
							return false, nil
						}
					}
				case 7:
					if strings.HasPrefix(datavalue, filtervalue) {
						if !IsAnd {
							return true, nil
						}
					} else {
						if IsAnd {
							return false, nil
						}
					}
				case 8:
					if strings.HasSuffix(datavalue, filtervalue) {
						if !IsAnd {
							return true, nil
						}
					} else {
						if IsAnd {
							return false, nil
						}
					}
				}

			}
			if strings.Contains(strings.ToLower(filter.Type.Name()), "float32") {
				if filter.Comparison == 6 || filter.Comparison == 7 || filter.Comparison == 8 {
					return false, errors.New("Invalid Comparison")
				}
				datavalue := float64(value.Interface().(float32))
				filtervalue := float64(reflect.ValueOf(filter.Value).Interface().(float32))
				switch filter.Comparison {
				case 0:
					if datavalue == filtervalue {
						if !IsAnd {
							return true, nil
						}
					} else {
						if IsAnd {
							return false, nil
						}
					}
				case 1:
					if datavalue < filtervalue {
						if !IsAnd {
							return true, nil
						}
					} else {
						if IsAnd {
							return false, nil
						}
					}
				case 2:
					if datavalue <= filtervalue {
						if !IsAnd {
							return true, nil
						}
					} else {
						if IsAnd {
							return false, nil
						}
					}
				case 3:
					if datavalue > filtervalue {
						if !IsAnd {
							return true, nil
						}
					} else {
						if IsAnd {
							return false, nil
						}
					}
				case 4:
					if datavalue >= filtervalue {
						if !IsAnd {
							return true, nil
						}
					} else {
						if IsAnd {
							return false, nil
						}
					}
				case 5:
					if datavalue != filtervalue {
						if !IsAnd {
							return true, nil
						}
					} else {
						if IsAnd {
							return false, nil
						}
					}
				}
			}
			if strings.Contains(strings.ToLower(filter.Type.Name()), "float64") {
				if filter.Comparison == 6 || filter.Comparison == 7 || filter.Comparison == 8 {
					return false, errors.New("Invalid Comparison")
				}
				datavalue := value.Interface().(float64)
				filtervalue := reflect.ValueOf(filter.Value).Interface().(float64)
				switch filter.Comparison {
				case 0:
					if datavalue == filtervalue {
						if !IsAnd {
							return true, nil
						}
					} else {
						if IsAnd {
							return false, nil
						}
					}
				case 1:
					if datavalue < filtervalue {
						if !IsAnd {
							return true, nil
						}
					} else {
						if IsAnd {
							return false, nil
						}
					}
				case 2:
					if datavalue <= filtervalue {
						if !IsAnd {
							return true, nil
						}
					} else {
						if IsAnd {
							return false, nil
						}
					}
				case 3:
					if datavalue > filtervalue {
						if !IsAnd {
							return true, nil
						}
					} else {
						if IsAnd {
							return false, nil
						}
					}
				case 4:
					if datavalue >= filtervalue {
						if !IsAnd {
							return true, nil
						}
					} else {
						if IsAnd {
							return false, nil
						}
					}
				case 5:
					if datavalue != filtervalue {
						if !IsAnd {
							return true, nil
						}
					} else {
						if IsAnd {
							return false, nil
						}
					}
				}
			}
			if strings.Contains(strings.ToLower(filter.Type.Name()), "time") {
				if filter.Comparison == 6 || filter.Comparison == 7 || filter.Comparison == 8 {
					return false, errors.New("Invalid Comparison")
				}
				datavalue := value.Interface().(time.Time)
				filtervalue := reflect.ValueOf(filter.Value).Interface().(time.Time)
				switch filter.Comparison {
				case 0:
					if datavalue.Equal(filtervalue) {
						if !IsAnd {
							return true, nil
						}
					} else {
						if IsAnd {
							return false, nil
						}
					}
				case 1:
					if datavalue.Before(filtervalue) {
						if !IsAnd {
							return true, nil
						}
					} else {
						if IsAnd {
							return false, nil
						}
					}
				case 2:
					if datavalue.Before(filtervalue) || datavalue.Equal(filtervalue) {
						if !IsAnd {
							return true, nil
						}
					} else {
						if IsAnd {
							return false, nil
						}
					}
				case 3:
					if datavalue.After(filtervalue) {
						if !IsAnd {
							return true, nil
						}
					} else {
						if IsAnd {
							return false, nil
						}
					}
				case 4:
					if datavalue.After(filtervalue) || datavalue.Equal(filtervalue) {
						if !IsAnd {
							return true, nil
						}
					} else {
						if IsAnd {
							return false, nil
						}
					}
				case 5:
					if !datavalue.Equal(filtervalue) {
						if !IsAnd {
							return true, nil
						}
					} else {
						if IsAnd {
							return false, nil
						}
					}
				}
			}
		}

	}
	if IsAnd {
		return true, nil
	}
	return false, nil

}

func takecountofdatas(datas interface{}, skip, take int) interface{} {
	var newdatas []interface{}
	if take == 0 {
		return datas
	}
	v := reflect.ValueOf(datas)
	if v.Len() <= skip {
		return nil
	}
	for i := 0; i < v.Len(); i++ {
		if i >= skip && i < skip+take {
			newdatas = append(newdatas, v.Index(i).Interface())
		}
	}
	return newdatas
}

func sortdata(datas interface{}, ColumnName string, IsAscending bool) (interface{}, error) {
	v := reflect.ValueOf(datas)

	datatype, err := getexiststype(datas, ColumnName)
	if err != nil {
		return nil, err
	}

	end := v.Len()

	for {

		if end == 0 {
			break
		}

		for i := 0; i < v.Len()-1; i++ {

			if strings.Contains(strings.ToLower(datatype.Name()), "int") {

				if (getdataoftype(v.Index(i), ColumnName).Int() < getdataoftype(v.Index(i+1), ColumnName).Int() &&
					!IsAscending) || (getdataoftype(v.Index(i), ColumnName).Int() > getdataoftype(v.Index(i+1), ColumnName).Int() && IsAscending) {
					intermediate := v.Index(i).Interface()
					v.Index(i).Set(v.Index(i + 1))
					v.Index(i + 1).Set(reflect.ValueOf(intermediate))
				}
			} else if strings.Contains(strings.ToLower(datatype.Name()), "float32") {

				if (getdataoftype(v.Index(i), ColumnName).Interface().(float32) < getdataoftype(v.Index(i+1), ColumnName).Interface().(float32) &&
					!IsAscending) || (getdataoftype(v.Index(i), ColumnName).Interface().(float32) > getdataoftype(v.Index(i+1), ColumnName).Interface().(float32) && IsAscending) {
					intermediate := v.Index(i).Interface()
					v.Index(i).Set(v.Index(i + 1))
					v.Index(i + 1).Set(reflect.ValueOf(intermediate))
				}
			} else if strings.Contains(strings.ToLower(datatype.Name()), "float64") {

				if (getdataoftype(v.Index(i), ColumnName).Interface().(float64) < getdataoftype(v.Index(i+1), ColumnName).Interface().(float64) &&
					!IsAscending) || (getdataoftype(v.Index(i), ColumnName).Interface().(float64) > getdataoftype(v.Index(i+1), ColumnName).Interface().(float64) && IsAscending) {
					intermediate := v.Index(i).Interface()
					v.Index(i).Set(v.Index(i + 1))
					v.Index(i + 1).Set(reflect.ValueOf(intermediate))
				}
			} else if strings.Contains(strings.ToLower(datatype.Name()), "string") {

				if (getdataoftype(v.Index(i), ColumnName).String() < getdataoftype(v.Index(i+1), ColumnName).String() &&
					!IsAscending) || (getdataoftype(v.Index(i), ColumnName).String() > getdataoftype(v.Index(i+1), ColumnName).String() && IsAscending) {
					intermediate := v.Index(i).Interface()
					v.Index(i).Set(v.Index(i + 1))
					v.Index(i + 1).Set(reflect.ValueOf(intermediate))
				}
			} else if strings.Contains(strings.ToLower(datatype.Name()), "time") {

				if (getdataoftype(v.Index(i), ColumnName).Interface().(time.Time).Before(getdataoftype(v.Index(i+1), ColumnName).Interface().(time.Time)) && !IsAscending) ||
					(getdataoftype(v.Index(i), ColumnName).Interface().(time.Time).After(getdataoftype(v.Index(i+1), ColumnName).Interface().(time.Time)) && IsAscending) {
					intermediate := v.Index(i).Interface()
					v.Index(i).Set(v.Index(i + 1))
					v.Index(i + 1).Set(reflect.ValueOf(intermediate))
				}
			}
		}

		end--

	}

	return datas, nil
}
