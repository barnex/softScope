#ifndef _INBOX_H_
#define _INBOX_H_

#include <stdint.h>

#define ADC_BUFSIZE	    2048
#define IR_PERIOD       128
#define MAX_NSAMPLES    (ADC_BUFSIZE/2)
#define MIN_NSAMPLES    (IR_PERIOD)

// Number of requested frames. If != 0, send frame and decrement.
volatile uint32_t reqFrames;

volatile uint32_t samples ;
volatile uint32_t timebase;
volatile uint32_t triglev ;

void init_inbox();

#endif
