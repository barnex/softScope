#ifndef _OUTBOX_H_
#define _OUTBOX_H_

#include <stdint.h>

// -- outbound communication
#define HEADER_WORDS  8

// Frame data header
typedef struct {
	uint32_t magic;                   // identifies start of header, 0xFFFFFFFF
	uint32_t samples;                 // number of samples
	uint32_t trigLev;
	uint32_t timeBase;
	uint32_t padding[HEADER_WORDS-4]; // unused space, needed for correct total size
} header_t;

header_t *outHeader;  // header written to software, first part of usartBuf
uint16_t *outData;    // data written to software, second part of usartBuf // TODO(a): bytes

void init_outbox();
void outbox_TX(uint32_t nDataBytes);

#endif
