package ua

import (
	"testing"

	"github.com/cloudwebrtc/go-sip-ua/pkg/auth"
	"github.com/ghettovoice/gosip/sip"
	"github.com/ghettovoice/gosip/sip/parser"
)

func TestBuildRequestAssignsManagedCSeq(t *testing.T) {
	ua := &UserAgent{config: &UserAgentConfig{}, cseqManager: NewCSeqManager(1)}
	fromURI := mustParseSipUriValue(t, "sip:alice@example.com")
	toURI := mustParseSipUriValue(t, "sip:bob@example.com")
	recipient := mustParseSipUriValue(t, "sip:bob@example.com")
	from := &sip.Address{
		Uri:    &fromURI,
		Params: sip.NewParams().Add("tag", sip.String{Str: "from-tag"}),
	}
	to := &sip.Address{Uri: &toURI}
	callID := sip.CallID("call-1")

	request1, err := ua.buildRequest(sip.MESSAGE, from, to, nil, recipient, nil, &callID)
	if err != nil {
		t.Fatalf("buildRequest 1 failed: %v", err)
	}
	request2, err := ua.buildRequest(sip.MESSAGE, from, to, nil, recipient, nil, &callID)
	if err != nil {
		t.Fatalf("buildRequest 2 failed: %v", err)
	}

	cseq1, ok := (*request1).CSeq()
	if !ok {
		t.Fatalf("request1 missing cseq")
	}
	cseq2, ok := (*request2).CSeq()
	if !ok {
		t.Fatalf("request2 missing cseq")
	}

	if cseq1.SeqNo != 1 {
		t.Fatalf("request1 cseq = %d, want 1", cseq1.SeqNo)
	}
	if cseq2.SeqNo != 2 {
		t.Fatalf("request2 cseq = %d, want 2", cseq2.SeqNo)
	}
}

func TestReplaceRequestCSeqUsesNextManagedValue(t *testing.T) {
	ua := &UserAgent{config: &UserAgentConfig{}, cseqManager: NewCSeqManager(1)}
	request := mustBuildObservedRequest(t, sip.INVITE, "call-observe", "from-observe", 7)

	ua.replaceRequestCSeq(request)

	cseq, ok := request.CSeq()
	if !ok {
		t.Fatalf("request missing cseq")
	}

	if cseq.SeqNo != 1 || cseq.MethodName != sip.INVITE {
		t.Fatalf("replaced cseq = %d %s, want 1 INVITE", cseq.SeqNo, cseq.MethodName)
	}
}

func TestAuthRetryConsumesNextManagedCSeq(t *testing.T) {
	ua := &UserAgent{config: &UserAgentConfig{}, cseqManager: NewCSeqManager(1)}
	fromURI := mustParseSipUriValue(t, "sip:alice@example.com")
	toURI := mustParseSipUriValue(t, "sip:bob@example.com")
	recipient := mustParseSipUriValue(t, "sip:bob@example.com")
	from := &sip.Address{
		Uri:    &fromURI,
		Params: sip.NewParams().Add("tag", sip.String{Str: "from-tag"}),
	}
	to := &sip.Address{Uri: &toURI}
	callID := sip.CallID("call-auth")

	request, err := ua.buildRequest(sip.MESSAGE, from, to, nil, recipient, nil, &callID)
	if err != nil {
		t.Fatalf("buildRequest failed: %v", err)
	}

	response := sip.NewResponseFromRequest("", *request, 401, "Unauthorized", "")
	response.AppendHeader(&sip.GenericHeader{
		HeaderName: "WWW-Authenticate",
		Contents:   `Digest realm="example.com",nonce="nonce-value"`,
	})

	authorizer := auth.NewClientAuthorizer("alice", "secret")
	if err := authorizer.AuthorizeRequest(*request, response); err != nil {
		t.Fatalf("AuthorizeRequest failed: %v", err)
	}
	ua.replaceRequestCSeq(*request)

	cseq, ok := (*request).CSeq()
	if !ok {
		t.Fatalf("request missing cseq after auth retry")
	}
	if cseq.SeqNo != 2 || cseq.MethodName != sip.MESSAGE {
		t.Fatalf("auth retry cseq = %d %s, want 2 MESSAGE", cseq.SeqNo, cseq.MethodName)
	}

	next := ua.nextRequestCSeq(sip.BYE)
	if next == nil || next.SeqNo != 3 || next.MethodName != sip.BYE {
		t.Fatalf("next managed cseq = %#v, want 3 BYE", next)
	}
}

func mustBuildObservedRequest(t *testing.T, method sip.RequestMethod, callID, fromTag string, seqNo uint32) sip.Request {
	t.Helper()

	builder := sip.NewRequestBuilder()
	fromURI := mustParseSipUriValue(t, "sip:alice@example.com")
	toURI := mustParseSipUriValue(t, "sip:bob@example.com")
	requestCallID := sip.CallID(callID)
	from := &sip.Address{
		Uri:    &fromURI,
		Params: sip.NewParams().Add("tag", sip.String{Str: fromTag}),
	}
	to := &sip.Address{Uri: &toURI}
	builder.SetMethod(method)
	builder.SetFrom(from)
	builder.SetTo(to)
	builder.SetRecipient(&toURI)
	builder.SetCallID(&requestCallID)
	builder.SetSeqNo(uint(seqNo))

	request, err := builder.Build()
	if err != nil {
		t.Fatalf("build observed request failed: %v", err)
	}

	return request
}

func mustParseSipUriValue(t *testing.T, raw string) sip.SipUri {
	t.Helper()

	uri, err := parser.ParseSipUri(raw)
	if err != nil {
		t.Fatalf("parse sip uri %q failed: %v", raw, err)
	}

	return uri
}
