package xtype

import (
	"database/sql/driver"
	"fmt"
	"github.com/kurzgesagtz/xgo/xerror"
	"github.com/nyaruka/phonenumbers"
	"strings"
)

var PhoneNumberDefaultRegion = "TH"

func init() {
	PhoneNumberDefaultRegion = strings.ToUpper("TH")
}

type Phone struct {
	phone *phonenumbers.PhoneNumber
}

func NewPhone(phone string, code string) (*Phone, error) {
	if phone == "" {
		return nil, xerror.NewError(xerror.ErrCodeInvalidRequest, xerror.WithMessage("Phone number is empty"))
	}
	num, err := phonenumbers.Parse(phone, code)
	if err != nil {
		return nil, xerror.NewError(xerror.ErrCodeInvalidRequest, xerror.WithMessage(err.Error()))
	}
	if !phonenumbers.IsValidNumber(num) {
		return nil, xerror.NewError(xerror.ErrCodeInvalidRequest, xerror.WithMessage("Invalid format"))
	}
	return &Phone{
		phone: num,
	}, nil
}

func (p *Phone) GormDataType() string {
	return "varchar(128)"
}

func (p *Phone) Scan(value interface{}) (err error) {
	if p == nil {
		return nil
	}
	var phone *phonenumbers.PhoneNumber
	switch v := value.(type) {
	case []byte:
		phone, err = phonenumbers.Parse(string(v), PhoneNumberDefaultRegion)
		if err != nil {
			return err
		}
	case string:
		phone, err = phonenumbers.Parse(v, PhoneNumberDefaultRegion)
		if err != nil {
			return err
		}
	default:
		return nil
	}
	*p = Phone{
		phone: phone,
	}
	return nil
}

func (p *Phone) Value() (driver.Value, error) {
	if p == nil || p.phone == nil {
		return nil, nil
	}
	return phonenumbers.Format(p.phone, phonenumbers.E164), nil
}

func (p *Phone) String() string {
	return phonenumbers.Format(p.phone, phonenumbers.E164)
}

func (p *Phone) FormatE164() string {
	return phonenumbers.Format(p.phone, phonenumbers.E164)
}

func (p *Phone) FormatInternational() string {
	return phonenumbers.Format(p.phone, phonenumbers.INTERNATIONAL)
}

func (p *Phone) FormatNational() string {
	return phonenumbers.Format(p.phone, phonenumbers.NATIONAL)
}

func (p *Phone) FormatRFC3966() string {
	return phonenumbers.Format(p.phone, phonenumbers.RFC3966)
}

func (p *Phone) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", p.FormatE164())), nil
}

func (p *Phone) UnmarshalJSON(data []byte) error {
	raw := strings.TrimSuffix(strings.TrimPrefix(string(data), "\""), "\"")
	phone, err := phonenumbers.Parse(raw, PhoneNumberDefaultRegion)
	if err != nil {
		return err
	}
	*p = Phone{
		phone: phone,
	}
	return nil
}
