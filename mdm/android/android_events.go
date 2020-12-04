package android

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mattrax/Mattrax/internal/db"
	"github.com/mattrax/Mattrax/mdm/android/policies"
	"github.com/mattrax/Mattrax/mdm/protocol"
)

func (p *Protocol) Events() protocol.EventHandlers {
	return protocol.EventHandlers{
		CreatePolicy: p.updatePolicy,
		UpdatePolicy: p.updatePolicy,
		DeletePolicy: p.deletePolicy,
	}
}

func (p *Protocol) updatePolicy(policy db.Policy) error {
	var policyPayload protocol.Policy
	if err := json.Unmarshal(policy.Payload, &policyPayload); err != nil {
		return fmt.Errorf("error decoding policy payload: %w", err)
	}

	androidPolicy, err := policies.GenerateAndroidPolicy(policyPayload)
	if err != nil {
		return fmt.Errorf("error generating Android policy: %w", err)
	}

	tenant, err := p.srv.DB.GetTenant(context.Background(), policy.TenantID)
	if err != nil {
		return fmt.Errorf("error retrieving tenant: %w", err)
	}

	if !tenant.AfwEnterpriseID.Valid {
		return nil
	}

	if _, err := p.ams.Enterprises.Policies.Patch(fmt.Sprintf("enterprises/%s/policies/%s", tenant.AfwEnterpriseID.String, policy.ID), &androidPolicy).Do(); err != nil {
		return fmt.Errorf("error patching policy to Google API: %w", err)
	}

	return nil
}

func (p *Protocol) deletePolicy(policy db.Policy) error {
	tenant, err := p.srv.DB.GetTenant(context.Background(), policy.TenantID)
	if err != nil {
		return fmt.Errorf("error retrieving tenant: %w", err)
	}

	if !tenant.AfwEnterpriseID.Valid {
		return nil
	}

	if _, err := p.ams.Enterprises.Policies.Delete(fmt.Sprintf("enterprises/%s/policies/%s", tenant.AfwEnterpriseID.String, policy.ID)).Do(); err != nil {
		return fmt.Errorf("error deleting policy on Google API: %w", err)
	}

	return nil
}
