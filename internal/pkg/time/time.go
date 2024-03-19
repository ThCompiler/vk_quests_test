package time

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/pkg/errors"
	"strings"
	"time"
)

const Format = "02.01.2006 - 15:04:05"

type FormattedTime struct {
	time.Time
}

func (ft *FormattedTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")

	tm, err := time.Parse(Format, s)
	if err != nil {
		return err
	}

	*ft = FormattedTime{Time: tm}

	return nil
}

func (ft *FormattedTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(ft.String())
}

func (ft *FormattedTime) String() string {
	return ft.Format(Format)
}

func (ft *FormattedTime) Scan(src any) error {
	if tm, ok := src.(time.Time); ok {
		ft.Time = tm
		return nil
	}
	return errors.Errorf("invalid type of data for FormattedTime %v", src)
}

func (ft *FormattedTime) Value() (driver.Value, error) {
	return driver.Value(ft.Time), nil
}

func Parse(value string) (FormattedTime, error) {
	tm, err := time.Parse(Format, value)
	bt := FormattedTime{Time: tm}

	return bt, err
}

func MustParse(value string) FormattedTime {
	bt, err := Parse(value)
	if err != nil {
		panic(err)
	}

	return bt
}
