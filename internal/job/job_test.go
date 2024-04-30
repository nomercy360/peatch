package job

import (
	"github.com/peatch-io/peatch/internal/db"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateImgURLWithCompleteUserData(t *testing.T) {
	user := db.User{
		FirstName: StringPointer("John"),
		LastName:  StringPointer("Doe"),
		Title:     StringPointer("Product"),
		AvatarURL: StringPointer("users/149/KO7uaU43.svg"),
		Badges: []db.Badge{
			{Text: "Mentor", Color: "17BEBB", Icon: "e8d3"},
			{Text: "Founder", Color: "FF8C42", Icon: "eb39"},
		},
	}

	expectedURL := "https://peatch-image-preview.vercel.app/api/image?title=John Doe&subtitle=Product&avatar=https://assets.peatch.io/users/149/KO7uaU43.svg&tags=Mentor,17BEBB,e8d3;Founder,FF8C42,eb39;"
	url, err := createImgURL("https://peatch-image-preview.vercel.app", &user)

	assert.NoError(t, err)
	assert.Equal(t, expectedURL, url)
}

func TestCreateImgURLWithIncompleteUserData(t *testing.T) {
	user := db.User{
		FirstName: StringPointer("John"),
		LastName:  StringPointer("Doe"),
		// Title and AvatarURL are missing
	}

	url, err := createImgURL("https://peatch-image-preview.vercel.app", &user)

	assert.Error(t, err)
	assert.Equal(t, "", url)
}

func StringPointer(s string) *string {
	return &s
}
