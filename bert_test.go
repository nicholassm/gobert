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
  out         := new(bytes.Buffer);
  Encode(out, val);
  assertBytesEqual(t, message, expected, out.Bytes());
}

func TestEncodeNil(t *testing.T) {
  expected := []byte {SMALL_TUPLE, 2, ATOM, 4, 'b', 'e', 'r', 't', ATOM, 3, 'n', 'i', 'l'};
  testEncode(t, "Encoding nil failed", nil, expected);
}

func TestEncodeTrue(t *testing.T) {
  expected := []byte {SMALL_TUPLE, 2, ATOM, 4, 'b', 'e', 'r', 't', ATOM, 4, 't', 'r', 'u', 'e'};
  testEncode(t, "Encoding true failed", true, expected);
}

func TestEncodeFalse(t *testing.T) {
  expected := []byte {SMALL_TUPLE, 2, ATOM, 4, 'b', 'e', 'r', 't', ATOM, 5, 'f', 'a', 'l', 's', 'e'};
  testEncode(t, "Encoding true failed", false, expected);
}

func TestEncodeFloat(t *testing.T) {
  expected := bytes.Add([]byte {FLOAT}, strings.Bytes(fmt.Sprintf("%.20e", 1.2e3)));
  testEncode(t, "Encoding float failed", 1.2e3, expected);
}
