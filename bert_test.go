package bert

import (
  "fmt";
  "strings";
  "bytes";
  "testing";
)

func assertBytesEqual(t *testing.T, message string, expected, actual []byte) {
  if !bytes.Equal(expected, actual) {
    t.Errorf("%s: '%s' != '%s'", message, expected, actual);
  }
}

func testEncode(t *testing.T, message string, val interface{}, expected []byte) {
  out := new(bytes.Buffer);
  Encode(out, val);
  assertBytesEqual(t, message, expected, out.Bytes());
}

func TestEncodeNil(t *testing.T) {
  expected := []byte {MAGIC, SMALL_TUPLE, 2, ATOM, 4, 'b', 'e', 'r', 't', ATOM, 3, 'n', 'i', 'l'};
  testEncode(t, "Encoding nil failed", nil, expected);
}

func TestEncodeTrue(t *testing.T) {
  expected := []byte {MAGIC, SMALL_TUPLE, 2, ATOM, 4, 'b', 'e', 'r', 't', ATOM, 4, 't', 'r', 'u', 'e'};
  testEncode(t, "Encoding true failed", true, expected);
}

func TestEncodeFalse(t *testing.T) {
  expected := []byte {MAGIC, SMALL_TUPLE, 2, ATOM, 4, 'b', 'e', 'r', 't', ATOM, 5, 'f', 'a', 'l', 's', 'e'};
  testEncode(t, "Encoding true failed", false, expected);
}

func TestEncodeFloat(t *testing.T) {
  expected := bytes.Add([]byte {MAGIC, FLOAT}, strings.Bytes(fmt.Sprintf("%.20e", float(1.2e3))));
  testEncode(t, "Encoding float failed", 1.2e3, expected);
}

func TestEncodeFloat32(t *testing.T) {
  expected := bytes.Add([]byte {MAGIC, FLOAT}, strings.Bytes(fmt.Sprintf("%.20e", float32(1.2e3))));
  testEncode(t, "Encoding float32 failed", 1.2e3, expected);
}

func TestEncodeFloat64(t *testing.T) {
  expected := bytes.Add([]byte {MAGIC, FLOAT}, strings.Bytes(fmt.Sprintf("%.20e", float64(1.2e3))));
  testEncode(t, "Encoding float64 failed", 1.2e3, expected);
}

func TestEncodeString(t *testing.T) {
  expected := bytes.Add([]byte {MAGIC, BIN, 0, 0, 0, 3}, strings.Bytes("foo"));
  testEncode(t, "Encoding string failed", "foo", expected);
}

func TestEncodeArray(t *testing.T) {
  expected := []byte {MAGIC, LIST, 0, 0, 0, 2, SMALL_INT, 1, SMALL_INT, 2, NIL};
  testEncode(t, "Encoding array failed", []byte {1, 2}, expected);
}

func TestEncodeMap(t *testing.T) {
  expected := []byte {MAGIC, SMALL_TUPLE, 3, ATOM, 4, 'b', 'e', 'r', 't', ATOM, 4, 'd', 'i', 'c', 't', LIST, 0, 0, 0, 1, SMALL_TUPLE, 2, SMALL_INT, 1, SMALL_INT, 2, NIL};
  testEncode(t, "Encoding map failed", map[byte]byte {1: 2}, expected);
}
