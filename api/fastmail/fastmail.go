package fastmail

import (
	"errors"

	"git.sr.ht/~rockorager/go-jmap"
	"git.sr.ht/~rockorager/go-jmap/mail/email"
)

type FastmailClient struct {
	JmapClient *jmap.Client
	AccountId  jmap.ID
	MailboxId  jmap.ID
}

func NewEmailClient(
	sessionEndpoint string,
	accessToken string,
	accountId string,
	mailboxId string,
) *FastmailClient {
	client := &jmap.Client{
		SessionEndpoint: sessionEndpoint,
	}

	client.WithAccessToken(accessToken)

	return &FastmailClient{
		JmapClient: client,
		AccountId:  jmap.ID(accountId),
		MailboxId:  jmap.ID(mailboxId),
	}
}

func (client *FastmailClient) FetchLivetrackEmail() (string, error) {

	// Create a new request
	req := &jmap.Request{}

	// Invoke a method. The CallID of this method will be returned to be
	// used when chaining calls
	callId := req.Invoke(&email.Query{
		Account: client.AccountId,
		Filter: &email.FilterCondition{
			InMailbox: client.MailboxId,
		},
		Sort: []*email.SortComparator{
			{Property: "receivedAt", IsAscending: false},
		},
		Limit: 1,
	})

	// Invoke a result reference call

	req.Invoke(&email.Get{
		Account:             client.AccountId,
		FetchHTMLBodyValues: true,
		Properties:          []string{"htmlBody", "bodyValues"},
		ReferenceIDs: &jmap.ResultReference{
			ResultOf: callId,        // The CallID of the referenced method
			Name:     "Email/query", // The name of the referenced method
			Path:     "/ids/*",      // JSON pointer to the location of the reference
		},
	})

	resp, err := client.JmapClient.Do(req)
	if err != nil {
		return "", err
	}

	// Loop through the responses to invidividual invocations
	for _, inv := range resp.Responses {
		// Our result to individual calls is in the Args field of the
		// invocation
		switch r := inv.Args.(type) {
		case *email.GetResponse:
			val := r.List[0].BodyValues["1"].Value
			return val, nil
		case *jmap.MethodError:
			return "", errors.New(*r.Description)
		}
	}
	return "", errors.New("no emails returned")
}
