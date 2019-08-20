package models

import "strconv"

// ExperienceTable holds all xp repartition.
var ExperienceTable = []int64{
	0,
	525,
	1760,
	3781,
	7184,
	12186,
	19324,
	29377,
	43181,
	61693,
	85990,
	117506,
	157384,
	207736,
	269997,
	346462,
	439268,
	551295,
	685171,
	843709,
	1030734,
	1249629,
	1504995,
	1800847,
	2142652,
	2535122,
	2984677,
	3496798,
	4080655,
	4742836,
	5490247,
	6334393,
	7283446,
	8384398,
	9541110,
	10874351,
	12361842,
	14018289,
	15859432,
	17905634,
	20171471,
	22679999,
	25456123,
	28517857,
	31897771,
	35621447,
	39721017,
	44225461,
	49176560,
	54607467,
	60565335,
	67094245,
	74247659,
	82075627,
	90631041,
	99984974,
	110197515,
	121340161,
	133497202,
	146749362,
	161191120,
	176922628,
	194049893,
	212684946,
	232956711,
	255001620,
	278952403,
	304972236,
	333233648,
	363906163,
	397194041,
	433312945,
	472476370,
	514937180,
	560961898,
	610815862,
	664824416,
	723298169,
	786612664,
	855129128,
	929261318,
	1009443795,
	1096169525,
	1189918242,
	1291270350,
	1400795257,
	1519130326,
	1646943474,
	1784977296,
	1934009687,
	2094900291,
	2268549086,
	2455921256,
	2658074992,
	2876116901,
	3111280300,
	3364828162,
	3638186694,
	3932818530,
	4250334444,
	4250334444, // Level 100 no more xp needed.
}

// XpToNextLevel retrieves xp needed for the next level.
func XpToNextLevel(level int) int64 {
	if level <= 0 {
		level = 0
	} else if level > 100 {
		level = 100
	}
	return ExperienceTable[level]
}

// CurrentXp returns the xp gain for this current level.
func CurrentXp(xp int64, level int) int64 {
	return xp - XpToNextLevel(level-1)
}

// XpNeeded returns the remaining XP needed for the next level.
func XpNeeded(level int) int64 {
	if level <= 0 {
		level = 1
	} else if level > 100 {
		level = 100
	}
	return ExperienceTable[level] - ExperienceTable[level-1]
}

// PrettyPrint write string with space separators.
func PrettyPrint(n int64) string {
	in := strconv.FormatInt(n, 10)
	out := make([]byte, len(in)+(len(in)-2+int(in[0]/'0'))/3)
	if in[0] == '-' {
		in, out[0] = in[1:], '-'
	}

	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			return string(out)
		}
		if k++; k == 3 {
			j, k = j-1, 0
			out[j] = ' '
		}
	}
}
