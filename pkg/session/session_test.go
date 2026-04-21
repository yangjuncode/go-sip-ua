package session

import (
	"testing"

	"github.com/ghettovoice/gosip/sip"
	"github.com/ghettovoice/gosip/sip/parser"
)

func TestSessionMakeRequestUsesManagedIncrementingCSeq(t *testing.T) {
	contactURI := mustParseSessionSipUri(t, "sip:alice@example.com")
	contact := &sip.ContactHeader{
		Address: contactURI,
		Params:  sip.NewParams(),
	}
	inviteRequest := mustBuildInviteRequest(t, 41)
	inviteResponse := sip.NewResponseFromRequest("", inviteRequest, 200, "OK", "")

	nextSeq := uint32(100)
	session := NewInviteSession(nil, "UAC", func(method sip.RequestMethod) *sip.CSeq {
		cseq := &sip.CSeq{SeqNo: nextSeq, MethodName: method}
		nextSeq++
		return cseq
	}, contact, inviteRequest, sip.CallID("call-1"), nil, Outgoing, nil)
	session.response = inviteResponse

	infoRequest := session.makeRequest(session.uaType, sip.INFO, sip.MessageID(session.callID), session.request, session.response)
	byeRequest := session.makeRequest(session.uaType, sip.BYE, sip.MessageID(session.callID), session.request, session.response)

	infoCSeq, _ := infoRequest.CSeq()
	byeCSeq, _ := byeRequest.CSeq()

	if infoCSeq.SeqNo != 100 || infoCSeq.MethodName != sip.INFO {
		t.Fatalf("info cseq = %d %s, want 100 INFO", infoCSeq.SeqNo, infoCSeq.MethodName)
	}
	if byeCSeq.SeqNo != 101 || byeCSeq.MethodName != sip.BYE {
		t.Fatalf("bye cseq = %d %s, want 101 BYE", byeCSeq.SeqNo, byeCSeq.MethodName)
	}
}

func mustBuildInviteRequest(t *testing.T, seqNo uint32) sip.Request {
	t.Helper()

	builder := sip.NewRequestBuilder()
	fromURI := mustParseSessionUri(t, "sip:alice@example.com")
	toURI := mustParseSessionUri(t, "sip:bob@example.com")
	contactURI := mustParseSessionUri(t, "sip:alice@example.com")
	callID := sip.CallID("call-1")
	builder.SetMethod(sip.INVITE)
	builder.SetFrom(&sip.Address{
		Uri:    fromURI,
		Params: sip.NewParams().Add("tag", sip.String{Str: "from-tag"}),
	})
	builder.SetTo(&sip.Address{Uri: toURI})
	builder.SetContact(&sip.Address{Uri: contactURI})
	builder.SetRecipient(toURI)
	builder.SetCallID(&callID)
	builder.SetSeqNo(uint(seqNo))

	request, err := builder.Build()
	if err != nil {
		t.Fatalf("build invite request failed: %v", err)
	}

	return request
}

func mustParseSessionSipUri(t *testing.T, raw string) sip.ContactUri {
	t.Helper()

	uri, err := parser.ParseUri(raw)
	if err != nil {
		t.Fatalf("parse sip uri %q failed: %v", raw, err)
	}

	return uri.(sip.ContactUri)
}

func mustParseSessionUri(t *testing.T, raw string) sip.Uri {
	t.Helper()

	uri, err := parser.ParseUri(raw)
	if err != nil {
		t.Fatalf("parse uri %q failed: %v", raw, err)
	}

	return uri
}
