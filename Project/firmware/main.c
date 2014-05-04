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

// ADC data buffer _____________________________________________________________________________________

uint16_t volatile *adc1Buf = NULL;   // ADC data buffer, to be accessed via startChunk, nextChunk.
uint16_t volatile *adc2Buf = NULL;   
volatile int       adcPos  = 0; // ADC chunk number currently being written.

// Called at the end of TIM3_IRQHandler.
void TIM3_IRQHook() {
	adcPos += ADC_CHUNKSIZE;
	if (adcPos >= ADC_BUFSIZE) {
		adcPos = 0;
		LEDToggle(LED_OK);
	}
}

// Return pointer to usable adc data chunk, length ADC_CHUNKSIZE
// Pointer will be in chunk right behind ADC position.
uint16_t volatile *startChunk() {
	int c = adcPos - (adcPos % ADC_CHUNKSIZE) - ADC_CHUNKSIZE;
	if(c < 0) {
		c = ADC_BUFSIZE - ADC_CHUNKSIZE;
	}
	return &adc1Buf[c];
}

// Return pointer to usable adc data chunk, length ADC_CHUNKSIZE,
// given the pointer to the previous chunk. Waits until ADC has
// advanced enough.
uint16_t volatile* nextChunk(uint16_t volatile* current) {
	uint16_t volatile* a = &current[ADC_CHUNKSIZE];
	if(a >= &adc1Buf[ADC_BUFSIZE]) {
		a = &adc1Buf[0];
	}
	uint16_t volatile* b = &a[ADC_CHUNKSIZE];
	while(&adc1Buf[adcPos] >= a && &adc1Buf[adcPos] < b) {
		//wait for ADC to exit upcoming chunk
	}
	return a;
}


// check that ADC write position is not in [a, a+ADC_CHUNKSIZE[,
// which would be timing error if using that data.
void checkTiming(uint16_t volatile *chunk) {
	if(&adc1Buf[adcPos] >= &chunk[0] && &adc1Buf[adcPos] < &chunk[ADC_CHUNKSIZE]) {
		error(UNMET_TIMING, 0); // value: by how much samples timing was missed
	}
}

// init ________________________________________________________________________________________________

void init() {
	NVIC_PriorityGroupConfig( NVIC_PriorityGroup_4 );
	init_LEDs();

	// initial settings
	timebase = 12000;
	nSamples = 512;

	// ADC
	adc1Buf = emalloc(ADC_BUFSIZE*sizeof(adc1Buf[0]));
	adc2Buf = emalloc(ADC_BUFSIZE*sizeof(adc2Buf[0]));
	memset((void*)adc1Buf, 0, ADC_BUFSIZE*sizeof(adc1Buf[0]));
	memset((void*)adc2Buf, 0, ADC_BUFSIZE*sizeof(adc2Buf[0]));

	init_clock(timebase, ADC_CHUNKSIZE);
	clock_TIM3_IRQHook = TIM3_IRQHook;  // Register TIM3_IRQHook to be called at the end of TIM3_IRQHandler

	init_analogIn();
	init_ADC(adc1Buf, adc2Buf, ADC_BUFSIZE);

	init_USART1(115200);
	init_outbox();
	init_inbox();

	enable_clock();
}

// main ___________________________________________________________________________________________________



int main(void) {
	init();

	for(;;) {


		while(reqFrames == 0) {
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

		// trigger
		uint16_t volatile* c = startChunk();
		uint16_t last = c[0];

		// for the moment only look in one chunk
		int i = 0;  // i is trigger point
		for(i=0; i<ADC_CHUNKSIZE; i++) {
			uint16_t current = c[i];
			if (last < triglev && current >= triglev) {
				break;
			}
			last = current;
		}


		// copy partial chunk
		int o = 0;  // output numbers copied
		memcpy((void*)(&outData[o]), (void*)(&c[i]), (ADC_CHUNKSIZE-i)*sizeof(adc1Buf[0]));
		o += (ADC_CHUNKSIZE-i);

		// copy the rest
		for(; o < hdr->nsamples; o+=ADC_CHUNKSIZE) {
			c = nextChunk(c);
			memcpy((void*)(&outData[o]), (void*)(c), ADC_CHUNKSIZE*sizeof(adc1Buf[0])); // may copy a bit too much
		}

		outbox_TX(hdr->nbytes);
	}
}



