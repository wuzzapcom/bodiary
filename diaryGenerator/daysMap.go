package diaryGenerator

//DaysMap map from day of week to number of string in table
type DaysMap map[int][]string

func (p *DaysMap) contains(day int) bool {

	pickedDays := *p

	_, ok := pickedDays[day]

	*p = pickedDays

	return ok

}

func (p *DaysMap) addDay(day int) {

	pickedDays := *p

	if !p.contains(day) {

		pickedDays[day] = make([]string, 2)

	}

	*p = pickedDays

}

func (p *DaysMap) addMessage(day int, numberOfWeek int, message string) {

	pickedDays := *p

	array := pickedDays[day]

	array[numberOfWeek] = message

	*p = pickedDays

}
