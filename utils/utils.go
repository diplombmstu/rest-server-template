package utils

import "errors"

func StringToError(string *string) error {
    if string != nil {
        return errors.New(*string)
    } else {
        return nil
    }
}
