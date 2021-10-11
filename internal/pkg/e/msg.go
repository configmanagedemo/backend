package e

// Msgs definition
var Msgs = map[int]string{
	Success:   "ok",
	Error:     "fail",
	NotLogin:  "not login",
	LoginFail: "login fail",
	NoAuth:    "no auth",
	LostParam: "lost param",
}

// GetMsg get msg by code
func GetMsg(code int) string {
	msg, ok := Msgs[code]
	if ok {
		return msg
	}

	return Msgs[Error]
}
