#include "adc.h"
#include "clock.h"
#include "error.h"
#include "inbox.h"
#include "leds.h"
#include "outbox.h"
#include "usart.h"


static message_t incoming;
static int rxByte = 0;

void setSamples(uint32_t s) {
	if(s > MAX_NSAMPLES) {
		s = MAX_NSAMPLES;
	}
	if(s < MIN_NSAMPLES) {
		s = MIN_NSAMPLES;
	}
	nSamples = s;
}

void setTriglev(uint32_t t) {
	if(t > 4095) {
		t = 4095;
	}
	if(t < 1) {
		t = 1;
	}
	triglev = t;
}

void setTimebase(uint32_t p) {
	if(p < MIN_CLOCK_PERIOD){
		p = MIN_CLOCK_PERIOD;
	}
	init_clock(p, ADC_CHUNKSIZE);
	enable_clock();
}

static void handleIncoming() {
	LEDOff(LED_ERR);

	// check magic number and ignore (but signal) bad message
	if(incoming.magic != MSG_MAGIC){
		error(BAD_MAGIC, incoming.magic);
		return;
	}

	switch(incoming.command) {
	default:
		error(BAD_COMMAND, incoming.command);
		break;
	case SET_SAMPLES:
		setSamples(incoming.value);
		break;
	case SET_TIMEBASE:
		setTimebase(incoming.value);
		break;
	case SET_TRIGLEV:
		setTriglev(incoming.value);
		break;
	case CLEAR_ERR:
		hdr->errno = 0;
		hdr->errval = 0;
	case REQ_FRAMES:
		reqFrames = incoming.value;
		break;
	}
}

// called upon usart RX to handle incoming byte
// writes to _inbuf until full, then copies to inbox
static void RXHandler(uint8_t data) {
	LEDOn(LED_RX);

	uint8_t *inArr = (uint8_t*)(&incoming);
	inArr[rxByte] = data;
	rxByte++;

	if (rxByte == sizeof(message_t)) {
		rxByte = 0;
		LEDOff(LED_RX);
		handleIncoming();
	}
}

void init_inbox() {
	USART1_RXHandler = RXHandler;
}
