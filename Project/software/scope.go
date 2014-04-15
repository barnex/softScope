package softscope

import(
	"log"
)

var(
	tty TTY
)


func Init(ttyDev string, baud int) {
	log.SetFlags(0)
	log.Println("Using", ttyDev, "@", baud, "baud")

	_, err := OpenTTY(ttyDev, baud)
	if err != nil{
		log.Fatal(err)
	}
}
