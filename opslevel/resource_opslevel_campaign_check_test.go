package opslevel_test

import (
	"testing"

	"github.com/opslevel/opslevel-go/v2026"
	opsleveltf "github.com/opslevel/terraform-provider-opslevel/opslevel"
)

// The copy mutation reports exactly which checks it created. Anything other than
// one means we cannot say which check this resource owns, and guessing risks
// managing - and later deleting - a check we did not create.
func TestSelectCreatedCheck_ReturnsTheSingleCreatedCheck(t *testing.T) {
	created := []opslevel.Check{{Name: "Arize"}}
	created[0].Id = "copy-z"

	check, err := opsleveltf.SelectCreatedCheck(created)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if check.Id != "copy-z" {
		t.Errorf("expected copy-z, got %q", check.Id)
	}
}

func TestSelectCreatedCheck_ErrorsWhenNothingWasCreated(t *testing.T) {
	if _, err := opsleveltf.SelectCreatedCheck([]opslevel.Check{}); err == nil {
		t.Error("expected an error when the copy created no checks")
	}
}

func TestSelectCreatedCheck_ErrorsWhenMoreThanOneWasCreated(t *testing.T) {
	created := []opslevel.Check{{Name: "a"}, {Name: "b"}}
	if _, err := opsleveltf.SelectCreatedCheck(created); err == nil {
		t.Error("expected an error when the copy created more than one check")
	}
}

// A check on an ended campaign fails the API's active? gate, so GetCheck reports
// it as missing even though it exists and is still attached. Removing it from
// state would make Terraform plan a re-copy, which is not the remedy.
func TestCampaignCheckWasDeleted_EndedCampaignHidesItsChecks(t *testing.T) {
	if opsleveltf.CampaignCheckWasDeleted(opslevel.CampaignStatusEnumEnded) {
		t.Error("a check missing from an ended campaign must not be treated as deleted")
	}
}

func TestCampaignCheckWasDeleted_ActiveCampaignMeansReallyDeleted(t *testing.T) {
	for _, status := range []opslevel.CampaignStatusEnum{
		opslevel.CampaignStatusEnumInProgress,
		opslevel.CampaignStatusEnumDraft,
		opslevel.CampaignStatusEnumScheduled,
		opslevel.CampaignStatusEnumDelayed,
	} {
		if !opsleveltf.CampaignCheckWasDeleted(status) {
			t.Errorf("a check missing from a %q campaign really was deleted", status)
		}
	}
}
