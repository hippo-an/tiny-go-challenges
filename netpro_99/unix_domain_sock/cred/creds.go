package cred

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
)

// init 에서는 클라이언트의 연결을 허용하는 그룹 ID의 목록을 flag 로 전달받습니다.
func init() {
	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage:\n\t%s <group names>\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

}

// GroupID 를 사용하여 user 패키지에서 그룹을 조회하고 group id 여부에 대해서 map 으로 저장합니다.
func parseGroupNames(args []string) map[string]struct{} {
	groups := make(map[string]struct{})
	for _, arg := range args {
		grp, err := user.LookupGroup(arg)
		if err != nil {
			log.Println(err)
			continue
		}
		groups[grp.Gid] = struct{}{}

	}
	return groups
}

func main() {
	flag.Parse()

	groups := parseGroupNames(flag.Args())
	socketAddr := filepath.Join(os.TempDir(), "creds.sock")
	addr, err := net.ResolveUnixAddr("unix", socketAddr)
	if err != nil {
		log.Fatal(err)
	}

	s, err := net.ListenUnix("unix", addr)

	if err != nil {
		log.Fatal(err)
	}
	// 시그널 대기 및 고루틴을 통한 우아한 리스너 종료
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		_ = s.Close()
	}()

	fmt.Printf("Listening on %s ...\n", socketAddr)
	for {
		conn, err := s.AcceptUnix()
		if err != nil {
			break
		}

		if Allowed(conn, groups) {
			_, err = conn.Write([]byte("Welcome\n"))
			if err == nil {
				continue
			}
		} else {
			_, err = conn.Write([]byte("Access denied\n"))
		}
		if err != nil {
			log.Println(err)
		}
		_ = conn.Close()
	}
}
