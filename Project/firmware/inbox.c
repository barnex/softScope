#include "inbox.h"

#include "adc.h"
#include "clock.h"
#include "leds.h"
#include "usart.h"

volatile uint32_t samples    = 512;    // TODO(a): keep below MAX_NSAMPLES
volatile uint32_t timebase   = 4200;
volatile uint32_t trigLev    = (1<<10);


typedef struct {
	uint32_t magic;
	uint32_t command;
	uint32_t value;
} message_t;

static message_t incoming;
static int rxByte = 0;

enum {
    SAMPLES=  1,
    TIMEBASE= 2,
    TRIGLEV=  3,
    SOFTGAIN= 4
};

void setSamples(uint32_t s) {
	if(s > MAX_NSAMPLES) {
		s = MAX_NSAMPLES;
	}
	if(s < MIN_NSAMPLES) {
		s = MIN_NSAMPLES;
	}
	samples = s;
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
	if(p < 42){
		p = 42; // 1MHz
	}
	if(p>42000){
		p = 42000;
	}
	init_clock(p, IR_PERIOD);
	enable_clock();
}

static void handleIncoming() {
	LEDOff(LED_ERR);
	switch(incoming.command) {
	default:
		LEDOn(LED_ERR);
		break;
	case SAMPLES:
		setSamples(incoming.value);
	case TIMEBASE:
		setTimebase(incoming.value);
	case TRIGLEV:
		setTriglev(incoming.value);
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
