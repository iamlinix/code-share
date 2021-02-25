#ifndef __UTIL_H__
#define __UTIL_H__


#define NVR_LOG		"/var/log/nvr/nvrconn.log"

#define NVR_ROOT_PATH	"/opt/tav/"
#define NVR_BIN_PATH	NVR_ROOT_PATH"bin"
#define NVR_LOG_CONF	NVR_ROOT_PATH"conf/slog.conf"
#define NVR_PIC_PATH	NVR_ROOT_PATH"pic/"

///////////////////////////////////////////////////////////////////////////////
//
int code_convert(char *inbuf, size_t inlen, char *outbuf, size_t outlen);

#endif
