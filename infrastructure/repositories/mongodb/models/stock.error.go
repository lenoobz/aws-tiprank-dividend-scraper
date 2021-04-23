package models

// SectorModeError struct
type SectorModelError struct {
	Message string
}

func (e *SectorModelError) Error() string { return e.Message }

func (e *SectorModelError) As(target error) bool {
	_, ok := target.(*SectorModelError)
	return ok
}

// DividendModelError struct
type DividendModelError struct {
	Message string
}

func (e *DividendModelError) Error() string { return e.Message }

func (e *DividendModelError) As(target error) bool {
	_, ok := target.(*DividendModelError)
	return ok
}

// CountryModelError struct
type CountryModelError struct {
	Message string
}

func (e *CountryModelError) Error() string { return e.Message }

func (e *CountryModelError) As(target error) bool {
	_, ok := target.(*CountryModelError)
	return ok
}
