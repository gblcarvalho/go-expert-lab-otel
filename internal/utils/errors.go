package utils

import "errors"


var ErrInvalidCEP = errors.New("Invalid CEP")
var ErrCEPNotFound = errors.New("CEP Not Found")
var ErrWeather = errors.New("Error Weather API")
