#include "business.h"
#include "db-util.h"
#include <ibase.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include "loguru.hpp"

#define FIELD_MAX 128

ISC_STATUS status[64];

void init_milestones(void* conn, void* milestones) {
	if (!my_execute((MYSQL *)conn, "INSERT IGNORE INTO milestone (ind, milestone) VALUES (0, -1), (1, -1) "
			", (2, -1), (3, -1), (4, -1), (5, -1), (6, -1), (7, -1), (8, -1), (9, -1), (10, -1)")) {
		LOG_F(ERROR, "failed to initiate milestones");
	} else {
		int rows, cols, rn;
		int ind, stone;
		char** results = my_execute_query((MYSQL *)conn, "SELECT ind, milestone FROM milestone ORDER BY ind",
				&rows, &cols);
		if (results && rows > 0 && cols > 0) {
			LOG_F(INFO, "%d milestones restored", rows);
			if (rows != GMS_MSTONE_MAX) {
				LOG_F(WARNING, "milestone number mismatch: %d vs %d", rows, GMS_MSTONE_MAX);
			}

			rn = 0;
			for (int i = 0; i < rows; ++ i) {
				ind = strtol(results[rn], NULL, 10);
				stone = strtol(results[rn + 1], NULL, 10);
				if (ind < GMS_MSTONE_MAX)
					((int *)milestones)[ind] = stone;
				else
					LOG_F(ERROR, "milestone index overflow: %d vs %d", ind, GMS_MSTONE_MAX);
				rn += cols;
			}
			my_free_results(results, rows, cols);
		} else {
			LOG_F(ERROR, "failed to load milestones");
		}
	}
}

bool item_handler(void* p1, void* p2, void* p3, void* p4, void* p5) {
    XSQLDA* sqlda = (XSQLDA *)p1;
	isc_stmt_handle* pstmt = (isc_stmt_handle *)p2;
    MYSQL* conn = (MYSQL *)p3;
    int* stone = (int *)p4;
	char* query = (char *)p5;
	bool success = true;

    if (sqlda && pstmt && conn && stone && query) {
		char val[FIELD_MAX];
		int num_cols = sqlda->sqld, acc, total = 0;
		int latest_id = *stone;

		sprintf(query, "%s", "REPLACE INTO item (ITEM_ID, ITEM_CODE, ITEM_NAME) VALUES ");
		acc = strlen(query);

		while (isc_dsql_fetch(status, pstmt, SQL_DIALECT_V6, sqlda) == 0) {
			if (total > 0) {
				query[acc++] = ',';
			}
			++total;

			latest_id = *(int *) sqlda->sqlvar[0].sqldata;

			query[acc++] = '(';
			for (int i = 0; i < num_cols; ++i) {
				fb_fetch_column(&sqlda->sqlvar[i], val, FIELD_MAX);
				if (i != 0)
					query[acc++] = ',';
				sprintf(query + acc, "%s", val);
				acc += strlen(val);
			}
			query[acc++] = ')';
		}

		if (total > 0) {
			if (!my_execute(conn, query)) {
				LOG_F(ERROR, "failed to update item");
				success = false;
				latest_id = -1;
			}
			else {
				if (latest_id > 0 && latest_id > *stone) {
					sprintf(query, "UPDATE milestone SET milestone = %d WHERE ind = %d", latest_id, GMS_MSTONE_ITEM);
					if (my_execute_update(conn, query) <= 0) {
						LOG_F(ERROR, "failed to update item milestone");
						success = false;
					}
					else {
						*stone = latest_id;
					}
				}
			}
		}

		if (success)
			LOG_F(INFO, "[%d] items updated, latest id is: %d", total, latest_id);
	} else {
		LOG_F(ERROR, "invalid parameters in items handler");
	}

	return success;
}

bool promo_handler(void* p1, void* p2, void* p3, void* p4, void* p5) {
    XSQLDA* sqlda = (XSQLDA *)p1;
    isc_stmt_handle* pstmt = (isc_stmt_handle *)p2;
    MYSQL* conn = (MYSQL *)p3;
    int* stone = (int *)p4;
	char* query = (char *)p5;
	bool success = true;

    if (sqlda && pstmt && conn && stone && query) {
		char promo_id[FIELD_MAX], promo_name[FIELD_MAX], sdate[FIELD_MAX], stime[FIELD_MAX], edate[FIELD_MAX], etime[FIELD_MAX];
		int num_cols = sqlda->sqld, acc, total = 0;
		int latest_id = *stone;

		sprintf(query, "%s", "REPLACE INTO promo (PROMO_ID, PROMO_NAME, START_DATE, END_DATE) VALUES ");
		acc = strlen(query);

		while (isc_dsql_fetch(status, pstmt, SQL_DIALECT_V6, sqlda) == 0) {
			if (total > 0) {
				query[acc++] = ',';
			}
			++total;

			latest_id = *(int *) sqlda->sqlvar[0].sqldata;

			fb_fetch_column(&sqlda->sqlvar[0], promo_id, FIELD_MAX);
			fb_fetch_column(&sqlda->sqlvar[1], promo_name, FIELD_MAX);
			fb_fetch_column(&sqlda->sqlvar[2], sdate, FIELD_MAX);
			fb_fetch_column(&sqlda->sqlvar[3], stime, FIELD_MAX);
			fb_fetch_column(&sqlda->sqlvar[4], edate, FIELD_MAX);
			fb_fetch_column(&sqlda->sqlvar[5], etime, FIELD_MAX);

			// combine date and time into datetime
			// a datetime field looks like '2018-01-01 00:00:00', INCLUDING the single quote
			memcpy(sdate + 12, stime + 12, 8);
			memcpy(edate + 12, etime + 12, 8);

			sprintf(query + acc, "(%s,%s,%s,%s)", promo_id, promo_name, sdate, edate);
			acc += strlen(query + acc);
			/*
			query[acc++] = '(';
			for (int i = 0; i < num_cols; ++i) {
				fb_fetch_column(&sqlda->sqlvar[i], val, FIELD_MAX);
				if (i != 0)
					query[acc++] = ',';
				sprintf(query + acc, "%s", val);
				acc += strlen(val);
			}
			query[acc++] = ')';
			 */
		}

		if (total > 0) {
			if (!my_execute(conn, query)) {
				LOG_F(ERROR, "failed to update pormo");
				success = false;
				latest_id = -1;
			}
			else {
				if (latest_id > 0 && latest_id > *stone) {
					sprintf(query, "UPDATE milestone SET milestone = %d WHERE ind = %d", latest_id, GMS_MSTONE_PROMO);
					if (my_execute_update(conn, query) <= 0) {
						LOG_F(ERROR, "failed to update pormo milestone");
						success = false;
					}
					else {
						*stone = latest_id;
					}
				}
			}
		}

		if (success)
			LOG_F(INFO, "[%d] promos updated, latest id is: %d", total, latest_id);
	} else {
		LOG_F(ERROR, "invalid parameters in promos handler");
	}

	return success;
}

bool pmnt_handler(void* p1, void* p2, void* p3, void* p4, void* p5) {
    XSQLDA* sqlda = (XSQLDA *)p1;
    isc_stmt_handle* pstmt = (isc_stmt_handle *)p2;
    MYSQL* conn = (MYSQL *)p3;
    int* stone = (int *)p4;
	char* query = (char *)p5;
	bool success = true;

    if (sqlda && pstmt && conn && stone && query) {
		char val[FIELD_MAX];
		int num_cols = sqlda->sqld, acc, total = 0;
		int latest_id = *stone;

		sprintf(query, "%s", "REPLACE INTO pmnt (PMNT_ID, PMNT_CODE, PMNT_SUB_CODE, PMNT_NAME) VALUES ");
		acc = strlen(query);

		while (isc_dsql_fetch(status, pstmt, SQL_DIALECT_V6, sqlda) == 0) {
			if (total > 0) {
				query[acc++] = ',';
			}
			++total;

			latest_id = *(int *) sqlda->sqlvar[0].sqldata;

			query[acc++] = '(';
			for (int i = 0; i < num_cols; ++i) {
				fb_fetch_column(&sqlda->sqlvar[i], val, FIELD_MAX);
				if (i != 0)
					query[acc++] = ',';
				sprintf(query + acc, "%s", val);
				acc += strlen(val);
			}
			query[acc++] = ')';
		}

		if (total > 0) {
			if (!my_execute(conn, query)) {
				LOG_F(ERROR, "failed to update pmnt");
				success = false;
				latest_id = -1;
			}
			else {
				if (latest_id > 0 && latest_id > *stone) {
					sprintf(query, "UPDATE milestone SET milestone = %d WHERE ind = %d", latest_id, GMS_MSTONE_PMNT);
					if (my_execute_update(conn, query) <= 0) {
						LOG_F(ERROR, "failed to update pmnt milestone");
						success = false;
					}
					else {
						*stone = latest_id;
					}
				}
			}
		}

		if (success)
			LOG_F(INFO, "[%d] pmnts updated, latest id is: %d", total, latest_id);
	} else {
		LOG_F(ERROR, "invalid parameters in pmnts handler");
	}

	return success;
}

bool fuel_handler(void* p1, void* p2, void* p3, void* p4, void* p5) {
    XSQLDA* sqlda = (XSQLDA *)p1;
    isc_stmt_handle* pstmt = (isc_stmt_handle *)p2;
    MYSQL* conn = (MYSQL *)p3;
    int* stone = (int *)p4;
	char* query = (char *)p5;
	bool success = true;

	if (sqlda && pstmt && conn && stone && query) {
		char val[FIELD_MAX];
		int num_cols = sqlda->sqld, acc, total = 0;
		int latest_sernum = *stone;

		sprintf(query, "%s",
			"REPLACE INTO fuel_sales (ID, POST_TIME, LITERS, TOTAL_PRICE, UNIT_PRICE, PUMP_ID, PROMO_REF, ITEM_ID, TILL_NUM) VALUES ");
		acc = strlen(query);

		while (isc_dsql_fetch(status, pstmt, SQL_DIALECT_V6, sqlda) == 0) {
			if (total > 0) {
				query[acc++] = ',';
			}
			++total;

			latest_sernum = *(int *) sqlda->sqlvar[0].sqldata;

			query[acc++] = '(';
			for (int i = 0; i < num_cols; ++i) {
				fb_fetch_column(&sqlda->sqlvar[i], val, FIELD_MAX);
				if (i != 0)
					query[acc++] = ',';
				sprintf(query + acc, "%s", val);
				acc += strlen(val);
			}
			query[acc++] = ')';
		}

		if (total > 0) {
			if (!my_execute(conn, query)) {
				LOG_F(ERROR, "failed to update fuel till item");
				success = false;
				latest_sernum = -1;
			}
			else {
				if (latest_sernum > 0 && latest_sernum > *stone) {
					sprintf(query, "UPDATE milestone SET milestone = %d WHERE ind = %d", latest_sernum, GMS_MSTONE_FUEL);
					if (my_execute_update(conn, query) <= 0) {
						LOG_F(ERROR, "failed to update fuel item milestone");
						success = false;
					}
					else {
						*stone = latest_sernum;
					}
				}
			}
		}

		if (success)
			LOG_F(INFO, "[%d] fuel items updated, latest sernum is: %d", total, latest_sernum);
	} else {
		LOG_F(ERROR, "invalid parameters in fuel handler");
	}

	return success;
}

bool non_fuel_handler(void* p1, void* p2, void* p3, void* p4, void* p5) {
    XSQLDA* sqlda = (XSQLDA *)p1;
    isc_stmt_handle* pstmt = (isc_stmt_handle *)p2;
    MYSQL* conn = (MYSQL *)p3;
    int* stone = (int *)p4;
	char* query = (char *)p5;
	bool success = true;

	if (sqlda && pstmt && conn && stone && query) {
		char val[FIELD_MAX];
		int num_cols = sqlda->sqld, acc, total = 0;
		int latest_sernum = *stone;

		sprintf(query, "%s",
			"REPLACE INTO non_fuel_sales (ID, POST_TIME, QUANTITY, TOTAL_PRICE, UNIT_PRICE, PROMO_REF, ITEM_ID, TILL_NUM) VALUES ");
		acc = strlen(query);

		while (isc_dsql_fetch(status, pstmt, SQL_DIALECT_V6, sqlda) == 0) {
			if (total > 0) {
				query[acc++] = ',';
			}
			++total;

			latest_sernum = *(int *) sqlda->sqlvar[0].sqldata;

			query[acc++] = '(';
			for (int i = 0; i < num_cols; ++i) {
				fb_fetch_column(&sqlda->sqlvar[i], val, FIELD_MAX);
				if (i != 0)
					query[acc++] = ',';
				sprintf(query + acc, "%s", val);
				acc += strlen(val);
			}
			query[acc++] = ')';
		}

		if (total > 0) {
			if (!my_execute(conn, query)) {
				LOG_F(ERROR, "failed to update fuel till item");
				success = false;
				latest_sernum = -1;
			}
			else {
				if (latest_sernum > 0 && latest_sernum > *stone) {
					sprintf(query, "UPDATE milestone SET milestone = %d WHERE ind = %d", latest_sernum, GMS_MSTONE_NON_FUEL);
					if (my_execute_update(conn, query) <= 0) {
						LOG_F(ERROR, "failed to update non fuel item milestone");
						success = false;
					}
					else {
						*stone = latest_sernum;
					}
				}
			}
		}

		if (success)
			LOG_F(INFO, "[%d] non fuel items updated, latest sernum is: %d", total, latest_sernum);
	} else {
		LOG_F(ERROR, "invalid parameters in non fuel handler");
	}

	return success;
}

bool till_promo_handler(void* p1, void* p2, void* p3, void* p4, void* p5) {
	XSQLDA* sqlda = (XSQLDA *)p1;
    isc_stmt_handle* pstmt = (isc_stmt_handle *)p2;
    MYSQL* conn = (MYSQL *)p3;
    int* stone = (int *)p4, latest_num = *stone;
	char* query = (char *)p5;
	bool success = true;

	if (sqlda && pstmt && conn && stone && query) {
		char val[FIELD_MAX];
		int num_cols = sqlda->sqld, acc, total = 0, inserted = 0;

		sprintf(query, "%s",
			"INSERT INTO till_promo (PROMO_ID, TILL_NUM, PROMO_QUANTITY, TOTAL_SAVINGS) VALUES ");
		acc = strlen(query);

		while (isc_dsql_fetch(status, pstmt, SQL_DIALECT_V6, sqlda) == 0) {
			if (total > 0) {
				query[acc++] = ',';
			}
			++total;

            latest_num = *(int *) sqlda->sqlvar[1].sqldata;

			query[acc++] = '(';
			for (int i = 0; i < num_cols; ++i) {
				fb_fetch_column(&sqlda->sqlvar[i], val, FIELD_MAX);
				if (i != 0)
					query[acc++] = ',';
				sprintf(query + acc, "%s", val);
				acc += strlen(val);
			}
			query[acc++] = ')';
		}

		if (total > 0) {
			if (!my_execute(conn, query)) {
				LOG_F(ERROR, "failed to update till promos");
				success = false;
			}
			else {
				if (total > 0) {
					sprintf(query, "UPDATE milestone SET milestone = milestone + %d WHERE ind = %d", latest_num, GMS_MSTONE_TILL_PROMO);
					if (my_execute_update(conn, query) <= 0) {
						LOG_F(ERROR, "failed to update till promo milestone");
						success = false;
					}
					else {
						*stone = latest_num;
					}
				}
			}
		}

		if (success)
			LOG_F(INFO, "[%d] till promos updated, latest size: %d", total, *stone);
	} else {
		LOG_F(ERROR, "invalid parameters in till promo handler");
	}

	return success;
}

bool till_pmnt_handler(void* p1, void* p2, void* p3, void* p4, void* p5) {
	XSQLDA* sqlda = (XSQLDA *)p1;
    isc_stmt_handle* pstmt = (isc_stmt_handle *)p2;
    MYSQL* conn = (MYSQL *)p3;
    int* stone = (int *)p4;
	char* query = (char *)p5;
	bool success = true;

	if (sqlda && pstmt && conn && stone && query) {
		char val[FIELD_MAX];
		int num_cols = sqlda->sqld, acc, total = 0;
		int latest_id = *stone;

		sprintf(query, "%s",
			"REPLACE INTO till_pmnt (PMNT_ID, TILL_NUM, POST_TIME, PMNT_CODE, PMNT_SUB_CODE, DISCOUNT_VAL) VALUES ");
		acc = strlen(query);

		while (isc_dsql_fetch(status, pstmt, SQL_DIALECT_V6, sqlda) == 0) {
			if (total > 0) {
				query[acc++] = ',';
			}
			++total;

			latest_id = *(int *) sqlda->sqlvar[0].sqldata;

			query[acc++] = '(';
			for (int i = 0; i < num_cols; ++i) {
				fb_fetch_column(&sqlda->sqlvar[i], val, FIELD_MAX);
				if (i != 0)
					query[acc++] = ',';
				sprintf(query + acc, "%s", val);
				acc += strlen(val);
			}
			query[acc++] = ')';
		}

		if (total > 0) {
			if (!my_execute(conn, query)) {
				LOG_F(ERROR, "failed to update till pmnts");
				success = false;
				latest_id = -1;
			}
			else {
				if (latest_id > 0 && latest_id > *stone) {
					sprintf(query, "UPDATE milestone SET milestone = %d WHERE ind = %d", latest_id, GMS_MSTONE_TILL_PMNT);
					if (my_execute_update(conn, query) <= 0) {
						LOG_F(ERROR, "failed to update till pmnt milestone");
						success = false;
					}
					else {
						*stone = latest_id;
					}
				}
			}
		}

		if (success)
			LOG_F(INFO, "[%d] till pmnts updated, latest id is: %d", total, latest_id);
	} else {
		LOG_F(ERROR, "invalid parameters in till pmnt handler");
	}

	return success;
}

bool credit_handler(void* p1, void* p2, void* p3, void* p4, void* p5) {
	XSQLDA* sqlda = (XSQLDA *)p1;
    isc_stmt_handle* pstmt = (isc_stmt_handle *)p2;
    MYSQL* conn = (MYSQL *)p3;
    int* stone = (int *)p4;
    char* query = (char *)p5;
	bool success = true;

	if (sqlda && pstmt && conn && stone && query) {
		char *p, card_num[FIELD_MAX], pmnt_id[FIELD_MAX];
		int acc, total = 0, affected = 0, i;
		int latest_id = *stone;
		const char* format = "SET CARD_NUM = %s WHERE PMNT_ID = %s";

		if (!my_execute(conn, "START TRANSACTION")) {
			LOG_F(ERROR, "failed to start transaction for updating fuel credit");
			success = false;
		}
		else {
			sprintf(query, "%s", "UPDATE till_pmnt ");
			acc = strlen(query);

			while (isc_dsql_fetch(status, pstmt, SQL_DIALECT_V6, sqlda) == 0) {
				p = query + acc;
                fb_fetch_column(&sqlda->sqlvar[0], card_num, FIELD_MAX);
                fb_fetch_column(&sqlda->sqlvar[1], pmnt_id,  FIELD_MAX);

                // clean card num field: all card num characters have to be printable
                affected = strlen(card_num);
                for (i = 1; i < affected; ++ i) {
                    if (card_num[i] < 32 || card_num[i] > 126) {
                        if (i < affected - 1) {
                            card_num[i] = '\'';
                            card_num[i + 1] = '\0';
                        } else {
                            LOG_F(WARNING, "card num not printable: %s", card_num);
                        }
                        break;
                    }
                }

                sprintf(p, format, card_num, pmnt_id);

				affected = my_execute_update(conn, query);

				if (affected == 0) {
					LOG_F(INFO, "undocumented credit entries retrieved, latest vs retrieved: %d vs %s", latest_id, pmnt_id);
					break;
				}
				else if (affected < 0) {
					LOG_F(ERROR, "failed to update fuel credit, id = %d", latest_id);
					success = false;
					break;
				}

				++total;
				latest_id = *(int ISC_FAR *) sqlda->sqlvar[1].sqldata;
			}

			if (!my_execute(conn, "COMMIT")) {
				LOG_F(ERROR, "failed to commit transaction for updating fuel credit");
				success = false;
			}
			else {
				if (latest_id > 0 && latest_id > *stone) {
					sprintf(query, "UPDATE milestone SET milestone = %d WHERE ind = %d", latest_id, GMS_MSTONE_CREDIT);
					if (!my_execute(conn, query)) {
						LOG_F(ERROR, "failed to update fuel credit milestone");
						success = false;
					}
					else {
						*stone = latest_id;
					}
				}
			}
		}

		if (success)
			LOG_F(INFO, "[%d] fuel payment credit entries updated, latest id is: %d", total, latest_id);
	} else {
		LOG_F(ERROR, "invalid parameters in credit handler");
	}

	return success;
}

bool pump_handler(void* p1, void* p2, void* p3, void* p4, void* p5) {
    XSQLDA* sqlda = (XSQLDA *)p1;
    isc_stmt_handle* pstmt = (isc_stmt_handle *)p2;
    MYSQL* conn = (MYSQL *)p3;
    int* stone = (int *)p4;
    char* query = (char *)p5;
    bool success = true;

    if (sqlda && pstmt && conn && stone && query) {
        char *p, id[FIELD_MAX], fuel_start[FIELD_MAX], fuel_end[FIELD_MAX];
        int acc, total = 0, affected = 0;
        int latest_id = *stone;

        const char* format = "SET PUMP_UP = %s, PUMP_DOWN = %s WHERE ID = %s";

        if (!my_execute(conn, "START TRANSACTION")) {
            LOG_F(ERROR, "failed to start transaction for updating pump timestamp");
            success = false;
        }
        else {
            sprintf(query, "%s", "UPDATE fuel_sales ");
            acc = strlen(query);

            while (isc_dsql_fetch(status, pstmt, SQL_DIALECT_V6, sqlda) == 0) {
                p = query + acc;
                fb_fetch_column(&sqlda->sqlvar[0], id,         FIELD_MAX);
                fb_fetch_column(&sqlda->sqlvar[1], fuel_start, FIELD_MAX);
                fb_fetch_column(&sqlda->sqlvar[2], fuel_end,  FIELD_MAX);
                sprintf(p, format, fuel_start, fuel_end, id);

                affected = my_execute_update(conn, query);

                if (affected == 0) {
                    LOG_F(INFO, "undocumented pump entries retrieved, latest vs retrieved: %d vs %s", latest_id, id);
                    break;
                }
                else if (affected < 0) {
                    LOG_F(ERROR, "failed to update fuel credit, id = %d", latest_id);
                    success = false;
                    break;
                }

                ++total;
                latest_id = *(int ISC_FAR *) sqlda->sqlvar[0].sqldata;
            }

            if (!my_execute(conn, "COMMIT")) {
                LOG_F(ERROR, "failed to commit transaction for updating pump timestamp");
                success = false;
            }
            else {
                if (latest_id > 0 && latest_id > *stone) {
                    sprintf(query, "UPDATE milestone SET milestone = %d WHERE ind = %d", latest_id, GMS_MSTONE_PUMP_TIME);
                    if (!my_execute(conn, query)) {
                        LOG_F(ERROR, "failed to update pump milestone");
                        success = false;
                    }
                    else {
                        *stone = latest_id;
                    }
                }
            }
        }

        if (success)
            LOG_F(INFO, "[%d] pump timestamp entries updated, latest id is: %d", total, latest_id);
    } else {
        LOG_F(ERROR, "invalid parameters in pump handler");
    }

    return success;
}

void bind_till_vehicle(void* p1, void* p2, void* p3, void* p4) {
	MYSQL* conn = (MYSQL *)p1;
	int* stone_till = (int *)p2, *stone_car = (int *)p3, latest_till = *stone_till, latest_car = *stone_car;
	char* query = (char *)p4;
	int r1, c1, r2, c2, i, total = 0;
	char **s1, **s2;

	sprintf(query, "SELECT DISTINCT TILL_NUM, PUMP_UP, PUMP_DOWN FROM fuel_sales WHERE PUMP_UP IS NOT NULL AND TILL_NUM > %d ORDER BY TILL_NUM", *stone_till);
	if ((s1 = my_execute_query(conn, query, &r1, &c1))) {
		for (i = 0; i < r1; ++ i) {
			sprintf(query, "SELECT MAX(ID) AS MAXD, PLATE_NUM, MIN(POST_TIME) AS MINT, MAX(POST_TIME) AS MAXT, COUNT(*) as CNT "
				  "FROM vehicle WHERE CHAR_LENGTH(PLATE_NUM) > 4 AND ID > %d GROUP BY PLATE_NUM HAVING CNT > 1 AND "
	              "TIME_TO_SEC(TIMEDIFF(MAXT, MINT)) < 1800 AND TIME_TO_SEC(TIMEDIFF(MAXT, MINT)) > 30 "
			      "AND MINT < '%s' AND TIME_TO_SEC(TIMEDIFF('%s', MINT)) < 1200 ORDER BY MINT LIMIT 1", latest_car,
			      s1[i * c1 + 1], s1[i * c1 + 1]);
			if ((s2 = my_execute_query(conn, query, &r2, &c2))) {
				sprintf(query, "REPLACE INTO till_vehicle (TILL_NUM, PLATE_NUM) VALUES (%s, '%s')", s1[i * c1], s2[1]);
				if (my_execute_update(conn, query) <= 0) {
					LOG_F(ERROR, "failed to update till vehicle mapping: %s | %s", s1[i * c1], s2[1]);
					my_free_results(s2, r2, c2);
					break;
				}

				latest_car = atoi(s2[0]);
				++ total;

				my_free_results(s2, r2, c2);
			} else {
				LOG_F(WARNING, "no matching car can be found for till: %s", s1[i * c1]);
			}
			latest_till = atoi(s1[i * c1]);
		}

		my_free_results(s1, r1, c1);
	}

	if (latest_till > *stone_till) {
		sprintf(query, "UPDATE milestone SET milestone = %d WHERE ind = %d", latest_till, GMS_MSTONE_MAP_TILL);
		if (!my_execute(conn, query)) {
			LOG_F(ERROR, "failed to update map till milestone");
		}
		else {
			*stone_till = latest_till;
		}
	}

	if (latest_car > *stone_car) {
		sprintf(query, "UPDATE milestone SET milestone = %d WHERE ind = %d", latest_car, GMS_MSTONE_MAP_CAR);
		if (!my_execute(conn, query)) {
			LOG_F(ERROR, "failed to update map car milestone");
		}
		else {
			*stone_car = latest_car;
		}
	}

	LOG_F(INFO, "[%d] till vehicle mappings updated, latest till number is: %d", total, latest_till);
}