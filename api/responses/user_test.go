package responses

import "testing"

func TestUserResponse_Render(t *testing.T) {
	resp := UserResponse{}
	err := resp.Render(nil, nil)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
}

func TestUsersResponse_Render(t *testing.T) {
	resp := UsersResponse{}
	err := resp.Render(nil, nil)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
}
