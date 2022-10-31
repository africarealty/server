package errors

const (
	ErrCodeUserIdEmpty                            = "AUTH-001"
	ErrCodeUserEmailEmpty                         = "AUTH-004"
	ErrCodeUserNoValidEmail                       = "AUTH-005"
	ErrCodeUserNameNotUnique                      = "AUTH-006"
	ErrCodeUserNotFound                           = "AUTH-007"
	ErrCodeUserNotActive                          = "AUTH-008"
	ErrCodeUserLocked                             = "AUTH-009"
	ErrCodeLogoutNoSID                            = "AUTH-010"
	ErrCodeUserStorageCreate                      = "AUTH-011"
	ErrCodeUserStorageClearCache                  = "AUTH-012"
	ErrCodeUserStorageUpdate                      = "AUTH-013"
	ErrCodeUserStorageAeroKey                     = "AUTH-014"
	ErrCodeUserStorageGetCache                    = "AUTH-015"
	ErrCodeUserStoragePutCache                    = "AUTH-016"
	ErrCodeUserStorageGetDb                       = "AUTH-017"
	ErrCodeUserStorageCreateIndex                 = "AUTH-018"
	ErrCodeUserStorageGetCacheByUsername          = "AUTH-019"
	ErrCodeUserStorageGetByIds                    = "AUTH-020"
	ErrCodeUserStorageDelete                      = "AUTH-021"
	ErrCodeSessionStorageAeroKey                  = "AUTH-022"
	ErrCodeSessionStorageGetCache                 = "AUTH-023"
	ErrCodeSessionStoragePutCache                 = "AUTH-024"
	ErrCodeSessionStorageClearCache               = "AUTH-025"
	ErrCodeSessionStorageGetDb                    = "AUTH-026"
	ErrCodeSessionGetByUser                       = "AUTH-027"
	ErrCodeSessionStorageUpdateLastActivity       = "AUTH-028"
	ErrCodeSessionStorageUpdateLogout             = "AUTH-029"
	ErrCodeSessionStorageCreateSession            = "AUTH-030"
	ErrCodeSessionNotFound                        = "AUTH-031"
	ErrCodeSessionLoggedOut                       = "AUTH-032"
	ErrCodeSecurityPermissionsDenied              = "AUTH-033"
	ErrCodeSessionAuthorizationInvalidResource    = "AUTH-034"
	ErrCodeSidEmpty                               = "AUTH-035"
	ErrCodeNoAuthHeader                           = "AUTH-036"
	ErrCodeAuthHeaderInvalid                      = "AUTH-037"
	ErrCodeUserInvalidPassword                    = "AUTH-038"
	ErrCodeNoUID                                  = "AUTH-039"
	ErrCodeNotAllowed                             = "AUTH-040"
	ErrCodeUserStorageSetToken                    = "AUTH-041"
	ErrCodeUserStorageGetToken                    = "AUTH-042"
	ErrCodeUserRegEmptyRq                         = "AUTH-043"
	ErrCodeUserRegPasswordTooSimple               = "AUTH-044"
	ErrCodeUserActivationTokenEmpty               = "AUTH-045"
	ErrCodeUserActivationNotExistedOnInvalidToken = "AUTH-046"
	ErrCodeUserActivationInvalidOperation         = "AUTH-047"
	ErrCodeUserRegPasswordNotSpecified            = "AUTH-048"
	ErrCodeUserRegPasswordConfirmationNotEqual    = "AUTH-049"
)
