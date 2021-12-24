package server

import (
	"log"
	"net/http"

	"github.com/LarsFox/pow-tcp/src/quotes"
)

const (
	hashcashHeader = "X-Hashcash"
)

//go:generate mockgen -destination=server_validator_mock_test.go -source=server.go -package=server
type validator interface {
	Challenge(ip string) (string, error)
	Validate(solution string) bool
}

// POWServer is a simple Proof-of-Work server implementation.
type POWServer struct {
	book      *quotes.Book
	validator validator
}

// NewPOWServer returns new Proof-of-Work server.
func NewPOWServer(book *quotes.Book, validator validator) *POWServer {
	return &POWServer{
		book:      book,
		validator: validator,
	}
}

// Listen listens on TCP network.
func (s *POWServer) Listen(addr string) {
	http.ListenAndServe(addr, http.HandlerFunc(s.handle))
}

// TODO: get proper IP.
func (s *POWServer) getIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		return xff
	}
	return "1.2.3.4"
}

func (s *POWServer) handle(w http.ResponseWriter, r *http.Request) {
	hashcash := r.Header.Get(hashcashHeader)
	if hashcash == "" {
		s.handleNewChallenge(w, r)
		return
	}

	val := s.validator.Validate(hashcash)
	if !val {
		s.sendErr(w, http.StatusForbidden)
		return
	}

	quote, err := s.book.Random()
	if err != nil {
		s.sendErr(w, http.StatusInternalServerError)
		return
	}

	log.Println("serving quote", string(quote))
	w.Write(quote)
}

func (s *POWServer) handleNewChallenge(w http.ResponseWriter, r *http.Request) {
	ip := s.getIP(r)
	challenge, err := s.validator.Challenge(ip)
	if err != nil {
		s.sendErr(w, http.StatusInternalServerError)
		return
	}

	log.Println("serving challenge", challenge)
	w.Header().Set(hashcashHeader, challenge)
	w.WriteHeader(http.StatusOK)
}

func (s *POWServer) sendErr(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}
