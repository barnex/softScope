#include <misc.h>
#include <stdint.h>
#include <string.h>
#include <stdlib.h>

#include "adc.h"
#include "clock.h"
#include "leds.h"
#include "usart.h"
#include "utils.h"
#include "settings.h"

volatile uint16_t *samplesBuffer;


// -- outbound communication
#define HEADER_WORDS  8

// Frame data header
typedef struct {
	uint32_t magic;                   // identifies start of header, 0xFFFFFFFF
	uint32_t samples;                 // number of samples
	uint32_t trigLev;
	uint32_t timeBase;
	uint32_t padding[HEADER_WORDS-4]; // unused space, needed for correct total size
} header_t;

static uint8_t *usartBuf;   // embeds outbox and outdata so they can be TX'ed in one go
static header_t *outbox;    // header written to software, first part of usartBuf
static uint16_t *outData;   // data written to software, second part of usartBuf



// Called at the end of TIM3_IRQHandler.
void TIM3_IRQHook() {
	LEDToggle(LED_OK);
}

void init() {
	NVIC_PriorityGroupConfig( NVIC_PriorityGroup_4 );

	// ADC
	samplesBuffer   = malloc(MAX_SAMPLES*sizeof(samplesBuffer[0]));
	memset((void*)samplesBuffer, 0, MAX_SAMPLES*sizeof(samplesBuffer[0]));

	// outbound communication
	int headerBytes = sizeof(header_t);
	int dataBytes   = MAX_SAMPLES*sizeof(outData[0]);
	usartBuf	    = malloc(headerBytes + dataBytes);
	memset(usartBuf, 0, headerBytes + dataBytes);

	outbox = (header_t*)(usartBuf);                      // header is embedded in beginning of usart buffer
	outData = (uint16_t*)(&usartBuf[headerBytes]);       // data is embedded next

	init_clock(timebase, IR_PERIOD);
	clock_TIM3_IRQHook = TIM3_IRQHook;  // Register TIM3_IRQHook to be called at the end of TIM3_IRQHandler
	init_ADC(samplesBuffer, MAX_SAMPLES);
	init_USART1(115200);
	init_inbox();
	init_analogIn();
	init_LEDs();

	enable_clock();
}

int main(void) {
	init();

	for(;;) {

		while(transmitting) {
			// wait
		}

		volatile int c = 2000000;
		while(c>0) {
			c--;
		}

		memcpy((void*)(outData), (void*)samplesBuffer, samples * sizeof(samplesBuffer[0]));

		outbox->magic = 0xFFFFFFFF;
		outbox->samples = samples;
		outbox->trigLev = triglev;
		outbox->timeBase = timebase;
		USART_asyncTX(usartBuf, sizeof(header_t) + MAX_SAMPLES * sizeof(samplesBuffer[0])); // todo: transfer samples

	}
}



