package entities

import (
	"fmt"
	"reflect"

	"github.com/tinylib/msgp/msgp"
)

var (
	// Mapping from string representation of a TickDataProvider concrete type to the concrete type itself.
	// The string representation is used as type discriminator when encoding/decoding TickDataProvider.
	tickDataProviderImplMap = map[string]reflect.Type{}
)

// RegisterTickDataProviderImpl registers the concrete types of an TickDataProvider.
// This function is not thread-safe and should be only call in init().
func RegisterTickDataProviderImpl(provider TickDataProvider) {
	if _, ok := provider.(msgp.Encodable); !ok {
		panic("expected provider to implement msgp.Encodable")
	}
	if _, ok := provider.(msgp.Decodable); !ok {
		panic("expected provider to implement msgp.Decodable")
	}
	if _, ok := provider.(msgp.Marshaler); !ok {
		panic("expected provider to implement msgp.Marshaler")
	}
	if _, ok := provider.(msgp.Unmarshaler); !ok {
		panic("expected provider to implement msgp.Unmarshaler")
	}
	if _, ok := provider.(msgp.Sizer); !ok {
		panic("expected provider to implement msgp.Sizer")
	}
	typ := reflect.ValueOf(provider).Elem().Type()
	tickDataProviderImplMap[typ.String()] = typ
}

// TickDataProviderWrapper is a wrapper of TickDataProvider and is implemented msgp.Encodable, msgp.Decodable, msgp.Marshaler, msgp.Unmarshaler, and msgp.Sizer
type TickDataProviderWrapper struct {
	TickDataProvider
}

func NewTickDataProviderWrapper(provider TickDataProvider) *TickDataProviderWrapper {
	if provider == nil {
		return nil
	}
	return &TickDataProviderWrapper{provider}
}

// Get the inner TickDataProvider, return nil if TickDataProviderWrapper is nil
func (p *TickDataProviderWrapper) Get() TickDataProvider {
	if p == nil {
		return nil
	}
	return p.TickDataProvider
}

// EncodeMsg implements msgp.Encodable
func (p *TickDataProviderWrapper) EncodeMsg(en *msgp.Writer) (err error) {
	typ := reflect.ValueOf(p.TickDataProvider).Elem().Type()
	err = en.WriteString(typ.String())
	if err != nil {
		return
	}

	if _, ok := tickDataProviderImplMap[typ.String()]; !ok {
		err = fmt.Errorf("unregistered type %s", typ.String())
		return
	}

	err = p.TickDataProvider.(msgp.Encodable).EncodeMsg(en)
	if err != nil {
		err = msgp.WrapError(err, "TickDataProvider")
		return
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (p *TickDataProviderWrapper) DecodeMsg(dc *msgp.Reader) (err error) {
	var typStr string
	typStr, err = dc.ReadString()
	if err != nil {
		return
	}

	typ, ok := tickDataProviderImplMap[typStr]
	if !ok {
		err = fmt.Errorf("unregistered type %s", typStr)
		return
	}

	providerVal := reflect.New(typ)
	err = providerVal.Interface().(msgp.Decodable).DecodeMsg(dc)
	if err != nil {
		err = msgp.WrapError(err, "TickDataProviderMsgp")
		return
	}
	p.TickDataProvider = providerVal.Interface().(TickDataProvider)
	return
}

// MarshalMsg implements msgp.Marshaler
func (p *TickDataProviderWrapper) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, p.Msgsize())
	typ := reflect.ValueOf(p.TickDataProvider).Elem().Type()
	o = msgp.AppendString(o, typ.String())

	if _, ok := tickDataProviderImplMap[typ.String()]; !ok {
		err = fmt.Errorf("unregistered type %s", typ.String())
		return
	}

	o, err = p.TickDataProvider.(msgp.Marshaler).MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "TickDataProvider")
		return
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (p *TickDataProviderWrapper) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var typStr string
	typStr, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		return
	}

	typ, ok := tickDataProviderImplMap[typStr]
	if !ok {
		err = fmt.Errorf("unregistered type %s", typStr)
		return
	}

	providerVal := reflect.New(typ)
	bts, err = providerVal.Interface().(msgp.Unmarshaler).UnmarshalMsg(bts)
	if err != nil {
		err = msgp.WrapError(err, "TickDataProvider")
		return
	}
	p.TickDataProvider = providerVal.Interface().(TickDataProvider)
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (p *TickDataProviderWrapper) Msgsize() int {
	typ := reflect.ValueOf(p.TickDataProvider).Elem().Type()
	return msgp.StringPrefixSize + len(typ.String()) + p.TickDataProvider.(msgp.Sizer).Msgsize()
}
