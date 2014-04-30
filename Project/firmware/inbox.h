#ifndef _INBOX_H_
#define _INBOX_H_

#include <stdint.h>

#define ADC_CHUNKSIZE   512
#define ADC_BUFSIZE	    (4 * ADC_CHUNKSIZE)
#define MAX_NSAMPLES    4096                 // not really hard limit
#define MIN_NSAMPLES    (ADC_CHUNKSIZE/4)    // not really hard limit

#define MIN_CLOCK_PERIOD 30   // 1.4 MSample/s limit

#define MSG_MAGIC 0xFAFBFCFD

// Incoming message
typedef struct {
	uint32_t magic;
	uint32_t command;
	uint32_t value;
} message_t;


// Value for message_t command
enum {
	INVALID   = 0,
    SET_SAMPLES   = 1,
    SET_TIMEBASE  = 2,
    SET_TRIGLEV   = 3,
    CLEAR_ERR = 1000,  // Clear errno
    REQ_FRAMES= 1001,  // Request a number of frames to be sent
};


// Number of requested frames. If != 0, send frame and decrement.
volatile uint32_t reqFrames;

volatile uint32_t nSamples ;
volatile uint32_t timebase;
volatile uint32_t triglev ;

void init_inbox();

#endif
