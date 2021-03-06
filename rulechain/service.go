package rulechain

import (
	"context"

	"github.com/cloustone/pandas/mainflux"
	"github.com/cloustone/pandas/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	RULE_STATUS_CREATED = "created"
	RULE_STATUS_STARTED = "started"
	RULE_STATUS_STOPPED = "stopped"
	RULE_STATUS_UNKNOWN = "unknown"
)

var (
	// ErrConflict indicates usage of the existing email during account
	// registration.
	ErrConflict = errors.New("email already taken")

	// ErrMalformedEntity indicates malformed entity specification
	// (e.g. invalid realmname or password).
	ErrMalformedEntity = errors.New("malformed entity specification")

	// ErrUnauthorizedAccess indicates missing or invalid credentials provided
	// when accessing a protected resource.
	ErrUnauthorizedAccess = errors.New("missing or invalid credentials provided")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = errors.New("non-existent entity")

	// ErrRuleChainNotFound indicates a non-existent realm request.
	ErrRuleChainNotFound = errors.New("non-existent rulechain")

	// ErrScanMetadata indicates problem with metadata in db.
	ErrScanMetadata = errors.New("Failed to scan metadata")

	// ErrMissingEmail indicates missing email for password reset request.
	ErrMissingEmail = errors.New("missing email for password reset")

	// ErrUnauthorizedPrincipal indicate the pricipal can not be recognized
	ErrUnauthorizedPrincipal = errors.New("unauthorized principal")
)

//Service service
type Service interface {
	AddNewRuleChain(context.Context, string, RuleChain) error
	GetRuleChainInfo(context.Context, string, string) (RuleChain, error)
	UpdateRuleChain(context.Context, string, RuleChain) (RuleChain, error)
	RevokeRuleChain(context.Context, string, string) error
	ListRuleChain(context.Context, string, uint64, uint64) (RuleChainPage, error)
	UpdateRuleChainStatus(context.Context, string, string, string) error
	SaveStates(*mainflux.Message) error
}

var _ Service = (*rulechainService)(nil)

type rulechainService struct {
	auth       mainflux.AuthNServiceClient
	rulechains RuleChainRepository
	//mutex      sync.RWMutex
	instanceManager instanceManager
	rulechainsCache RuleChainCache
}

//New new
func New(auth mainflux.AuthNServiceClient, rulechains RuleChainRepository, instancemanager instanceManager, rulechainscache RuleChainCache) Service {
	return &rulechainService{
		auth:            auth,
		rulechains:      rulechains,
		instanceManager: instancemanager,
		rulechainsCache: rulechainscache,
	}
}

func (svc rulechainService) AddNewRuleChain(ctx context.Context, token string, rulechain RuleChain) error {
	_, err := svc.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return err
	}
	return svc.rulechains.Save(ctx, rulechain)
}

func (svc rulechainService) GetRuleChainInfo(ctx context.Context, token string, RuleChainID string) (RuleChain, error) {
	res, err := svc.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return RuleChain{}, err
	}
	rulechain, err := svc.rulechains.Retrieve(ctx, res.GetValue(), RuleChainID)
	if err != nil {
		return RuleChain{}, errors.Wrap(ErrRuleChainNotFound, err)
	}

	return rulechain, nil
}

func (svc rulechainService) UpdateRuleChain(ctx context.Context, token string, rulechain RuleChain) (RuleChain, error) {

	res, err := svc.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return RuleChain{}, err
	}

	old_rulechain, err := svc.rulechains.Retrieve(ctx, res.GetValue(), rulechain.ID)
	if err != nil {
		return RuleChain{}, errors.Wrap(ErrRuleChainNotFound, err)
	}
	if old_rulechain.Status == RULE_STATUS_STARTED {
		return RuleChain{}, status.Error(codes.FailedPrecondition, "")
	}

	return svc.rulechains.Update(ctx, rulechain)
}

func (svc rulechainService) RevokeRuleChain(ctx context.Context, token string, RuleChainID string) error {

	res, err := svc.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return err
	}

	rulechain, err := svc.rulechains.Retrieve(ctx, res.GetValue(), RuleChainID)
	if err != nil {
		return errors.Wrap(ErrRuleChainNotFound, err)
	}
	if rulechain.Status == RULE_STATUS_STARTED {
		return status.Error(codes.FailedPrecondition, "")
	}

	return svc.rulechains.Revoke(ctx, res.GetValue(), RuleChainID)
}

func (svc rulechainService) ListRuleChain(ctx context.Context, token string, offset uint64, limit uint64) (RuleChainPage, error) {

	res, err := svc.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return RuleChainPage{}, err
	}

	return svc.rulechains.List(ctx, res.GetValue(), offset, limit)
}

func (svc rulechainService) UpdateRuleChainStatus(ctx context.Context, token string, RuleChainID string, updatestatus string) error {
	res, err := svc.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return err
	}
	rulechain, err := svc.rulechains.Retrieve(ctx, res.GetValue(), RuleChainID)
	if err != nil {
		return errors.Wrap(ErrRuleChainNotFound, err)
	}

	switch updatestatus {
	case UPDATE_RULE_STATUS_START:
		if rulechain.Status != RULE_STATUS_CREATED && rulechain.Status != RULE_STATUS_STOPPED {
			return status.Error(codes.FailedPrecondition, "")
		}

		return svc.instanceManager.startRuleChain(&rulechain)
	case UPDATE_RULE_STATUS_STOP:
		if rulechain.Status != RULE_STATUS_STARTED {
			return status.Error(codes.FailedPrecondition, "")
		}

		return svc.instanceManager.stopRuleChain(&rulechain)
	}
	return nil
}

func (svc rulechainService) SaveStates(msg *mainflux.Message) error {
	/*
		rulechainmessage := message.NewMessage()
		if err := rulechainmessage.UnmarshalBinary(msg.GetPayload()); err != nil {
			logrus.WithError(err).Errorf("rulechain instance receive message failed")
			return err
		}
		return svc.instanceManager.HandleMessage(rulechainmessage, msg)
	*/
	return nil
}
