package dto

import (
	"errors"
	"strings"
)

type CretaePostDTO struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (cretaePostDTO *CretaePostDTO) ValidateInput() error {
	errArr := []string{}
	if len(cretaePostDTO.Title) == 0 {
		errArr = append(errArr, "title is reqired")
	} else if len(cretaePostDTO.Title) > 300 {
		errArr = append(errArr, "title has maximum 300 characters long")
	}
	if len(cretaePostDTO.Content) == 0 {
		errArr = append(errArr, "content is reqired")
	}
	if len(errArr) > 0 {
		return errors.New(strings.Join(errArr, ", "))
	}
	return nil
}
