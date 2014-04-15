/*
 * Utility for reading/writing raw tty device
 */
#include <errno.h>
#include <fcntl.h>
#include <stdio.h>
#include <string.h>
#include <termios.h>
#include <unistd.h>

// pair of baudrate + speed code
typedef struct {
	int x, y;
} int2;

// maps human-readable baud rate on speed code
int2 baudTable[18] = {
	{50     ,B50},
	{75     ,B75},
	{110    ,B110},
	{134    ,B134},
	{150    ,B150},
	{200    ,B200},
	{300    ,B300},
	{600    ,B600},
	{1200   ,B1200},
	{1800   ,B1800},
	{2400   ,B2400},
	{4800   ,B4800},
	{9600   ,B9600},
	{19200  ,B19200},
	{38400  ,B38400},
	{57600  ,B57600},
	{115200 ,B115200},
	{230400 ,B230400}
};

// return speed code for baud rate, or -1 if baud rate not in baudTable
int decodeBaud(int rate) {
	int NBAUD = sizeof(baudTable)/sizeof(baudTable[0]);
	int i = 0;

	for(i=0; i<NBAUD; i++) {
		if (baudTable[i].x == rate) {
			return baudTable[i].y;
		}
	}
	return -1; // not found
}

#define ERRLEN 4096
static char TTYErr[ERRLEN+1];
char *TTYerr = &TTYErr[0];

int openTTY(char* file, int baud) {

	// parse baud rate
	int baudRate = decodeBaud(baud);
	if(baudRate <= 0) {
		sprintf(&TTYErr[0], "invalid baud rate: %d", baud); // TODO: snprintf
		return -1;
	}

	// open tty
	int fd = open(file, O_RDWR | O_NOCTTY | O_SYNC);
	if(fd == -1) {
		sprintf(&TTYErr[0], "error opening %s: %s", file, strerror(errno));
		return -1;
	}
	if(!isatty(fd)) {
		sprintf(&TTYErr[0], "not a tty: %s", file);
		return -1;
	}

	// setup tty
	struct termios config;
	if(tcgetattr(fd, &config) < 0) {
		sprintf(&TTYErr[0], "tcgetattr %s: %s", file, strerror(errno));
		return -1;
	}
	cfmakeraw(&config);
	config.c_cflag &= ~(CSIZE | PARENB);
	config.c_cflag |= CS8;
	config.c_cc[VMIN]  = 128; // buffer
	config.c_cc[VTIME] = 1;   // return as quickly as possible

	// communication speed
	if(cfsetispeed(&config, baudRate) < 0 || cfsetospeed(&config, baudRate) < 0) {
		sprintf(&TTYErr[0], "tcsetispeed %s: %s", file, strerror(errno));
		return -1;
	}

	if(tcsetattr(fd, TCSANOW, &config) != 0) {
		sprintf(&TTYErr[0], "tcsetattr %s: %s", file, strerror(errno));
		return -1;
	}
	
	return fd;
}

int readTTY(int fd, void *buf, int N) {
	int n = read(fd, buf, N);
	if(n<0){
		sprintf(&TTYErr[0], "readTTY: %s", strerror(errno));
	}
	return n;
}

int writeTTY(int fd, void *buf, int N) {
	int n = write(fd, buf, N);
	if(n<N){
		sprintf(&TTYErr[0], "writeTTY: %s", strerror(errno));
	}
	return n;
}

void closeTTY(int fd) {
	close(fd);
}
