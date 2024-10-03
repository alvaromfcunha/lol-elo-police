package usecase

import (
	"errors"
	"os"

	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity/enum"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/repository"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/service"
)

type PolicePatrol struct {
	PlayerRepository     repository.IPlayerRepository
	RankedInfoRepository repository.IRankedInfoRepository
	LolService           service.ILolService
	WhatsappService      service.IWhatsappService
	TemplateService      service.ITemplateService
}

func (u PolicePatrol) Execute() error {
	players, err := u.PlayerRepository.GetAll()
	if err != nil {
		return errors.New("cannot get all players on database")
	}

	for _, player := range players {
		leagues, err := u.LolService.GetLeaguesBySummonerId(player.SummonerId)
		if err != nil {
			return errors.New("cannot player leagues by summoner id")
		}

		var soloQueueInfo *entity.RankedInfo
		if rankedInfo, err := u.RankedInfoRepository.GetByPlayerAndQueueType(player, enum.Solo); err == nil {
			soloQueueInfo = &rankedInfo
		}

		var flexQueueInfo *entity.RankedInfo
		if rankedInfo, err := u.RankedInfoRepository.GetByPlayerAndQueueType(player, enum.Flex); err == nil {
			flexQueueInfo = &rankedInfo
		}

		for _, league := range leagues {
			switch league.QueueType {
			case enum.Solo:
				if soloQueueInfo != nil {
					return u.checkRankedQueueUpdate(player, *soloQueueInfo, league)
				} else {
					return u.notifyNewRankedQueueEntry(player, enum.Solo, league)
				}
			case enum.Flex:
				if flexQueueInfo != nil {
					return u.checkRankedQueueUpdate(player, *flexQueueInfo, league)
				} else {
					return u.notifyNewRankedQueueEntry(player, enum.Flex, league)
				}
			}
		}
	}

	return nil
}

func (u PolicePatrol) checkRankedQueueUpdate(
	player entity.Player,
	rankedInfo entity.RankedInfo,
	leagueEntry service.LeagueEntry,
) error {
	isQueueUpdate := leagueEntry.LeaguePoints != rankedInfo.LeaguePoints ||
		leagueEntry.Wins != rankedInfo.Wins ||
		leagueEntry.Losses != rankedInfo.Losses ||
		leagueEntry.Rank != rankedInfo.Rank ||
		leagueEntry.Tier != rankedInfo.Tier

	if isQueueUpdate {
		queueUpdateData := service.QueueUpdateData{
			Player:         player,
			RankedType:     getReadableQueueType(rankedInfo.QueueType),
			NewLeagueEntry: leagueEntry,
			OldRankedInfo:  rankedInfo,
		}

		message, err := u.TemplateService.ExecuteQueueUpdateMessageTemplate(queueUpdateData)
		if err != nil {
			return errors.New("cannot execute queue update template")
		}

		whatsappGroupUser := os.Getenv("WPP_GROUP_USER")
		_, err = u.WhatsappService.SendMessageToGroup(message, whatsappGroupUser)
		if err != nil {
			return errors.New("cannot send message to whatsapp group")
		}

		rankedInfo.LeaguePoints = leagueEntry.LeaguePoints
		rankedInfo.Wins = leagueEntry.Wins
		rankedInfo.Losses = leagueEntry.Losses
		rankedInfo.Rank = leagueEntry.Rank
		rankedInfo.Tier = leagueEntry.Tier

		err = u.RankedInfoRepository.Update(rankedInfo)
		if err != nil {
			return errors.New("cannot update ranked info on database")
		}
	}

	return nil
}

func (u PolicePatrol) notifyNewRankedQueueEntry(
	player entity.Player,
	queueType enum.QueueType,
	leagueEntry service.LeagueEntry,
) error {
	queueNewEntryData := service.QueueNewEntryData{
		Player:      player,
		RankedType:  getReadableQueueType(queueType),
		LeagueEntry: leagueEntry,
	}

	message, err := u.TemplateService.ExecuteQueueNewEntryMessageTemplate(queueNewEntryData)
	if err != nil {
		return errors.New("cannot execute queue update template")
	}

	whatsappGroupUser := os.Getenv("WPP_GROUP_USER")
	_, err = u.WhatsappService.SendMessageToGroup(message, whatsappGroupUser)
	if err != nil {
		return errors.New("cannot send message to whatsapp group")
	}

	playerRankedInfo := entity.NewRankedInfo(
		player,
		queueType,
		leagueEntry.Tier,
		leagueEntry.Rank,
		leagueEntry.LeaguePoints,
		leagueEntry.Wins,
		leagueEntry.Losses,
	)

	err = u.RankedInfoRepository.Create(playerRankedInfo, player)
	if err != nil {
		return errors.New("cannot create new ranked info entry on database")
	}

	return nil
}

func getReadableQueueType(queueType enum.QueueType) string {
	switch queueType {
	case enum.Solo:
		return "Solo Queue"
	case enum.Flex:
		return "Flex"
	default:
		return "-"
	}
}
