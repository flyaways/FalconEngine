package tools

import (
	"fmt"
	"sync"
)

type FalconSearchEncoder interface {
	FalconEncoding() ([]byte, error)
	ToString() string
}

type FalconSearchDecoder interface {
	FalconDecoding(bytes []byte) error
	ToString() string
	// FalconPrepare(length int64) (int64,error)
}

type FalconCoder interface {
	FalconEncoding() ([]byte, error)
	FalconDecoding(bytes []byte) error
	ToString() string
}

//
//type DictValue struct {
//	Val uint64
//	ExtVal uint64
//}
//
//func NewDicValue() *DictValue{
//	return &DictValue{}
//}
//
//func (dv *DictValue) FalconEncoding() ([]byte,error) {
//	b:=make([]byte,24)
//	binary.LittleEndian.PutUint64(b[:8],uint64(16))
//	binary.LittleEndian.PutUint64(b[8:16],dv.Val)
//	binary.LittleEndian.PutUint64(b[16:],dv.ExtVal)
//	return b,nil
//
//}
//
//func (dv *DictValue) FalconDecoding(src []byte) error {
//	if len(src)!=24{
//		return fmt.Errorf("Length is not 24 byte")
//	}
//	dv.Val=binary.LittleEndian.Uint64(src[8:16])
//	dv.ExtVal=binary.LittleEndian.Uint64(src[16:])
//	return nil
//}
//
//func (dv *DictValue) ToString() string {
//	return fmt.Sprintf(`{ "Val": %d , "ExtVal"：%d }`,dv.Val,dv.ExtVal)
//}
//

//
//type DocId struct{
//	DocID uint32
//	Weight uint32
//}
//
//func (di *DocId) ToString() string {
//	return fmt.Sprintf(`{"id":%d,"weight":%d}`,di.DocID,di.Weight)
//}


// 字段类型
type FalconFieldType uint32

const (
	// 字符串类型
	TFalconString FalconFieldType = 0x0001
)


// 字符串信息
type FalconFieldInfo struct {
	Name string
	Type FalconFieldType
	Offset int64

}

func (ffi *FalconFieldInfo) ToString() string {
	return fmt.Sprintf(`{"name":"%s","type":"%d","offset":%d}`,ffi.Name,ffi.Type,ffi.Offset)
}


// 读写模式
type FalconMode uint32

const (
	TWriteMode FalconMode = 0x0001
	TReadMode FalconMode = 0x0002
	TRWMode FalconMode = 0x0003
)


const (
	TKeywordType string = "keyword"
	TTextType string = "text"
)




type FalconMapping struct{
	FieldName string
	FieldType string
}

func (fm *FalconMapping) GetFieldInfo() (*FalconFieldInfo,error) {

	switch fm.FieldType {
	case TKeywordType:
		fallthrough
	case TTextType:
		return &FalconFieldInfo{Name:fm.FieldName,Type:TFalconString},nil
	default:
		return nil,fmt.Errorf("Mapping is not right [ %s ]",fm.FieldType)
	}

}


type FalconIndexMappings struct {

	mappingLocker *sync.RWMutex
	Mappings map[string]*FalconMapping
}

func NewFalconIndexMappings() *FalconIndexMappings {
	return &FalconIndexMappings{mappingLocker:new(sync.RWMutex),Mappings:make(map[string]*FalconMapping)}
}

func (fim *FalconIndexMappings) GetMappings() []*FalconMapping {
	mappings := make([]*FalconMapping,0)
	fim.mappingLocker.RLock()
	defer fim.mappingLocker.RUnlock()
	for _,v:=range fim.Mappings {
		mappings=append(mappings,v)
	}
	return mappings

}

func (fim *FalconIndexMappings) GetFieldMapping(name string) (*FalconMapping,bool) {
	fim.mappingLocker.RLock()
	defer fim.mappingLocker.RUnlock()
	v,ok:= fim.Mappings[name]
	return v,ok
}

func (fim *FalconIndexMappings) AddFieldMapping(fieldMapping *FalconMapping) error {
	fim.mappingLocker.Lock()
	defer fim.mappingLocker.Unlock()
	if _,ok:=fim.Mappings[fieldMapping.FieldName];!ok {
		fim.Mappings[fieldMapping.FieldName] = fieldMapping
		return nil
	}

	return fmt.Errorf("mapping [ %s ] is already exist",fieldMapping.FieldName)

}