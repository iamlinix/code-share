#ifndef TAV_DB_UTIL_H_
#define TAV_DB_UTIL_H_

#include <mysql.h>
#include <ibase.h>

typedef enum _gms_ops {
    GMS_CLEANUP_COMMIT      = 0x01,
    GMS_CLEANUP_STMT        = 0x02,
    GMS_CLEANUP_SQLDA       = 0x04,
    GMS_CLEANUP_CLOSE       = 0x08,
    GMS_CLEANUP_ALL         = 0xFF
} gms_ops;

MYSQL* my_connect_database(const char* host, const char* user, const char* pass, const char* db, int port);
void my_close_connection(MYSQL** conn);
char** my_execute_query(MYSQL* conn, const char* query, int* rows, int* cols);
bool my_execute(MYSQL* conn, const char* query);
int my_execute_update(MYSQL* conn, const char* query);
void my_free_results(char** results, long row, long col);

bool fb_connect_database(const char* host, const char* user, const char* pass, const char* db, int port, isc_db_handle* pdb, isc_tr_handle* ptrans, isc_stmt_handle* pstmt, XSQLDA** psqlda);
bool fb_execute(ISC_STATUS* status, isc_tr_handle* ptrans, isc_stmt_handle* pstmt, XSQLDA* sqlda, const char* query, char* buffer);
void fb_fetch_column(XSQLVAR* var, char* buffer, int buffer_size);
void fb_cleanup(isc_db_handle* pdb, isc_stmt_handle* pstmt, isc_tr_handle* ptrans, XSQLDA** psqlda, int ops);

#endif