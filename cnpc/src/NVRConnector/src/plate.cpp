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
#include <string.h>
#include <sys/stat.h>

#include "HCNetSDK.h"
#include "iniparser.h"
#include "log.h"

#define NVR_LOG         "/var/log/nvr/plate.log"
#define NVR_LOG_CONF    "/opt/tav/conf/slog.conf"

#define OUT_PATH        "/opt/tav/plate"
#define INI_PATH        "plate.ini"


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

    LOG_DBG(strTemp);
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


int SavePicFile(int userId, char *src, char *dest, int *dcnt)
{
    bool res;

    res = NET_DVR_GetPicture(userId, src, dest);
    if (res == false)
    {
        LOG_ERR("save picture %s failed, error [%d]", src,
		NET_DVR_GetLastError());
        return -1;
    }

    (*dcnt)++;

    LOG_DBG(">>> saving picture(%d) to: %s.", *dcnt, dest);

    return 0;
}


int PackageDayUnit(int month, int day, int dcnt)
{
    char cmd[256] = {0};
    char prefix[256] = {0};

    if (dcnt <= 0) {
    	LOG_DBG("Stop package, day counter invalid: %d.", dcnt);
        return 0;
    }

    sprintf(prefix, "%s/%02d%02d", OUT_PATH, month, day);
    sprintf(cmd, "tar -czvf %s.tar.gz %s/*", prefix, prefix);

    int re = system(cmd);
    if (re < 0)
    {
    	LOG_ERR("excute tar cmd failed: %d.", re);
    } else
    {
        LOG_DBG("package unit %s success with counter (%d).",
		prefix, dcnt);
    }

    return 0;
}


int SearchUnit(LONG userID, int month, int day, int ch, int *dcnt)
{
    //create dir first
    if (CreateUnitDir(month, day) < 0) {
        LOG_ERR("create unit dir failed.");
        return -1;
    }

    NET_DVR_SMART_SEARCH_PIC_PARA picPara = {0};

    //picPara.dwSize = sizeof(NET_DVR_FIND_PICTURE_PARAM);
    picPara.dwChanNo = ch;

    picPara.struStartTime.wYear	   = 2019;
    picPara.struStartTime.byMonth  = month;
    picPara.struStartTime.byDay	   = day;
    picPara.struStartTime.byHour   = 0;
    picPara.struStartTime.byMinute = 0;
    picPara.struStartTime.bySecond = 0;

    picPara.struEndTime.wYear	   = 2019;
    picPara.struEndTime.byMonth	   = month;
    picPara.struEndTime.byDay	   = day;
    picPara.struEndTime.byHour	   = 23;
    picPara.struEndTime.byMinute   = 59;
    picPara.struEndTime.bySecond   = 59;

    // 0 -> vechile detection
    picPara.wSearchType = 0;

    LOG_DBG(">>> start search 2019-%02d-%02d, channel %d", month, day, ch);

    int h = NET_DVR_SmartSearchPicture(userID, &picPara);
    if(h < 0)
    {
        LOG_ERR("search file in %02d-%02d (%d) failed, last error %d",
		month, day, ch, NET_DVR_GetLastError());
        return -1;
    }

    NET_DVR_SMART_SEARCH_PIC_RET picInfo;

    while(true)
    {
        LONG result = NET_DVR_FindNextSmartPicture(h, &picInfo);
        if(result == NET_DVR_ISFINDING)
        {
            LOG_DBG(">>> file is searching...");
            continue;
        }
        else if(result == NET_DVR_FILE_SUCCESS)
        {
            LOG_DBG(">>> search file %s success", picInfo.sFileName);
            //LOG_DBG(">>> file size : %d", picInfo.dwFileSize);

            NET_DVR_PLATE_INFO *v = &picInfo.uPicFeature.struPlateInfo;
            LOG_DBG(">>> plate type: %d", v->byPlateType);
            //LOG_DBG(">>> plate color: %d", v->byColor);
            //LOG_DBG(">>> plate license: %s", v->sLicense);

            char destName[256] = {0};
            sprintf(destName, "%s/%02d%02d/%s-%s.jpg", OUT_PATH, month, day,
		    picInfo.sFileName, v->sLicense);

            SavePicFile(userID, picInfo.sFileName, destName, dcnt);
            usleep(5000);
        }
        else if(result == NET_DVR_FILE_NOFIND || result == NET_DVR_NOMOREFILE)
        {	
            LOG_DBG("<<< file not find");
            break;
        }
        else
        {
            //LOG_ERR("<<< search file failed: %ld", result);
            //break;
        }
    }

    NET_DVR_CloseSmartSearchPicture(h);

    return 0;
}


int SmartSearchPic(LONG userID)
{
    int month = 0;
    int day   = 0;
    int ch    = 0;
    int dayCounter = 0;

    for (month = cfg.smonth; month <= cfg.emonth; month++) {
    
        for (day = cfg.sday; day <= cfg.eday; day++) {
	
            dayCounter = 0;
            for (ch = cfg.schannel; ch <= cfg.echannel; ch++) {

	    	SearchUnit(userID, month, day, ch, &dayCounter);
	    }

            /* search a day file finish, do package now*/
            PackageDayUnit(month, day, dayCounter);
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

    /* LOGIN */
    LONG userID = NET_DVR_Login_V40(&userLoginInfo, &deviceInfo);
    if (userID < 0)
    {
        LOG_ERR("<<< Login failed with code %d", NET_DVR_GetLastError());
        NET_DVR_Cleanup();
        return -1;
    }

    LOG_DBG(">>> login success with user id (%ld)", userID);

    SmartSearchPic(userID);

    NET_DVR_Logout_V30(userID);

    NET_DVR_Cleanup();

    return 0;
}


