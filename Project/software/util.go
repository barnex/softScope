package softscope

import (
	"strconv"
	"log"
)

func debug(msg...interface{}){
	if *flag_debug{
		log.Println(msg...)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func atouint32(a string) uint32 {
	i, err := strconv.Atoi(a)
	check(err)
	if i < 0 {
		panic("not an uint32")
	}
	return uint32(i)
}
