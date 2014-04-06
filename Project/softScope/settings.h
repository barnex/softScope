#ifndef _SETTINGS_H_
#define _SETTINGS_H_

#include <stdint.h>

#define ADC_BUFSIZE	    2048
#define IR_PERIOD       128
#define MAX_NSAMPLES    (ADC_BUFSIZE/2)
#define MIN_NSAMPLES    (IR_PERIOD)

volatile uint32_t samples ;
volatile uint32_t timebase;
volatile uint32_t triglev ;

void init_inbox();

#endif
