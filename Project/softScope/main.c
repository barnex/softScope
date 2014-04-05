#include <misc.h>
#include <stdint.h>
#include <string.h>
#include <stdlib.h>

#include "adc.h"
#include "clock.h"
#include "leds.h"
#include "outbox.h"
#include "settings.h"
#include "usart.h"
#include "utils.h"

volatile uint16_t *samplesBuffer;


// Called at the end of TIM3_IRQHandler.
void TIM3_IRQHook() {
	LEDToggle(LED_OK);
}

void init() {
	NVIC_PriorityGroupConfig( NVIC_PriorityGroup_4 );

	// ADC
	samplesBuffer   = malloc(MAX_SAMPLES*sizeof(samplesBuffer[0]));
	memset((void*)samplesBuffer, 0, MAX_SAMPLES*sizeof(samplesBuffer[0]));

	init_outbox();
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

		outHeader->magic = 0xFFFFFFFF;
		outHeader->samples = samples;
		outHeader->trigLev = triglev;
		outHeader->timeBase = timebase;
		outbox_TX(MAX_SAMPLES*sizeof(samplesBuffer[0]));


	}
}



