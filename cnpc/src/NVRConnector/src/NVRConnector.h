#ifndef __NVRCONNECTOR_H__
#define __NVRCONNECTOR_H__


#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <time.h>
#include <mysql.h>
#include <queue>

#include "HCNetSDK.h"
#include "log.h"


typedef struct vehicle_info
{
    int station_id;
    int channel_id;
    NET_DVR_TIME_V30 tm;
    //struct tm t;
    time_t post_time;
    char plate_num[32];
    int plate_color;
    int veh_type;
    int veh_brand;
    int veh_sub_brand;
    int veh_model;
    char *pic_buff;
    int pic_buff_size;
    int direction;
} veh_info_t;

typedef struct message_queue
{
    std::queue<veh_info_t *> msg_queue;
    pthread_mutex_t qlock;
    pthread_cond_t  cond;
} message_queue_t;


///////////////////////////////////////////////////////////////////////////////
//
MYSQL *DB_Init();
int    DB_Connector(MYSQL *conn, veh_info_t *vehi);
void   DB_Close(MYSQL *conn);

#endif
