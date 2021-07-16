package usecase

import (
	"log"
	"time"

	"github.com/ardafirdausr/discuss-server/internal"
	"github.com/ardafirdausr/discuss-server/internal/entity"
)

type DiscussionUsecase struct {
	discussionRepo internal.DiscussionRepository
}

func NewDiscussionUsecase(discussionRepo internal.DiscussionRepository) *DiscussionUsecase {
	return &DiscussionUsecase{discussionRepo: discussionRepo}
}

func (du DiscussionUsecase) GetAllUserDiscussions(userID interface{}) ([]*entity.Discussion, error) {
	discussions, err := du.discussionRepo.GetDiscussionsByUserID(userID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return discussions, nil
}

func (du DiscussionUsecase) GetDiscussionByID(discussionID interface{}) (*entity.Discussion, error) {
	discussion, err := du.discussionRepo.GetDiscussionsByID(discussionID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return discussion, nil
}

func (du DiscussionUsecase) GetDiscussionByCode(code string) (*entity.Discussion, error) {

	discussion, err := du.discussionRepo.GetDiscussionsByCode(code)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return discussion, nil
}

func (du DiscussionUsecase) Create(param entity.CreateDiscussionParam) (*entity.Discussion, error) {
	discussion, err := du.GetDiscussionByCode(param.Code)
	_, isErrNF := err.(entity.ErrNotFound)
	if err != nil && !isErrNF {
		return nil, err
	}

	if discussion != nil {
		err = entity.ErrInvalidData{
			Message: "Discussion code already exists",
			Err:     nil,
		}
		return nil, err
	}

	param.CreatedAt = time.Now()
	param.UpdatedAt = time.Now()
	param.Members = make([]interface{}, 0)
	param.Members = append(param.Members, param.CreatorID)

	discussion, err = du.discussionRepo.Create(param)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return discussion, nil
}

func (du DiscussionUsecase) JoinDiscussion(discussionID, userID interface{}) (*entity.Discussion, error) {
	err := du.discussionRepo.AddMember(discussionID, userID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return du.GetDiscussionByID(discussionID)
}

func (du DiscussionUsecase) LeaveDiscussion(discussionID, userID interface{}) error {
	err := du.discussionRepo.RemoveMember(discussionID, userID)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (du DiscussionUsecase) Update(discussionID interface{}, param entity.UpdateDiscussionParam) error {
	discussion, err := du.GetDiscussionByID(discussionID)
	if _, isErrNF := err.(entity.ErrNotFound); isErrNF {
		err = entity.ErrNotFound{
			Message: "Discussion not found",
			Err:     nil,
		}
		return err
	}

	if err != nil {
		return err
	}

	existDiscussion, err := du.GetDiscussionByCode(param.Code)
	_, isErrNF := err.(entity.ErrNotFound)
	if err != nil && !isErrNF {
		return err
	}

	if existDiscussion != nil && existDiscussion.Code == param.Code && existDiscussion.ID != discussion.ID {
		err = entity.ErrInvalidData{
			Message: "Discussion code already exists",
			Err:     nil,
		}
		return err
	}

	param.UpdatedAt = time.Now()
	err = du.discussionRepo.UpdateByID(discussionID, param)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (du DiscussionUsecase) UpdatePassword(discussionID interface{}, param entity.UpdateDiscussionPassword) error {
	hashedPassword := hashString(param.Password)
	err := du.discussionRepo.UpdatePasswordByID(discussionID, hashedPassword)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (du DiscussionUsecase) UpdatePhoto(discussionID interface{}, param string) error {
	err := du.discussionRepo.UpdatePhotoByID(discussionID, param)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (du DiscussionUsecase) Delete(discussionID interface{}) error {
	err := du.discussionRepo.DeleteByID(discussionID)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
