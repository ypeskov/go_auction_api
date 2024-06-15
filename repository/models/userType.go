package models

type UserType struct {
	Id              int    `json:"id"`
	TypeName        string `json:"typeName" db:"type_name"`
	TypeDescription string `json:"typeDescription" db:"type_description"`
	TypeCode        string `json:"typeCode" db:"type_code"`
}
