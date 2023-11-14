#include <security/pam_appl.h>

#ifdef __linux__
#include <security/pam_ext.h>
#endif

// Wrapper to call into go
int pam_sm_authenticate_go(pam_handle_t *pamh, int flags, int argc, char **argv);
int pam_sm_authenticate(pam_handle_t *pamh, int flags, int argc, const char **argv) {
  return pam_sm_authenticate_go(pamh, flags, argc, (char**)argv);
}

// Wrapper to call into go
int pam_sm_acct_mgmt_go(pam_handle_t *pamh, int flags, int argc, char **argv);
int pam_sm_acct_mgmt(pam_handle_t *pamh, int flags, int argc, const char **argv) {
  return pam_sm_acct_mgmt_go(pamh, flags, argc, (char**)argv);
}

int pam_sm_setcred(pam_handle_t *pamh, int flags, int argc, char **argv) {
  return PAM_IGNORE;
}

char* argv_i(char **argv, int i) {
  return argv[i];
}

void pam_syslog_str(pam_handle_t *pamh, int priority, const char *str) {
#ifdef __linux__
  pam_syslog(pamh, priority, "%s", str);
#endif
}
