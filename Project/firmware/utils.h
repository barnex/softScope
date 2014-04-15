#ifndef _UTILS_H
#define _UTILS_H

#include <stdint.h>
#include "leds.h"


// ACHTUNG: Only works if:
//  * nBytes is a multiple of 4
//  * source and destination are word aligned
void memcpy32(uint32_t *dest, uint32_t *src, uint32_t nBytes);

// Goes into an infinite loop flashing the error led.
void panic();

void* emalloc(int nBytes);

#endif
