package shortuuid

import "testing"

func TestShortUUID(t *testing.T) {
	uuid := ShortUUID(6)
	if len(uuid) != 6 {
		t.Error("UUID length is not 6")
	}
	uuid = ShortUUID(7)
	if len(uuid) != 7 {
		t.Error("UUID length is not 7")
	}
	uuid = ShortUUID(8)
	if len(uuid) != 8 {
		t.Error("UUID length is not 8")
	}
	uuid = ShortUUID(9)
	if len(uuid) != 9 {
		t.Error("UUID length is not 9")
	}
	uuid = ShortUUID(10)
	if len(uuid) != 10 {
		t.Error("UUID length is not 10")
	}
	uuid = ShortUUID(11)
	if len(uuid) != 11 {
		t.Error("UUID length is not 11")
	}
	uuid = ShortUUID(12)
	if len(uuid) != 12 {
		t.Error("UUID length is not 12")
	}
	uuid = ShortUUID(13)
	if len(uuid) != 13 {
		t.Error("UUID length is not 13")
	}
}
