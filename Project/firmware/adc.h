#ifndef _ADC_H_
#define _ADC_H_

#include <stdint.h>

// Initialize the scope's analog input pin PA1
void init_analogIn();

// Initialize the ADC to write to samplesBuffer with size samples
void init_ADC(volatile uint16_t *samplesBuffer, int samples);

#endif
