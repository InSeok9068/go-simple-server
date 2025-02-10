package services

type UserRepository interface {
	Name() string
}

type UserService struct {
	repo UserRepository
}

type StringUserRepository struct {
}

func (r *StringUserRepository) Name() string {
	return "User!"
}

func (u *UserService) Name() string {
	return u.repo.Name()
}

func plus(a int, b int) int {
	return a + b
}
