#ifndef _ERROR_H_
#define _ERROR_H_

#include <stdint.h>

enum ErrorCode{
	NO_ERROR         = 0,
	BAD_MAGIC        = 1,
	BAD_COMMAND      = 2,
	BAD_CLOCK_PERIOD = 3
};

void error(uint32_t code, uint32_t value);

#endif
