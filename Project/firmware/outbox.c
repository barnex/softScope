#include <stdlib.h>
#include <string.h>

#include "outbox.h"
#include "inbox.h"
#include "usart.h"
#include "utils.h"

header_t *outHeader;
uint16_t *outData;

static uint8_t *usartBuf;   // embeds outHeader and outdata so they can be TX'ed in one go

void init_outbox() {
	// sanity check in case we forgot to update header_t.padding after changing the type
	if (sizeof(header_t) != HEADER_WORDS*sizeof(uint32_t)){
		panic();
	}

	int headerBytes = sizeof(header_t);
	int dataBytes   = MAX_NSAMPLES*sizeof(outData[0]);
	usartBuf	    = malloc(headerBytes + dataBytes);
	memset(usartBuf, 0, headerBytes + dataBytes);
	header = (header_t*)(usartBuf);                      // header is embedded in beginning of usart buffer
	outData = (uint16_t*)(&usartBuf[headerBytes]);       // data is embedded next
	header->magic = MSG_MAGIC;
}

void outbox_TX(uint32_t nDataBytes) {
	if(nDataBytes > MAX_NSAMPLES * sizeof(outData[0])) {
		panic(); // ask to send more than the size of usartBuf
	}
	USART_asyncTX(usartBuf, sizeof(header_t) + nDataBytes);
}
