package accord

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/mistsys/accord/protocol"
)

func ToHostCA(public CAPublic) *protocol.HostCA {
	validFrom, _ := ptypes.TimestampProto(public.ValidFrom)
	validUntil, _ := ptypes.TimestampProto(public.ValidUntil)
	return &protocol.HostCA{
		Id:         uint64(public.Id),
		PublicKey:  public.PublicKey,
		ValidFrom:  validFrom,
		ValidUntil: validUntil,
	}
}

func ToUserCA(public CAPublic) *protocol.UserCA {
	validFrom, _ := ptypes.TimestampProto(public.ValidFrom)
	validUntil, _ := ptypes.TimestampProto(public.ValidUntil)
	return &protocol.UserCA{
		Id:         uint64(public.Id),
		PublicKey:  public.PublicKey,
		ValidFrom:  validFrom,
		ValidUntil: validUntil,
	}
}
