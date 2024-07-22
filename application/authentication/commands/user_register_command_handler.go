package commands

import (
	"errors"
	"fmt"

	"github.com/beka-birhanu/finance-go/application/authentication/common"
	"github.com/beka-birhanu/finance-go/application/common/interfaces/jwt"
	"github.com/beka-birhanu/finance-go/application/common/interfaces/persistance"
	hash "github.com/beka-birhanu/finance-go/domain/common/authentication"
	"github.com/beka-birhanu/finance-go/domain/domain_errors"
	"github.com/beka-birhanu/finance-go/domain/models"
)

type UserRegisterCommandHandler struct {
	userRepository persistance.IUserRepository
	jwtService     jwt.IJwtService
	hashService    hash.IHashService
}

func NewRegisterCommandHandler(repository persistance.IUserRepository, jwtService jwt.IJwtService, hashService hash.IHashService) *UserRegisterCommandHandler {
	return &UserRegisterCommandHandler{userRepository: repository, jwtService: jwtService, hashService: hashService}
}

func (h *UserRegisterCommandHandler) Handle(command *UserRegisterCommand) (*common.AuthResult, error) {
	user, err := fromRegisterCommand(command, h.hashService)
	if err != nil {
		return nil, err
	}

	err = h.userRepository.CreateUser(user)
	if errors.Is(err, domain_errors.ErrUsernameConflict) {
		return nil, err
	} else if err != nil {
		return nil, fmt.Errorf("server error")
	}

	token, err := h.jwtService.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("server error")
	}

	return common.NewAuthResult(user.ID(), user.Username(), token), nil
}

func fromRegisterCommand(command *UserRegisterCommand, hashService hash.IHashService) (*models.User, error) {
	return models.NewUser(command.Username, command.Password, hashService)
}
