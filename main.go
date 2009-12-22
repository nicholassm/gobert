package main

import (
  "os";
  "bert";
)

func main() {
  bert.Encode(os.Stdout, int8(1));
  bert.Encode(os.Stdout, int32(1));
  bert.Encode(os.Stdout, true);
  bert.Encode(os.Stdout, 1.2e3);
}
