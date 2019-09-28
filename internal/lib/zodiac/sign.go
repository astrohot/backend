package zodiac

import "time"

// Sign ...
type Sign uint8

// Zodiac signs
const (
	Aries Sign = iota + 1
	Taurus
	Gemini
	Cancer
	Leo
	Virgo
	Libra
	Scorpio
	Sagittarius
	Capricorn
	Aquarius
	Pisces
)

var signMap = map[Sign]string{
	Capricorn:   "Capricórnio",
	Aquarius:    "Aquário",
	Pisces:      "Peixes",
	Taurus:      "Touro",
	Gemini:      "Gêmeos",
	Cancer:      "Câncer",
	Leo:         "Leão",
	Virgo:       "Virgem",
	Libra:       "Libra",
	Scorpio:     "Escorpião",
	Sagittarius: "Sagitário",
}

func (s Sign) String() string {
	return signMap[s]
}

// GetSign returns a sign accordingly with a birth date (timestamp).
func GetSign(birth time.Time) (s Sign) {
	day, month := birth.Day(), birth.Month()

	switch month {
	case time.March:
		if day < 21 {
			s = Pisces
		} else {
			s = Aries
		}

	case time.April:
		if day < 20 {
			s = Aries
		} else {
			s = Taurus
		}

	case time.May:
		if day < 21 {
			s = Taurus
		} else {
			s = Gemini
		}

	case time.June:
		if day < 22 {
			s = Gemini
		} else {
			s = Cancer
		}

	case time.July:
		if day < 23 {
			s = Cancer
		} else {
			s = Leo
		}

	case time.August:
		if day < 23 {
			s = Leo
		} else {
			s = Virgo
		}

	case time.September:
		if day < 23 {
			s = Virgo
		} else {
			s = Libra
		}

	case time.October:
		if day < 23 {
			s = Libra
		} else {
			s = Scorpio
		}

	case time.November:
		if day < 22 {
			s = Scorpio
		} else {
			s = Sagittarius
		}

	case time.December:
		if day < 22 {
			s = Sagittarius
		} else {
			s = Capricorn
		}

	case time.January:
		if day < 20 {
			s = Capricorn
		} else {
			s = Aquarius
		}

	case time.February:
		if day < 19 {
			s = Aquarius
		} else {
			s = Pisces
		}
	}

	return
}
