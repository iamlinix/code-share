#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <unistd.h>
#include <getopt.h>

#include "NVRConnector.h"
#include "util.h"


///////////////////////////////////////////////////////////////////////////////
//
#define NVR_IP		"192.0.0.200"
#define NVR_PORT	8000
#define NVR_USERNAME	"admin"
#define NVR_PASSWORD	"Beijing_2018"


/* global database handler */
MYSQL *conn;

message_queue_t g_msg_queue;

LONG IHandle = -1;
LONG IUserID;
int  iNum = 0;
int  is_daemon = 0;

///////////////////////////////////////////////////////////////////////////////
////

static void show_usage(void)
{
    printf("nvrconn <options>\n");
    printf("\t-D running in daemon mode\n");
    printf("\t-h\tshow help message\n");
}


static int parse_cmd(int argc, char **argv)
{
    char c;
    const char optstr[] = "Dh";

    opterr = 0;
    while ((c = getopt(argc, argv, optstr)) != -1)
    {
        switch (c) {
        case 'D':
            LOG_DBG(">>> running in daemon mode.");
            is_daemon = 1;
            break;

        case 'h':
            show_usage();
            exit(0);

        case ':':
            LOG_ERR(">>> option %c miss argument\n", optopt);
            return -1;

        case '?':
            LOG_ERR(">>> unknow option %c\n", optopt);
            return -1;
        }
    }

    if (optind != argc)
        return -1;

    return 0;
}


void LOG_Init()
{
    slog_init(NVR_LOG, NVR_LOG_CONF, 3, 3, 1);
}


void NVR_SDK_Version()
{
    unsigned int uiVersion = NET_DVR_GetSDKBuildVersion();

    char strTemp[1024] = { 0 };
    sprintf(strTemp, "HCNetSDK V%d.%d.%d.%d\n", \
            (0xff000000 & uiVersion) >> 24, \
            (0x00ff0000 & uiVersion) >> 16, \
            (0x0000ff00 & uiVersion) >> 8, \
            (0x000000ff & uiVersion));

    LOG_DBG(">>> NVR SDK Version: %s.", strTemp);
}


void NVR_Init()
{
    NET_DVR_Init();

    NVR_SDK_Version();
}


void NVR_Connect()
{
    NET_DVR_SetConnectTime(2000, 1);
    NET_DVR_SetReconnect(10000, true);
}


void NVR_Htime() 
{
    NET_DVR_TIME struParams = { 0 };
    DWORD dwReturnLen;
    bool  bRet;

    bRet = NET_DVR_GetDVRConfig(IUserID, NET_DVR_GET_TIMECFG, 1, \
                                &struParams, sizeof(NET_DVR_TIME), &dwReturnLen);
    if (!bRet)
    {
        LOG_ERR(">>> NET_DVR_GetDVRConfig -> NET_DVR_GET_TIMECFG error.");
        NET_DVR_Logout(IUserID);
        NET_DVR_Cleanup();
    }

    LOG_INFO(">>> NVR Time: %04d-%02d-%02d %02d:%02d:%02d\n",
              struParams.dwYear, struParams.dwMonth, struParams.dwDay,
              struParams.dwHour, struParams.dwMinute, struParams.dwSecond);
}


int NVR_Login()
{
    NET_DVR_DEVICEINFO_V30 struDeviceInfo;

    IUserID = NET_DVR_Login_V30((char *)NVR_IP, NVR_PORT, (char *)NVR_USERNAME,
                                (char *)NVR_PASSWORD, &struDeviceInfo);
    if (IUserID < 0)
    {
        LOG_ERR(">>> Login Failed! Error number: %d", NET_DVR_GetLastError());
        NET_DVR_Cleanup();
        return -1;
    }

    LOG_INFO(">>> Login Successfully!");
    return 0;
}


void CALLBACK ExceptionCallBack(DWORD dwType, LONG lUserID, LONG lHandle, void *pUser)
{
    switch (dwType)
    {
    case EXCEPTION_RECONNECT:
        LOG_INFO(">>> Reconnect... @%lld\n", time(NULL));
        break;
    default:
        break;
    }
}


void CALLBACK MessageCallback(LONG lCommand, NET_DVR_ALARMER *pAlarmer,
                              char *pAlarmInfo, DWORD dwBufLen, void* pUser)
{
    FILE *fSnapPicPlate = NULL;
    FILE *fSnapPic = NULL;
    char  filename[100] = {0};
    unsigned int i = 0;

    switch (lCommand)
    {
    case COMM_UPLOAD_PLATE_RESULT:
    {
        NET_DVR_PLATE_RESULT struPlateResult = { 0 };
        memcpy(&struPlateResult, pAlarmInfo, sizeof(struPlateResult));

        LOG_INFO(">>> COMM_UPLOAD_PLATE_RESULT\n");
        LOG_INFO(">>> plate: %s\n", struPlateResult.struPlateInfo.sLicense);

        //sscanf(struPlateResult.byAbsTime, "%04d%02d%02d%02d%02d%02d",
        //       &vehi->t.tm_year, &vehi->t.tm_mon, &vehi->t.tm_mday,
        //       &vehi->t.tm_hour, &vehi->t.tm_min, &vehi->t.tm_sec);

        if (struPlateResult.dwPicLen != 0 && struPlateResult.byResultType == 1)
        {
            sprintf(filename, "%s/%s.jpg", NVR_PIC_PATH, struPlateResult.struPlateInfo.sLicense);
            fSnapPic = fopen(filename, "wb");
            fwrite(struPlateResult.pBuffer1, struPlateResult.dwPicLen, 1, fSnapPic);
            iNum++;
            fclose(fSnapPic);
        }

        if (struPlateResult.dwPicPlateLen != 0 && struPlateResult.byResultType == 1)
        {
            sprintf(filename, "%s/%s.jpg", NVR_PIC_PATH, struPlateResult.struPlateInfo.sLicense);
            fSnapPicPlate = fopen(filename, "wb");
            fwrite(struPlateResult.pBuffer1, struPlateResult.dwPicLen, 1, fSnapPicPlate);
            iNum++;
            fclose(fSnapPicPlate);
        }
        break;
    }
    case COMM_ITS_PLATE_RESULT:
    {
        NET_ITS_PLATE_RESULT struITSPlateResult = { 0 };
        memcpy(&struITSPlateResult, pAlarmInfo, sizeof(struITSPlateResult));

        LOG_INFO(">>> Got %d pictures...", struITSPlateResult.dwPicNum);

        for (i = 0; i < struITSPlateResult.dwPicNum; i++)
        {
            LOG_INFO(">>> COMM_ITS_PLATE_RESULT [%d]", i);

            veh_info_t *vehi = new veh_info_t;
            memset(vehi, 0x0, sizeof(veh_info_t));

            /* just for test */
            vehi->station_id = 100;

            /* channel id */
            vehi->channel_id = struITSPlateResult.byChanIndex;
            LOG_INFO(">>> channel_id = %d", vehi->channel_id);

            /* timestamp in alarm */
            vehi->tm = struITSPlateResult.struSnapFirstPicTime;
            LOG_INFO(">>> capture time = %04d-%02d-%02d %02d:%02d:%02d",
                     vehi->tm.wYear, vehi->tm.byMonth, vehi->tm.byDay,
                     vehi->tm.byHour, vehi->tm.byMinute, vehi->tm.bySecond);

            /* local timestamp */
            time(&vehi->post_time);

            /* plate number, MUST convert to UTF-8 code */
            char *instr = struITSPlateResult.struPlateInfo.sLicense;
            //char *instr = (char *)"京QC65B3";
            unsigned char *uinstr = (unsigned char *)instr;
            if (uinstr[0] == 0xce && uinstr[1] == 0xde)
            {
                LOG_DBG(">>> got nono plate, skip it...");
                break;
            }

            size_t inlen = strlen(instr);
            char   outstr[32] = { 0 };
            size_t outlen = 32;
            int re = code_convert(instr, inlen, outstr, outlen);
            if (re < 0)
            {
                LOG_ERR(">>> code convert failed");
                return;
            }

            strcpy(vehi->plate_num, outstr);
            LOG_INFO(">>> plate number: %s", vehi->plate_num);

            /* plate color */
            vehi->plate_color = struITSPlateResult.struPlateInfo.byColor;
            LOG_INFO(">>> plate color: %d", vehi->plate_color);

            /* veh type */
            vehi->veh_type = struITSPlateResult.struVehicleInfo.byVehicleType;
            LOG_INFO(">>> veh type: %d", vehi->veh_type);

            /* veh brand */
            vehi->veh_brand = struITSPlateResult.struVehicleInfo.byVehicleLogoRecog;
            vehi->veh_sub_brand = struITSPlateResult.struVehicleInfo.byVehicleSubLogoRecog;
            vehi->veh_model = struITSPlateResult.struVehicleInfo.byVehicleModel;
            LOG_INFO(">>> veh brand: %d, sub brand: %d, model: %d",
                      vehi->veh_brand, vehi->veh_sub_brand, vehi->veh_model);

            /* direction */
            vehi->direction = struITSPlateResult.byDir;
            LOG_INFO(">>> direction: %d", vehi->direction);

            if ((struITSPlateResult.struPicInfo[i].dwDataLen != 0) &&
                ((struITSPlateResult.struPicInfo[i].byType == 1) ||
                 (struITSPlateResult.struPicInfo[i].byType == 2)))
            {
                vehi->pic_buff_size = struITSPlateResult.struPicInfo[i].dwDataLen;
                vehi->pic_buff = new char[vehi->pic_buff_size];
                memcpy(vehi->pic_buff, struITSPlateResult.struPicInfo[i].pBuffer, vehi->pic_buff_size);

               LOG_INFO(">>> abs time: %s", struITSPlateResult.struPicInfo[i].byAbsTime);
            }
            else
            {
                LOG_INFO(">>> perhaps no buffer...");
            }

            /* en-queue */
            message_queue_t *mq = &g_msg_queue;

            pthread_mutex_lock(&mq->qlock);
            mq->msg_queue.push(vehi);
            pthread_mutex_unlock(&mq->qlock);
            pthread_cond_signal(&mq->cond);
        }
        break;
    }
    default:
        break;
    }

    return;
}


void SetMessageCallBack()
{
    NET_DVR_SetDVRMessageCallBack_V30(MessageCallback, NULL);
}


int SetupAlarm()
{
    NET_DVR_SETUPALARM_PARAM struSetupParam = { 0 };
    struSetupParam.dwSize = sizeof(NET_DVR_SETUPALARM_PARAM);

    struSetupParam.byAlarmInfoType = 1;
    struSetupParam.byLevel = 1;

    IHandle = NET_DVR_SetupAlarmChan_V41(IUserID, &struSetupParam);
    if (IHandle < 0)
    {
        LOG_ERR(">>> NET_DVR_SetupAlarmChan_V41 Failed! Error number: %d", NET_DVR_GetLastError());
        NET_DVR_Logout(IUserID);
        NET_DVR_Cleanup();
        return -1;
    }

    return 0;
}


void CloseAlarm()
{
    if (!NET_DVR_CloseAlarmChan_V30(IHandle))
    {
        LOG_INFO(">>> NET_DVR_CloseAlarmChan_V30 Failed! Error number: %d", NET_DVR_GetLastError());
        NET_DVR_Logout(IUserID);
        NET_DVR_Cleanup();
        return;
    }

    IHandle = -1;
}


void OnExit(void)
{
    CloseAlarm();

    NET_DVR_Logout(IUserID);
    NET_DVR_Cleanup();
}


void *message_handler(void *para)
{
    LOG_INFO(">>> [DB] thread message handler start...");

    message_queue_t *mq = &g_msg_queue;

    while (1)
    {
        pthread_mutex_lock(&mq->qlock);
        if(!mq->msg_queue.empty())
        {
            veh_info_t *vehi = mq->msg_queue.front();
            mq->msg_queue.pop();
            pthread_mutex_unlock(&mq->qlock);

            /* handle current item */
            LOG_INFO(">>> [DB] Got a item with station id (%d).", vehi->station_id);
            LOG_INFO(">>> [DB] plate number (%s).", vehi->plate_num);

            DB_Connector(conn, vehi);

            if (vehi)
            {
                if (vehi->pic_buff)
                {
                    delete vehi->pic_buff;
                }
                delete vehi;
                vehi = NULL;
            }
        }
        else
        {
            pthread_mutex_unlock(&mq->qlock);
            //pthread_cond_wait(&mq->cond, &mq->qlock);
            sleep(1);
        }
    }

    return NULL;
}


//#define __OFFLINE_TEST__  1
#ifdef __OFFLINE_TEST__
void offline_test()
{
    LOG_DBG(">>> do offline test.");

    NET_ITS_PLATE_RESULT res;
    memset(&res, 0x0, sizeof(NET_ITS_PLATE_RESULT));
    res.dwPicNum = 1;

    MessageCallback(COMM_ITS_PLATE_RESULT, NULL, (char *)&res, 0, NULL);
}
#endif


int main(int argc, char **argv)
{
    int ret = 0;

    LOG_Init();

    if (parse_cmd(argc, argv))
    {
        show_usage();
        return -1;
    }

    if (is_daemon)
    {
        if (daemon(0, 0) == -1)
        {
            LOG_ERR("nvrconn daemon failed.");
            return -1;
        }
        printf(">>> after daemon...\n");
        chdir(NVR_BIN_PATH);
    }

    LOG_DBG("nvrconnd is running in %s mode.", is_daemon ?
            "background" : "foreground");

    pthread_mutex_init(&g_msg_queue.qlock, NULL);
    pthread_cond_init(&g_msg_queue.cond, NULL);

    conn = DB_Init();
    if (!conn)
    {
        LOG_ERR(">>> DB Init failed.");
        //return -1;
    }

    NVR_Init();

    NVR_Connect();

connect:
    ret = NVR_Login();
    if (ret < 0)
    {
        LOG_ERR(">>> NVR Login failed...");
        sleep(5);

        LOG_ERR(">>> After sleep 5s, retrying...");
        goto connect;
    }

    NVR_Htime();

    SetupAlarm();

    SetMessageCallBack();

    pthread_t recv;
    ret = pthread_create(&recv, NULL, &message_handler, NULL);
    if (ret)
    {
        LOG_ERR(">>> Create pthread error!");
        return -1;
    }

    for (;;)
    {
#ifdef __OFFLINE_TEST__
        offline_test();
#endif

        SetMessageCallBack();
        sleep(1);
    }

    if (conn) {
        DB_Close(conn);
    }

    //atexit(OnExit);
    return 0;
}


