package main

import (
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"
)

func fromXshell() (*url.URL, error){
	u, err:= url.Parse(os.Args[2])
	if err != nil {
		return nil, err
	}
	u.Fragment =os.Args[4]

	return u, nil
}

func fromSecureCRT() (*url.URL, error){
	s := fmt.Sprintf("ssh://%s:%s@%s:%s#%s",
		os.Args[11], os.Args[13], os.Args[7], os.Args[9], os.Args[4])
	return url.Parse(s)
}


func wait() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	<-ch
	os.Exit(-1)
}
