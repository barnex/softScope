#ifndef _SETTINGS_H_
#define _SETTINGS_H_

#include <stdint.h>


#define MAX_SAMPLES	    1024              // Number of samples for each acquisition/frame
#define MIN_SAMPLES	     128
#define IR_PERIOD      (MAX_SAMPLES/4)

volatile uint32_t samples ;
volatile uint32_t timebase;
volatile uint32_t triglev ;

void init_inbox();

#endif
