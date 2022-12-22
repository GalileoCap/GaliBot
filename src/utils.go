package main;

import (
  "io"
  "net/http"
);

var CurrIP string;

func contains(array []int64, value int64) bool {
	for _, v := range array {
		if v == value {
			return true;
		}
	}
	return false;
}

func max(x int, y int) int {
  if x > y {
    return x;
  }
  return y;
}

func min(x int, y int) int {
  if x < y {
    return x;
  }
  return y;
}

func updateIP() (bool, error) {
  resp, err := http.Get("https://api.ipify.org?format=text");
  if err != nil {
    return false, err;
  }
  defer resp.Body.Close();

  newIP_b, err := io.ReadAll(resp.Body);
  if err != nil {
    return false, err;
  }

  newIP, prevIP := string(newIP_b), CurrIP;
  CurrIP = newIP;

  return CurrIP != prevIP, nil;
}
