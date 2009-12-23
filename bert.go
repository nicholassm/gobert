package bert

import (
  "fmt";
	"io";
  "math";
	"os";
	"reflect";
	"strings";
	"encoding/binary";
)

const (
  SMALL_INT    = 97;
  INT          = 98;
  FLOAT        = 99;
  ATOM         = 100;
  SMALL_TUPLE  = 104;
  LARGE_TUPLE  = 105;
  NIL          = 106;
  STRING       = 107;
  LIST         = 108;
  BIN          = 109;
  SMALL_BIGNUM = 110;
  LARGE_BIGNUM = 111;
  MAGIC        = 131;
)

func writeByte(w io.Writer, data byte) os.Error {
  _, err := w.Write([]byte {data});
  return err;
}

func writeAtom(w io.Writer, atom string) os.Error {
  writeByte(w, ATOM);
  writeByte(w, uint8(len(atom)));
  io.WriteString(w, atom);
  
  return nil;
}

func writeTuple(w io.Writer, tuple []interface{}) os.Error {
  if len(tuple) <= 256 {
    writeByte(w, SMALL_TUPLE);
  } else {
    writeByte(w, LARGE_TUPLE);
  }
  
  binary.Write(w, binary.BigEndian, len(tuple));
	for i := 0; i < len(tuple); i++ {
		if err := writeValue(w, reflect.NewValue(tuple[i])); err != nil {
			return err
		}
	}
	
	return nil;
}

func writeNil(w io.Writer) os.Error {
  w.Write([]byte {SMALL_TUPLE, 2});
  writeAtom(w, "bert");
  writeAtom(w, "nil");
  
  return nil
}

func writeBool(w io.Writer, v bool) os.Error {
  w.Write([]byte {SMALL_TUPLE, 2});
  writeAtom(w, "bert");
  
  if v {
    writeAtom(w, "true")
  } else {
    writeAtom(w, "false")
  }
  
  return nil
}


func writeUint8(w io.Writer, v uint8) (err os.Error) {
  _, err = w.Write([]byte {SMALL_INT, v});
  return
}

func writeUint16(w io.Writer, v uint16) (err os.Error) {
  _, err = w.Write([] byte { SMALL_INT,
    byte(v >> 8),
    byte(v)
  });
  return
}

func writeUint32(w io.Writer, v uint32) (err os.Error) {
  _, err = w.Write([] byte {
    SMALL_INT,
    byte(v >> 24),
    byte(v >> 16),
    byte(v >> 8),
    byte(v)
  });
  return
}

func writeUint64(w io.Writer, v uint64) (err os.Error) {
  _, err = w.Write([] byte {
    SMALL_BIGNUM,
    8,
    0,
    byte(v >> 56),
    byte(v >> 48),
    byte(v >> 40),
    byte(v >> 32),
    byte(v >> 24),
    byte(v >> 16),
    byte(v >> 8),
    byte(v)
  });
  return
}

func writeInt64(w io.Writer, v int64) (err os.Error) {
  var sign byte;
  
  if v > 0 {
    sign = 0
  } else {
    sign = 1;
    v    = -v
  };
  
  _, err = w.Write([] byte {
    SMALL_BIGNUM,
    8,
    sign,
    byte(v),
    byte(v >> 8),
    byte(v >> 16),
    byte(v >> 24),
    byte(v >> 32),
    byte(v >> 40),
    byte(v >> 48),
    byte(v >> 56)
  });
  return
}

func writeBinary(w io.Writer, data []byte) os.Error {
  l := len(data);
  
  if l > math.MaxInt32 {
    return os.NewError(fmt.Sprintf("bert.Encode: binary too large: %d (max is %d)", l, math.MaxInt32));
  }
	
  writeByte(w, BIN);
  binary.Write(w, binary.BigEndian, uint32(l));
  w.Write(data);
  
  return nil
}

func writeFloat(w io.Writer, v string) os.Error {
  writeByte(w, FLOAT);
  fmt.Fprintf(w, v);
  
  return nil
}

func writeArrayOrSlice(w io.Writer, val reflect.ArrayOrSliceValue) os.Error {
  writeByte(w, LIST);
  binary.Write(w, binary.BigEndian, uint32(val.Len()));
  
	for i := 0; i < val.Len(); i++ {
		if err := writeValue(w, val.Elem(i)); err != nil {
			return err
		}
	}

	writeByte(w, NIL);
	return nil;
}

func writeMap(w io.Writer, val *reflect.MapValue) os.Error {
  writeByte(w, SMALL_TUPLE);
  writeByte(w, 3);
  writeAtom(w, "bert");
  writeAtom(w, "dict");
  
	keys := val.Keys();
	writeByte(w, LIST);
  binary.Write(w, binary.BigEndian, uint32(len(keys)));
	for i := 0; i < len(keys); i++ {
	  w.Write([]byte {SMALL_TUPLE, 2});
		writeValue(w, keys[i]);
		writeValue(w, val.Elem(keys[i]));
	}
	writeByte(w, NIL);

	return nil;
}

// [ag] Peculiarity: binary.Write(io.Writer, interface{}) can only write "fixed-
// size values". A fixed-size value is either a fixed-size integer (int8, uint8,
// int16, uint16, ...) or an array or struct containing only fixed-size values.
// Thus in the big switch, matching on reflect.IntValue is no good. Workaround?
func writeValue(w io.Writer, val reflect.Value) (err os.Error) {
	switch v := val.(type) {
	case nil:                   return writeNil(w)
	case *reflect.BoolValue:    return writeBool(w, v.Get())
	case *reflect.Uint8Value:	  return writeUint8(w, v.Get())
	case *reflect.Int8Value:	  return writeUint8(w, uint8(v.Get()))
	case *reflect.Uint16Value:  return writeUint16(w, v.Get())
	case *reflect.Int16Value:	  return writeUint16(w, uint16(v.Get()))
	case *reflect.Uint32Value:  return writeUint32(w, v.Get())
	case *reflect.Int32Value:	  return writeUint32(w, uint32(v.Get()))
	case *reflect.Uint64Value:  return writeUint64(w, v.Get())
	case *reflect.Int64Value:	  return writeInt64(w, v.Get())
	case *reflect.FloatValue:   return writeFloat(w, fmt.Sprintf("%.20e", v.Get()))
	case *reflect.Float32Value: return writeFloat(w, fmt.Sprintf("%.20e", v.Get()))
	case *reflect.Float64Value: return writeFloat(w, fmt.Sprintf("%.20e", v.Get()))
  case *reflect.StringValue:  return writeBinary(w, strings.Bytes(v.Get()))
	case *reflect.ArrayValue:   return writeArrayOrSlice(w, v)
	case *reflect.SliceValue:   return writeArrayOrSlice(w, v)
	case *reflect.MapValue:     return writeMap(w, v)
	default:
		return os.NewError("bert.Encode: invalid type " + v.Type().String())
	}
	
	return nil;
}

func Encode(w io.Writer, val interface{}) os.Error {
  writeByte(w, MAGIC);
	return writeValue(w, reflect.NewValue(val))
}

func Decode(r io.Reader) (data interface{}, err os.Error) {
  return nil, nil;
}
