package services

import "testing"

func TestPlus(t *testing.T) {
	result := plus(1, 2)
	if result != 3 {
		t.Errorf("1 + 2 = %d, want 3", result)
	}
}

type TestUserRepository struct {
}

func (r *TestUserRepository) Name() string {
	return "User"
}

func TestUser(t *testing.T) {
	u := &UserService{
		// repo: &StringUserRepository{},
		repo: &TestUserRepository{},
	}

	if u.Name() != "User" {
		t.Errorf("Name() = %s, want User", u.Name())
	}
}
