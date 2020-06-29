package service

import (
	"github.com/Alex27Khalupka/Go-course-task/pkg/apiserver"
	"log"
	"testing"
)

func TestService_GetGroups(t *testing.T){
	var s *apiserver.APIServer

	if err := s.Store.Open(); err != nil {
		log.Fatal(err)
		return
	}

	answer := GetGroups(s.Store.GetDB())

	s.Store.GetDB().QueryRow()
}
