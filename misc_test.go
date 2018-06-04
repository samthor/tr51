package tr51

import (
	"reflect"
	"testing"
)

func TestUnqualify(t *testing.T) {
	qualified := "👨‍⚕️"
	unqualified := "👨‍⚕"

	if qualified == unqualified {
		t.Errorf("qualified should not equal unqualified")
	}
	if unqualified != Unqualify(qualified) {
		t.Errorf("expected fe0f removed")
	}
}

func TestOrigins(t *testing.T) {
	expected := []string{"Origin_ZDings", "Origin_JCarrier"}
	if actual := Origins("z j"); !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected z j origins to equal=%v, was=%v", expected, actual)
	}
	if len(Origins("q p #")) != 0 {
		t.Errorf("expected zero origins for bogus data")
	}
}
