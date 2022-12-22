package main;

func contains(array []int64, value int64) bool {
  for _, v := range array {
    if v == value {
      return true;
    }
  }
  return false;
}
