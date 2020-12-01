package main

import "strings"

// Re-implementation of the type parsing logic from the driver.
// Unlike the driver, this implementation differentiates frozen types
// from non-frozen ones.

type Type int

const (
	TypeCustom    Type = 0x0000
	TypeAscii     Type = 0x0001
	TypeBigInt    Type = 0x0002
	TypeBlob      Type = 0x0003
	TypeBoolean   Type = 0x0004
	TypeCounter   Type = 0x0005
	TypeDecimal   Type = 0x0006
	TypeDouble    Type = 0x0007
	TypeFloat     Type = 0x0008
	TypeInt       Type = 0x0009
	TypeText      Type = 0x000A
	TypeTimestamp Type = 0x000B
	TypeUUID      Type = 0x000C
	TypeVarchar   Type = 0x000D
	TypeVarint    Type = 0x000E
	TypeTimeUUID  Type = 0x000F
	TypeInet      Type = 0x0010
	TypeDate      Type = 0x0011
	TypeTime      Type = 0x0012
	TypeSmallInt  Type = 0x0013
	TypeTinyInt   Type = 0x0014
	TypeDuration  Type = 0x0015
	TypeList      Type = 0x0020
	TypeMap       Type = 0x0021
	TypeSet       Type = 0x0022
	TypeUDT       Type = 0x0030
	TypeTuple     Type = 0x0031
)

func (t Type) IsCollection() bool {
	switch t {
	case TypeList, TypeMap, TypeSet:
		return true
	default:
		return false
	}
}

type TypeInfo interface {
	Type() Type
	IsFrozen() bool
}

type FrozenType struct {
	Inner TypeInfo
}

func (ft *FrozenType) Type() Type {
	return ft.Inner.Type()
}

func (ft *FrozenType) IsFrozen() bool {
	return true
}

type MapType struct {
	Key   TypeInfo
	Value TypeInfo
}

func (mt *MapType) Type() Type {
	return TypeMap
}

func (mt *MapType) IsFrozen() bool {
	return false
}

type ListType struct {
	Element TypeInfo
}

func (lt *ListType) Type() Type {
	return TypeList
}

func (lt *ListType) IsFrozen() bool {
	return false
}

type SetType struct {
	Element TypeInfo
}

func (st *SetType) Type() Type {
	return TypeSet
}

func (st *SetType) IsFrozen() bool {
	return false
}

type TupleType struct {
	Elements []TypeInfo
}

func (tt *TupleType) Type() Type {
	return TypeTuple
}

func (tt *TupleType) IsFrozen() bool {
	return false
}

type NativeType struct {
	RealType Type
}

func (nt *NativeType) Type() Type {
	return nt.RealType
}

func (nt *NativeType) IsFrozen() bool {
	return false
}

func parseType(str string) TypeInfo {
	if strings.HasPrefix(str, "frozen<") {
		innerStr := strings.TrimSuffix(strings.TrimPrefix(str, "frozen<"), ">")
		return &FrozenType{parseType(innerStr)}
	}
	if strings.HasPrefix(str, "list<") {
		innerStr := strings.TrimSuffix(strings.TrimPrefix(str, "list<"), ">")
		return &ListType{parseType(innerStr)}
	}
	if strings.HasPrefix(str, "set<") {
		innerStr := strings.TrimSuffix(strings.TrimPrefix(str, "set<"), ">")
		return &SetType{parseType(innerStr)}
	}
	if strings.HasPrefix(str, "map<") {
		innerStr := strings.TrimSuffix(strings.TrimPrefix(str, "map<"), ">")
		list := parseTypeList(innerStr)
		return &MapType{Key: list[0], Value: list[1]}
	}
	if strings.HasPrefix(str, "tuple<") {
		innerStr := strings.TrimSuffix(strings.TrimPrefix(str, "tuple<"), ">")
		list := parseTypeList(innerStr)
		return &TupleType{Elements: list}
	}
	return &NativeType{RealType: parseNativeType(str)}
}

func parseTypeList(str string) []TypeInfo {
	var ret []TypeInfo
	var level int
	var builder strings.Builder
	for _, r := range str {
		if r == ',' && level == 0 {
			s := strings.TrimSpace(builder.String())
			ret = append(ret, parseType(s))
			builder.Reset()
			continue
		}

		if r == '<' {
			level++
		} else if r == '>' {
			level--
		}
		builder.WriteRune(r)
	}
	if builder.Len() != 0 {
		s := strings.TrimSpace(builder.String())
		ret = append(ret, parseType(s))
	}
	return ret
}

func parseNativeType(str string) Type {
	switch str {
	case "ascii":
		return TypeAscii
	case "bigint":
		return TypeBigInt
	case "blob":
		return TypeBlob
	case "boolean":
		return TypeBoolean
	case "counter":
		return TypeCounter
	case "date":
		return TypeDate
	case "decimal":
		return TypeDecimal
	case "double":
		return TypeDouble
	case "duration":
		return TypeDuration
	case "float":
		return TypeFloat
	case "int":
		return TypeInt
	case "smallint":
		return TypeSmallInt
	case "tinyint":
		return TypeTinyInt
	case "time":
		return TypeTime
	case "timestamp":
		return TypeTimestamp
	case "uuid":
		return TypeUUID
	case "varchar":
		return TypeVarchar
	case "text":
		return TypeText
	case "varint":
		return TypeVarint
	case "timeuuid":
		return TypeTimeUUID
	case "inet":
		return TypeInet
	default:
		return TypeCustom
	}
}