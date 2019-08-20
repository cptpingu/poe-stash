package models

import "testing"

// assertEqual checks that two values are equals
func assertEqual(t *testing.T, expected, got int64) {
	if expected != got {
		t.Errorf("Expected %d but got %d", expected, got)
	}
}

// assertStringEqual checks that two strings are equals
func assertStringEqual(t *testing.T, expected, got string) {
	if expected != got {
		t.Errorf("Expected %q but got %q", expected, got)
	}
}

// TestXpToNextLevel tests that xp for next level is correct.
func TestXpToNextLevel(t *testing.T) {
	assertEqual(t, 0, XpToNextLevel(0))
	assertEqual(t, 525, XpToNextLevel(1))
	assertEqual(t, 1760, XpToNextLevel(2))
	assertEqual(t, 161191120, XpToNextLevel(60))
	assertEqual(t, 4250334444, XpToNextLevel(99))
	assertEqual(t, 4250334444, XpToNextLevel(100))
	assertEqual(t, 4250334444, XpToNextLevel(101))
}

// TestPrettyPrint tests that pretty printing works.
func TestPrettyPrint(t *testing.T) {
	assertStringEqual(t, "0", PrettyPrint(0))

	assertStringEqual(t, "1", PrettyPrint(1))
	assertStringEqual(t, "11", PrettyPrint(11))
	assertStringEqual(t, "111", PrettyPrint(111))
	assertStringEqual(t, "2 111", PrettyPrint(2111))
	assertStringEqual(t, "22 111", PrettyPrint(22111))
	assertStringEqual(t, "222 111", PrettyPrint(222111))
	assertStringEqual(t, "3 222 111", PrettyPrint(3222111))
	assertStringEqual(t, "33 222 111", PrettyPrint(33222111))
	assertStringEqual(t, "333 222 111", PrettyPrint(333222111))
	assertStringEqual(t, "4 333 222 111", PrettyPrint(4333222111))
	assertStringEqual(t, "44 333 222 111", PrettyPrint(44333222111))
	assertStringEqual(t, "444 333 222 111", PrettyPrint(444333222111))
	assertStringEqual(t, "5 444 333 222 111", PrettyPrint(5444333222111))
	assertStringEqual(t, "55 444 333 222 111", PrettyPrint(55444333222111))
	assertStringEqual(t, "555 444 333 222 111", PrettyPrint(555444333222111))
	assertStringEqual(t, "6 555 444 333 222 111", PrettyPrint(6555444333222111))
	assertStringEqual(t, "66 555 444 333 222 111", PrettyPrint(66555444333222111))
	assertStringEqual(t, "666 555 444 333 222 111", PrettyPrint(666555444333222111))

	assertStringEqual(t, "-1", PrettyPrint(-1))
	assertStringEqual(t, "-11", PrettyPrint(-11))
	assertStringEqual(t, "-111", PrettyPrint(-111))
	assertStringEqual(t, "-2 111", PrettyPrint(-2111))
	assertStringEqual(t, "-22 111", PrettyPrint(-22111))
	assertStringEqual(t, "-222 111", PrettyPrint(-222111))
	assertStringEqual(t, "-3 222 111", PrettyPrint(-3222111))
	assertStringEqual(t, "-33 222 111", PrettyPrint(-33222111))
	assertStringEqual(t, "-333 222 111", PrettyPrint(-333222111))
	assertStringEqual(t, "-4 333 222 111", PrettyPrint(-4333222111))
	assertStringEqual(t, "-44 333 222 111", PrettyPrint(-44333222111))
	assertStringEqual(t, "-444 333 222 111", PrettyPrint(-444333222111))
	assertStringEqual(t, "-5 444 333 222 111", PrettyPrint(-5444333222111))
	assertStringEqual(t, "-55 444 333 222 111", PrettyPrint(-55444333222111))
	assertStringEqual(t, "-555 444 333 222 111", PrettyPrint(-555444333222111))
	assertStringEqual(t, "-6 555 444 333 222 111", PrettyPrint(-6555444333222111))
	assertStringEqual(t, "-66 555 444 333 222 111", PrettyPrint(-66555444333222111))
	assertStringEqual(t, "-666 555 444 333 222 111", PrettyPrint(-666555444333222111))
}

// TestCurrentXp tests that we return the xp we have.
func TestCurrentXp(t *testing.T) {
	assertEqual(t, 0, CurrentXp(0, 0))
	assertEqual(t, 55, CurrentXp(55, 0))
	assertEqual(t, 0, CurrentXp(0, 1))
	assertEqual(t, 55, CurrentXp(55, 1))
	assertEqual(t, 0, CurrentXp(525, 2))
	assertEqual(t, 0, CurrentXp(7184, 5))
	assertEqual(t, 475, CurrentXp(1000, 2))
	assertEqual(t, 0, CurrentXp(3932818530, 99))
	assertEqual(t, 10, CurrentXp(3932818540, 99))
	assertEqual(t, 0, CurrentXp(4250334444, 100))
	assertEqual(t, 0, CurrentXp(4250334444, 101))
	assertEqual(t, 1, CurrentXp(4250334445, 100))
}

// TestXpNeeded tests that we returns the right amount of xp needed for next level.
func TestXpNeeded(t *testing.T) {
	assertEqual(t, 525, XpNeeded(0))
	assertEqual(t, 525, XpNeeded(1))
	assertEqual(t, 1235, XpNeeded(2))
	assertEqual(t, 14441758, XpNeeded(60))
	assertEqual(t, 317515914, XpNeeded(99))
	assertEqual(t, 0, XpNeeded(100))
	assertEqual(t, 0, XpNeeded(101))
}
