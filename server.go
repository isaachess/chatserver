package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type Server struct {
	host   string
	port   int
	logger Logger

	chatters_m sync.RWMutex
	chatters   map[string]net.Conn
}

func NewServer(host string, port int, l Logger) *Server {
	return &Server{
		host:     host,
		port:     port,
		logger:   l,
		chatters: make(map[string]net.Conn),
	}
}

func (s *Server) Serve() error {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.host, s.port))
	if err != nil {
		return err
	}

	log.Println("Successfully listening on", ln.Addr())

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		go func() {
			if err := s.handle(conn); err != nil {
				log.Println(err)
			}
		}()
	}
}

func (s *Server) handle(conn net.Conn) error {
	defer conn.Close()

	if err := s.writeMsg("Please enter your name: ", conn); err != nil {
		return err
	}

	sc := bufio.NewScanner(conn)

	// scan once for name response
	if !sc.Scan() {
		return sc.Err()
	}

	name := sc.Text()
	if s.haveChatter(name) {
		return s.writeMsg("Chatter with that name already exists.\n", conn)
	}

	s.addChatter(name, conn)
	// defer cleanups
	defer s.removeChatter(name)
	defer s.sendLogout(name)

	// tell new chatter welcome
	if err := s.writeMsg(s.welcomeMsg(name), conn); err != nil {
		return err
	}

	// alert other chatters a new friend has arrived!
	if err := s.broadcastMsg(name, loginMsg(name)); err != nil {
		return err
	}

	// scan for all new messages
	for sc.Scan() {
		msg := sc.Text()
		if err := s.broadcastMsg(name, chatMsg(name, msg)); err != nil {
			return err
		}
	}

	return sc.Err()
}

func (s *Server) addChatter(name string, conn net.Conn) {
	s.chatters_m.Lock()
	defer s.chatters_m.Unlock()
	s.chatters[name] = conn
}

func (s *Server) removeChatter(name string) {
	s.chatters_m.Lock()
	defer s.chatters_m.Unlock()
	delete(s.chatters, name)
}

func (s *Server) numChatters() int {
	s.chatters_m.RLock()
	defer s.chatters_m.RUnlock()
	return len(s.chatters)
}

func (s *Server) haveChatter(name string) bool {
	s.chatters_m.RLock()
	defer s.chatters_m.RUnlock()
	_, ok := s.chatters[name]
	return ok
}

func (s *Server) writeMsg(msg string, conn net.Conn) error {
	_, err := conn.Write([]byte(msg))
	return err
}

func (s *Server) broadcastMsg(fromName, msg string) error {
	s.chatters_m.RLock()
	defer s.chatters_m.RUnlock()
	toSend := tsMsg(msg)
	for name, conn := range s.chatters {
		if name == fromName {
			// don't want to send to yourself, that would be weird
			continue
		}
		if err := s.writeMsg(toSend, conn); err != nil {
			return err
		}
	}

	// log any message we broadcast to everyone
	s.logger.Log(toSend)

	return nil
}

func (s *Server) sendLogout(name string) {
	s.broadcastMsg(name, fmt.Sprintf("%s has logged out\n", name))
}

func (s *Server) welcomeMsg(name string) string {
	var others string
	total := s.numChatters()
	if total > 1 {
		others = fmt.Sprintf("There are %d total people here. Happy chatting!", total)
	} else {
		others = "You are alone here. We hope that changes soon."
	}

	return fmt.Sprintf("Welcome %s! %s\n", name, others)
}

func chatMsg(name, msg string) string {
	return fmt.Sprintf("%s: %s\n", name, msg)
}

func loginMsg(name string) string {
	return fmt.Sprintf("%s has logged in!\n", name)
}

func tsMsg(msg string) string {
	return fmt.Sprintf("(%s) %s", time.Now().Format(time.Kitchen), msg)
}
