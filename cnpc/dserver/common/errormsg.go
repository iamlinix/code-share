package common


const (
	ERR_CODE_NO_ERROR	= 0
	ERR_CODE_USER_NOT_FOUND = 1000
	ERR_CODE_USER_WRONG_PASS = 1001

	ERR_CODE_TOKEN_GEN_FAILED = 1050

	ERR_CODE_SYS_INVALID_PARA = 1101
	ERR_CODE_SYS_OPEN_READER = 1102
	ERR_CODE_SYS_GET_ROWS = 1103

	ERR_CODE_INVALID_CLASS	= 2000
	ERR_CODE_INVALID_ORG	= 2001
	ERR_CODE_INVALID_MCODE	= 2010
	ERR_CODE_INVALID_PCODE	= 2020

	ERR_CODE_CFG_GET_FAILED	= 3001
	ERR_CODE_CFG_NOT_FOUND	= 3002
	ERR_CODE_CFG_EXIST	= 3003
	ERR_CODE_CFG_ADD_FAILED	= 3004
	ERR_CODE_CFG_UPD_FAILED	= 3005
	ERR_CODE_CFG_DEL_FAILED	= 3006
	ERR_CODE_CFG_GEN_XLS_FAILED	= 3007
	ERR_CODE_CFG_IMPORT_FAILED	= 3008
)

const (
	ERR_MSG_NO_ERROR	= "success"
	ERR_MSG_USER_NOT_FOUND  = "user not found"
	ERR_MSG_USER_WRONG_PASS = "wrong password"

	ERR_MSG_TOKEN_GEN_FAILED = "generate token failed"

	ERR_MSG_SYS_INVALID_PARA = "invalid parameter"
	ERR_MSG_SYS_OPEN_READER = "open reader failed"
	ERR_MSG_SYS_GET_ROWS    = "reader get rows failed"

	ERR_MSG_INVALID_CLASS	= "invalid class level"
	ERR_MSG_INVALID_ORG	= "invalid org level"
	ERR_MSG_INVALID_MCODE	= "invalid material code"
	ERR_MSG_INVALID_PCODE	= "invalid plant code"

	ERR_MSG_CFG_GET_FAILED	= "config get failed"
	ERR_MSG_CFG_NOT_FOUND	= "config not found"
	ERR_MSG_CFG_EXIST	= "config already exist"
	ERR_MSG_CFG_ADD_FAILED	= "config add failed"
	ERR_MSG_CFG_UPD_FAILED	= "config update failed"
	ERR_MSG_CFG_DEL_FAILED	= "config delete failed"
	ERR_MSG_CFG_GEN_XLS_FAILED	= "config gen xls failed"
	ERR_MSG_CFG_IMPORT_FAILED	= "config import failed"
)

