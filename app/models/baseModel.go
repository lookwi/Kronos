package models

import (
	"fmt"
	"strings"
	"time"
)

type BaseModel struct {
	ID        uint64     `gorm:"primary_key; index:id;" json:"id" structs:"id"`
	CreatedAt time.Time  `json:"createdAt" structs:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt" structs:"updatedAt"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt" structs:"deletedAt"`
}

type NullType byte

const (
	_ NullType = iota
	// IsNull the same as `is null`
	IsNull
	// IsNotNull the same as `is not null`
	IsNotNull
)

// sql build where
func WhereBuild(where map[string]interface{}) (whereSQL string, vals []interface{}, err error) {
	vals = make([]interface{}, 0)
	for k, v := range where {
		ks := strings.Split(k, " ")
		if len(ks) > 2 {
			return "", nil, fmt.Errorf("Error in query condition: %s. ", k)
		}

		if whereSQL != "" {
			whereSQL += " AND "
		}

		//fmt.Println(strings.Join(ks, ","))
		switch len(ks) {
		case 1:
			//fmt.Println(reflect.TypeOf(v))
			switch v := v.(type) {
			case NullType:
				if v == IsNotNull {
					whereSQL += fmt.Sprint(k, " IS NOT NULL")
				} else {
					whereSQL += fmt.Sprint(k, " IS NULL")
				}
			default:
				whereSQL += fmt.Sprint(k, "=?")
				vals = append(vals, v)
			}

		case 2:
			k = ks[0]
			switch ks[1] {
			case "=":
				whereSQL += fmt.Sprint(k, "=?")
				vals = append(vals, v)
			case ">":
				whereSQL += fmt.Sprint(k, ">?")
				vals = append(vals, v)

			case ">=":
				whereSQL += fmt.Sprint(k, ">=?")
				vals = append(vals, v)

			case "<":
				whereSQL += fmt.Sprint(k, "<?")
				vals = append(vals, v)

			case "<=":
				whereSQL += fmt.Sprint(k, "<=?")
				vals = append(vals, v)

			case "!=":
				whereSQL += fmt.Sprint(k, "!=?")
				vals = append(vals, v)

			case "<>":
				whereSQL += fmt.Sprint(k, "!=?")
				vals = append(vals, v)

			case "in":
				whereSQL += fmt.Sprint(k, " in (?) ")
				vals = append(vals, v)

			case "like":
				whereSQL += fmt.Sprint(k, " like ? ")
				vals = append(vals, v)
			}

		}
	}
	return
}
