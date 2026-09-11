package zain

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
)

func (c *Client) CreateTicket(ctx context.Context, req *SubmitTicketReq) (*SubmitTicketResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("zain: submit ticket request cannot be nil")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if req.MSISDN == "" {
		req.MSISDN = c.GetMSISDN()
	}
	if req.Source == "" {
		req.Source = "MyZain"
	}
	if req.TicketLanguage == "" {
		c.mu.RLock()
		req.TicketLanguage = c.language
		c.mu.RUnlock()
	}

	_ = writer.WriteField("city", req.City)
	_ = writer.WriteField("ticket_language", req.TicketLanguage)
	_ = writer.WriteField("msisdn", req.MSISDN)
	_ = writer.WriteField("source", req.Source)
	_ = writer.WriteField("governorate", req.Governorate)
	_ = writer.WriteField("description", req.Description)
	_ = writer.WriteField("summary_id", req.SummaryID)
	_ = writer.WriteField("question_answers", req.QuestionAnswers)

	for i, att := range req.Attachments {
		fieldName := fmt.Sprintf("attachment_%d", i)
		fileName := att.Name
		if fileName == "" {
			fileName = fmt.Sprintf("file_%d.jpg", i)
		}
		part, err := writer.CreateFormFile(fieldName, fileName)
		if err == nil {
			_, _ = part.Write([]byte(att.Data))
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("zain: closing multipart writer: %w", err)
	}

	fullURL := c.baseURL + "/api/complaints/create-ticket"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("zain: creating multipart request: %w", err)
	}

	c.mu.RLock()
	httpReq.Header.Set("Skel-Accept-Language", c.language)
	httpReq.Header.Set("Skel-Platform", c.platform)
	httpReq.Header.Set("Skel-OS-Version", c.osVersion)
	httpReq.Header.Set("Skel-Fix-Version", c.appVersion)
	httpReq.Header.Set("Skel-Installation-Id", c.installationID)
	httpReq.Header.Set("User-Agent", c.userAgent)
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())
	if c.accessToken != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.accessToken)
	}
	if c.msisdn != "" {
		httpReq.Header.Set("gal-msisdn", c.msisdn)
	}
	c.mu.RUnlock()

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("zain: multipart request failed: %w", err)
	}
	defer resp.Body.Close()

	var apiResp APIResponse[SubmitTicketResponse]
	if err := decodeResponse(resp, &apiResp); err != nil {
		return nil, err
	}

	return &apiResp.Data, nil
}

func (c *Client) GetQuestions(ctx context.Context, summaryCode string) (*QuestionsData, error) {
	if summaryCode == "" {
		return nil, fmt.Errorf("zain: summaryCode cannot be empty")
	}

	path := "api/complaints/questions?txt_summary_code=" + url.QueryEscape(summaryCode)
	return doAndDecode[QuestionsData](ctx, c, http.MethodGet, path, nil, false, false)
}

func (c *Client) GetReOpenReasons(ctx context.Context) (*ReOpenReasonsData, error) {
	return doAndDecode[ReOpenReasonsData](ctx, c, http.MethodGet, "api/complaints/reopen_reason", nil, false, false)
}

func (c *Client) GetSummaries(ctx context.Context, categorizationCode string, sourceOfCreation ...string) (*SummariesData, error) {
	source := "MyZain"
	if len(sourceOfCreation) > 0 && sourceOfCreation[0] != "" {
		source = sourceOfCreation[0]
	}

	params := url.Values{}
	params.Set("categorization_code", categorizationCode)
	params.Set("source_of_creation", source)

	path := "api/complaints/summary?" + params.Encode()
	return doAndDecode[SummariesData](ctx, c, http.MethodGet, path, nil, false, false)
}

func (c *Client) GetTemplates(ctx context.Context, sourceOfCreation, issueType, category, compType string) (*TemplatesData, error) {
	if sourceOfCreation == "" {
		sourceOfCreation = "MyZain"
	}
	if issueType == "" {
		issueType = "Complaint"
	}
	if compType == "" {
		compType = "Zain App-Self Service"
	}

	params := url.Values{}
	params.Set("source_of_creation", sourceOfCreation)
	params.Set("issue_type", issueType)
	if category != "" {
		params.Set("category", category)
	}
	params.Set("type", compType)

	path := "api/complaints/items?" + params.Encode()
	return doAndDecode[TemplatesData](ctx, c, http.MethodGet, path, nil, false, false)
}

func (c *Client) GetTickets(ctx context.Context, qualification string) (*TicketsData, error) {
	path := "api/complaints/list-ticket"
	if qualification != "" {
		path += "?qualification=" + url.QueryEscape(qualification)
	}

	return doAndDecode[TicketsData](ctx, c, http.MethodGet, path, nil, false, false)
}

func (c *Client) ReOpenTicket(ctx context.Context, req *ReOpenTicketsReq) (*ReOpenTicketsResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("zain: reopen tickets request cannot be nil")
	}
	return doAndDecode[ReOpenTicketsResponse](ctx, c, http.MethodPost, "api/complaints/reopen-ticket", req, false, false)
}
