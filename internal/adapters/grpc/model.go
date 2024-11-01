package grpc

import (
	"net/mail"

	"github.com/fernandoocampo/players/internal/players"
	pb "github.com/fernandoocampo/players/pkg/pb/players"
)

func toNewPlayer(pbPlayer *pb.CreatePlayerRequest) players.NewPlayer {
	return players.NewPlayer{
		FirstName: pbPlayer.GetFirstname(),
		LastName:  pbPlayer.GetLastname(),
		Nickname:  pbPlayer.GetNickname(),
		Email:     mail.Address{Address: pbPlayer.GetEmail()},
		Password:  pbPlayer.GetPassword(),
		Country:   pbPlayer.GetCountry(),
	}
}

func toSearchCriteria(request *pb.SearchPlayersRequest) players.SearchCriteria {
	return players.SearchCriteria{
		Country: players.NewString(request.GetCountry()),
		Limit:   uint16(request.GetLimit()),
		Offset:  uint16(request.GetOffset()),
	}
}

func toUpdatePlayer(pbPlayer *pb.UpdatePlayerRequest, playerID *players.PlayerID) players.UpdatePlayer {
	var updatePlayer players.UpdatePlayer

	updatePlayer.ID = *playerID

	if pbPlayer.GetCountry() != "" {
		updatePlayer.Country = players.NewString(pbPlayer.GetCountry())
	}

	if pbPlayer.GetEmail() != "" {
		updatePlayer.Email = &mail.Address{Address: pbPlayer.GetEmail()}
	}

	if pbPlayer.GetFirstname() != "" {
		updatePlayer.FirstName = players.NewString(pbPlayer.GetFirstname())
	}

	if pbPlayer.GetLastname() != "" {
		updatePlayer.LastName = players.NewString(pbPlayer.GetLastname())
	}

	if pbPlayer.GetNickname() != "" {
		updatePlayer.Nickname = players.NewString(pbPlayer.GetNickname())
	}

	if pbPlayer.GetPassword() != "" {
		updatePlayer.Password = players.NewString(pbPlayer.GetPassword())
	}

	return updatePlayer
}

func toCreatePlayerReply(playerID *players.PlayerID) *pb.CreatePlayerReply {
	newReply := pb.CreatePlayerReply{
		PlayerId: playerID.String(),
	}

	return &newReply
}

func updatePlayerReplyOK() *pb.UpdatePlayerReply {
	return &pb.UpdatePlayerReply{
		Ok: true,
	}
}

func toSearchPlayerReply(result *players.SearchResult) *pb.SearchPlayersReply {
	newReply := pb.SearchPlayersReply{
		Total:       int64(result.Total),
		Limit:       uint32(result.Limit),
		Offset:      uint32(result.Offset),
		PlayerItems: make([]*pb.PlayerItem, 0, len(result.Items)),
	}

	for _, item := range result.Items {
		newReply.PlayerItems = append(newReply.PlayerItems, toPBPlayerItem(item))
	}

	return &newReply
}

func toPBPlayerItem(item players.PlayerItem) *pb.PlayerItem {
	newPBPlayerItem := pb.PlayerItem{
		Id:        item.ID.String(),
		Firstname: item.FirstName,
		Lastname:  item.LastName,
		Nickname:  item.Nickname,
		Country:   item.Country,
	}

	return &newPBPlayerItem
}

func deletePlayerReplyOK() *pb.DeletePlayerReply {
	return &pb.DeletePlayerReply{
		Ok: true,
	}
}
