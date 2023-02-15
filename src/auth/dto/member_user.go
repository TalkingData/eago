package dto

import (
	"eago/auth/model"
	authpb "eago/auth/proto"
)

func CopyMemberUserGrpc(frm []*model.MemberUser, to *authpb.MemberUsers) {
	to.Users = make([]*authpb.MemberUsers_MemberUser, len(frm))
	for idx, u := range frm {
		to.Users[idx] = &authpb.MemberUsers_MemberUser{
			Id:       u.Id,
			Username: u.Username,
			IsOwner:  u.IsOwner,
			JoinedAt: u.JoinedAt.String(),
		}
	}
}
