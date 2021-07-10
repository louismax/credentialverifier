package credentialverifier

import "testing"

func TestIsValid(t *testing.T) {
	t.Log(IsValid("610122198310134420", false))
	t.Log(IsValid("150000199703191282", true))
}

func TestGetInfo(t *testing.T) {
	t.Log(GetInfo("150000199703191282", false))
	as, err := GetInfo("830000199505245608", true)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", as)
}

func TestFakeId(t *testing.T) {
	t.Log(FakeId())
}

func TestFakeRequireId(t *testing.T) {
	t.Log(FakeRequireId(true, "台湾省", "199505", 0))
	t.Log(FakeRequireId(true, "香港特别行政区", "199505", 0))
}

func TestUpgradeId(t *testing.T) {
	t.Log(UpgradeId("610104620927690"))
}
