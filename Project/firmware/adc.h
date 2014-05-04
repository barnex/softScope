#ifndef _ADC_H_
#define _ADC_H_

#include <stdint.h>

// Initialize the scope's analog input pin PA1
void init_analogIn();

// Initialize the ADC to write to adcXbuf with size samples
void init_ADC(uint16_t volatile *adc1buf, uint16_t volatile *adc2buf, int samples);

#endif
