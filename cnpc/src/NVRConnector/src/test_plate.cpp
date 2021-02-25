#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <unistd.h>
#include <getopt.h>

#include "HCNetSDK.h"
#include "NVRConnector.h"
#include "util.h"


///////////////////////////////////////////////////////////////////////////////
//
#define NVR_IP		"192.168.1.210"
#define NVR_PORT	8000
#define NVR_USERNAME	"admin"
#define NVR_PASSWORD	"Ys123456"


LONG IHandle = -1;
LONG IUserID;
int  iNum = 0;

///////////////////////////////////////////////////////////////////////////////
////


void NVR_SDK_Version()
{
    unsigned int uiVersion = NET_DVR_GetSDKBuildVersion();

    char strTemp[1024] = { 0 };
    sprintf(strTemp, "HCNetSDK V%d.%d.%d.%d\n", \
            (0xff000000 & uiVersion) >> 24, \
            (0x00ff0000 & uiVersion) >> 16, \
            (0x0000ff00 & uiVersion) >> 8, \
            (0x000000ff & uiVersion));

    printf(">>> NVR SDK Version: %s.\n", strTemp);
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
        printf(">>> NET_DVR_GetDVRConfig -> NET_DVR_GET_TIMECFG error.\n");
        NET_DVR_Logout(IUserID);
        NET_DVR_Cleanup();
    }

    printf(">>> NVR Time: %04d-%02d-%02d %02d:%02d:%02d\n",
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
        printf(">>> Login Failed! Error number: %d\n", NET_DVR_GetLastError());
        NET_DVR_Cleanup();
        return -1;
    }

    printf(">>> Login Successfully !\n");
    printf(">>> sSerialNumber = %s\n", struDeviceInfo.sSerialNumber);
    printf(">>> byDVRType = %d\n", struDeviceInfo.byDVRType);
    printf(">>> wDevType = %d\n", struDeviceInfo.wDevType);
    return 0;
}


void CALLBACK ExceptionCallBack(DWORD dwType, LONG lUserID, LONG lHandle, void *pUser)
{
    switch (dwType)
    {
    case EXCEPTION_RECONNECT:
        printf(">>> Reconnect... @%lld\n", time(NULL));
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

        printf(">>> COMM_UPLOAD_PLATE_RESULT\n");
        printf(">>> plate: %s\n", struPlateResult.struPlateInfo.sLicense);

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

        for (i = 0; i < struITSPlateResult.dwPicNum; i++)
        {
            printf(">>> COMM_ITS_PLATE_RESULT\n");

            veh_info_t *vehi = new veh_info_t;
            memset(vehi, 0x0, sizeof(veh_info_t));

            /* just for test */
            vehi->station_id = 100;

            /* channel id */
            vehi->channel_id = struITSPlateResult.byChanIndex;
            printf(">>> channel_id = %d\n", vehi->channel_id);

            //NET_DVR_TIME_V30  tm = struITSPlateResult.struSnapFirstPicTime;
            //cout << tm.wYear << "-" << tm.byMonth << "-" << tm.byDay << " " << tm.byHour << ":" << tm.byMinute << ":" << tm.bySecond << endl;
            time(&vehi->post_time);
            printf(">>> post time = %lld\n", vehi->post_time);

            /* plate number, MUST convert to UTF-8 code */
            unsigned char *instr = (unsigned char *)struITSPlateResult.struPlateInfo.sLicense;
            printf(">>> plate: 0x%02x 0x%02x 0x%02x\n", instr[0], instr[1], instr[2]);
            if (instr[0] == 0xce && instr[1] == 0xde)
            {
                printf(">>> none plate, skip it...\n");
            }
            //char *instr = (char *)"京QC65B3";
            size_t inlen = strlen((char *)instr);
            char outstr[32] = { 0 };
            size_t outlen = 32;
            int re = code_convert((char *)instr, inlen, outstr, outlen);
            if (re < 0)
            {
                printf(">>> code convert failed\n");
                return;
            }

            strcpy(vehi->plate_num, outstr);
            printf(">>> plate number: %s\n", vehi->plate_num);

            /* plate color */
            vehi->plate_color = struITSPlateResult.struPlateInfo.byColor;
            printf(">>> plate color: %d\n", vehi->plate_color);

            /* veh type */
            vehi->veh_type = struITSPlateResult.struVehicleInfo.byVehicleType;
            printf(">>> veh type: %d\n", vehi->veh_type);

            /* veh brand */
            vehi->veh_brand = struITSPlateResult.struVehicleInfo.byVehicleLogoRecog;
            vehi->veh_sub_brand = struITSPlateResult.struVehicleInfo.byVehicleSubLogoRecog;
            vehi->veh_model = struITSPlateResult.struVehicleInfo.byVehicleModel;

            if ((struITSPlateResult.struPicInfo[i].dwDataLen != 0) &&
                ((struITSPlateResult.struPicInfo[i].byType == 1) ||
                 (struITSPlateResult.struPicInfo[i].byType == 2)))
            {
                vehi->pic_buff_size = struITSPlateResult.struPicInfo[i].dwDataLen;
                vehi->pic_buff = new char[vehi->pic_buff_size];
                memcpy(vehi->pic_buff, struITSPlateResult.struPicInfo[i].pBuffer, vehi->pic_buff_size);
            }

            /*
            if ((struITSPlateResult.struPicInfo[i].dwDataLen != 0) && (struITSPlateResult.struPicInfo[i].byType == 0))
            {
                sprintf(filename, "./pic/1/%s_%d.jpg", struITSPlateResult.struPlateInfo.sLicense, i);
                fSnapPicPlate = fopen(filename, "wb");
                fwrite(struITSPlateResult.struPicInfo[i].pBuffer, struITSPlateResult.struPicInfo[i].dwDataLen, 1, fSnapPicPlate);
                iNum++;
                fclose(fSnapPicPlate);
            }
            */
            //车牌小图片
            if ((struITSPlateResult.struPicInfo[i].dwDataLen != 0) && (struITSPlateResult.struPicInfo[i].byType == 1))
            {
                printf(">>> test small picture\n");
                printf(">>> abs time = %s\n", struITSPlateResult.struPicInfo[i].byAbsTime);
                printf(">>> positon = (%f, %f, %f, %f)\n",
                        struITSPlateResult.struPicInfo[i].struPlateRect.fX,
                        struITSPlateResult.struPicInfo[i].struPlateRect.fY,
                        struITSPlateResult.struPicInfo[i].struPlateRect.fWidth,
                        struITSPlateResult.struPicInfo[i].struPlateRect.fHeight);
            }
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

    struSetupParam.byAlarmInfoType = 0;
    struSetupParam.byLevel = 1;

    IHandle = NET_DVR_SetupAlarmChan_V41(IUserID, &struSetupParam);
    if (IHandle < 0)
    {
        printf(">>> NET_DVR_SetupAlarmChan_V41 Failed! Error number: %d\n", NET_DVR_GetLastError());
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
        printf(">>> NET_DVR_CloseAlarmChan_V30 Failed! Error number: %d\n", NET_DVR_GetLastError());
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


int main(int argc, char **argv)
{
    int ret = 0;

    NVR_Init();

    NVR_Connect();

    ret = NVR_Login();
    if (ret < 0)
    {
        printf(">>> NVR Login failed...\n");
        return -1;
    }

    NVR_Htime();

    SetupAlarm();

    SetMessageCallBack();

    for (;;)
    {
        SetMessageCallBack();
        sleep(1);
    }

    return 0;
}


