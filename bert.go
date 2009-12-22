package bert

import (
  "fmt";
	"io";
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
  MAX_INT      = (1 << 27) -1;
  MIN_INT      = -(1 << 27);
)

func writeByte(w io.Writer, data byte) os.Error {
  _, err := w.Write([]byte {data});
  return err;
}

// Write string as atom to w.
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

func writeArrayOrSlice(w io.Writer, val reflect.ArrayOrSliceValue) os.Error {
  writeByte(w, LIST);

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
	for i := 0; i < len(keys); i++ {
		writeValue(w, keys[i]);
		writeValue(w, val.Elem(keys[i]));
	}

	return nil;
}

// [ag] Peculiarity: binary.Write(io.Writer, interface{}) can only write "fixed-
// size values". A fixed-size value is either a fixed-size integer (int8, uint8,
// int16, uint16, ...) or an array or struct containing only fixed-size values.
// Thus in the big switch, matching on reflect.IntValue is no good. Workaround?
func writeValue(w io.Writer, val reflect.Value) (err os.Error) {
	switch v := val.(type) {
	case *reflect.Int8Value:
	  // SMALL_INT
	  writeByte(w, SMALL_INT);
	  writeByte(w, uint8(v.Get()));
	case *reflect.Int32Value:
	  // INT
	  writeByte(w, INT);
	  binary.Write(w, binary.BigEndian, v.Get());
	case nil:
	  // NIL
    writeByte(w, SMALL_TUPLE);
    writeByte(w, 2);
	  writeAtom(w, "bert");
	  writeAtom(w, "nil")
	case *reflect.FloatValue:
	  // FLOAT
	  writeByte(w, FLOAT);
	  fmt.Fprintf(w, "%.20e", v.Get());
  case *reflect.StringValue:
	  // STRING
	  bytes := strings.Bytes(v.Get());
	  writeByte(w, BIN);
	  writeByte(w, uint8(len(bytes)));
	  w.Write(bytes);
	case *reflect.ArrayValue:
	  // LIST
	case *reflect.Int64Value:
	  // SMALL_BIGNUM
	case *reflect.MapValue:
	  // {bert, dict, keysAndValues}
    writeByte(w, SMALL_TUPLE);
    writeByte(w, 3);
	  writeAtom(w, "bert");
	  writeMap(w, v);
	case *reflect.BoolValue:
	  // {bert, true} or {bert, false}
    writeByte(w, SMALL_TUPLE);
    writeByte(w, 2);
	  writeAtom(w, "bert");
	  if v.Get() {
	    writeAtom(w, "true");
	  } else {
	    writeAtom(w, "false");
	  }
	default:
		return os.NewError("bert.Encode: invalid type " + v.Type().String())
	}
	
	return nil;
}

func Encode(w io.Writer, val interface{}) os.Error {
	return writeValue(w, reflect.NewValue(val))
}

func Decode(r io.Reader) (data interface{}, err os.Error) {
  return nil, nil;
}
