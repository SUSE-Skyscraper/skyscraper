package responses

import "testing"

func TestTagResponse_Render(t *testing.T) {
	resp := TagResponse{}
	err := resp.Render(nil, nil)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
}

func TestTagsResponse_Render(t *testing.T) {
	resp := TagsResponse{}
	err := resp.Render(nil, nil)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
}
