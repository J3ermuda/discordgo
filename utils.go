package discordgo

import (
	"sort"
	"strconv"
	"time"
)

// Contains checks if a slice of strings contains the string to search for
// haystack      : slice of strings to search in
// needle        : string to search for
func Contains(haystack []string, needle string) bool {
	sort.Strings(haystack)
	pos := sort.SearchStrings(haystack, needle)

	if pos == len(haystack) {
		return false
	}
	return haystack[pos] == needle
}

// ContainsIDObject checks if the haystack IDGettable contains the needle IDGettable
// haystack      : slice of IDGettables to search in
// needle        : IDGettable to search for
func ContainsIDObject(haystack []IDGettable, needle IDGettable) (contains bool) {
	if len(haystack) < 1 {
		return false
	}

	for _, item := range haystack {
		if item.GetID() == needle.GetID() {
			return true
		}
	}

	return false
}

// SnowflakeToTime converts a snowflake ID to a Time object
// snowflake      : the snowflake ID to convert
func SnowflakeToTime(snowflake string) (returnTime time.Time, err error) {
	n, err := strconv.ParseInt(snowflake, 10, 64)
	if err != nil {
		return
	}

	timestamp := ((n >> 22) + 1420070400000) * 1000000
	returnTime = time.Unix(0, timestamp).UTC()
	return
}

// IDWrapper is a struct used to make a bare ID string into an IDGettable object,
// this can then be used to pass it to a function that accepts any IDGettable object
type IDWrapper struct {
	ID string
}

// GetID returns the ID inside the IDWrapper
func (i IDWrapper) GetID() string {
	return i.ID
}

// CreatedAt returns the creation time in UTC of the object that the snowflake ID represents
func (i IDWrapper) CreatedAt() (creation time.Time, err error) {
	return SnowflakeToTime(i.ID)
}

// TimeSorter is a struct for allowing sorting on objects with a CreatedAt method
type TimeSorter []TimeSortable

func (t TimeSorter) Len() int {
	return len(t)
}

func (t TimeSorter) Less(i, j int) bool {
	iCreation, _ := t[i].CreatedAt()
	jCreation, _ := t[j].CreatedAt()
	return iCreation.After(jCreation)
}

func (t TimeSorter) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
