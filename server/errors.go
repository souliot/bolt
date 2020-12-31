package server

import (
	"errors"
)

var (
	ErrEmptyKey      = errors.New("server: key is not provided")
	ErrKeyNotFound   = errors.New("server: key not found")
	ErrValueProvided = errors.New("server: value is provided")
	ErrLeaseProvided = errors.New("server: lease is provided")
	ErrTooManyOps    = errors.New("server: too many operations in txn request")
	ErrDuplicateKey  = errors.New("server: duplicate key given in txn request")
	ErrCompacted     = errors.New("server: mvcc: required revision has been compacted")
	ErrFutureRev     = errors.New("server: mvcc: required revision is a future revision")
	ErrNoSpace       = errors.New("server: mvcc: database space exceeded")

	ErrLeaseNotFound    = errors.New("server: requested lease not found")
	ErrLeaseExist       = errors.New("server: lease already exists")
	ErrLeaseTTLTooLarge = errors.New("server: too large lease TTL")

	ErrRequestTooLarge        = errors.New("server: request is too large")
	ErrRequestTooManyRequests = errors.New("server: too many requests")

	ErrRootUserNotExist     = errors.New("server: root user does not exist")
	ErrRootRoleNotExist     = errors.New("server: root user does not have root role")
	ErrUserAlreadyExist     = errors.New("server: user name already exists")
	ErrUserEmpty            = errors.New("server: user name is empty")
	ErrUserNotFound         = errors.New("server: user name not found")
	ErrRoleAlreadyExist     = errors.New("server: role name already exists")
	ErrRoleNotFound         = errors.New("server: role name not found")
	ErrAuthFailed           = errors.New("server: authentication failed, invalid user ID or password")
	ErrPermissionDenied     = errors.New("server: permission denied")
	ErrRoleNotGranted       = errors.New("server: role is not granted to the user")
	ErrPermissionNotGranted = errors.New("server: permission is not granted to the role")
	ErrAuthNotEnabled       = errors.New("server: authentication is not enabled")
	ErrInvalidAuthToken     = errors.New("server: invalid auth token")
	ErrInvalidAuthMgmt      = errors.New("server: invalid auth management")

	ErrNotCapable                 = errors.New("server: not capable")
	ErrStopped                    = errors.New("server: server stopped")
	ErrTimeout                    = errors.New("server: request timed out")
	ErrTimeoutDueToLeaderFail     = errors.New("server: request timed out, possibly due to previous leader failure")
	ErrTimeoutDueToConnectionLost = errors.New("server: request timed out, possibly due to connection lost")
	ErrUnhealthy                  = errors.New("server: unhealthy cluster")
	ErrCorrupt                    = errors.New("server: corrupt cluster")
)
