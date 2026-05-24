package services

import (
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/repositories"
)

func buildUserProfile(userID string) (*dto.UserProfileResponse, error) {
	user, err := repositories.GetUserProfileUser(userID)
	if err != nil {
		return nil, err
	}

	summary, err := repositories.GetUserProfileSummary(userID)
	if err != nil {
		return nil, err
	}

	batting, err := repositories.GetUserBattingStats(userID)
	if err != nil {
		return nil, err
	}

	bowling, err := repositories.GetUserBowlingStats(userID)
	if err != nil {
		return nil, err
	}

	fielding, err := repositories.GetUserFieldingStats(userID)
	if err != nil {
		return nil, err
	}

	recentMatches, err := repositories.GetUserRecentMatches(userID)
	if err != nil {
		return nil, err
	}

	return &dto.UserProfileResponse{
		User:          *user,
		Summary:       *summary,
		Batting:       *batting,
		Bowling:       *bowling,
		Fielding:      *fielding,
		RecentMatches: recentMatches,
	}, nil
}

func GetUserProfile(userID string) (*dto.UserProfileResponse, error) {
	return buildUserProfile(userID)
}

func GetPublicUserProfile(userID string) (*dto.UserProfileResponse, error) {
	return buildUserProfile(userID)
}

func UpdateUserProfile(userID string, req dto.UpdateProfileRequest) (*dto.UserProfileResponse, error) {
	_, err := repositories.UpdateUserProfile(userID, req)
	if err != nil {
		return nil, err
	}

	return GetUserProfile(userID)
}
