#include "ibase.h"
#include "utf8.h"
#include "loguru.hpp"
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <signal.h>
#include <getopt.h>
#include <iostream>
#include <string>
#include "mysql.h"
#include "db-util.h"
#include "business.h"

const int BUFFER_SIZE_MAX = 64 * 1024 * 1024;
const int ONEK = 1024;
const int ONEM = ONEK * ONEK;

const char* business_queries[] = {
        // item
        "SELECT ITEMID, ITEMCODE, ITEMNAME FROM ITEM WHERE ITEMID > %d ORDER BY ITEMID ROWS %d",

        // promo
        "SELECT PROMO_ID, PROMO_NAME, SDATE, STIME, EDATE, ETIME FROM PROMO WHERE PROMO_ID > %d ORDER BY PROMO_ID ROWS %d",

        // pmnt
        "SELECT PMNT_ID, PMNTCODE_ID, PMSUBCODE_ID, PMNT_NAME FROM PMNT WHERE PMNT_ID > %d ORDER BY PMNT_ID ROWS %d",

        // fuel sales
        "SELECT SERNUM, POSTED, WEIGHT, TOTAL, PRICE, PUMP_ID, PROMO_REF, PLU, TILLNUM FROM TILLITEM WHERE HOSE_ID = 1 AND SERNUM > %d ORDER BY SERNUM ROWS %d",

        // non fuel sales
        "SELECT SERNUM, POSTED, WEIGHT, TOTAL, PRICE, PROMO_REF, PLU, TILLNUM FROM TILLITEM WHERE HOSE_ID = 0 AND SERNUM > %d ORDER BY SERNUM ROWS %d",

        // till pmnt
        "SELECT SERNUM, BILLNUM, POSTED, PMCODE, PMSUBCODE, DISCOUNT_VAL FROM TILLPMNT WHERE SERNUM > %d ORDER BY SERNUM ROWS %d",

        // till pmnt credit
        "SELECT CARDNUM, SERNUM FROM TILLPMNT_CRDT_TRL WHERE SERNUM > %d ORDER BY SERNUM ROWS %d",

        // till promo
        "SELECT PROMO_ID, TILLNUM, PROMO_QUANTITY, TOTAL_SAVINGS FROM TILLPROMO WHERE TILLNUM > %d ORDER BY TILLNUM",

        // pump timestamp
        "SELECT TILLITEMNUM, FUEL_START_TIME, FUEL_END_TIME FROM TILLITEMSTATUS WHERE FUEL_START_TIME IS NOT NULL AND FUEL_START_TIME > '1900-01-01 00:00:00' AND TILLITEMNUM > %d ORDER BY TILLITEMNUM ROWS %d"
};

const query_hander business_handlers[] = {
        item_handler,
        promo_handler,
        pmnt_handler,
        fuel_handler,
        non_fuel_handler,
        till_pmnt_handler,
        credit_handler,
        till_promo_handler,
        pump_handler
};

const struct option long_options[] = {
    { "daemon",         no_argument,        NULL,   'd' },
    { "interval",       required_argument,  NULL,   'i' },
    { "rows",           required_argument,  NULL,   'r' },
    { "server-addr",    required_argument,  NULL,   'a' },
    { "server-port",    required_argument,  NULL,   'p' },
    { "server-user",    required_argument,  NULL,   'u' },
    { "server-pass",    required_argument,  NULL,   'c' },
    { "server-db",      required_argument,  NULL,   'b' },
    { "db-addr",        required_argument,  NULL,   'A' },
    { "db-port",        required_argument,  NULL,   'P' },
    { "db-user",        required_argument,  NULL,   'U' },
    { "db-pass",        required_argument,  NULL,   'C' },
    { "db-db",          required_argument,  NULL,   'B' },
    { NULL,             0,                  NULL,   0   }
};
bool keep_running = true;

void signal_handler(int signum) {
    LOG_F(INFO, "signal received: %d", signum);
    if (signum == SIGINT || signum == SIGKILL || signum == SIGHUP || signum == SIGTERM)
        keep_running = false;
}

int main(int argc, char* argv[]) {
    std::string bos_server = "10.94.185.225";
    int bos_port = 3050;
    std::string bos_user = "SYSDBA";
    std::string bos_pass = "masterkey";
    std::string bos_db = "C:\\Office\\Db\\OFFICE.GDB";
    std::string mysql_server = "localhost";
    int mysql_port = 3306;
    std::string mysql_user = "tav";
    std::string mysql_pass = "Qwer_1234";
    std::string mysql_db = "tav";
    bool daemonize = false;
    int interval = 3600, rows = 1000;
    int c = 0;
    while ((c = getopt_long(argc, argv, "di:r:a:p:u:c:b:A:P:U:C:B:", long_options, NULL)) != -1) {
        switch(c) {
            case 'd':
                daemonize = true;
                break;

            case 'i':
                interval = atoi(optarg);
                break;

            case 'r':
                rows = atoi(optarg);
                break;

            case 'a':
                bos_server = optarg;
                break;

            case 'p':
                bos_port = atoi(optarg);
                break;

            case 'u':
                bos_user = optarg;
                break;

            case 'b':
                bos_db = optarg;
                break;

            case 'c':
                bos_pass = optarg;
                break;

            case 'A':
                mysql_server = optarg;
                break;

            case 'P':
                mysql_port = atoi(optarg);
                break;

            case 'U':
                mysql_user = optarg;
                break;

            case 'C':
                mysql_pass = optarg;
                break;

            case 'B':
                mysql_db = optarg;
                break;

            default:
                std::cout << "unknown argument: " << (char)c << std::endl;
                break;
        }
    }

    if (daemonize) {
        pid_t pid = fork();

        if (pid < 0) {
            std::cout << "fork error\n";
            exit(EXIT_FAILURE);
        }

        if (pid > 0)
            exit(EXIT_SUCCESS);

        if (setsid() < 0) {
            std::cout << "setsid error\n";
            exit(EXIT_FAILURE);
        }

        pid = fork();

        if (pid < 0) {
            std::cout << "second fork error\n";
            exit(EXIT_FAILURE);
        }

        if (pid > 0)
            exit(EXIT_SUCCESS);

        umask(0);

        chdir("/");

        int d;
        for (d = sysconf(_SC_OPEN_MAX); d > 0; -- d) {
            close(d);
        }
    }

    loguru::init(argc, argv);
    if (daemonize)
        loguru::g_stderr_verbosity = loguru::Verbosity_OFF;
    loguru::add_file("/var/log/petro.log", loguru::Append, loguru::Verbosity_MAX);

    int counter = 0, i;
    char* buffer = new char[BUFFER_SIZE_MAX];
    int query_num = sizeof(business_queries) / sizeof(char *);
    int handler_num = sizeof(business_handlers) / sizeof(query_hander);
    if (query_num != handler_num) {
        LOG_F(ERROR, "query num and handler num mismatch: %d vs %d", query_num, handler_num);
        exit(EXIT_FAILURE);
    }

    MYSQL* conn = my_connect_database(mysql_server.c_str(), mysql_user.c_str(), mysql_pass.c_str(), mysql_db.c_str(), mysql_port);
    int* milestones = new int[GMS_MSTONE_MAX];
    if (conn) {
        init_milestones(conn, milestones);
        for (i = 0; i < GMS_MSTONE_MAX; ++ i) {
            LOG_F(INFO, "milestone %02d:%d", i, milestones[i]);
        }
        my_close_connection(&conn);
    } else {
        LOG_F(ERROR, "failed to connect to mysql server: %s:%d", mysql_server.c_str(), mysql_port);
        exit(EXIT_FAILURE);
    }

    signal(SIGTERM, signal_handler);
    signal(SIGKILL, signal_handler);
    signal(SIGINT, signal_handler);
    signal(SIGHUP, signal_handler);

    isc_db_handle fb_db = 0;
    isc_tr_handle fb_trans = 0;
    isc_stmt_handle fb_stmt = 0;
    XSQLDA* sqlda = 0;
    char query_query[ONEK], *data_buffer = new char[ONEM];
    int round = 1, stone;
    ISC_STATUS status[64];
    while (keep_running) {
        if (counter == 0) {
            LOG_F(INFO, "******** round [%07d] starts running ********", round);

            if (fb_connect_database(bos_server.c_str(), bos_user.c_str(), bos_pass.c_str(), bos_db.c_str(),
                    bos_port, &fb_db, &fb_trans, &fb_stmt, &sqlda)) {
                conn = my_connect_database(mysql_server.c_str(), mysql_user.c_str(), mysql_pass.c_str(), mysql_db.c_str(), mysql_port);
                if (conn) {
                    for (i = 0; i < query_num; ++ i) {
                        stone = milestones[i];
                        if (i != GMS_MSTONE_TILL_PROMO)
                            sprintf(query_query, business_queries[i], stone, rows);
                        else
                            sprintf(query_query, business_queries[i], stone);
                        if (!fb_execute(status, &fb_trans, &fb_stmt, sqlda, query_query, data_buffer)) {
                            LOG_F(ERROR, "main loop: error running query [%d]", i);
                            isc_dsql_free_statement(status, &fb_stmt, DSQL_close);
                            break;
                        }

                        memset(buffer, 0, BUFFER_SIZE_MAX);
                        if (!business_handlers[i](sqlda, &fb_stmt, conn, &milestones[i], buffer)) {
                            LOG_F(ERROR, "main loop: handler error: %d", i);
                            break;
                        }

                        isc_dsql_free_statement(status, &fb_stmt, DSQL_close);
                    }

                    fb_cleanup(&fb_db, &fb_stmt, &fb_trans, &sqlda, GMS_CLEANUP_SQLDA | GMS_CLEANUP_COMMIT | GMS_CLEANUP_CLOSE);
                    fb_stmt = 0;

                    bind_till_vehicle(conn, &milestones[GMS_MSTONE_MAP_TILL], &milestones[GMS_MSTONE_MAP_CAR], buffer);

                    my_close_connection(&conn);
                } else
                    LOG_F(ERROR, "failed to connect to local server");

                fb_cleanup(&fb_db, &fb_stmt, &fb_trans, &sqlda, GMS_CLEANUP_ALL);
            } else {
                LOG_F(ERROR, "failed to connect to remote server");
            }

            LOG_F(INFO, "******** round [%07d] completes ********", round);
            ++ round;
        } else {
            sleep(1);
        }

        ++ counter;
        if (counter == interval)
            counter = 0;
    }

    // clean up
    delete[] buffer;
    delete[] milestones;
    delete[] data_buffer;

    LOG_F(INFO, "petro terminates");

    return 0;
}
