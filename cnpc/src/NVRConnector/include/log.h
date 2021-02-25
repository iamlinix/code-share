#ifndef __LOG_H__
#define __LOG_H__

#include "slog.h"


/////////////////////for daemon log////////////////////
#define LOG_DBG(MSG, ...) { \
    slog_debug(0,  MSG, ##__VA_ARGS__); \
}

#define LOG_INFO(MSG, ...) { \
    slog_info(0, MSG, ##__VA_ARGS__); \
}

#define LOG_ERR(MSG, ...) { \
    slog_error(0, MSG, ##__VA_ARGS__); \
}


#endif

