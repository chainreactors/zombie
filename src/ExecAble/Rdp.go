package ExecAble

/*
//linux编译环境配置
#cgo darwin CFLAGS: -DCGO_OS_DARWIN=1
#cgo darwin CFLAGS: -I../../pkg/lib/mac/freerdp2/include/freerdp2
#cgo darwin CFLAGS: -I../../pkg/lib/mac/freerdp2/include/winpr2
#cgo darwin LDFLAGS: -L${SRCDIR}/../../pkg/lib/mac/freerdp2/lib
#cgo darwin LDFLAGS: -lfreerdp2 -lwinpr2 -lssl -lcrypto

   #if defined(CGO_OS_WINDOWS)
   	//static char* os = "windows";
   	#define uint u_int
   #endif



   #if defined(CGO_OS_DARWIN)
   	//static char* os = "darwin";
   #endif


   #if defined(CGO_OS_LINUX)
   	//static char* os = "linux";
   	#define uint int
   #endif


   #include <freerdp/freerdp.h>

   uint rdp_connect(char *server, uint port, char *domain, char *login, char *password) {
       uint err;
   	err = 500;
       freerdp* instance;
       instance = freerdp_new();
       if (instance == NULL || freerdp_context_new(instance) == FALSE) {
           return err;
       }
       instance->settings->Username = login;
       instance->settings->Password = password;
       instance->settings->IgnoreCertificate = TRUE;
       instance->settings->AuthenticationOnly = TRUE;
       instance->settings->ServerHostname = server;
       instance->settings->ServerPort = port;
       instance->settings->Domain = domain;
       freerdp_connect(instance);
       err = freerdp_get_last_error(instance->context);
   	if (err == 0){
   		freerdp_disconnect(instance);
   		freerdp_free(instance);
   		err = 200;
   		return err;
   	}
   	if (err == 0x00020014){
   		freerdp_free(instance);
   		err = 404;
   		// login failure
   		return err;
   	}
   	if (err == 0x00020015){
   		freerdp_free(instance);
   		err = 404;
   		// login failure
   		return err;
   	}
   	if (err == 0x0002000c){
   		freerdp_free(instance);
   		err = 501;
   		// cannot establish rdp connection, either the port is not opened or it's
   		// no rdp
   		return err;
   	}
   	freerdp_free(instance);
       return err;
       //switch (err) {
       //    case 0:
   	//
       //    case 0x00020009:
       //    case 0x0002000d:
       //    case 0x00020006:
       //    case 0x00020008:
       //    case :
   	//		freerdp_free(instance);
   	//		return err;
       //        // cannot establish rdp connection, either the port is not opened or it's
       //        // not rdp
       //}
   }

   uint check_rdp(char *ip, uint port, char *domain, char *login, char *password) {
       uint login_result = 0;
       wLog *root = WLog_GetRoot();
       WLog_SetStringLogLevel(root, "OFF");
       login_result = rdp_connect(ip, port, domain, login, password);
       return login_result;
   }
*/
import "C"
import (
	"Zombie/src/Utils"
	"errors"
	"strings"
	"sync"
)

var mtx sync.Mutex

type RdpService struct {
	Utils.IpInfo
	Username string `json:"username"`
	Password string `json:"password"`
	Input    string
}

func (s *RdpService) Query() bool {
	return false
}

func (s *RdpService) GetInfo() bool {
	return false
}

func (s *RdpService) Connect() bool {
	err, res := RdpConnectTest(s.Username, s.Password, s.IpInfo)
	if err == nil && res {
		return true
	}
	return false

}

func (s *RdpService) DisConnect() bool {
	return false
}

func (s *RdpService) SetQuery(query string) {
	s.Input = query
}

func (s *RdpService) Output(res interface{}) {

}

func RdpConnectTest(User string, Password string, info Utils.IpInfo) (err error, result bool) {

	var UserName, DoaminName string

	if strings.Contains(User, "/") {
		UserName = strings.Split(User, "/")[1]
		DoaminName = strings.Split(User, "/")[0]
	} else {
		UserName = User
		DoaminName = ""
	}

	res, _ := RdpConnect(info.Ip, DoaminName, UserName, Password, info.Port)

	if res == true {
		result = true
	}

	return err, result
}

func RdpConnect(ip, domain, login, password string, port int) (bool, error) {
	mtx.Lock()
	defer mtx.Unlock()

	var nIp *C.char = C.CString(ip)
	var nDomain *C.char = C.CString(domain)
	var nLogin *C.char = C.CString(login)
	var nPassword *C.char = C.CString(password)
	var nPort C.uint = C.uint(port)

	rInt := uint(C.check_rdp(nIp, nPort, nDomain, nLogin, nPassword))
	switch rInt {
	case 200:
		return true, nil
	case 500:
		return false, errors.New("freerdp init failed")
	case 501:
		return false, errors.New("no rdp")
	default:
		return false, nil
	}

}
