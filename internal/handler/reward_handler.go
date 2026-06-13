package handler

import (
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"github.com/indim/bali-resik-backend/internal/domain/errors"
	"github.com/indim/bali-resik-backend/internal/dto/request"
	"github.com/indim/bali-resik-backend/internal/helper"
	"github.com/indim/bali-resik-backend/internal/middleware"
	ucase "github.com/indim/bali-resik-backend/internal/interfaces/usecase"
)

type RewardHandler struct {
	rewardUseCase ucase.RewardUseCase
	log           *logrus.Logger
}

func NewRewardHandler(rewardUseCase ucase.RewardUseCase, log *logrus.Logger) *RewardHandler {
	return &RewardHandler{
		rewardUseCase: rewardUseCase,
		log:           log,
	}
}

func (h *RewardHandler) Create(c echo.Context) error {
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		return helper.Unauthorized(c, "Tenant context required")
	}

	var req request.CreateRewardRequest
	if err := helper.ValidateRequest(c, &req); err != nil {
		return err
	}

	resp, err := h.rewardUseCase.CreateReward(tenantID, &req)
	if err != nil {
		return h.handleRewardError(c, err)
	}

	return helper.SuccessCreated(c, resp, "Reward created")
}

func (h *RewardHandler) List(c echo.Context) error {
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		return helper.Unauthorized(c, "Tenant context required")
	}

	resp, err := h.rewardUseCase.ListRewards(tenantID)
	if err != nil {
		return h.handleRewardError(c, err)
	}

	return helper.SuccessOK(c, resp, "Rewards retrieved")
}

func (h *RewardHandler) Update(c echo.Context) error {
	rewardID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.BadRequest(c, "Invalid reward ID")
	}

	var req request.UpdateRewardRequest
	if err := helper.ValidateRequest(c, &req); err != nil {
		return err
	}

	resp, err := h.rewardUseCase.UpdateReward(rewardID, &req)
	if err != nil {
		return h.handleRewardError(c, err)
	}

	return helper.SuccessOK(c, resp, "Reward updated")
}

func (h *RewardHandler) Delete(c echo.Context) error {
	rewardID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return helper.BadRequest(c, "Invalid reward ID")
	}

	if err := h.rewardUseCase.DeleteReward(rewardID); err != nil {
		return h.handleRewardError(c, err)
	}

	return helper.SuccessOK(c, nil, "Reward deleted")
}

func (h *RewardHandler) Redeem(c echo.Context) error {
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		return helper.Unauthorized(c, "Tenant context required")
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return helper.Unauthorized(c, "User authentication required")
	}

	var req request.RedeemRewardRequest
	if err := helper.ValidateRequest(c, &req); err != nil {
		return err
	}

	rewardID, err := uuid.Parse(req.RewardID)
	if err != nil {
		return helper.BadRequest(c, "Invalid reward ID")
	}

	resp, err := h.rewardUseCase.RedeemReward(tenantID, userID, rewardID)
	if err != nil {
		return h.handleRewardError(c, err)
	}

	return helper.SuccessCreated(c, resp, "Reward redeemed successfully")
}

func (h *RewardHandler) GetPoints(c echo.Context) error {
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		return helper.Unauthorized(c, "Tenant context required")
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return helper.Unauthorized(c, "User authentication required")
	}

	resp, err := h.rewardUseCase.GetUserPoints(tenantID, userID)
	if err != nil {
		return h.handleRewardError(c, err)
	}

	return helper.SuccessOK(c, resp, "Points retrieved")
}

func (h *RewardHandler) GetTransactions(c echo.Context) error {
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		return helper.Unauthorized(c, "Tenant context required")
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return helper.Unauthorized(c, "User authentication required")
	}

	page, perPage := parsePageParams(c)

	transactions, total, err := h.rewardUseCase.GetTransactionHistory(tenantID, userID, page, perPage)
	if err != nil {
		return h.handleRewardError(c, err)
	}

	meta := &helper.PaginationMeta{
		Page:    page,
		PerPage: perPage,
		Total:   total,
	}

	return helper.SuccessPaginated(c, transactions, meta, "Transaction history retrieved")
}

func (h *RewardHandler) handleRewardError(c echo.Context, err error) error {
	domainErr, ok := err.(*errors.DomainError)
	if !ok {
		h.log.WithError(err).Error("unhandled reward error")
		return helper.InternalError(c, "An unexpected error occurred")
	}

	switch domainErr.Code {
	case "NOT_FOUND":
		return helper.NotFound(c, domainErr.Message)
	case "REWARD_INACTIVE", "OUT_OF_STOCK":
		return helper.BadRequest(c, domainErr.Message)
	case "INSUFFICIENT_POINTS":
		return helper.BadRequest(c, domainErr.Message)
	default:
		return helper.InternalError(c, domainErr.Message)
	}
}

func parsePageParams(c echo.Context) (int, int) {
	page := 1
	perPage := 20
	if p := c.QueryParam("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if p := c.QueryParam("per_page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 && v <= 100 {
			perPage = v
		}
	}
	return page, perPage
}
