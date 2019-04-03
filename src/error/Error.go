package error

import "errors"

type ApiError struct {
	Status       string `json:"status"`
	ErrorCode    int32  `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	Err          error  `json:"-"`
}

func new(code int32, message string) ApiError {
	return ApiError{"failure", code, message, errors.New(message)}
}

var (

	// ErrNoToken :
	ErrNoToken = new(1046, "no token")
	// ErrTokenExpire :
	ErrTokenExpire = new(1045, "expire/connection fail")
	// ErrInvalidToken :
	ErrInvalidToken = new(1048, "invalid token")
	// ErrAuthFail :
	ErrAuthFail = new(1044, "auth fail(login twice)")

	// ErrNotFound : target not found
	ErrNotFound = new(2000, "not found")

	// ErrSQLConn : raised when sql connections
	ErrSQLConn = new(3000, "SQL connection fail")
	// ErrSQLExec : raised when sql execution
	ErrSQLExec = new(3000, "SQL execution fail")
	// ErrSQLScan : raised when sql row scan fail
	ErrSQLScan = new(3001, "SQL scan fail")

	// ErrRequestParam :
	ErrRequestParam = new(5000, "param check fail")

	// ErrUUIDGen :
	ErrUUIDGen = new(8000, "uuid gen fail")
	// ErrFileOpen :
	ErrFileOpen = new(8001, "file open fail")
	// ErrFileRead :
	ErrFileRead = new(8002, "file read fail")
	// ErrFileSize :
	ErrFileSize = new(8003, "file size")
	// ErrFileWrite :
	ErrFileWrite = new(8004, "file write fail")

	// ErrEsConn :
	ErrEsConn = new(9000, "elasticsearch connection fail")
	// ErrEsExec :
	ErrEsExec = new(9001, "elasticsearch execution fail")
	// ErrEsIndexNotExist :
	ErrEsIndexNotExist = new(9002, "elasticsearch index not exist")

	// ErrRedisConn :
	ErrRedisConn = new(3100, "redis connection fail")
	// ErrRedisExec :
	ErrRedisExec = new(3101, "redis exection fail")
	// ErrTokenGen :
	ErrTokenGen = new(1047, "token gen tail")
	// ErrLogin :
	ErrLogin = new(5001, "login fail")

	// ErrWechatBrandCreateFail :
	ErrWechatBrandCreateFail = new(5002, "create fail")
	// ErrWechatBrandUpdateFail :
	ErrWechatBrandUpdateFail = new(5003, "update fail")
	// ErrWechatBrandDeleteFail :
	ErrWechatBrandDeleteFail = new(5004, "delete fail")
)
