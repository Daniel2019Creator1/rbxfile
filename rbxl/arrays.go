package rbxl

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
)

// Encodes and decodes a Value based on its fields
type fielder interface {
	// Value.Type
	Type() Type
	// Length of each field
	fieldLen() []int
	// Set bytes of nth field
	fieldSet(int, []byte) error
	// Get bytes of nth field
	fieldGet(int) []byte
}

// Encodes Values that implement the fielder interface.
func interleaveFields(id Type, a []Value) (b []byte, err error) {
	if len(a) == 0 {
		return b, nil
	}

	af := make([]fielder, len(a))
	for i, v := range a {
		af[i] = v.(fielder)
		if af[i].Type() != id {
			return nil, fmt.Errorf("element %d is of type %s where %s is expected", i, af[i].Type().String(), id.String())
		}
	}

	// list is assumed to contain the same kinds of values

	// Number of bytes per field
	nbytes := af[0].fieldLen()
	// Number fields per value
	nfields := len(nbytes)
	// Number of values
	nvalues := len(af)

	// Total bytes per value
	tbytes := 0
	// Offset of each field slice
	ofields := make([]int, len(nbytes)+1)
	for i, n := range nbytes {
		tbytes += n
		ofields[i+1] = ofields[i] + n*nvalues
	}

	b = make([]byte, tbytes*nvalues)

	// List of each field slice
	fields := make([][]byte, nfields)
	for i := range fields {
		// Each field slice affects the final array
		fields[i] = b[ofields[i]:ofields[i+1]]
	}

	for i, v := range af {
		for f, field := range fields {
			fb := v.fieldGet(f)
			if len(fb) != nbytes[f] {
				panic("length of field's bytes does not match given field length")
			}
			copy(field[i*nbytes[f]:], fb)
		}
	}

	// Interleave each field slice independently
	for i, field := range fields {
		if err = interleave(field, nbytes[i]); err != nil {
			return nil, err
		}
	}

	return b, nil
}

// Decodes Values that implement the fielder interface.
func deinterleaveFields(id Type, b []byte) (a []Value, err error) {
	if len(b) == 0 {
		return a, nil
	}

	newValue := valueGenerators[id]
	if newValue == nil {
		return nil, fmt.Errorf("type identifier 0x%X is not a valid Type.", id)
	}

	// Number of bytes per field
	nbytes := newValue().(fielder).fieldLen()
	// Number fields per value
	nfields := len(nbytes)

	// Total bytes per value
	tbytes := 0
	for _, n := range nbytes {
		tbytes += n
	}

	if len(b)%tbytes != 0 {
		return nil, fmt.Errorf("length of array (%d) is not divisible by value byte size (%d)", len(b), tbytes)
	}

	// Number of values
	nvalues := len(b) / tbytes
	// Offset of each field slice
	ofields := make([]int, len(nbytes)+1)
	for i, n := range nbytes {
		ofields[i+1] = ofields[i] + n*nvalues
	}

	a = make([]Value, nvalues)

	// List of each field slice
	fields := make([][]byte, nfields)
	for i := range fields {
		fields[i] = b[ofields[i]:ofields[i+1]]
	}

	// Deinterleave each field slice independently
	for i, field := range fields {
		if err = deinterleave(field, nbytes[i]); err != nil {
			return nil, err
		}
	}

	for i := range a {
		v := newValue()
		vf := v.(fielder)
		for f, field := range fields {
			n := nbytes[f]
			fb := field[i*n : i*n+n]
			vf.fieldSet(f, fb)
		}
		a[i] = v
	}

	return a, nil
}

// Interleave transforms an array of bytes by interleaving them based on a
// given size. The size must be a divisor of the array length.
//
// The array is divided into groups, each `length` in size. The nth elements
// of each group are then moved so that they are group together. For example:
//
//     Original:    abcd1234
//     Interleaved: a1b2c3d4
func interleave(bytes []byte, length int) error {
	if length <= 0 {
		return errors.New("length must be greater than 0")
	}
	if len(bytes)%length != 0 {
		return errors.New("length must be a divisor of array length")
	}

	// Matrix transpose algorithm
	cols := length
	rows := len(bytes) / length
	if rows == cols {
		for r := 0; r < rows; r++ {
			for c := 0; c < r; c++ {
				bytes[r*cols+c], bytes[c*cols+r] = bytes[c*cols+r], bytes[r*cols+c]
			}
		}
	} else {
		tmp := make([]byte, len(bytes))
		for r := 0; r < rows; r++ {
			for c := 0; c < cols; c++ {
				tmp[c*rows+r] = bytes[r*cols+c]
			}
		}
		for i, b := range tmp {
			bytes[i] = b
		}
	}

	return nil
}

func deinterleave(bytes []byte, size int) error {
	if size <= 0 {
		return errors.New("size must be greater than 0")
	}
	if len(bytes)%size != 0 {
		return errors.New("size must be a divisor of array length")
	}

	return interleave(bytes, len(bytes)/size)
}

// Appends the bytes of a list of Values into a byte array.
func appendValueBytes(id Type, a []Value) (b []byte, err error) {
	for i, v := range a {
		if v.Type() != id {
			return nil, fmt.Errorf("element %d is of type `%s` where `%s` is expected", i, v.Type().String(), id.String())
		}

		b = append(b, v.Bytes()...)
	}

	return b, nil
}

// Reads a byte array as an array of Values of a certain type. Size is the
// byte size of each Value. If size is less than 0, then values are assumed to
// be of variable length. The first 4 bytes of a value is read as length N of
// the value. Field then indicates the size of each field in the value, so the
// next N*field bytes are read as the full value.
func appendByteValues(id Type, b []byte, size int, field int) (a []Value, err error) {
	gen := valueGenerators[id]
	if size < 0 {
		// Variable length; get size from first 4 bytes.
		ba := b
		for len(ba) > 0 {
			if len(ba) < 4 {
				return nil, errors.New("expected 4 more bytes in array")
			}
			size := int(binary.LittleEndian.Uint32(ba))
			if len(ba[4:]) < size*field {
				return nil, fmt.Errorf("expected %d more bytes in array", size*field)
			}

			v := gen()
			if err := v.FromBytes(ba[:4+size*field]); err != nil {
				return nil, err
			}
			a = append(a, v)

			ba = ba[4+size*field:]
		}
	} else {
		for i := 0; i+size <= len(b); i += size {
			v := gen()
			if err := v.FromBytes(b[i : i+size]); err != nil {
				return nil, err
			}
			a = append(a, v)
		}
	}
	return a, nil
}

// ValuesToBytes encodes a slice of values into binary form, according to t.
// Returns an error if a value cannot be encoded as t.
func ValuesToBytes(t Type, a []Value) (b []byte, err error) {
	if !t.Valid() {
		return nil, fmt.Errorf("invalid type (%02X)", t)
	}
	for i, v := range a {
		if v.Type() != t {
			return nil, fmt.Errorf("element %d is of type `%s` where `%s` is expected", i, v.Type().String(), t.String())
		}
	}

	switch t {
	case TypeString,
		TypeBool,
		TypeDouble,
		TypeRay,
		TypeFaces,
		TypeAxes,
		TypeVector3int16,
		TypeNumberSequence,
		TypeColorSequence,
		TypeNumberRange:
		// Append each value as bytes with no further operation.
		b, err = appendValueBytes(t, a)
	case TypeInt,
		TypeFloat,
		TypeBrickColor,
		TypeToken,
		TypeSharedString:
		// Append each value a bytes, then interleave to improve compression.
		if b, err = appendValueBytes(t, a); err != nil {
			break
		}
		err = interleave(b, 4)
	case TypeInt64:
		// Append each value a bytes, then interleave to improve compression.
		if b, err = appendValueBytes(t, a); err != nil {
			break
		}
		err = interleave(b, 8)
	case TypeUDim,
		TypeUDim2,
		TypeColor3,
		TypeVector2,
		TypeVector3,
		TypeRect2D,
		TypeColor3uint8:
		// Interleave fields.
		return interleaveFields(t, a)
	case TypeCFrame:
		// The bytes of each value can vary in length.
		p := make([]Value, len(a))
		for i, cf := range a {
			cf := cf.(*ValueCFrame)
			// Build matrix part.
			b = append(b, cf.Special)
			if cf.Special == 0 {
				// Write all components.
				r := make([]byte, len(cf.Rotation)*4)
				for i, f := range cf.Rotation {
					binary.LittleEndian.PutUint32(r[i*4:i*4+4], math.Float32bits(f))
				}
				b = append(b, r...)
			}
			// Prepare position part.
			p[i] = &cf.Position
		}
		// Build position part.
		pb, _ := interleaveFields(TypeVector3, p)
		b = append(b, pb...)
	case TypeReference:
		// Because values are generated in sequence, they are likely to be
		// relatively close to each other. Subtracting each value from the
		// previous will likely produce small values that compress well.
		if len(a) == 0 {
			break
		}
		const size = 4
		b = make([]byte, len(a)*size)
		var prev ValueReference
		for i, ref := range a {
			ref := ref.(*ValueReference)
			if i == 0 {
				copy(b[i*size:i*size+size], ref.Bytes())
			} else {
				// Convert absolute ref to relative ref.
				copy(b[i*size:i*size+size], (*ref - prev).Bytes())
			}
			prev = *ref
		}
		err = interleave(b, size)
	case TypePhysicalProperties:
		// The bytes of each value can vary in length.
		q := make([]byte, 20)
		for _, pp := range a {
			pp := pp.(*ValuePhysicalProperties)
			b = append(b, pp.CustomPhysics)
			if pp.CustomPhysics != 0 {
				// Write all fields.
				binary.LittleEndian.PutUint32(q[0*4:0*4+4], math.Float32bits(pp.Density))
				binary.LittleEndian.PutUint32(q[1*4:1*4+4], math.Float32bits(pp.Friction))
				binary.LittleEndian.PutUint32(q[2*4:2*4+4], math.Float32bits(pp.Elasticity))
				binary.LittleEndian.PutUint32(q[3*4:3*4+4], math.Float32bits(pp.FrictionWeight))
				binary.LittleEndian.PutUint32(q[4*4:4*4+4], math.Float32bits(pp.ElasticityWeight))
				b = append(b, q...)
			}
		}
	case TypeVector2int16:
		err = errors.New("not implemented")
	}

	return
}

// ValuesFromBytes decodes b according to t, into a slice of values, the type of
// each corresponding to t.
func ValuesFromBytes(t Type, b []byte) (a []Value, err error) {
	if !t.Valid() {
		return nil, fmt.Errorf("invalid type (%02X)", t)
	}

	switch t {
	case TypeBool,
		TypeFaces,
		TypeAxes:
		// Append from constant size 1.
		a, err = appendByteValues(t, b, 1, 0)
	case TypeVector3int16:
		// Append from constant size 6.
		a, err = appendByteValues(t, b, 6, 0)
	case TypeDouble,
		TypeNumberRange:
		// Append from constant size 8.
		a, err = appendByteValues(t, b, 8, 0)
	case TypeRay:
		// Append from constant size 24.
		a, err = appendByteValues(t, b, 24, 0)
	case TypeString:
		// Append from variable size.
		a, err = appendByteValues(t, b, -1, 1)
	case TypeNumberSequence:
		// Append from variable size.
		a, err = appendByteValues(t, b, -1, sizeNSK)
	case TypeColorSequence:
		// Append from variable size.
		a, err = appendByteValues(t, b, -1, sizeCSK)
	case TypeInt,
		TypeFloat,
		TypeBrickColor,
		TypeToken,
		TypeSharedString:
		// Deinterleave, then append from size 4.
		bc := make([]byte, len(b))
		copy(bc, b)
		if err = deinterleave(bc, 4); err != nil {
			return nil, err
		}
		a, err = appendByteValues(t, bc, 4, 0)
	case TypeInt64:
		// Deinterleave, then append from size 8.
		bc := make([]byte, len(b))
		copy(bc, b)
		if err = deinterleave(bc, 8); err != nil {
			return nil, err
		}
		a, err = appendByteValues(t, bc, 8, 0)
	case TypeUDim,
		TypeUDim2,
		TypeColor3,
		TypeVector2,
		TypeVector3,
		TypeRect2D,
		TypeColor3uint8:
		// Deinterleave fields.
		a, err = deinterleaveFields(t, b)
	case TypeCFrame:
		cfs := make([]*ValueCFrame, 0)
		// This loop reads the matrix data. i is the current position in the
		// byte array. n is the expected size of the position data, which
		// increases every time another CFrame is read. As long as the number of
		// remaining bytes is greater than n, then the next byte can be assumed
		// to be matrix data. By the end, the number of remaining bytes should
		// be exactly equal to n.
		i := 0
		for n := 0; len(b)-i > n; n += 12 {
			cf := new(ValueCFrame)
			cf.Special = b[i]
			i++
			if cf.Special == 0 {
				q := len(cf.Rotation) * 4
				r := b[i:]
				if len(r) < q {
					return nil, fmt.Errorf("expected %d more bytes in array", q)
				}
				for i := range cf.Rotation {
					cf.Rotation[i] = math.Float32frombits(binary.LittleEndian.Uint32(r[i*4 : i*4+4]))
				}
				i += q
			}
			cfs = append(cfs, cf)
		}
		// Read remaining position data using the Position field, which is a
		// ValueVector3.
		a, err = deinterleaveFields(TypeVector3, b[i:])
		if err != nil {
			return
		}
		if len(a) != len(cfs) {
			return nil, errors.New("number of positions does not match number of matrices")
		}
		// Hack: use 'a' variable to receive Vector3 values, then replace them
		// with CFrames. This lets us avoid needing to copy 'cfs' to 'a', and
		// needing to create a second array.
		for i, p := range a {
			cfs[i].Position = *p.(*ValueVector3)
			a[i] = cfs[i]
		}
	case TypeReference:
		if len(b) == 0 {
			return
		}
		const size = 4
		if len(b)%size != 0 {
			return nil, fmt.Errorf("array must be divisible by %d", size)
		}
		bc := make([]byte, len(b))
		copy(bc, b)
		if err = deinterleave(bc, size); err != nil {
			return nil, err
		}
		a = make([]Value, len(bc)/size)
		for i := 0; i < len(bc)/size; i++ {
			ref := new(ValueReference)
			ref.FromBytes(bc[i*size : i*size+size])
			if i > 0 {
				// Convert relative ref to absolute ref.
				r := *a[i-1].(*ValueReference)
				*ref = r + *ref
			}
			a[i] = ref
		}
	case TypePhysicalProperties:
		for i := 0; i < len(b); {
			pp := new(ValuePhysicalProperties)
			pp.CustomPhysics = b[i]
			i++
			if pp.CustomPhysics != 0 {
				const size = 5 * 4
				p := b[i:]
				if len(p) < size {
					return nil, fmt.Errorf("expected %d more bytes in array", size)
				}
				pp.Density = math.Float32frombits(binary.LittleEndian.Uint32(p[0*4 : 0*4+4]))
				pp.Friction = math.Float32frombits(binary.LittleEndian.Uint32(p[1*4 : 1*4+4]))
				pp.Elasticity = math.Float32frombits(binary.LittleEndian.Uint32(p[2*4 : 2*4+4]))
				pp.FrictionWeight = math.Float32frombits(binary.LittleEndian.Uint32(p[3*4 : 3*4+4]))
				pp.ElasticityWeight = math.Float32frombits(binary.LittleEndian.Uint32(p[4*4 : 4*4+4]))
				i += size
			}
			a = append(a, pp)
		}
	case TypeVector2int16:
		err = errors.New("not implemented")
	}

	return
}
