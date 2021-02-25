#include <iconv.h>
#include <string.h>

#include "log.h"
#include "util.h"

///////////////////////////////////////////////////////////////////////////////
//
int code_convert(char *inbuf, size_t inlen, char *outbuf, size_t outlen)
{
    char **pin  = &inbuf;
    char **pout = &outbuf;
    iconv_t cd;

    cd = iconv_open("UTF-8", "GB2312");
    if (cd==0)
    {
        LOG_ERR(">>> iconv open failed.");
        return -1;
    }

    memset(outbuf, 0x0, outlen);

    if (iconv(cd, pin, &inlen, pout, &outlen) == (size_t)-1)
    {
        LOG_ERR(">>> iconv failed.");
        iconv_close(cd);
        return -1;
    }

    iconv_close(cd);

    return 0;
}

