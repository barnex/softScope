#include "error.h"
#include "leds.h"
#include "outbox.h"

void error(uint32_t code, uint32_t value){
	LEDOn(LED_ERR);
	hdr->errno = code;
	hdr->errval = value;	
}


