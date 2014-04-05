#include "inbox.h"
#include "usart.h"
#include "leds.h"

volatile inbox_t inbox;   

//enum{
	//FRAMESYNC = 1;
//};

//int rxState = FRAMESYNC;

typedef struct{
	uint32_t magic;
	uint32_t command;
	uint32_t value;
}message_t;

static message_t incoming;
static int rxByte = 0;

// called upon usart RX to handle incoming byte
// writes to _inbuf until full, then copies to inbox
static void RXHandler(uint8_t data){
	LEDOn(LED1);

	uint8_t *inArr = (uint8_t*)(&incoming);	
	inArr[rxByte] = data;

	if (rxByte == sizeof(message_t)){
		rxByte = 0;
		LEDOff(LED1);
	}
}

void init_inbox(){
	USART1_RXHandler = RXHandler;
}
