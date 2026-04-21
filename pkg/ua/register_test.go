package ua

import (
	"testing"

	"github.com/cloudwebrtc/go-sip-ua/pkg/account"
	"github.com/ghettovoice/gosip/sip"
	"github.com/ghettovoice/gosip/sip/parser"
)

func TestRegisterPrepareRequestUsesManagedIncrementingCSeq(t *testing.T) {
	ua := &UserAgent{config: &UserAgentConfig{}, cseqManager: NewCSeqManager(1)}
	profile := &account.Profile{
		URI:        mustParseRegisterUri(t, "sip:alice@example.com"),
		InstanceID: "nil",
	}
	recipient := mustParseRegisterSipUriValue(t, "sip:registrar.example.com")
	register := NewRegister(ua, profile, recipient, nil)

	if err := register.prepareRequest(300); err != nil {
		t.Fatalf("prepareRequest 1 failed: %v", err)
	}
	cseq1, _ := (*register.request).CSeq()

	if err := register.prepareRequest(300); err != nil {
		t.Fatalf("prepareRequest 2 failed: %v", err)
	}
	cseq2, _ := (*register.request).CSeq()

	if cseq1.SeqNo != 1 {
		t.Fatalf("first register cseq = %d, want 1", cseq1.SeqNo)
	}
	if cseq2.SeqNo != 2 {
		t.Fatalf("second register cseq = %d, want 2", cseq2.SeqNo)
	}
}

func mustParseRegisterSipUriValue(t *testing.T, raw string) sip.SipUri {
	t.Helper()

	uri, err := parser.ParseSipUri(raw)
	if err != nil {
		t.Fatalf("parse sip uri %q failed: %v", raw, err)
	}

	return uri
}

func mustParseRegisterUri(t *testing.T, raw string) sip.Uri {
	t.Helper()

	uri, err := parser.ParseUri(raw)
	if err != nil {
		t.Fatalf("parse uri %q failed: %v", raw, err)
	}

	return uri
}
