package xtype

import (
	"context"
	"fmt"
	"reflect"

	"github.com/kurzgesagtz/xgo/xerror"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/schema"
)

type HashString struct {
	hash bool
	cost int
	str  string
}

func NewHashString(str string) HashString {
	cost, err := bcrypt.Cost([]byte(str))
	if err != nil {
		return HashString{
			str: str,
		}
	}
	return HashString{
		str:  str,
		hash: true,
		cost: cost,
	}
}

func (hs *HashString) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) (err error) {
	if hs == nil {
		return nil
	}

	switch value := dbValue.(type) {
	case []byte:
		*hs = HashString{
			hash: true,
			str:  string(value),
		}
	case string:
		*hs = HashString{
			hash: true,
			str:  value,
		}
	case nil:
		*hs = HashString{
			hash: false,
			str:  "",
		}
	default:
		return xerror.NewError(xerror.ErrCodeInternalError, xerror.WithMessage(fmt.Sprintf("unsupported data %#v", dbValue)))
	}
	if hs.str == "" {
		return nil
	}
	hs.cost, err = bcrypt.Cost([]byte(hs.str))
	return err
}

func (hs *HashString) Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error) {
	if hs == nil {
		return nil, nil
	}

	if hs.hash {
		return hs.str, nil
	}
	if hs.str == "" {
		return nil, nil
	}
	bytes, _ := bcrypt.GenerateFromPassword([]byte(hs.str), 10)
	return string(bytes), nil
}

func (hs *HashString) Equal(pwd string) bool {
	if hs == nil {
		return false
	}

	if hs.hash {
		err := bcrypt.CompareHashAndPassword([]byte(hs.str), []byte(pwd))
		return err == nil
	}
	return hs.str == pwd
}

func (hs *HashString) String() string {
	if hs == nil {
		return ""
	}
	return hs.str
}
