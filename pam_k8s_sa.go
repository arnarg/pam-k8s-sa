package main

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -lpam -fPIC

#include <stdlib.h>
#include <security/pam_appl.h>
#include <security/pam_modules.h>

#ifdef __linux__
#include <security/pam_ext.h>
#endif

char* argv_i(const char **argv, int i);
void pam_syslog_str(pam_handle_t *pamh, int priority, const char *str);
*/
import "C"

import (
	"fmt"
	"log/syslog"
	"unsafe"
)

// Required to compile
func main() {}

// Logger
func pamStrError(pamh *C.pam_handle_t, errnum C.int) string {
	return C.GoString(C.pam_strerror(pamh, errnum))
}

func pamSyslog(pamh *C.pam_handle_t, priority syslog.Priority, format string, a ...interface{}) {
	cstr := C.CString(fmt.Sprintf(format, a...))
	defer C.free(unsafe.Pointer(cstr))

	C.pam_syslog_str(pamh, C.int(priority), cstr)
}

type pamLogger struct {
	pamh *C.pam_handle_t
}

func (l *pamLogger) Infof(format string, a ...interface{}) {
	pamSyslog(l.pamh, syslog.LOG_INFO, format, a...)
}

func (l *pamLogger) Warnf(format string, a ...interface{}) {
	pamSyslog(l.pamh, syslog.LOG_WARNING, format, a...)
}

func (l *pamLogger) Errf(format string, a ...interface{}) {
	pamSyslog(l.pamh, syslog.LOG_ERR, format, a...)
}

//export pam_sm_authenticate_go
func pam_sm_authenticate_go(pamh *C.pam_handle_t, flags C.int, argc C.int, argv **C.char) C.int {
	// Create pam logger
	l := &pamLogger{pamh}

	// Copy args to Go strings
	args := make([]string, int(argc))
	for i := 0; i < int(argc); i++ {
		args[i] = C.GoString(C.argv_i(argv, C.int(i)))
	}

	// Parse args
	conf, err := parseConfig(args)
	if err != nil {
		l.Errf("failed to parse arguments: %s", err)
		return C.PAM_SYSTEM_ERR
	}

	// Get user
	var cUser *C.char
	if errnum := C.pam_get_user(pamh, &cUser, nil); errnum != C.PAM_SUCCESS {
		l.Errf("failed to get user: %s", pamStrError(pamh, errnum))
		return errnum
	}

	user := C.GoString(cUser)
	if len(user) == 0 {
		l.Warnf("empty user")
		return C.PAM_USER_UNKNOWN
	}

	// Get password (token)
	var cToken *C.char
	if errnum := C.pam_get_authtok(pamh, C.PAM_AUTHTOK, &cToken, nil); errnum != C.PAM_SUCCESS {
		l.Errf("failed to get token: %s", pamStrError(pamh, errnum))
		return errnum
	}
	token := C.GoString(cToken)

	if err := pamAuthenticate(l, user, token, conf); err != nil {
		l.Errf("failed to verify token: %s", err)
		return C.PAM_AUTH_ERR
	}
	return C.PAM_SUCCESS
}

//export pam_sm_acct_mgmt_go
func pam_sm_acct_mgmt_go(pamh *C.pam_handle_t, flags C.int, argc C.int, argv **C.char) C.int {
	// TODO actually do some checks?
	return C.PAM_SUCCESS
}
