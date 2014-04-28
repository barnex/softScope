#ifndef _INBOX_H_
#define _INBOX_H_

#include <stdint.h>

#define ADC_BUFSIZE	    2048
#define IR_PERIOD       32
#define MAX_NSAMPLES    ADC_BUFSIZE
#define MIN_NSAMPLES    (IR_PERIOD)

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
