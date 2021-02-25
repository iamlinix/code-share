#include "db-util.h"
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <time.h>
#include "loguru.hpp"
#include "utf8.h"

typedef PARAMVARY VARY2;

#ifndef FB_ALIGN
#define FB_ALIGN(n,b) ((n+1) & ~1)
#endif

MYSQL* my_connect_database(const char* host, const char* user, const char* pass, const char* db, int port) {
    MYSQL* conn = mysql_init(NULL);
    if (conn) {
        if (NULL == mysql_real_connect(conn, host, user, pass, db, port, NULL, 0)) {
            LOG_F(ERROR, "failed to connect to database: %s:%d", host, port);
            mysql_close(conn);
            conn = NULL;
        }

        if (mysql_set_character_set(conn, "utf8")) {
            LOG_F(ERROR, "failed to to set utf8 encoding");
            mysql_close(conn);
            conn = NULL;
        }
    } else {
        LOG_F(ERROR, "failed to initialize mysql connection");
    }

    return conn;
}

void my_close_connection(MYSQL** conn) {
    if (conn && *conn) {
        mysql_close(*conn);
        *conn = NULL;
    }
}

char** my_execute_query(MYSQL* conn, const char* query, int* rows, int* cols) {
    char** results = NULL;

    if (!conn || !query || !rows || !cols) {
        LOG_F(ERROR, "'execute_query': conn & query & rows & cols cannot be null");
    } else {
        if (mysql_real_query(conn, query, strlen(query))) {
            LOG_F(ERROR, "failed to execute query: [%s]", query);
        } else {
            MYSQL_RES* res = mysql_store_result(conn);
            if (res) {
                MYSQL_ROW row;
                int r, c, rn, cn;

                *rows = rn = mysql_num_rows(res);
                *cols = cn = mysql_num_fields(res);
                if (rn > 0) {
                    results = new char *[rn * cn];

                    r = 0;
                    while ((row = mysql_fetch_row(res))) {
                        for (c = 0; c < cn; ++c) {
                            if (row[c]) {
                                results[r + c] = strdup(row[c]);
                            } else {
                                results[r + c] = strdup("NULL");
                            }
                        }
                        r += cn;
                    }
                }

                mysql_free_result(res);
            }
        }
    }

    return results;
}

bool my_execute(MYSQL* conn, const char* query) {
    bool success = true;
    int ret = 0;
    if (!conn || !query) {
        LOG_F(ERROR, "'execute': conn & query cannot be null");
        success = false;
    } else {
        if (ret = mysql_query(conn, query)) {
            LOG_F(ERROR, "failed to execute query: [%s], error: %d, %s", query, ret, mysql_error(conn));
            success = false;
        }
    }
    return success;
}

int my_execute_update(MYSQL* conn, const char* query) {
    int affected = 0;
    int ret = 0;
    if (!conn || !query) {
        LOG_F(ERROR, "'execute_update': conn & query cannot be null");
    } else {
        if ((ret = mysql_query(conn, query))) {
            LOG_F(ERROR, "failed to execute query: [%s], error: %d", query, ret);
        } else {
            affected = mysql_affected_rows(conn);
        }
    }
    return affected;
}

void my_free_results(char** results, long row, long col) {
    if (results) {
        long total = row * col, i;
        for (i = 0; i < total; ++ i) {
            free(results[i]);
        }
        delete results;
    }
}

bool fb_connect_database(const char* host, const char* user, const char* pass, const char* db, int port, isc_db_handle* pdb, isc_tr_handle* ptrans, isc_stmt_handle* pstmt, XSQLDA** psqlda) {
    bool success = true;
    ISC_SCHAR* dpb = NULL;
    short dpb_length = 0;
    ISC_STATUS status[64];
    char db_full_path[1024];

    if (port > 0)
        sprintf(db_full_path, "%s/%d:%s", host, port, db);
    else
        sprintf(db_full_path, "%s:%s", host, db);
    isc_modify_dpb(&dpb, &dpb_length, isc_dpb_user_name, user, strlen(user));
    isc_modify_dpb(&dpb, &dpb_length, isc_dpb_password, pass, strlen(pass));
    isc_modify_dpb(&dpb, &dpb_length, isc_dpb_sql_dialect, "\3", 1);

    if (isc_attach_database(status, 0, db_full_path, pdb, dpb_length, dpb)) {
        isc_print_status(status);
        LOG_F(ERROR, "failed to connect to database: %s:%d", host, port);
        success = false;
    } else {
        if (isc_dsql_allocate_statement(status, pdb, pstmt)) {
            isc_print_status(status);
            LOG_F(ERROR, "failed to allocate statement");
            fb_cleanup(pdb, pstmt, ptrans, psqlda, GMS_CLEANUP_CLOSE);
            success = false;
        } else {
            if (isc_start_transaction(status, ptrans, 1, pdb, 0, NULL)) {
                isc_print_status(status);
                LOG_F(ERROR, "failed to start transaction");
                fb_cleanup(pdb, pstmt, ptrans, psqlda, GMS_CLEANUP_CLOSE | GMS_CLEANUP_STMT);
                success = false;
            } else {
                if (psqlda) {
                    if (*psqlda)
                        LOG_F(WARNING, "sql data pointer is not null");
                    *psqlda = (XSQLDA ISC_FAR *) malloc(XSQLDA_LENGTH(64));
                    (*psqlda)->sqln = 64;
                    (*psqlda)->version = 1;
                }
            }
        }
    }

    if (dpb)
        isc_free(dpb);

    return success;
}

bool fb_execute(ISC_STATUS* status, isc_tr_handle* ptrans, isc_stmt_handle* pstmt, XSQLDA* sqlda, const char* query, char* buffer) {
    bool success = true;
    int num_cols, offset, length, alignment, type, i;
    XSQLVAR* var;

    if (!status || !ptrans || !pstmt || !sqlda || !query) {
        LOG_F(ERROR, "remote execute: parameters cannot be null");
        success = false;
    } else {
        if (isc_dsql_prepare(status, ptrans, pstmt, 0, query, SQL_DIALECT_V6, sqlda)) {
            LOG_F(ERROR, "remote execute: failed to prepare statement");
            isc_print_status(status);
            success = false;
        } else {
            num_cols = sqlda->sqld;
            for (var = sqlda->sqlvar, offset = 0, i = 0; i < num_cols; ++var, ++i) {
                length = alignment = var->sqllen;
                type = var->sqltype & ~1;

                if (type == SQL_TEXT)
                    alignment = 1;
                else if (type == SQL_VARYING)
                {
                    length += sizeof (short) + 1;
                    alignment = sizeof (short);
                }
                /*  RISC machines are finicky about word alignment
                **  So the output buffer values must be placed on
                **  word boundaries where appropriate
                */
                offset = FB_ALIGN(offset, alignment);
                var->sqldata = (char *) buffer + offset;
                offset += length;
                offset = FB_ALIGN(offset, sizeof(short));
                var->sqlind = (short*) ((char *) buffer + offset);
                offset += sizeof  (short);
            }

            if (isc_dsql_execute(status, ptrans, pstmt, SQL_DIALECT_V6, NULL)) {
                LOG_F(ERROR, "remote execute: failed to run query: [%s]", query);
                success = false;
            }
        }
    }

    return success;
}

void fb_fetch_column(XSQLVAR* var, char* buffer, int buffer_size) {
    short       dtype;
    char        *p = buffer, escape[128];
    char        blob_s[20], date_s[25];
    VARY2        *vary2;
    struct tm   times;
    ISC_QUAD    bid;

    if (!buffer)
        return;

    memset(buffer, 0, buffer_size);

    dtype = var->sqltype & ~1;

    /* Null handling.  If the column is nullable and null */
    if ((var->sqltype & 1) && (*var->sqlind < 0))
    {
        sprintf(p, "%s", "NULL");
    }
    else
    {
        switch (dtype)
        {
            case SQL_TEXT:
                sprintf(p, "%s", var->sqldata);
                gb_to_utf8(p, p, 128);
                mysql_escape_string(escape, p, strlen(p));
                sprintf(p, "'%s'", escape);
                break;

            case SQL_VARYING:
                vary2 = (VARY2*) var->sqldata;
                vary2->vary_string[vary2->vary_length] = '\0';
                sprintf(p, "%s", vary2->vary_string);
                gb_to_utf8(p, p, 128);
                mysql_escape_string(escape, p, strlen(p));
                sprintf(p, "'%s'", escape);
                break;

            case SQL_SHORT:
            case SQL_LONG:
            case SQL_INT64:
            {
                ISC_INT64	value;
                short		dscale;
                switch (dtype)
                {
                    case SQL_SHORT:
                        value = (ISC_INT64) *(short *) var->sqldata;
                        break;
                    case SQL_LONG:
                        value = (ISC_INT64) *(int *) var->sqldata;
                        break;
                    case SQL_INT64:
                        value = (ISC_INT64) *(ISC_INT64 *) var->sqldata;
                        break;
                }
                dscale = var->sqlscale;
                if (dscale < 0)
                {
                    ISC_INT64	tens;
                    short	i;

                    tens = 1;
                    for (i = 0; i > dscale; i--)
                        tens *= 10;

                    if (value >= 0)
                        sprintf (p, "%lld.%lld",
                                 (ISC_INT64) value / tens,
                                 (ISC_INT64) value % tens);
                    else if ((value / tens) != 0)
                        sprintf (p, "%lld.%lld",
                                 (ISC_INT64) (value / tens),
                                 (ISC_INT64) -(value % tens));
                    else
                        sprintf (p, "-0.%lld",
                                 (ISC_INT64) -(value % tens));
                }
                else if (dscale)
                    sprintf (p, "%lld%0*d",
                             (ISC_INT64) value, dscale, 0);
                else
                    sprintf (p, "%lld",
                             (ISC_INT64) value);
            }
                break;

            case SQL_FLOAT:
                sprintf(p, "%g", *(float *) (var->sqldata));
                break;

            case SQL_DOUBLE:
                sprintf(p, "%f", *(double *) (var->sqldata));
                break;

            case SQL_TIMESTAMP:
                isc_decode_timestamp((ISC_TIMESTAMP *)var->sqldata, &times);
                sprintf(date_s, "%04d-%02d-%02d %02d:%02d:%02d",
                        times.tm_year + 1900,
                        times.tm_mon+1,
                        times.tm_mday,
                        times.tm_hour,
                        times.tm_min,
                        times.tm_sec);
                sprintf(p, "'%s'", date_s);
                break;

            case SQL_TYPE_DATE:
                isc_decode_sql_date((ISC_DATE *)var->sqldata, &times);
                sprintf(date_s, "%04d-%02d-%02d",
                        times.tm_year + 1900,
                        times.tm_mon+1,
                        times.tm_mday);
                sprintf(p, "'%s'", date_s);
                break;

            case SQL_TYPE_TIME:
                isc_decode_sql_time((ISC_TIME *)var->sqldata, &times);
                sprintf(date_s, "%02d:%02d:%02d",
                        times.tm_hour,
                        times.tm_min,
                        times.tm_sec);
                sprintf(p, "'%s'", date_s);
                break;

            case SQL_BLOB:
            case SQL_ARRAY:
                /* Print the blob id on blobs or arrays */
                bid = *(ISC_QUAD *) var->sqldata;
                sprintf(blob_s, "%08x:%08x", bid.gds_quad_high, bid.gds_quad_low);
                sprintf(p, "%17s ", blob_s);
                break;

            default:
                break;
        }
    }
}

void fb_cleanup(isc_db_handle* pdb, isc_stmt_handle* pstmt, isc_tr_handle* ptrans, XSQLDA** psqlda, int ops) {
    ISC_STATUS status[64];

    if ((ops & GMS_CLEANUP_STMT) && pstmt && *pstmt) {
        if (isc_dsql_free_statement(status, pstmt, DSQL_close)) {
            isc_print_status(status);
            LOG_F(ERROR, "failed to free statement");
        }
        *pstmt = 0;
    }

    if ((ops & GMS_CLEANUP_COMMIT) && ptrans && *ptrans) {
        if (isc_commit_transaction(status, ptrans)) {
            isc_print_status(status);
            LOG_F(ERROR, "failed to commit transaction");
        }
        *ptrans = 0;
    }

    if ((ops & GMS_CLEANUP_SQLDA) && psqlda && *psqlda) {
        free(*psqlda);
        *psqlda = 0;
    }

    if ((ops & GMS_CLEANUP_CLOSE) && pdb && *pdb) {
        if (isc_detach_database(status, pdb)) {
            isc_print_status(status);
            LOG_F(ERROR, "failed to close database connection");
        }
        *pdb = 0;
    }
}