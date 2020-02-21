package rbxl

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"

	"github.com/robloxapi/rbxfile"
)

// Type represents a type that can be serialized.
type Type byte

const (
	TypeInvalid      Type = 0x0
	TypeString       Type = 0x1
	TypeBool         Type = 0x2
	TypeInt          Type = 0x3
	TypeFloat        Type = 0x4
	TypeDouble       Type = 0x5
	TypeUDim         Type = 0x6
	TypeUDim2        Type = 0x7
	TypeRay          Type = 0x8
	TypeFaces        Type = 0x9
	TypeAxes         Type = 0xA
	TypeBrickColor   Type = 0xB
	TypeColor3       Type = 0xC
	TypeVector2      Type = 0xD
	TypeVector3      Type = 0xE
	TypeVector2int16 Type = 0xF
	TypeCFrame       Type = 0x10
	//TypeCFrameQuat Type = 0x11
	TypeToken              Type = 0x12
	TypeReference          Type = 0x13
	TypeVector3int16       Type = 0x14
	TypeNumberSequence     Type = 0x15
	TypeColorSequence      Type = 0x16
	TypeNumberRange        Type = 0x17
	TypeRect2D             Type = 0x18
	TypePhysicalProperties Type = 0x19
	TypeColor3uint8        Type = 0x1A
	TypeInt64              Type = 0x1B
	TypeSharedString       Type = 0x1C
)

func (t Type) Valid() bool {
	return TypeString <= t && t <= TypeSharedString && t != 0x11
}

// String returns a string representation of the type. If the type is not
// valid, then the returned value will be "Invalid".
func (t Type) String() string {
	switch t {
	case TypeString:
		return "String"
	case TypeBool:
		return "Bool"
	case TypeInt:
		return "Int"
	case TypeFloat:
		return "Float"
	case TypeDouble:
		return "Double"
	case TypeUDim:
		return "UDim"
	case TypeUDim2:
		return "UDim2"
	case TypeRay:
		return "Ray"
	case TypeFaces:
		return "Faces"
	case TypeAxes:
		return "Axes"
	case TypeBrickColor:
		return "BrickColor"
	case TypeColor3:
		return "Color3"
	case TypeVector2:
		return "Vector2"
	case TypeVector3:
		return "Vector3"
	case TypeVector2int16:
		return "Vector2int16"
	case TypeCFrame:
		return "CFrame"
	// case TypeCFrameQuat:
	// 	return "CFrameQuat"
	case TypeToken:
		return "Token"
	case TypeReference:
		return "Reference"
	case TypeVector3int16:
		return "Vector3int16"
	case TypeNumberSequence:
		return "NumberSequence"
	case TypeColorSequence:
		return "ColorSequence"
	case TypeNumberRange:
		return "NumberRange"
	case TypeRect2D:
		return "Rect2D"
	case TypePhysicalProperties:
		return "PhysicalProperties"
	case TypeColor3uint8:
		return "Color3uint8"
	case TypeInt64:
		return "Int64"
	case TypeSharedString:
		return "SharedString"
	default:
		return "Invalid"
	}
}

// ValueType returns the rbxfile.Type that corresponds to the type.
func (t Type) ValueType() rbxfile.Type {
	switch t {
	case TypeString:
		return rbxfile.TypeString
	case TypeBool:
		return rbxfile.TypeBool
	case TypeInt:
		return rbxfile.TypeInt
	case TypeFloat:
		return rbxfile.TypeFloat
	case TypeDouble:
		return rbxfile.TypeDouble
	case TypeUDim:
		return rbxfile.TypeUDim
	case TypeUDim2:
		return rbxfile.TypeUDim2
	case TypeRay:
		return rbxfile.TypeRay
	case TypeFaces:
		return rbxfile.TypeFaces
	case TypeAxes:
		return rbxfile.TypeAxes
	case TypeBrickColor:
		return rbxfile.TypeBrickColor
	case TypeColor3:
		return rbxfile.TypeColor3
	case TypeVector2:
		return rbxfile.TypeVector2
	case TypeVector3:
		return rbxfile.TypeVector3
	case TypeVector2int16:
		return rbxfile.TypeVector2int16
	case TypeCFrame:
		return rbxfile.TypeCFrame
	case TypeToken:
		return rbxfile.TypeToken
	case TypeReference:
		return rbxfile.TypeReference
	case TypeVector3int16:
		return rbxfile.TypeVector3int16
	case TypeNumberSequence:
		return rbxfile.TypeNumberSequence
	case TypeColorSequence:
		return rbxfile.TypeColorSequence
	case TypeNumberRange:
		return rbxfile.TypeNumberRange
	case TypeRect2D:
		return rbxfile.TypeRect2D
	case TypePhysicalProperties:
		return rbxfile.TypePhysicalProperties
	case TypeColor3uint8:
		return rbxfile.TypeColor3uint8
	case TypeInt64:
		return rbxfile.TypeInt64
	case TypeSharedString:
		return rbxfile.TypeSharedString
	default:
		return rbxfile.TypeInvalid
	}
}

// FromValueType returns the Type corresponding to a given rbxfile.Type.
func FromValueType(t rbxfile.Type) Type {
	switch t {
	case rbxfile.TypeString:
		return TypeString
	case rbxfile.TypeBool:
		return TypeBool
	case rbxfile.TypeInt:
		return TypeInt
	case rbxfile.TypeFloat:
		return TypeFloat
	case rbxfile.TypeDouble:
		return TypeDouble
	case rbxfile.TypeUDim:
		return TypeUDim
	case rbxfile.TypeUDim2:
		return TypeUDim2
	case rbxfile.TypeRay:
		return TypeRay
	case rbxfile.TypeFaces:
		return TypeFaces
	case rbxfile.TypeAxes:
		return TypeAxes
	case rbxfile.TypeBrickColor:
		return TypeBrickColor
	case rbxfile.TypeColor3:
		return TypeColor3
	case rbxfile.TypeVector2:
		return TypeVector2
	case rbxfile.TypeVector3:
		return TypeVector3
	case rbxfile.TypeVector2int16:
		return TypeVector2int16
	case rbxfile.TypeCFrame:
		return TypeCFrame
	case rbxfile.TypeToken:
		return TypeToken
	case rbxfile.TypeReference:
		return TypeReference
	case rbxfile.TypeVector3int16:
		return TypeVector3int16
	case rbxfile.TypeNumberSequence:
		return TypeNumberSequence
	case rbxfile.TypeColorSequence:
		return TypeColorSequence
	case rbxfile.TypeNumberRange:
		return TypeNumberRange
	case rbxfile.TypeRect2D:
		return TypeRect2D
	case rbxfile.TypePhysicalProperties:
		return TypePhysicalProperties
	case rbxfile.TypeColor3uint8:
		return TypeColor3uint8
	case rbxfile.TypeInt64:
		return TypeInt64
	case rbxfile.TypeSharedString:
		return TypeSharedString
	default:
		return TypeInvalid
	}
}

// Value represents a value of a certain Type.
type Value interface {
	// Type returns an identifier indicating the type.
	Type() Type

	// FromBytes receives the value of the type from a byte array.
	FromBytes([]byte) error

	// Bytes returns the encoded value of the type as a byte array.
	Bytes() []byte
}

// NewValue returns new Value of the given Type. The initial value will not
// necessarily be the zero for the type. If the given type is invalid, then a
// nil value is returned.
func NewValue(typ Type) Value {
	newValue, ok := valueGenerators[typ]
	if !ok {
		return nil
	}
	return newValue()
}

type valueGenerator func() Value

var valueGenerators = map[Type]valueGenerator{
	TypeString:       newValueString,
	TypeBool:         newValueBool,
	TypeInt:          newValueInt,
	TypeFloat:        newValueFloat,
	TypeDouble:       newValueDouble,
	TypeUDim:         newValueUDim,
	TypeUDim2:        newValueUDim2,
	TypeRay:          newValueRay,
	TypeFaces:        newValueFaces,
	TypeAxes:         newValueAxes,
	TypeBrickColor:   newValueBrickColor,
	TypeColor3:       newValueColor3,
	TypeVector2:      newValueVector2,
	TypeVector3:      newValueVector3,
	TypeVector2int16: newValueVector2int16,
	TypeCFrame:       newValueCFrame,
	//TypeCFrameQuat: newValueCFrameQuat,
	TypeToken:              newValueToken,
	TypeReference:          newValueReference,
	TypeVector3int16:       newValueVector3int16,
	TypeNumberSequence:     newValueNumberSequence,
	TypeColorSequence:      newValueColorSequence,
	TypeNumberRange:        newValueNumberRange,
	TypeRect2D:             newValueRect2D,
	TypePhysicalProperties: newValuePhysicalProperties,
	TypeColor3uint8:        newValueColor3uint8,
	TypeInt64:              newValueInt64,
	TypeSharedString:       newValueSharedString,
}

////////////////////////////////////////////////////////////////

// Encodes signed integers so that the bytes of negative numbers are more
// similar to positive numbers, making them more compressible.
//
// https://developers.google.com/protocol-buffers/docs/encoding#types
func encodeZigzag32(n int32) uint32 {
	return uint32((n << 1) ^ (n >> 31))
}

func decodeZigzag32(n uint32) int32 {
	return int32((n >> 1) ^ uint32((int32(n&1)<<31)>>31))
}

func encodeZigzag64(n int64) uint64 {
	return uint64((n << 1) ^ (n >> 63))
}

func decodeZigzag64(n uint64) int64 {
	return int64((n >> 1) ^ uint64((int64(n&1)<<63)>>63))
}

// Encodes a Binary32 float with sign at LSB instead of MSB.
func encodeRobloxFloat(f float32) uint32 {
	n := math.Float32bits(f)
	return (n << 1) | (n >> 31)
}

func decodeRobloxFloat(n uint32) float32 {
	f := (n >> 1) | (n << 31)
	return math.Float32frombits(f)
}

////////////////////////////////////////////////////////////////

type ValueString []byte

func newValueString() Value {
	return new(ValueString)
}

func (ValueString) Type() Type {
	return TypeString
}

func (v ValueString) Bytes() []byte {
	b := make([]byte, len(v)+4)
	binary.LittleEndian.PutUint32(b, uint32(len(v)))
	copy(b[4:], v)
	return b
}

func (v *ValueString) FromBytes(b []byte) error {
	if len(b) < 4 {
		return errors.New("array length must be greater than or equal to 4")
	}

	length := binary.LittleEndian.Uint32(b[:4])
	str := b[4:]
	if uint32(len(str)) != length {
		return fmt.Errorf("string length (%d) does not match integer length (%d)", len(str), length)
	}

	*v = make(ValueString, len(str))
	copy(*v, str)

	return nil
}

////////////////////////////////////////////////////////////////

type ValueBool bool

func newValueBool() Value {
	return new(ValueBool)
}

func (ValueBool) Type() Type {
	return TypeBool
}

func (v ValueBool) Bytes() []byte {
	if v {
		return []byte{1}
	}
	return []byte{0}
}

func (v *ValueBool) FromBytes(b []byte) error {
	if len(b) != 1 {
		return errors.New("array length must be 1")
	}

	*v = b[0] != 0

	return nil
}

////////////////////////////////////////////////////////////////

type ValueInt int32

func newValueInt() Value {
	return new(ValueInt)
}

func (ValueInt) Type() Type {
	return TypeInt
}

func (v ValueInt) Bytes() []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, encodeZigzag32(int32(v)))
	return b
}

func (v *ValueInt) FromBytes(b []byte) error {
	if len(b) != 4 {
		return errors.New("array length must be 4")
	}

	*v = ValueInt(decodeZigzag32(binary.BigEndian.Uint32(b)))

	return nil
}

////////////////////////////////////////////////////////////////

type ValueFloat float32

func newValueFloat() Value {
	return new(ValueFloat)
}

func (ValueFloat) Type() Type {
	return TypeFloat
}

func (v ValueFloat) Bytes() []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, encodeRobloxFloat(float32(v)))
	return b
}

func (v *ValueFloat) FromBytes(b []byte) error {
	if len(b) != 4 {
		return errors.New("array length must be 4")
	}

	*v = ValueFloat(decodeRobloxFloat(binary.BigEndian.Uint32(b)))

	return nil
}

////////////////////////////////////////////////////////////////

type ValueDouble float64

func newValueDouble() Value {
	return new(ValueDouble)
}

func (ValueDouble) Type() Type {
	return TypeDouble
}

func (v ValueDouble) Bytes() []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, math.Float64bits(float64(v)))
	return b
}

func (v *ValueDouble) FromBytes(b []byte) error {
	if len(b) != 8 {
		return errors.New("array length must be 8")
	}

	*v = ValueDouble(math.Float64frombits(binary.LittleEndian.Uint64(b)))

	return nil
}

////////////////////////////////////////////////////////////////

type ValueUDim struct {
	Scale  ValueFloat
	Offset ValueInt
}

func newValueUDim() Value {
	return new(ValueUDim)
}

func (ValueUDim) Type() Type {
	return TypeUDim
}

func (v ValueUDim) Bytes() []byte {
	b := make([]byte, 8)

	copy(b[0:4], v.Scale.Bytes())
	copy(b[4:8], v.Offset.Bytes())

	return b
}

func (v *ValueUDim) FromBytes(b []byte) error {
	if len(b) != 8 {
		return errors.New("array length must be 8")
	}

	v.Scale.FromBytes(b[0:4])
	v.Offset.FromBytes(b[4:8])

	return nil
}

func (ValueUDim) fieldLen() []int {
	return []int{4, 4}
}

func (v *ValueUDim) fieldSet(i int, b []byte) (err error) {
	switch i {
	case 0:
		err = v.Scale.FromBytes(b)
	case 1:
		err = v.Offset.FromBytes(b)
	}
	return
}

func (v ValueUDim) fieldGet(i int) (b []byte) {
	switch i {
	case 0:
		return v.Scale.Bytes()
	case 1:
		return v.Offset.Bytes()
	}
	return
}

////////////////////////////////////////////////////////////////

type ValueUDim2 struct {
	ScaleX  ValueFloat
	ScaleY  ValueFloat
	OffsetX ValueInt
	OffsetY ValueInt
}

func newValueUDim2() Value {
	return new(ValueUDim2)
}

func (ValueUDim2) Type() Type {
	return TypeUDim2
}

func (v ValueUDim2) Bytes() []byte {
	b := make([]byte, 16)
	copy(b[0:4], v.ScaleX.Bytes())
	copy(b[4:8], v.ScaleY.Bytes())
	copy(b[8:12], v.OffsetX.Bytes())
	copy(b[12:16], v.OffsetY.Bytes())
	return b
}

func (v *ValueUDim2) FromBytes(b []byte) error {
	if len(b) != 16 {
		return errors.New("array length must be 16")
	}

	v.ScaleX.FromBytes(b[0:4])
	v.ScaleY.FromBytes(b[4:8])
	v.OffsetX.FromBytes(b[8:12])
	v.OffsetY.FromBytes(b[12:16])

	return nil
}

func (ValueUDim2) fieldLen() []int {
	return []int{4, 4, 4, 4}
}

func (v *ValueUDim2) fieldSet(i int, b []byte) (err error) {
	switch i {
	case 0:
		err = v.ScaleX.FromBytes(b)
	case 1:
		err = v.ScaleY.FromBytes(b)
	case 2:
		err = v.OffsetX.FromBytes(b)
	case 3:
		err = v.OffsetY.FromBytes(b)
	}
	return
}

func (v ValueUDim2) fieldGet(i int) (b []byte) {
	switch i {
	case 0:
		return v.ScaleX.Bytes()
	case 1:
		return v.ScaleY.Bytes()
	case 2:
		return v.OffsetX.Bytes()
	case 3:
		return v.OffsetY.Bytes()
	}
	return
}

////////////////////////////////////////////////////////////////

type ValueRay struct {
	OriginX    float32
	OriginY    float32
	OriginZ    float32
	DirectionX float32
	DirectionY float32
	DirectionZ float32
}

func newValueRay() Value {
	return new(ValueRay)
}

func (ValueRay) Type() Type {
	return TypeRay
}

func (v ValueRay) Bytes() []byte {
	b := make([]byte, 24)
	binary.LittleEndian.PutUint32(b[0:4], math.Float32bits(v.OriginX))
	binary.LittleEndian.PutUint32(b[4:8], math.Float32bits(v.OriginY))
	binary.LittleEndian.PutUint32(b[8:12], math.Float32bits(v.OriginZ))
	binary.LittleEndian.PutUint32(b[12:16], math.Float32bits(v.DirectionX))
	binary.LittleEndian.PutUint32(b[16:20], math.Float32bits(v.DirectionY))
	binary.LittleEndian.PutUint32(b[20:24], math.Float32bits(v.DirectionZ))
	return b
}

func (v *ValueRay) FromBytes(b []byte) error {
	if len(b) != 24 {
		return errors.New("array length must be 24")
	}

	v.OriginX = math.Float32frombits(binary.LittleEndian.Uint32(b[0:4]))
	v.OriginY = math.Float32frombits(binary.LittleEndian.Uint32(b[4:8]))
	v.OriginZ = math.Float32frombits(binary.LittleEndian.Uint32(b[8:12]))
	v.DirectionX = math.Float32frombits(binary.LittleEndian.Uint32(b[12:16]))
	v.DirectionY = math.Float32frombits(binary.LittleEndian.Uint32(b[16:20]))
	v.DirectionZ = math.Float32frombits(binary.LittleEndian.Uint32(b[20:24]))

	return nil
}

////////////////////////////////////////////////////////////////

type ValueFaces struct {
	Right, Top, Back, Left, Bottom, Front bool
}

func newValueFaces() Value {
	return new(ValueFaces)
}

func (ValueFaces) Type() Type {
	return TypeFaces
}

func (v ValueFaces) Bytes() []byte {
	flags := [6]bool{v.Right, v.Top, v.Back, v.Left, v.Bottom, v.Front}
	var b byte
	for i, flag := range flags {
		if flag {
			b = b | (1 << uint(i))
		}
	}

	return []byte{b}
}

func (v *ValueFaces) FromBytes(b []byte) error {
	if len(b) != 1 {
		return errors.New("array length must be 1")
	}

	v.Right = b[0]&(1<<0) != 0
	v.Top = b[0]&(1<<1) != 0
	v.Back = b[0]&(1<<2) != 0
	v.Left = b[0]&(1<<3) != 0
	v.Bottom = b[0]&(1<<4) != 0
	v.Front = b[0]&(1<<5) != 0

	return nil
}

////////////////////////////////////////////////////////////////

type ValueAxes struct {
	X, Y, Z bool
}

func newValueAxes() Value {
	return new(ValueAxes)
}

func (ValueAxes) Type() Type {
	return TypeAxes
}

func (v ValueAxes) Bytes() []byte {
	flags := [3]bool{v.X, v.Y, v.Z}
	var b byte
	for i, flag := range flags {
		if flag {
			b = b | (1 << uint(i))
		}
	}

	return []byte{b}
}

func (v *ValueAxes) FromBytes(b []byte) error {
	if len(b) != 1 {
		return errors.New("array length must be 1")
	}

	v.X = b[0]&(1<<0) != 0
	v.Y = b[0]&(1<<1) != 0
	v.Z = b[0]&(1<<2) != 0

	return nil
}

////////////////////////////////////////////////////////////////

type ValueBrickColor uint32

func newValueBrickColor() Value {
	return new(ValueBrickColor)
}

func (ValueBrickColor) Type() Type {
	return TypeBrickColor
}

func (v ValueBrickColor) Bytes() []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(v))
	return b
}

func (v *ValueBrickColor) FromBytes(b []byte) error {
	if len(b) != 4 {
		return errors.New("array length must be 4")
	}

	*v = ValueBrickColor(binary.BigEndian.Uint32(b))

	return nil
}

////////////////////////////////////////////////////////////////

type ValueColor3 struct {
	R, G, B ValueFloat
}

func newValueColor3() Value {
	return new(ValueColor3)
}

func (ValueColor3) Type() Type {
	return TypeColor3
}

func (v ValueColor3) Bytes() []byte {
	b := make([]byte, 12)
	copy(b[0:4], v.R.Bytes())
	copy(b[4:8], v.G.Bytes())
	copy(b[8:12], v.B.Bytes())
	return b
}

func (v *ValueColor3) FromBytes(b []byte) error {
	if len(b) != 12 {
		return errors.New("array length must be 12")
	}

	v.R.FromBytes(b[0:4])
	v.G.FromBytes(b[4:8])
	v.B.FromBytes(b[8:12])

	return nil
}

func (ValueColor3) fieldLen() []int {
	return []int{4, 4, 4}
}

func (v *ValueColor3) fieldSet(i int, b []byte) (err error) {
	switch i {
	case 0:
		err = v.R.FromBytes(b)
	case 1:
		err = v.G.FromBytes(b)
	case 2:
		err = v.B.FromBytes(b)
	}
	return
}

func (v ValueColor3) fieldGet(i int) (b []byte) {
	switch i {
	case 0:
		return v.R.Bytes()
	case 1:
		return v.G.Bytes()
	case 2:
		return v.B.Bytes()
	}
	return
}

////////////////////////////////////////////////////////////////

type ValueVector2 struct {
	X, Y ValueFloat
}

func newValueVector2() Value {
	return new(ValueVector2)
}

func (ValueVector2) Type() Type {
	return TypeVector2
}

func (v ValueVector2) Bytes() []byte {
	b := make([]byte, 8)
	copy(b[0:4], v.X.Bytes())
	copy(b[4:8], v.Y.Bytes())
	return b
}

func (v *ValueVector2) FromBytes(b []byte) error {
	if len(b) != 8 {
		return errors.New("array length must be 8")
	}

	v.X.FromBytes(b[0:4])
	v.Y.FromBytes(b[4:8])

	return nil
}

func (ValueVector2) fieldLen() []int {
	return []int{4, 4}
}

func (v *ValueVector2) fieldSet(i int, b []byte) (err error) {
	switch i {
	case 0:
		err = v.X.FromBytes(b)
	case 1:
		err = v.Y.FromBytes(b)
	}
	return
}

func (v ValueVector2) fieldGet(i int) (b []byte) {
	switch i {
	case 0:
		return v.X.Bytes()
	case 1:
		return v.Y.Bytes()
	}
	return
}

////////////////////////////////////////////////////////////////

type ValueVector3 struct {
	X, Y, Z ValueFloat
}

func newValueVector3() Value {
	return new(ValueVector3)
}

func (ValueVector3) Type() Type {
	return TypeVector3
}

func (v ValueVector3) Bytes() []byte {
	b := make([]byte, 12)
	copy(b[0:4], v.X.Bytes())
	copy(b[4:8], v.Y.Bytes())
	copy(b[8:12], v.Z.Bytes())
	return b
}

func (v *ValueVector3) FromBytes(b []byte) error {
	if len(b) != 12 {
		return errors.New("array length must be 12")
	}

	v.X.FromBytes(b[0:4])
	v.Y.FromBytes(b[4:8])
	v.Z.FromBytes(b[8:12])

	return nil
}

func (ValueVector3) fieldLen() []int {
	return []int{4, 4, 4}
}

func (v *ValueVector3) fieldSet(i int, b []byte) (err error) {
	switch i {
	case 0:
		err = v.X.FromBytes(b)
	case 1:
		err = v.Y.FromBytes(b)
	case 2:
		err = v.Z.FromBytes(b)
	}
	return
}

func (v ValueVector3) fieldGet(i int) (b []byte) {
	switch i {
	case 0:
		return v.X.Bytes()
	case 1:
		return v.Y.Bytes()
	case 2:
		return v.Z.Bytes()
	}
	return
}

////////////////////////////////////////////////////////////////

type ValueVector2int16 struct {
	X, Y int16
}

func newValueVector2int16() Value {
	return new(ValueVector2int16)
}

func (ValueVector2int16) Type() Type {
	return TypeVector2int16
}

func (v ValueVector2int16) Bytes() []byte {
	b := make([]byte, 4)

	binary.LittleEndian.PutUint16(b[0:2], uint16(v.X))
	binary.LittleEndian.PutUint16(b[2:4], uint16(v.Y))

	return b
}

func (v *ValueVector2int16) FromBytes(b []byte) error {
	if len(b) != 4 {
		return errors.New("array length must be 4")
	}

	v.X = int16(binary.LittleEndian.Uint16(b[0:2]))
	v.Y = int16(binary.LittleEndian.Uint16(b[2:4]))

	return nil
}

////////////////////////////////////////////////////////////////

type ValueCFrame struct {
	Special  uint8
	Rotation [9]float32
	Position ValueVector3
}

func newValueCFrame() Value {
	return new(ValueCFrame)
}

func (ValueCFrame) Type() Type {
	return TypeCFrame
}

func (v ValueCFrame) Bytes() []byte {
	var b []byte
	if v.Special == 0 {
		b = make([]byte, 49)
		r := b[1:]
		for i, f := range v.Rotation {
			binary.LittleEndian.PutUint32(r[i*4:i*4+4], math.Float32bits(f))
		}
	} else {
		b = make([]byte, 13)
		b[0] = v.Special
	}

	copy(b[len(b)-12:], v.Position.Bytes())

	return b
}

func (v *ValueCFrame) FromBytes(b []byte) error {
	if b[0] == 0 && len(b) != 49 {
		return errors.New("array length must be 49")
	} else if b[0] != 0 && len(b) != 13 {
		return errors.New("array length must be 13")
	}

	v.Special = b[0]

	if b[0] == 0 {
		r := b[1:]
		for i := range v.Rotation {
			v.Rotation[i] = math.Float32frombits(binary.LittleEndian.Uint32(r[i*4 : i*4+4]))
		}
	} else {
		for i := range v.Rotation {
			v.Rotation[i] = 0
		}
	}

	v.Position.FromBytes(b[len(b)-12:])

	return nil
}

////////////////////////////////////////////////////////////////

type ValueToken uint32

func newValueToken() Value {
	return new(ValueToken)
}

func (ValueToken) Type() Type {
	return TypeToken
}

func (v ValueToken) Bytes() []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(v))
	return b
}

func (v *ValueToken) FromBytes(b []byte) error {
	if len(b) != 4 {
		return errors.New("array length must be 4")
	}

	*v = ValueToken(binary.BigEndian.Uint32(b))

	return nil
}

////////////////////////////////////////////////////////////////

type ValueReference int32

func newValueReference() Value {
	return new(ValueReference)
}

func (ValueReference) Type() Type {
	return TypeReference
}

func (v ValueReference) Bytes() []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, encodeZigzag32(int32(v)))
	return b
}

func (v *ValueReference) FromBytes(b []byte) error {
	if len(b) != 4 {
		return errors.New("array length must be 4")
	}

	*v = ValueReference(decodeZigzag32(binary.BigEndian.Uint32(b)))

	return nil
}

////////////////////////////////////////////////////////////////

type ValueVector3int16 struct {
	X, Y, Z int16
}

func newValueVector3int16() Value {
	return new(ValueVector3int16)
}

func (ValueVector3int16) Type() Type {
	return TypeVector3int16
}

func (v ValueVector3int16) Bytes() []byte {
	b := make([]byte, 6)

	binary.LittleEndian.PutUint16(b[0:2], uint16(v.X))
	binary.LittleEndian.PutUint16(b[2:4], uint16(v.Y))
	binary.LittleEndian.PutUint16(b[4:6], uint16(v.Z))

	return b
}

func (v *ValueVector3int16) FromBytes(b []byte) error {
	if len(b) != 6 {
		return errors.New("array length must be 6")
	}

	v.X = int16(binary.LittleEndian.Uint16(b[0:2]))
	v.Y = int16(binary.LittleEndian.Uint16(b[2:4]))
	v.Z = int16(binary.LittleEndian.Uint16(b[4:6]))

	return nil
}

////////////////////////////////////////////////////////////////

const sizeNSK = 3 * 4

type ValueNumberSequenceKeypoint struct {
	Time, Value, Envelope float32
}

type ValueNumberSequence []ValueNumberSequenceKeypoint

func newValueNumberSequence() Value {
	return new(ValueNumberSequence)
}

func (ValueNumberSequence) Type() Type {
	return TypeNumberSequence
}

func (v ValueNumberSequence) Bytes() []byte {
	b := make([]byte, 4+len(v)*sizeNSK)

	binary.LittleEndian.PutUint32(b, uint32(len(v)))
	ba := b[4:]

	for i, nsk := range v {
		bk := ba[i*sizeNSK:]
		binary.LittleEndian.PutUint32(bk[0:4], math.Float32bits(nsk.Time))
		binary.LittleEndian.PutUint32(bk[4:8], math.Float32bits(nsk.Value))
		binary.LittleEndian.PutUint32(bk[8:12], math.Float32bits(nsk.Envelope))
	}

	return b
}

func (v *ValueNumberSequence) FromBytes(b []byte) error {
	if len(b) < 4 {
		return errors.New("array length must be at least 4")
	}

	length := int(binary.LittleEndian.Uint32(b))
	ba := b[4:]
	if len(ba) != sizeNSK*length {
		return fmt.Errorf("expected array length of %d (4 + %d * %d)", 4+sizeNSK*length, sizeNSK, length)
	}

	a := make(ValueNumberSequence, length)
	for i := 0; i < length; i++ {
		bk := ba[i*sizeNSK:]
		a[i] = ValueNumberSequenceKeypoint{
			Time:     math.Float32frombits(binary.LittleEndian.Uint32(bk[0:4])),
			Value:    math.Float32frombits(binary.LittleEndian.Uint32(bk[4:8])),
			Envelope: math.Float32frombits(binary.LittleEndian.Uint32(bk[8:12])),
		}
	}

	*v = a

	return nil
}

////////////////////////////////////////////////////////////////

const sizeCSK = 4 + 3*4 + 4

type ValueColorSequenceKeypoint struct {
	Time     float32
	Value    ValueColor3
	Envelope float32
}

type ValueColorSequence []ValueColorSequenceKeypoint

func newValueColorSequence() Value {
	return new(ValueColorSequence)
}

func (ValueColorSequence) Type() Type {
	return TypeColorSequence
}

func (v ValueColorSequence) Bytes() []byte {
	b := make([]byte, 4+len(v)*sizeCSK)

	binary.LittleEndian.PutUint32(b, uint32(len(v)))
	ba := b[4:]

	for i, csk := range v {
		bk := ba[i*sizeCSK:]
		binary.LittleEndian.PutUint32(bk[0:4], math.Float32bits(csk.Time))
		binary.LittleEndian.PutUint32(bk[4:8], math.Float32bits(float32(csk.Value.R)))
		binary.LittleEndian.PutUint32(bk[8:12], math.Float32bits(float32(csk.Value.G)))
		binary.LittleEndian.PutUint32(bk[12:16], math.Float32bits(float32(csk.Value.B)))
		binary.LittleEndian.PutUint32(bk[16:20], math.Float32bits(csk.Envelope))
	}

	return b
}

func (v *ValueColorSequence) FromBytes(b []byte) error {
	if len(b) < 4 {
		return errors.New("array length must be at least 4")
	}

	length := int(binary.LittleEndian.Uint32(b))
	ba := b[4:]
	if len(ba) != sizeCSK*length {
		return fmt.Errorf("expected array length of %d (4 + %d * %d)", 4+sizeCSK*length, sizeCSK, length)
	}

	a := make(ValueColorSequence, length)
	for i := 0; i < length; i++ {
		bk := ba[i*sizeCSK:]
		c3 := *new(ValueColor3)
		c3.FromBytes(bk[4:16])
		a[i] = ValueColorSequenceKeypoint{
			Time: math.Float32frombits(binary.LittleEndian.Uint32(bk[0:4])),
			Value: ValueColor3{
				R: ValueFloat(math.Float32frombits(binary.LittleEndian.Uint32(bk[4:8]))),
				G: ValueFloat(math.Float32frombits(binary.LittleEndian.Uint32(bk[8:12]))),
				B: ValueFloat(math.Float32frombits(binary.LittleEndian.Uint32(bk[12:16]))),
			},
			Envelope: math.Float32frombits(binary.LittleEndian.Uint32(bk[16:20])),
		}
	}

	*v = a

	return nil
}

////////////////////////////////////////////////////////////////

type ValueNumberRange struct {
	Min, Max float32
}

func newValueNumberRange() Value {
	return new(ValueNumberRange)
}

func (ValueNumberRange) Type() Type {
	return TypeNumberRange
}

func (v ValueNumberRange) Bytes() []byte {
	b := make([]byte, 8)

	binary.LittleEndian.PutUint32(b[0:4], math.Float32bits(v.Min))
	binary.LittleEndian.PutUint32(b[4:8], math.Float32bits(v.Max))

	return b
}

func (v *ValueNumberRange) FromBytes(b []byte) error {
	if len(b) != 8 {
		return errors.New("array length must be 8")
	}

	v.Min = math.Float32frombits(binary.LittleEndian.Uint32(b[0:4]))
	v.Max = math.Float32frombits(binary.LittleEndian.Uint32(b[4:8]))

	return nil
}

////////////////////////////////////////////////////////////////

type ValueRect2D struct {
	Min, Max ValueVector2
}

func newValueRect2D() Value {
	return new(ValueRect2D)
}

func (ValueRect2D) Type() Type {
	return TypeRect2D
}

func (v ValueRect2D) Bytes() []byte {
	b := make([]byte, 16)

	copy(b[0:8], v.Min.Bytes())
	copy(b[8:16], v.Max.Bytes())

	return b
}

func (v *ValueRect2D) FromBytes(b []byte) error {
	if len(b) != 16 {
		return errors.New("array length must be 16")
	}

	v.Min.FromBytes(b[0:8])
	v.Max.FromBytes(b[8:16])

	return nil
}

func (ValueRect2D) fieldLen() []int {
	return []int{4, 4, 4, 4}
}

func (v *ValueRect2D) fieldSet(i int, b []byte) (err error) {
	switch i {
	case 0:
		err = v.Min.X.FromBytes(b)
	case 1:
		err = v.Min.Y.FromBytes(b)
	case 2:
		err = v.Max.X.FromBytes(b)
	case 3:
		err = v.Max.Y.FromBytes(b)
	}
	return
}

func (v ValueRect2D) fieldGet(i int) (b []byte) {
	switch i {
	case 0:
		return v.Min.X.Bytes()
	case 1:
		return v.Min.Y.Bytes()
	case 2:
		return v.Max.X.Bytes()
	case 3:
		return v.Max.Y.Bytes()
	}
	return
}

////////////////////////////////////////////////////////////////

type ValuePhysicalProperties struct {
	CustomPhysics    byte
	Density          float32
	Friction         float32
	Elasticity       float32
	FrictionWeight   float32
	ElasticityWeight float32
}

func newValuePhysicalProperties() Value {
	return new(ValuePhysicalProperties)
}

func (ValuePhysicalProperties) Type() Type {
	return TypePhysicalProperties
}

func (v ValuePhysicalProperties) Bytes() []byte {
	if v.CustomPhysics != 0 {
		b := make([]byte, 21)
		b[0] = v.CustomPhysics
		q := b[1:]
		binary.LittleEndian.PutUint32(q[0*4:0*4+4], math.Float32bits(v.Density))
		binary.LittleEndian.PutUint32(q[1*4:1*4+4], math.Float32bits(v.Friction))
		binary.LittleEndian.PutUint32(q[2*4:2*4+4], math.Float32bits(v.Elasticity))
		binary.LittleEndian.PutUint32(q[3*4:3*4+4], math.Float32bits(v.FrictionWeight))
		binary.LittleEndian.PutUint32(q[4*4:4*4+4], math.Float32bits(v.ElasticityWeight))
		return b
	}
	return make([]byte, 1)
}

func (v *ValuePhysicalProperties) FromBytes(b []byte) error {
	if b[0] == 0 && len(b) != 21 {
		return errors.New("array length must be 21")
	} else if b[0] != 0 && len(b) != 1 {
		return errors.New("array length must be 1")
	}

	v.CustomPhysics = b[0]
	if v.CustomPhysics != 0 {
		p := b[1:]
		v.Density = math.Float32frombits(binary.LittleEndian.Uint32(p[0*4 : 0*4+4]))
		v.Friction = math.Float32frombits(binary.LittleEndian.Uint32(p[1*4 : 1*4+4]))
		v.Elasticity = math.Float32frombits(binary.LittleEndian.Uint32(p[2*4 : 2*4+4]))
		v.FrictionWeight = math.Float32frombits(binary.LittleEndian.Uint32(p[3*4 : 3*4+4]))
		v.ElasticityWeight = math.Float32frombits(binary.LittleEndian.Uint32(p[4*4 : 4*4+4]))
	} else {
		v.Density = 0
		v.Friction = 0
		v.Elasticity = 0
		v.FrictionWeight = 0
		v.ElasticityWeight = 0
	}

	return nil
}

////////////////////////////////////////////////////////////////

type ValueColor3uint8 struct {
	R, G, B byte
}

func newValueColor3uint8() Value {
	return new(ValueColor3uint8)
}

func (ValueColor3uint8) Type() Type {
	return TypeColor3uint8
}

func (v ValueColor3uint8) Bytes() []byte {
	b := make([]byte, 3)
	b[0] = v.R
	b[1] = v.G
	b[2] = v.B
	return b
}

func (v *ValueColor3uint8) FromBytes(b []byte) error {
	if len(b) != 3 {
		return errors.New("array length must be 3")
	}

	v.R = b[0]
	v.G = b[1]
	v.B = b[2]

	return nil
}

func (ValueColor3uint8) fieldLen() []int {
	return []int{1, 1, 1}
}

func (v *ValueColor3uint8) fieldSet(i int, b []byte) (err error) {
	switch i {
	case 0:
		v.R = b[0]
	case 1:
		v.G = b[0]
	case 2:
		v.B = b[0]
	}
	return
}

func (v ValueColor3uint8) fieldGet(i int) (b []byte) {
	switch i {
	case 0:
		return []byte{v.R}
	case 1:
		return []byte{v.G}
	case 2:
		return []byte{v.B}
	}
	return
}

////////////////////////////////////////////////////////////////

type ValueInt64 int64

func newValueInt64() Value {
	return new(ValueInt64)
}

func (ValueInt64) Type() Type {
	return TypeInt64
}

func (v ValueInt64) Bytes() []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, encodeZigzag64(int64(v)))
	return b
}

func (v *ValueInt64) FromBytes(b []byte) error {
	if len(b) != 8 {
		return errors.New("array length must be 8")
	}

	*v = ValueInt64(decodeZigzag64(binary.BigEndian.Uint64(b)))

	return nil
}

////////////////////////////////////////////////////////////////

type ValueSharedString uint32

func newValueSharedString() Value {
	return new(ValueSharedString)
}

func (ValueSharedString) Type() Type {
	return TypeSharedString
}

func (v ValueSharedString) Bytes() []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(v))
	return b
}

func (v *ValueSharedString) FromBytes(b []byte) error {
	if len(b) != 4 {
		return errors.New("array length must be 4")
	}

	*v = ValueSharedString(binary.BigEndian.Uint32(b))

	return nil
}

////////////////////////////////////////////////////////////////
