#include <stdio.h>
#include <string.h>
#include <errno.h>
#include <errmsg.h>
#include <sys/stat.h>

#include "NVRConnector.h"
#include "util.h"


///////////////////////////////////////////////////////////////////////////////
//
MYSQL *DB_Init()
{
    MYSQL *conn;

    /* db configure*/
    char *server = (char *)"127.0.0.1";
    char *user   = (char *)"tav";
    char *password = (char *)"Qwer_1234";
    char *database = (char *)"tav";
    int   port = 3306;

    /* do connect */
    conn = mysql_init(NULL);

    /* Connect to database */
    if (!mysql_real_connect(conn, server, user, password, database, port, NULL, 0))
    {
        LOG_ERR(">>> Connect DB failed: %s\n", mysql_error(conn));
        return NULL;
    }

    if (mysql_set_character_set(conn, "utf8"))
    {
        LOG_ERR(">>> MySQL set character set [utf8] failed: %s", mysql_error(conn));
        mysql_close(conn);
        return NULL;
    }

    LOG_INFO(">>> DB Connected.");

    return conn;
}


MYSQL *DB_Reconnect()
{
    LOG_INFO(">>> DB Re-Connected.");

    return DB_Init();
}


void DB_Close(MYSQL *conn)
{
    LOG_INFO(">>> DB Closed.");

    if (conn)
    {
        mysql_close(conn);
    }
}


static int dump_picture_file(veh_info_t *vehi, char *name, int size)
{
    char full_name[256] = { 0 };
    char full_dir[256] = { 0 };

    if (!vehi->pic_buff || vehi->pic_buff_size <= 0)
    {
        return 0;
    }
	
    //struct tm *fmt_tm = &vehi->t;
    struct tm *fmt_tm;
    fmt_tm = localtime(&vehi->post_time);

    if (vehi->plate_num[0] == 0x0)
    {
        snprintf(name, size, "%04d%02d%02d_%02d%02d%02d_none.jpg",
                 //vehi->tm.wYear, vehi->tm.byMonth, vehi->tm.byDay,
                 //vehi->tm.byHour, vehi->tm.byMinute, vehi->tm.bySecond);
                 1900 + fmt_tm->tm_year, fmt_tm->tm_mon + 1, fmt_tm->tm_mday,
                 fmt_tm->tm_hour, fmt_tm->tm_min, fmt_tm->tm_sec);
    }
    else
    {
        snprintf(name, size, "%04d%02d%02d_%02d%02d%02d_%s.jpg",
                 //vehi->tm.wYear, vehi->tm.byMonth, vehi->tm.byDay,
                 //vehi->tm.byHour, vehi->tm.byMinute, vehi->tm.bySecond,
                 1900 + fmt_tm->tm_year, fmt_tm->tm_mon + 1, fmt_tm->tm_mday,
                 fmt_tm->tm_hour, fmt_tm->tm_min, fmt_tm->tm_sec,
                 vehi->plate_num);
    }

    /* check dir first */
    snprintf(full_dir, 256, "%s/%04d%02d%02d", NVR_PIC_PATH,
             1900 + fmt_tm->tm_year, fmt_tm->tm_mon + 1, fmt_tm->tm_mday);

    if (0 == access(full_dir, 0))
    {

    }
    else
    {
        if (mkdir(full_dir, 0755) == -1)
        {
            LOG_ERR(">>> mkdir %s failed.", full_dir);
            return -1;
        }
        else
        {
            LOG_DBG(">>> mkdir %s success.", full_dir);
        }
    }

    snprintf(full_name, 256, "%s/%s", full_dir, name);
    LOG_INFO(">>> picture file name: %s", full_name);

    FILE *fp = fopen(full_name, "wb");
    if (fp)
    {
        fwrite(vehi->pic_buff, vehi->pic_buff_size, sizeof(char), fp);
        fclose(fp);
        LOG_INFO(">>> dump picture to file.");
    }
    else 
    {
        LOG_INFO(">>> dump picture to file failed: %d(%s).", errno, strerror(errno));
    }

    return 0;
}


int DB_Connector(MYSQL *conn, veh_info_t *vehi)
{
    char sql[1024] = {0};

    /* Dump picture file */
    char out_name[256] = { 0 };
    dump_picture_file(vehi, out_name, 256);

    if (!conn || !vehi)
    {
        LOG_ERR("para is null (conn: %p, vehi: %p).", conn, vehi);
        return -1;
    }

    /* convert time */
    char tbuff[20] = { 0 };
    struct tm *fmt_tm;
    LOG_DBG(">>> post time: %lld", vehi->post_time);
    fmt_tm = localtime(&vehi->post_time);
    strftime(tbuff, 20, "%Y-%m-%d %H:%M:%S", fmt_tm);
    //snprintf(tbuff, 20, "%04d-%02d-%02d %02d:%02d:%02d",
    //         vehi->tm.wYear, vehi->tm.byMonth, vehi->tm.byDay,
    //         vehi->tm.byHour, vehi->tm.byMinute, vehi->tm.bySecond);
    LOG_INFO(">>> time buff: %s", tbuff);

    /* convert time */
    char nbuff[20] = { 0 };
    struct tm *now_tm;
    time_t now;
    time(&now);
    now_tm = localtime(&now);
    strftime(nbuff, 20, "%Y-%m-%d %H:%M:%S", now_tm);
    LOG_INFO(">>> now time buff: %s", nbuff);

    //if (vehi->plate_num[0] == 0x65 && vehi->plate_num[1] == 0xe0)
    //{
    //    LOG_DBG(">>> Got a none plate, skip it.");
    //    return 0;
    //}

    /* format sql command */
    snprintf(sql, 1024, "INSERT INTO vehicle (STATION_ID, CHANNEL_ID, POST_TIME, \
			PLATE_NUM, PLATE_COLOR, VEH_TYPE, VEH_BRAND, VEH_SUB_BRAND, \
                        VEH_MODEL, VEH_PIC_FNAME, CREATE_TIME) \
			VALUES (%d, %d, \"%s\", \"%s\", %d, %d, %d, %d, %d, \"%s\", \"%s\")",
                        vehi->station_id, vehi->channel_id, tbuff, vehi->plate_num, \
                        vehi->plate_color, vehi->veh_type, vehi->veh_brand, \
                        vehi->veh_sub_brand, vehi->veh_model, out_name, nbuff);

    /* send SQL query */
    LOG_INFO(">>> [DB] begin excute sql...");
    LOG_DBG(">>> [DB] command: %s", sql);

    int status = mysql_query(conn, sql);
    if (status != 0)
    {
        if (status == CR_SERVER_LOST || status == CR_SERVER_GONE_ERROR)
        {
            LOG_INFO(">>> [DB] mysql conn lost, reconnect...");

            /* do re-connect */
            DB_Close(conn);
            conn = DB_Reconnect();
            if (conn)
            {
                if (mysql_query(conn, sql))
                {
                    LOG_ERR("mysql query failed.");
                    return -1;
                }
            }
        }
        else
        {
            LOG_ERR(">>> mysql query failed with code %d(%s).",
                     status, mysql_error(conn));
        }
    }
    LOG_INFO(">>> [DB] finish excute sql.");

    //res = mysql_use_result(conn);

    return 0;
}


