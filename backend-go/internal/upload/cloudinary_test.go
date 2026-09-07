package upload

import "testing"

func TestExtractPublicID(t *testing.T) {
	cases := []struct {
		name    string
		url     string
		want    string
		wantErr bool
	}{
		{
			name: "folder + versioned URL (typical upload)",
			url:  "https://res.cloudinary.com/demo/image/upload/v1699999999/uploads/abcd1234.png",
			want: "uploads/abcd1234",
		},
		{
			name: "nested folder (voice messages)",
			url:  "https://res.cloudinary.com/demo/video/upload/v1699999999/uploads/voice-messages/xyz789.webm",
			want: "uploads/voice-messages/xyz789",
		},
		{
			name: "no version segment",
			url:  "https://res.cloudinary.com/demo/image/upload/uploads/abcd1234.jpg",
			want: "uploads/abcd1234",
		},
		{
			name:    "not a Cloudinary upload URL",
			url:     "https://example.com/some/image.png",
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ExtractPublicID(tc.url)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got public_id %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}
