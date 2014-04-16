#ifndef _OUTBOX_H_
#define _OUTBOX_H_

#include <stdint.h>

// size of header_t in words.
#define HEADER_WORDS  16

// Frame data header
typedef struct {
	uint32_t magic;                   // identifies start of header, 0xFFFFFFFF

	uint32_t errno;                   // code of last error, e.g.: BAD_COMMAND
	uint32_t errval;                  // value that caused the error, e.g., the value of the bad command

	uint32_t nbytes;                  // number of bytes in payload (sent after header)

	uint32_t nchans;                  // number of channels
	uint32_t nsamples;                // number of samples per channel
	uint32_t bitdepth;                // number of bits per sample

	uint32_t trigLev;
	uint32_t timeBase;

	uint32_t padding[HEADER_WORDS-10]; // unused space, needed for correct total size, should be HEADER_WORDS minus number of words in the struct!
} header_t;

header_t *header;     // header written to software, first part of usartBuf
uint16_t *outData;    // data written to software, second part of usartBuf // TODO(a): bytes

void init_outbox();
void outbox_TX(uint32_t nDataBytes);

#endif
