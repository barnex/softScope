#ifndef _INBOX_H_
#define _INBOX_H_

#include <stdint.h>


// -- I/O protocol
#define HEADER_WORDS  8


// Frame data header
typedef struct{
	uint32_t magic;                   // identifies start of header, 0xFFFFFFFF
	uint32_t samples;                 // number of samples
	uint32_t trigLev;
	uint32_t timeBase;
	uint32_t padding[HEADER_WORDS-4]; // unused space, needed for correct total size
} header_t;


volatile header_t inbox;   // last header sent from software, values can be used

void init_inbox();


#endif
