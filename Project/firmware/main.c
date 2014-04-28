#include <misc.h>
#include <stdint.h>
#include <string.h>
#include <stdlib.h>

#include "adc.h"
#include "clock.h"
#include "error.h"
#include "inbox.h"
#include "leds.h"
#include "outbox.h"
#include "usart.h"
#include "utils.h"

uint16_t volatile *_samplesBuffer = NULL;  // to be accessed via nextChunk
volatile int       _adcPos        = 0;

// Called at the end of TIM3_IRQHandler.
void TIM3_IRQHook() {
	_adcPos += IR_PERIOD;
	if (_adcPos >= ADC_BUFSIZE) {
		_adcPos = 0;
		LEDToggle(LED_OK);
	}
}

// return index for samplesBuffer where usable chunk starts (up to chunk + IR_PERIOD)
int currentChunk(){
	int c = _adcPos - (_adcPos % IR_PERIOD) - IR_PERIOD;
	if(c < 0){
		c = ADC_BUFSIZE - IR_PERIOD;
	}
	return c;
}

// given an index in samplesbuffer, return index of next chunk (current + IR_PERIOD),
// but wait until ADC is not writing there anymore.
int nextChunk(int current){
	int a = current + IR_PERIOD;
	if(a >= ADC_BUFSIZE){
		a = 0;
	}
	int b = a + IR_PERIOD;
	while(_adcPos >= a && _adcPos < b){
		//wait for ADC to exit upcoming chunk
	}
	return a;
}

void init() {
	NVIC_PriorityGroupConfig( NVIC_PriorityGroup_4 );
	init_LEDs();

	// initial settings
	timebase = 12000;
	nSamples = 512;

	// ADC
	_samplesBuffer   = emalloc(ADC_BUFSIZE*sizeof(_samplesBuffer[0]));
	memset((void*)_samplesBuffer, 0, ADC_BUFSIZE*sizeof(_samplesBuffer[0]));

	init_clock(timebase, IR_PERIOD);
	clock_TIM3_IRQHook = TIM3_IRQHook;  // Register TIM3_IRQHook to be called at the end of TIM3_IRQHandler

	init_analogIn();
	init_ADC(_samplesBuffer, ADC_BUFSIZE);

	init_USART1(115200);
	init_outbox();
	init_inbox();

	enable_clock();
}

// check that ADC write position is not in [a, a+IR_PERIOD[,
// which would be timing error if using that data.
void checkTiming(int a){
	if(_adcPos >= a && _adcPos < a + IR_PERIOD){
		error(UNMET_TIMING, _adcPos-a); // value: by how much samples timing was missed
	}
}

int main(void) {
	init();

	for(;;) {


		while(reqFrames == 0){
			// wait until frame is requested
		}
		reqFrames--; // TODO: is not atomic

		while(transmitting) {
			// wait until previous transmission finished
		}

		// Make non-volatile copies of settings
		hdr->nsamples = nSamples;
		hdr->bitdepth = 16; // TODO: variable
		hdr->nchans = 1;    // TODO: variable
		hdr->nbytes = sizeof(uint16_t) * hdr->nchans * hdr->nsamples;
		hdr->trigLev = triglev;
		hdr->timeBase = timebase;

		// go through samplesbuffer in small chunks,
		// trailing behind the ADC
		int c = currentChunk();
		for(int n = 0; n < hdr->nsamples; n+=IR_PERIOD){
			checkTiming(c);
			memcpy((void*)(&outData[n]), (void*)(&_samplesBuffer[c]), IR_PERIOD*sizeof(_samplesBuffer[0]));
			checkTiming(c);
			c = nextChunk(c);
		}

		outbox_TX(hdr->nbytes);
	}
}



