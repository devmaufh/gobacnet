package encoding

import (
	bactype "github.com/devmaufh/gobacnet/types"
	"math"
)

func (e *Encoder) WriteProperty(invokeID uint8, wp bactype.WritePropertyData) {
	// ConfirmedRequest PDU
	e.write(uint8(0x00)) // PDU Type: ConfirmedRequest
	e.write(invokeID)    // Invoke ID
	e.write(uint8(15))   // Service Choice: WriteProperty (0x0F)

	// Object Identifier (Context Tag 0)
	e.contextObjectID(0, wp.Object.ID.Type, wp.Object.ID.Instance)

	// Property Identifier (Context Tag 1)
	e.contextUnsigned(1, uint32(wp.Property.Type))

	// Array Index (Context Tag 2) - optional
	if wp.Property.ArrayIndex != 0xFFFFFFFF {
		e.contextUnsigned(2, wp.Property.ArrayIndex)
	}

	// Value (Context Tag 3 - opening)
	e.openingTag(3)

	switch v := wp.Value.(type) {
	case float32:
		// Tag 4 = REAL
		e.tag(tagInfo{ID: 4, Context: false, Value: 4})
		e.write(realToBytes(v))
	case int:
		// Tag 2 = Signed Int
		e.tag(tagInfo{ID: 2, Context: false, Value: 1})
		e.write(uint8(v))
	case uint:
		// Tag 2 = Unsigned Int
		e.tag(tagInfo{ID: 2, Context: false, Value: 1})
		e.write(uint8(v))
	case bool:
		// Tag 1 = Boolean
		e.tag(tagInfo{ID: 1, Context: false, Value: 1})
		if v {
			e.write(uint8(1))
		} else {
			e.write(uint8(0))
		}
	default:
		// fallback if el tipo no está soportado
		e.err = bactype.ErrUnsupportedWriteValue
		return
	}

	// Value (Context Tag 3 - closing)
	e.closingTag(3)

	// Priority (optional, Context Tag 4)
	if wp.Priority != 0 {
		e.contextUnsigned(4, uint32(wp.Priority))
	}
}

// Conversión de float32 a bytes (IEEE 754)
func realToBytes(val float32) []byte {
	bits := math.Float32bits(val)
	return []byte{
		byte(bits >> 24),
		byte(bits >> 16),
		byte(bits >> 8),
		byte(bits),
	}
}
