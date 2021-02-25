/*
 * Copyright(C) 2010,Custom Co., Ltd 
 * FileName: playback.cpp
 * Description: 
 * Version: 1.0
 * Author: panyadong
 * Create Date: 2010-4-24
 * Modification History：
 */

#include <unistd.h> 
#include <stdio.h>
#include <time.h>
#include <stdlib.h>
#include <memory.h>
#include <sys/stat.h>

#include "HCNetSDK.h"
#include "iniparser.h"
#include "log.h"


#define NVR_LOG         "/var/log/nvr/video.log"
#define NVR_LOG_CONF    "/opt/tav/conf/slog.conf"

#define OUT_PATH        "/opt/tav/video"
#define INI_PATH        "video.ini"


struct config {

    int smonth;
    int sday;
    int schannel;

    int emonth;
    int eday;
    int echannel;
};

struct config cfg;

///////////////////////////////////////////////////////////////////////////////
//
void LOG_Init()
{
    slog_init(NVR_LOG, NVR_LOG_CONF, 3, 3, 1);
}


int Parse_INI_File(const char * ini_name)
{
    dictionary  *ini;

    /* Some temporary variables to hold query results */
    ini = iniparser_load(ini_name);
    if (ini == NULL) {
        LOG_ERR("cannot parse file: %s", ini_name);
        return -1 ;
    }

    //iniparser_dump(ini, stderr);

    /* Get start attributes */
    LOG_DBG("start:");
    cfg.smonth = iniparser_getint(ini, "start:month", -1);
    LOG_DBG("month:   [%d]", cfg.smonth);
    cfg.sday = iniparser_getint(ini, "start:day", -1);
    LOG_DBG("day:     [%d]", cfg.sday);
    cfg.schannel = iniparser_getint(ini, "start:channel", -1);
    LOG_DBG("channel: [%d]", cfg.schannel);

    /* Get end attributes */
    LOG_DBG("end:");
    cfg.emonth = iniparser_getint(ini, "end:month", -1);
    LOG_DBG("month:   [%d]", cfg.emonth);
    cfg.eday = iniparser_getint(ini, "end:day", -1);
    LOG_DBG("day:     [%d]", cfg.eday);
    cfg.echannel = iniparser_getint(ini, "end:channel", -1);
    LOG_DBG("channel: [%d]", cfg.echannel);

    iniparser_freedict(ini);

    LOG_DBG(">>> load ini file finish");
    return 0 ;
}


void GET_SDK_Version()
{
    unsigned int uiVersion = NET_DVR_GetSDKBuildVersion();

    char strTemp[1024] = {0};
    sprintf(strTemp, "HCNetSDK V%d.%d.%d.%d\n", \
        (0xff000000 & uiVersion)>>24, \
        (0x00ff0000 & uiVersion)>>16, \
        (0x0000ff00 & uiVersion)>>8, \
        (0x000000ff & uiVersion));
    printf(strTemp);
}


void Login_Init(NET_DVR_USER_LOGIN_INFO *u)
{
    strcpy(u->sDeviceAddress, "192.0.0.200");
    u->wPort = 8000;
    strcpy(u->sUserName,  "admin");
    strcpy(u->sPassword , "Beijing_2018");
}


int CreateUnitDir(int month, int day)
{
    char path[256] = {0};

    sprintf(path, "%s/%02d%02d", OUT_PATH, month, day);

    if (access(path, 0) != 0)
    {
	if (mkdir(path, 0755) == -1)
        {
            LOG_ERR("mkdir %s failed.", path);
            return -1;
        }

    	LOG_DBG("dir (%s) not exist, create now", path);
    } 

    return 0;
}


int SaveVideoFile(LONG userId, char *dest, LPNET_DVR_PLAYCOND cond)
{
    LONG hFile = 0;
    if ((hFile = NET_DVR_GetFileByTime_V40(userId, dest, cond)) < 0)
    {
        LOG_ERR("save video file (%s) failed. error [%d]", dest,
		NET_DVR_GetLastError());
        return -1;
    }

    if (!NET_DVR_PlayBackControl_V40(hFile, NET_DVR_PLAYSTART, NULL, 0, NULL, 0))
    {
        LOG_ERR("play back control (%s) failed, error [%d]", dest,
		NET_DVR_GetLastError());
        return -1;
    }

    int pos = 0;
    for (pos = 0;  pos < 100 && pos >= 0; pos = NET_DVR_GetDownloadPos(hFile))
    {
        LOG_DBG("++++ downloading, progress %d%", pos);
        sleep(5);   //second
    }

    LOG_DBG(">>>> download finish", pos);

    if (!NET_DVR_StopGetFile(hFile))
    {
        LOG_ERR("failed to stop get file (%s), error [%d]", dest,
		NET_DVR_GetLastError());
        return -1;
    }

    if (pos < 0 || pos > 100)
    {
        LOG_ERR("download video (%s) failed, error [%d]", dest,
		NET_DVR_GetLastError());
        return -1;
    }

    LOG_DBG(">>> saving video to: %s.", dest);
    return 0;
}


int SmartSearchVideo(LONG userID)
{
    int month = 0;
    int day   = 0;
    int ch    = 0;
    int hour  = 0;

    NET_DVR_PLAYCOND cond = {0};

    for (month = cfg.smonth; month <= cfg.emonth; month++) {
    
        for (day = cfg.sday; day <= cfg.eday; day++) {
	
            if (CreateUnitDir(month, day) < 0) {
                LOG_ERR("create unit dir failed.");
                return -1;
            }

            for (ch = cfg.schannel; ch <= cfg.echannel; ch++) {

                /* set channel */
                cond.dwChannel = ch;

                for (hour = 0; hour < 24; hour++) {
                    /* set begin and end time */
                    NET_DVR_TIME startTime, stopTime;
                    startTime.dwYear   = 2019;
                    startTime.dwMonth  = month;
                    startTime.dwDay    = day;
                    startTime.dwHour   = hour;
                    startTime.dwMinute = 0;
                    startTime.dwSecond = 0;

                    stopTime.dwYear	= 2019;
                    stopTime.dwMonth	= month;
                    stopTime.dwDay	= day;
                    stopTime.dwHour	= hour + 1;
                    stopTime.dwMinute	= 0;
                    stopTime.dwSecond	= 0;

                    cond.struStartTime = startTime;
                    cond.struStopTime = stopTime;

                    char destName[256] = {0};
                    sprintf(destName, "%s/%02d%02d/ch%02d-%04d%02d%02d-%02d.mp4",
                            OUT_PATH, month, day, ch, startTime.dwYear,
			    month, day, hour);

                    SaveVideoFile(userID, destName, &cond);
                }
	    }
	}
    
    }

    return 0;
}


int main()
{
    LOG_Init();

    if (Parse_INI_File(INI_PATH) < 0)
    {
        LOG_ERR("<<< load ini file failed");
        return -1;
    }

    NET_DVR_Init();

    GET_SDK_Version();

    /* login info */
    NET_DVR_USER_LOGIN_INFO userLoginInfo = {0};

    Login_Init(&userLoginInfo);

    NET_DVR_DEVICEINFO_V40 deviceInfo = { 0 };

    LONG userID = NET_DVR_Login_V40(&userLoginInfo, &deviceInfo);
    if (userID < 0)
    {
        LOG_ERR("<<< Login failed with code %d", NET_DVR_GetLastError());
        NET_DVR_Cleanup();
        return -1;
    }

    LOG_DBG(">>> login success with user id (%ld)", userID);

    SmartSearchVideo(userID);

    NET_DVR_Logout_V30(userID);

    NET_DVR_Cleanup();

    return 0;
}



