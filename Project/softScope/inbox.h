#ifndef _INBOX_H_
#define _INBOX_H_

#include <stdint.h>

typedef struct{
	uint32_t samples;
	uint32_t timeBase;
	uint32_t trigLev;
}inbox_t;


volatile inbox_t inbox;

void init_inbox();


#endif
