#include <misc.h>
#include <stdint.h>
#include <string.h>
#include <stdlib.h>

#include "adc.h"
#include "clock.h"
#include "inbox.h"
#include "leds.h"
#include "outbox.h"
#include "usart.h"
#include "utils.h"

volatile uint16_t *samplesBuffer;

volatile int adcPos = 0;

// Called at the end of TIM3_IRQHandler.
void TIM3_IRQHook() {
	adcPos += IR_PERIOD;
	if (adcPos > ADC_BUFSIZE) {
		adcPos = 0;
		LEDToggle(LED_OK);
	}
}


void init() {
	NVIC_PriorityGroupConfig( NVIC_PriorityGroup_4 );
	init_LEDs();

	timebase = 420;

	// ADC
	samplesBuffer   = emalloc(ADC_BUFSIZE*sizeof(samplesBuffer[0]));
	memset((void*)samplesBuffer, 0, ADC_BUFSIZE*sizeof(samplesBuffer[0]));

	init_clock(timebase, IR_PERIOD);
	clock_TIM3_IRQHook = TIM3_IRQHook;  // Register TIM3_IRQHook to be called at the end of TIM3_IRQHandler

	init_analogIn();
	init_ADC(samplesBuffer, ADC_BUFSIZE);

	init_USART1(115200);
	init_outbox();
	init_inbox();

	enable_clock();
}

int main(void) {
	init();

	for(;;) {

		while(transmitting) {
			// wait until previous transmission finished
		}

		while(reqFrames == 0){
			// wait until frame is requested
		}
		reqFrames--; // TODO: is not atomic

		// Make non-volatile copies of settings
		hdr->nsamples = samples;
		hdr->bitdepth = 16; // TODO: variable
		hdr->nchans = 1; // TODO: variable
		hdr->nbytes = sizeof(uint16_t) * hdr->nchans * hdr->nsamples;
		
		//hdr->trigLev = triglev;
		//hdr->timeBase = timebase;

		memcpy((void*)(outData), (void*)samplesBuffer, hdr->nbytes);
		outbox_TX(hdr->nbytes);
	}
}



