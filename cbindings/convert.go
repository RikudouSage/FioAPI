package main

/*
#include <stdint.h>
*/
import "C"
import "time"

func dateToC(date time.Time) C.uint64_t {
	return C.uint64_t(date.UnixMilli())
}

func dateFromC(date C.uint64_t) time.Time {
	return time.UnixMilli(int64(date))
}
