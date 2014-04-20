package softscope

import(
	"strconv"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}


func atouint32(a string)uint32{
	i, err := strconv.Atoi(a)
	check(err)
	if i < 0 {panic("not an uint32")}
	return uint32(i)
}
