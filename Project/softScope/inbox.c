#include "inbox.h"
#include "usart.h"
#include "leds.h"

volatile header_t inbox;   

static volatile header_t _inbuf;  // receive buffer for incoming usart communication, not to be used
static volatile int _inpos = 0;   // position where next received byte should go in headerArray

// called upon usart RX to handle incoming byte
// writes to _inbuf until full, then copies to inbox
void myRXHandler(uint8_t data){
		LEDOn(LED1);
		uint8_t* arr = (uint8_t*)(&_inbuf);
 		arr[_inpos] = data;
		_inpos++;
		if (_inpos >= sizeof(header_t)){
			_inpos = 0;
			inbox = _inbuf; // header complete, copy to visible header

		}
		LEDOff(LED1);
}

void init_inbox(){
	USART1_RXHandler = myRXHandler;
}
