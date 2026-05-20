package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"ridehailing/backend/internal/model"
)

const (
	maxDriverAIPromptLength     = 1200
	maxPassengerAIMessages      = 8
	maxPassengerAIMessageLength = 500
)

var (
	monthDayPattern = regexp.MustCompile(`(\d{1,2})月(\d{1,2})[日号]?`)
)

// DriverTripAIDraft 是司机端 AI 产出的结构化班次草稿。
type DriverTripAIDraft struct {
	Prompt             string   `json:"prompt"`
	TripName           string   `json:"tripName"`
	StartCity          string   `json:"startCity"`
	EndCity            string   `json:"endCity"`
	DepartureTimeLocal string   `json:"departureTimeLocal"`
	ArrivalTimeLocal   string   `json:"arrivalTimeLocal"`
	SeatTotal          int      `json:"seatTotal"`
	PriceCent          int      `json:"priceCent"`
	VehicleType        string   `json:"vehicleType"`
	Stops              []string `json:"stops"`
	Remark             string   `json:"remark"`
	Suggestions        []string `json:"suggestions"`

	Origin        string `json:"origin"`
	Destination   string `json:"destination"`
	DepartureTime string `json:"departureTime"`
	ArrivalTime   string `json:"arrivalTime"`
	SeatCount     int    `json:"seatCount"`
	Seats         int    `json:"seats"`
	Response      string `json:"response"`
}

// AIChatMessage 是前端会话消息。
type AIChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type PassengerAIIntent string

const (
	PassengerAIIntentGeneral PassengerAIIntent = "general"
	PassengerAIIntentRoute   PassengerAIIntent = "route"
	PassengerAIIntentOrders  PassengerAIIntent = "orders"
	PassengerAIIntentRefund  PassengerAIIntent = "refund"
)

// PassengerAIChatResult 是返回给前端的统一结构。
type PassengerAIChatResult struct {
	Reply       string              `json:"reply"`
	Response    string              `json:"response,omitempty"`
	Suggestions []string            `json:"suggestions"`
	Intent      PassengerAIIntent   `json:"intent"`
	Context     *PassengerAIContext `json:"context,omitempty"`
}

type PassengerAIContext struct {
	Intent           PassengerAIIntent    `json:"intent"`
	RouteQuery       *AIRouteQuery        `json:"routeQuery,omitempty"`
	RouteResults     []*AITripCard        `json:"routeResults,omitempty"`
	OrderSummary     *AIOrderSummary      `json:"orderSummary,omitempty"`
	OrderResults     []*AIOrderCard       `json:"orderResults,omitempty"`
	RefundRules      []string             `json:"refundRules,omitempty"`
	SystemHints      []string             `json:"systemHints,omitempty"`
	KnowledgeSources []*AIKnowledgeSource `json:"knowledgeSources,omitempty"`
}

type AIKnowledgeSource struct {
	DocumentID  string `json:"documentId"`
	ChunkID     string `json:"chunkId"`
	Title       string `json:"title"`
	SectionPath string `json:"sectionPath"`
	Content     string `json:"content"`
	FinalScore  string `json:"finalScore"`
}

type AIRouteQuery struct {
	StartCity     string `json:"startCity"`
	EndCity       string `json:"endCity"`
	Date          string `json:"date"`
	AllowTransfer bool   `json:"allowTransfer,omitempty"`
}

type AITripLeg struct {
	TripID        uint   `json:"tripId"`
	Route         string `json:"route"`
	DepartureTime string `json:"departureTime"`
	ArrivalTime   string `json:"arrivalTime"`
	SeatAvailable int    `json:"seatAvailable"`
	PriceCent     int    `json:"priceCent"`
	VehicleType   string `json:"vehicleType"`
}

type AIRouteSuggestion struct {
	Route            string `json:"route"`
	TransferCity     string `json:"transferCity"`
	FirstLegCount    int    `json:"firstLegCount"`
	SecondLegCount   int    `json:"secondLegCount"`
	TotalOptionCount int    `json:"totalOptionCount"`
	Reason           string `json:"reason"`
}

type AITripCard struct {
	ID                 uint                 `json:"id"`
	Kind               string               `json:"kind"`
	Route              string               `json:"route"`
	DepartureTime      string               `json:"departureTime"`
	ArrivalTime        string               `json:"arrivalTime"`
	SeatAvailable      int                  `json:"seatAvailable"`
	PriceCent          int                  `json:"priceCent"`
	VehicleType        string               `json:"vehicleType"`
	TransferCity       string               `json:"transferCity,omitempty"`
	TransferWaitMinute int                  `json:"transferWaitMinute,omitempty"`
	MatchedCity        string               `json:"matchedCity,omitempty"`
	MatchRoles         []string             `json:"matchRoles,omitempty"`
	Stops              []string             `json:"stops,omitempty"`
	Legs               []*AITripLeg         `json:"legs,omitempty"`
	Suggestions        []*AIRouteSuggestion `json:"suggestions,omitempty"`
}

type AIOrderSummary struct {
	TotalCount               int `json:"totalCount"`
	PendingPaymentCount      int `json:"pendingPaymentCount"`
	PendingVerificationCount int `json:"pendingVerificationCount"`
	CompletedCount           int `json:"completedCount"`
	RefundRequestedCount     int `json:"refundRequestedCount"`
	RefundRejectedCount      int `json:"refundRejectedCount"`
	RefundedCount            int `json:"refundedCount"`
}

type AIOrderCard struct {
	ID               uint   `json:"id"`
	OrderNo          string `json:"orderNo"`
	Route            string `json:"route"`
	DepartureTime    string `json:"departureTime"`
	OrderStatus      string `json:"orderStatus"`
	PayStatus        string `json:"payStatus"`
	RefundStatus     string `json:"refundStatus"`
	RefundReviewNote string `json:"refundReviewNote,omitempty"`
	Amount           int    `json:"amount"`
}

// AIService 统一封装 AI 请求和本地工具执行。
type AIService struct {
	apiKey   string
	baseURL  string
	model    string
	client   *http.Client
	endpoint string

	ticketService     *TicketService
	orderService      *OrderService
	userService       *UserService
	knowledgeService  *KnowledgeService
	tokenUsageService *TokenUsageService
}

func NewAIService(
	apiKey string,
	baseURL string,
	model string,
	timeout time.Duration,
	ticketService *TicketService,
	orderService *OrderService,
	userService *UserService,
	knowledgeService *KnowledgeService,
	tokenUsageService *TokenUsageService,
) *AIService {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	if strings.TrimSpace(model) == "" {
		model = "gpt-5-mini"
	}
	if timeout <= 0 {
		timeout = 20 * time.Second
	}

	return &AIService{
		apiKey:            strings.TrimSpace(apiKey),
		baseURL:           baseURL,
		model:             strings.TrimSpace(model),
		client:            &http.Client{Timeout: timeout},
		endpoint:          baseURL + "/chat/completions",
		ticketService:     ticketService,
		orderService:      orderService,
		userService:       userService,
		knowledgeService:  knowledgeService,
		tokenUsageService: tokenUsageService,
	}
}

// GenerateDriverTripDraft 使用 function calling 让模型返回结构化班次参数。
func (s *AIService) GenerateDriverTripDraft(ctx context.Context, currentUserID uint, currentUserRole, prompt string) (*DriverTripAIDraft, error) {
	if currentUserRole != model.RoleDriver {
		return nil, errors.New("only driver can use driver AI")
	}
	if s.apiKey == "" {
		return nil, errors.New("OPENAI_API_KEY is not configured")
	}

	prompt = strings.TrimSpace(prompt)
	if prompt == "" {
		return nil, errors.New("prompt is required")
	}
	if len([]rune(prompt)) > maxDriverAIPromptLength {
		return nil, fmt.Errorf("prompt is too long, max %d characters", maxDriverAIPromptLength)
	}

	messages := []chatCompletionMessage{
		{
			Role: "system",
			Content: strings.Join([]string{
				"你是一个司机端班次草稿助手。",
				"你必须调用函数 build_driver_trip_draft 返回参数。",
				"不要直接输出自然语言。",
				"时间格式必须是 YYYY-MM-DDTHH:mm。",
				"priceCent 表示整数分。",
				"vehicleType 只能是：商务大巴、城际快线、拼车专线。",
			}, "\n"),
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	response, err := s.requestChatCompletion(ctx, aiUsageContext{
		UserID:  currentUserID,
		Role:    currentUserRole,
		Feature: model.TokenFeatureDriverAIDraft,
	}, chatCompletionRequest{
		Model:          s.model,
		Messages:       messages,
		Tools:          driverDraftTools(),
		ToolChoice:     forceToolChoice("build_driver_trip_draft"),
		EnableThinking: false,
	})
	if err != nil {
		return nil, err
	}

	message, err := response.firstMessage()
	if err != nil {
		return nil, err
	}

	call, err := firstToolCallByName(message.ToolCalls, "build_driver_trip_draft")
	if err != nil {
		return nil, err
	}

	draft, err := parseDriverTripAIDraft(call.Function.Arguments)
	if err != nil {
		return nil, fmt.Errorf("parse driver tool arguments failed: %w", err)
	}

	draft.Prompt = prompt
	reconcileDriverTripDraft(draft)
	normalizeDriverTripDraft(draft)
	if err := validateDriverTripDraft(draft); err != nil {
		return nil, err
	}

	return draft, nil
}

// ChatPassenger runs the tool-calling flow for passenger chat.
// 1. The first model call decides which local tool to invoke and returns JSON arguments.
// 2. The backend executes local tools such as ticket search, order lookup, and refund rule lookup.
// 3. If tool results are enough for a direct answer, the backend returns them immediately.
// 4. Otherwise, a second model call summarizes tool results through reply_directly JSON.
func (s *AIService) ChatPassenger(
	ctx context.Context,
	currentUserID uint,
	currentUserRole string,
	messages []AIChatMessage,
) (*PassengerAIChatResult, error) {
	if s.apiKey == "" {
		return nil, errors.New("OPENAI_API_KEY is not configured")
	}

	normalized := normalizeChatMessages(messages)
	if len(normalized) == 0 {
		return nil, errors.New("messages are required")
	}

	now := time.Now()
	lastMessage := lastPassengerMessage(normalized)
	if directResult, ok := buildDirectCurrentDateReply(lastMessage, now); ok {
		directResult.Intent = PassengerAIIntentGeneral
		return directResult, nil
	}

	if s.knowledgeService != nil && shouldUseRAG(lastMessage) {
		log.Printf("passenger ai rag triggered query=%q", lastMessage)
		ragResults, err := s.knowledgeService.SearchKnowledge(ctx, SearchKnowledgeInput{
			Query:    lastMessage,
			TopK:     4,
			Category: inferKnowledgeCategory(lastMessage),
			UserID:   currentUserID,
			Role:     currentUserRole,
			Feature:  model.TokenFeaturePassengerAI,
		})
		if err != nil {
			log.Printf("passenger ai rag search failed: err=%v", err)
		} else if len(ragResults) > 0 {
			log.Printf("passenger ai rag search success: topk=%d", len(ragResults))
			return s.chatPassengerWithRAG(ctx, currentUserID, currentUserRole, normalized, ragResults)
		}
	}

	plannerMessages := buildPassengerPlannerMessages(normalized, now)
	plannerResponse, err := s.requestChatCompletion(ctx, aiUsageContext{
		UserID:  currentUserID,
		Role:    currentUserRole,
		Feature: model.TokenFeaturePassengerAI,
	}, chatCompletionRequest{
		Model:          s.model,
		Messages:       plannerMessages,
		Tools:          passengerPlannerTools(),
		ToolChoice:     "auto",
		EnableThinking: false,
	})

	if err != nil {
		return nil, err
	}

	plannerMessage, err := plannerResponse.firstMessage()
	if err != nil {
		return nil, err
	}

	log.Printf("passenger ai planner content=%q tool_calls=%d", plannerMessage.Content, len(plannerMessage.ToolCalls))
	for _, call := range plannerMessage.ToolCalls {
		log.Printf("passenger ai planner tool=%s args=%s", call.Function.Name, call.Function.Arguments)
	}

	// 极端情况下模型没有走工具调用，直接按普通文本兜底。
	if len(plannerMessage.ToolCalls) == 0 {
		result, err := parsePassengerAIChatResult(plannerMessage.Content)
		if err != nil {
			return nil, fmt.Errorf("parse planner output failed: %w", err)
		}
		if strings.TrimSpace(result.Reply) == "" {
			result.Reply = "我可以先帮你整理问题，但还需要更具体的信息。你可以补充出发地、目的地、日期，或者直接询问订单和退款。"
		}
		if len(result.Suggestions) == 0 {
			result.Suggestions = []string{
				"帮我查明天杭州到苏州的票",
				"帮我看下我的订单",
				"我这笔订单能退款吗",
			}
		}
		result.Intent = PassengerAIIntentGeneral
		return result, nil
	}

	contextData := &PassengerAIContext{
		Intent:      PassengerAIIntentGeneral,
		SystemHints: []string{},
	}

	toolMessages := make([]chatCompletionMessage, 0, len(plannerMessage.ToolCalls)+1)
	toolMessages = append(toolMessages, chatCompletionMessage{
		Role:      "assistant",
		Content:   plannerMessage.Content,
		ToolCalls: plannerMessage.ToolCalls,
	})

	for _, call := range plannerMessage.ToolCalls {
		if call.Function.Name == "reply_directly" {
			log.Printf("passenger ai planner replied directly args=%s", call.Function.Arguments)
			result, err := parsePassengerReplyDirectlyArguments(call.Function.Arguments)
			if err != nil {
				return nil, fmt.Errorf("parse direct reply arguments failed: %w", err)
			}
			if result.Intent == "" {
				result.Intent = PassengerAIIntentGeneral
			}
			result.Context = contextData
			return result, nil
		}
		log.Printf("passenger ai executing tool=%s", call.Function.Name)
		toolMessage, err := s.executePassengerToolCall(ctx, currentUserID, currentUserRole, lastMessage, now, call, contextData)
		if err != nil {
			log.Printf("passenger ai tool failed tool=%s err=%v", call.Function.Name, err)
			return nil, err
		}
		log.Printf("passenger ai tool output tool=%s content=%s", call.Function.Name, toolMessage.Content)
		toolMessages = append(toolMessages, toolMessage)
	}

	log.Printf("passenger ai tool context: intent=%s routeResults=%d orderResults=%d refundRules=%d",
		contextData.Intent,
		len(contextData.RouteResults),
		len(contextData.OrderResults),
		len(contextData.RefundRules),
	)

	finalMessages := buildPassengerResponderMessages(normalized, toolMessages, now)
	finalResponse, err := s.requestChatCompletion(ctx, aiUsageContext{
		UserID:  currentUserID,
		Role:    currentUserRole,
		Feature: model.TokenFeaturePassengerAI,
	}, chatCompletionRequest{
		Model:          s.model,
		Messages:       finalMessages,
		Tools:          passengerReplyTools(),
		ToolChoice:     forceToolChoice("reply_directly"),
		EnableThinking: false,
	})
	if err != nil {
		return nil, err
	}

	finalMessage, err := finalResponse.firstMessage()
	if err != nil {
		return nil, err
	}
	log.Printf("passenger ai final content=%q tool_calls=%d", finalMessage.Content, len(finalMessage.ToolCalls))
	for _, call := range finalMessage.ToolCalls {
		log.Printf("passenger ai final tool=%s args=%s", call.Function.Name, call.Function.Arguments)
	}

	replyCall, err := firstToolCallByName(finalMessage.ToolCalls, "reply_directly")
	if err != nil {
		// 如果模型没有按工具返回，尝试从文本直接解析。
		result, parseErr := parsePassengerAIChatResult(finalMessage.Content)
		if parseErr != nil {
			return nil, err
		}
		result.Intent = contextData.Intent
		result.Context = contextData
		return result, nil
	}

	result, err := parsePassengerReplyDirectlyArguments(replyCall.Function.Arguments)
	if err != nil {
		return nil, fmt.Errorf("parse final reply arguments failed: %w", err)
	}
	if result.Intent == "" {
		result.Intent = contextData.Intent
	}
	result.Context = contextData
	log.Printf("passenger ai final reply intent=%s reply=%q suggestions=%d", result.Intent, result.Reply, len(result.Suggestions))
	return result, nil
}

func (s *AIService) executePassengerToolCall(
	ctx context.Context,
	currentUserID uint,
	currentUserRole string,
	lastMessage string,
	now time.Time,
	call chatCompletionToolCall,
	contextData *PassengerAIContext,
) (chatCompletionMessage, error) {
	switch call.Function.Name {
	case "search_tickets":
		var args struct {
			StartCity     string `json:"startCity"`
			EndCity       string `json:"endCity"`
			Date          string `json:"date"`
			AllowTransfer bool   `json:"allowTransfer"`
		}
		if err := json.Unmarshal([]byte(call.Function.Arguments), &args); err != nil {
			return chatCompletionMessage{}, fmt.Errorf("parse search_tickets arguments failed: %w", err)
		}

		args.StartCity = strings.TrimSpace(args.StartCity)
		args.EndCity = strings.TrimSpace(args.EndCity)
		args.Date = normalizePassengerSearchDate(strings.TrimSpace(args.Date), lastMessage, now)
		args.AllowTransfer = true
		if strings.Contains(lastMessage, "直达") || strings.Contains(lastMessage, "不要中转") {
			args.AllowTransfer = false
		}
		log.Printf("passenger ai calling SearchTickets: start=%s end=%s date=%s allowTransfer=%t", args.StartCity, args.EndCity, args.Date, args.AllowTransfer)

		dateValue, err := time.ParseInLocation("2006-01-02", args.Date, time.Local)
		if err != nil {
			output := mustMarshalToolOutput(map[string]interface{}{
				"error":   "invalid_date",
				"message": "date must use YYYY-MM-DD",
			})
			return buildToolMessage(call.ID, output), nil
		}

		trips, err := s.ticketService.SearchTickets(ctx, SearchTicketsInput{
			StartCity:     args.StartCity,
			EndCity:       args.EndCity,
			Date:          dateValue,
			AllowTransfer: args.AllowTransfer,
		})
		if err != nil {
			output := mustMarshalToolOutput(map[string]interface{}{
				"error":   "search_failed",
				"message": err.Error(),
			})
			return buildToolMessage(call.ID, output), nil
		}

		contextData.Intent = PassengerAIIntentRoute
		contextData.RouteQuery = &AIRouteQuery{
			StartCity:     args.StartCity,
			EndCity:       args.EndCity,
			Date:          args.Date,
			AllowTransfer: args.AllowTransfer,
		}
		contextData.RouteResults = mapTicketSearchResultsToAITripCards(trips)
		log.Printf("passenger ai search_tickets result_count=%d start=%s end=%s date=%s allowTransfer=%t",
			len(trips),
			args.StartCity,
			args.EndCity,
			args.Date,
			args.AllowTransfer,
		)

		output := mustMarshalToolOutput(map[string]interface{}{
			"routeQuery":   contextData.RouteQuery,
			"routeResults": contextData.RouteResults,
		})
		return buildToolMessage(call.ID, output), nil

	case "search_city_tickets":
		var args struct {
			City string `json:"city"`
			Date string `json:"date"`
			Role string `json:"role"`
		}
		if err := json.Unmarshal([]byte(call.Function.Arguments), &args); err != nil {
			return chatCompletionMessage{}, fmt.Errorf("parse search_city_tickets arguments failed: %w", err)
		}

		args.City = strings.TrimSpace(args.City)
		args.Role = strings.TrimSpace(args.Role)
		args.Date = normalizePassengerSearchDate(strings.TrimSpace(args.Date), lastMessage, now)
		log.Printf("passenger ai calling SearchCityTickets: city=%s date=%s role=%s", args.City, args.Date, args.Role)

		dateValue, err := time.ParseInLocation("2006-01-02", args.Date, time.Local)
		if err != nil {
			output := mustMarshalToolOutput(map[string]interface{}{
				"error":   "invalid_date",
				"message": "date must use YYYY-MM-DD",
			})
			return buildToolMessage(call.ID, output), nil
		}

		trips, err := s.ticketService.SearchCityTickets(ctx, SearchCityTicketsInput{
			City: args.City,
			Date: dateValue,
			Role: args.Role,
		})
		if err != nil {
			output := mustMarshalToolOutput(map[string]interface{}{
				"error":   "search_failed",
				"message": err.Error(),
			})
			return buildToolMessage(call.ID, output), nil
		}

		contextData.Intent = PassengerAIIntentRoute
		contextData.RouteQuery = &AIRouteQuery{
			StartCity: args.City,
			Date:      args.Date,
		}
		contextData.RouteResults = mapTicketSearchResultsToAITripCards(trips)
		log.Printf("passenger ai search_city_tickets result_count=%d city=%s date=%s role=%s",
			len(trips),
			args.City,
			args.Date,
			args.Role,
		)

		output := mustMarshalToolOutput(map[string]interface{}{
			"city":         args.City,
			"date":         args.Date,
			"role":         args.Role,
			"routeResults": contextData.RouteResults,
		})
		return buildToolMessage(call.ID, output), nil

	case "list_my_orders":
		if currentUserID == 0 || currentUserRole == "" {
			output := mustMarshalToolOutput(map[string]interface{}{
				"error":   "unauthorized",
				"message": "user is not logged in",
			})
			return buildToolMessage(call.ID, output), nil
		}

		orders, err := s.orderService.ListMyOrders(ctx, currentUserID, currentUserRole)
		if err != nil {
			output := mustMarshalToolOutput(map[string]interface{}{
				"error":   "list_orders_failed",
				"message": err.Error(),
			})
			return buildToolMessage(call.ID, output), nil
		}

		if contextData.Intent == PassengerAIIntentGeneral {
			contextData.Intent = PassengerAIIntentOrders
		}
		contextData.OrderResults = mapOrdersToAIOrderCards(orders)
		contextData.OrderSummary = buildAIOrderSummary(orders)
		log.Printf("passenger ai list_my_orders result_count=%d user=%d", len(orders), currentUserID)

		output := mustMarshalToolOutput(map[string]interface{}{
			"orderSummary": contextData.OrderSummary,
			"orderResults": contextData.OrderResults,
		})
		return buildToolMessage(call.ID, output), nil

	case "get_refund_rules":
		if contextData.Intent == PassengerAIIntentGeneral {
			contextData.Intent = PassengerAIIntentRefund
		}
		contextData.RefundRules = defaultRefundRules()
		log.Printf("passenger ai get_refund_rules rule_count=%d", len(contextData.RefundRules))
		output := mustMarshalToolOutput(map[string]interface{}{
			"refundRules": contextData.RefundRules,
		})
		return buildToolMessage(call.ID, output), nil

	default:
		output := mustMarshalToolOutput(map[string]interface{}{
			"error":   "unknown_tool",
			"message": call.Function.Name,
		})
		return buildToolMessage(call.ID, output), nil
	}
}

func buildPassengerPlannerMessages(messages []AIChatMessage, now time.Time) []chatCompletionMessage {
	cityToolHint := "When the user asks about one city only, including trips starting there, ending there, or passing through it, call search_city_tickets. When the user gives both start and end cities, call search_tickets and allow transfer by default."
	result := []chatCompletionMessage{
		{
			Role: "system",
			Content: strings.Join([]string{
				cityToolHint,
				"你是一个中文出行票务助手。",
				"你必须优先使用函数调用来处理用户请求。",
				"如果用户问路线、票、班次、余票，调用 search_tickets，并默认允许中转。",
				"如果用户问订单状态、我的订单，调用 list_my_orders。",
				"如果用户问退款、退票、退费，优先调用 get_refund_rules，必要时再调用 list_my_orders。",
				"如果不需要外部数据就能回答，请调用 reply_directly。",
				"不要直接输出普通文本。",
			}, "\n"),
		},
	}

	for _, item := range messages {
		result = append(result, chatCompletionMessage{
			Role:    item.Role,
			Content: item.Content,
		})
	}

	if len(result) > 0 {
		result[0].Content = fmt.Sprintf(
			"今天是 %s，请基于这个日期理解“今天、明天、后天、今年、某月某日”等时间表达。\n%s",
			now.Format("2006-01-02"),
			result[0].Content,
		)
	}

	return result
}

func buildPassengerResponderMessages(messages []AIChatMessage, toolMessages []chatCompletionMessage, now time.Time) []chatCompletionMessage {
	result := []chatCompletionMessage{
		{
			Role: "system",
			Content: strings.Join([]string{
				"你已经拿到了工具执行结果。",
				"现在不要再调用 search_tickets、list_my_orders、get_refund_rules。",
				"你必须调用 reply_directly 返回最终答案。",
				"reply 用简体中文直接回答用户，suggestions 返回简短建议数组。",
				"回答必须严格基于工具结果，不要编造不存在的数据。",
				"如果工具结果包含 direct、transfer、suggestion，多给出几个可选路线，并说明推荐理由。",
			}, "\n"),
		},
	}

	for _, item := range messages {
		result = append(result, chatCompletionMessage{
			Role:    item.Role,
			Content: item.Content,
		})
	}

	if len(result) > 0 {
		result[0].Content = fmt.Sprintf(
			"今天是 %s。最终回答里涉及日期、年份时，必须基于这个日期。\n%s",
			now.Format("2006-01-02"),
			result[0].Content,
		)
	}

	result = append(result, toolMessages...)
	return result
}

func passengerPlannerTools() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "search_tickets",
				"description": "按起点、终点和日期搜索可售班次；后端会同时生成直达、中转和模糊路线候选",
				"parameters": map[string]interface{}{
					"type":                 "object",
					"additionalProperties": false,
					"properties": map[string]interface{}{
						"startCity": map[string]interface{}{
							"type":        "string",
							"description": "起点城市或站点",
						},
						"endCity": map[string]interface{}{
							"type":        "string",
							"description": "终点城市或站点",
						},
						"date": map[string]interface{}{
							"type":        "string",
							"description": "出发日期，格式 YYYY-MM-DD",
						},
					},
					"required": []string{"startCity", "endCity", "date"},
				},
			},
		},
		{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "search_city_tickets",
				"description": "Search trips related to one city. The city can be the trip start, trip end, or any stop/via city.",
				"parameters": map[string]interface{}{
					"type":                 "object",
					"additionalProperties": false,
					"properties": map[string]interface{}{
						"city": map[string]interface{}{
							"type":        "string",
							"description": "City or stop name to search.",
						},
						"date": map[string]interface{}{
							"type":        "string",
							"description": "Departure date in YYYY-MM-DD.",
						},
						"role": map[string]interface{}{
							"type":        "string",
							"description": "Optional: any, start, end, or stop.",
							"enum":        []string{"any", "start", "end", "stop"},
						},
					},
					"required": []string{"city", "date"},
				},
			},
		},
		{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "list_my_orders",
				"description": "查询当前登录用户的订单列表和状态摘要",
				"parameters": map[string]interface{}{
					"type":                 "object",
					"additionalProperties": false,
					"properties":           map[string]interface{}{},
				},
			},
		},
		{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "get_refund_rules",
				"description": "获取平台退款规则摘要",
				"parameters": map[string]interface{}{
					"type":                 "object",
					"additionalProperties": false,
					"properties":           map[string]interface{}{},
				},
			},
		},
		{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "reply_directly",
				"description": "当不需要外部服务时，直接返回最终回答",
				"parameters": map[string]interface{}{
					"type":                 "object",
					"additionalProperties": false,
					"properties": map[string]interface{}{
						"reply": map[string]interface{}{
							"type":        "string",
							"description": "直接回复用户的自然语言答案",
						},
						"suggestions": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"type": "string",
							},
						},
					},
					"required": []string{"reply", "suggestions"},
				},
			},
		},
	}
}

func passengerReplyTools() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "reply_directly",
				"description": "返回最终回答",
				"parameters": map[string]interface{}{
					"type":                 "object",
					"additionalProperties": false,
					"properties": map[string]interface{}{
						"reply": map[string]interface{}{
							"type": "string",
						},
						"suggestions": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"type": "string",
							},
						},
					},
					"required": []string{"reply", "suggestions"},
				},
			},
		},
	}
}

func driverDraftTools() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "build_driver_trip_draft",
				"description": "根据司机自然语言描述生成结构化班次草稿",
				"parameters": map[string]interface{}{
					"type":                 "object",
					"additionalProperties": false,
					"properties": map[string]interface{}{
						"tripName":           map[string]interface{}{"type": "string"},
						"startCity":          map[string]interface{}{"type": "string"},
						"endCity":            map[string]interface{}{"type": "string"},
						"departureTimeLocal": map[string]interface{}{"type": "string"},
						"arrivalTimeLocal":   map[string]interface{}{"type": "string"},
						"seatTotal":          map[string]interface{}{"type": "integer", "minimum": 1},
						"priceCent":          map[string]interface{}{"type": "integer", "minimum": 1},
						"vehicleType": map[string]interface{}{
							"type": "string",
							"enum": []string{"商务大巴", "城际快线", "拼车专线"},
						},
						"stops": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"type": "string",
							},
						},
						"remark": map[string]interface{}{"type": "string"},
						"suggestions": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"type": "string",
							},
						},
					},
					"required": []string{
						"tripName",
						"startCity",
						"endCity",
						"departureTimeLocal",
						"arrivalTimeLocal",
						"seatTotal",
						"priceCent",
						"vehicleType",
						"stops",
						"remark",
						"suggestions",
					},
				},
			},
		},
	}
}

func forceToolChoice(name string) map[string]interface{} {
	return map[string]interface{}{
		"type": "function",
		"function": map[string]interface{}{
			"name": name,
		},
	}
}

func buildToolMessage(toolCallID, content string) chatCompletionMessage {
	return chatCompletionMessage{
		Role:       "tool",
		ToolCallID: toolCallID,
		Content:    content,
	}
}

func mustMarshalToolOutput(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		return `{"error":"marshal_tool_output_failed"}`
	}
	return string(data)
}

type chatCompletionRequest struct {
	Model          string                   `json:"model"`
	Messages       []chatCompletionMessage  `json:"messages"`
	Tools          []map[string]interface{} `json:"tools,omitempty"`
	ToolChoice     interface{}              `json:"tool_choice,omitempty"`
	EnableThinking bool                     `json:"enable_thinking"`
	Temperature    float64                  `json:"temperature,omitempty"`
}

type chatCompletionMessage struct {
	Role       string                   `json:"role"`
	Content    string                   `json:"content,omitempty"`
	ToolCalls  []chatCompletionToolCall `json:"tool_calls,omitempty"`
	ToolCallID string                   `json:"tool_call_id,omitempty"`
}

type chatCompletionToolCall struct {
	ID       string                     `json:"id"`
	Type     string                     `json:"type"`
	Function chatCompletionToolFunction `json:"function"`
}

type chatCompletionToolFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type aiUsageContext struct {
	UserID  uint
	Role    string
	Feature string
}

type chatCompletionResponse struct {
	Choices []chatCompletionChoice `json:"choices"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type chatCompletionChoice struct {
	Message      chatCompletionMessage `json:"message"`
	FinishReason string                `json:"finish_reason"`
}

func (r *chatCompletionResponse) firstMessage() (*chatCompletionMessage, error) {
	if len(r.Choices) == 0 {
		return nil, errors.New("model returned no choices")
	}
	return &r.Choices[0].Message, nil
}

func (s *AIService) requestChatCompletion(ctx context.Context, usageCtx aiUsageContext, requestBody chatCompletionRequest) (*chatCompletionResponse, error) {
	if requestBody.Model == "" {
		requestBody.Model = s.model
	}

	payload, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, parseOpenAIHTTPError(resp.StatusCode, body)
	}

	var result chatCompletionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	if s.tokenUsageService != nil {
		_ = s.tokenUsageService.Record(ctx, RecordTokenUsageInput{
			UserID:           usageCtx.UserID,
			Role:             usageCtx.Role,
			Feature:          usageCtx.Feature,
			RequestKind:      model.TokenUsageKindChat,
			Provider:         "openai-compatible",
			Model:            requestBody.Model,
			PromptTokens:     result.Usage.PromptTokens,
			CompletionTokens: result.Usage.CompletionTokens,
			TotalTokens:      result.Usage.TotalTokens,
		})
	}
	return &result, nil
}

func firstToolCallByName(toolCalls []chatCompletionToolCall, name string) (*chatCompletionToolCall, error) {
	for i := range toolCalls {
		if toolCalls[i].Function.Name == name {
			return &toolCalls[i], nil
		}
	}
	return nil, fmt.Errorf("tool call %s not found", name)
}

func parseOpenAIHTTPError(statusCode int, body []byte) error {
	var envelope struct {
		Error struct {
			Message string `json:"message"`
			Type    string `json:"type"`
			Code    string `json:"code"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &envelope); err == nil && strings.TrimSpace(envelope.Error.Message) != "" {
		return fmt.Errorf("openai api error (%d): %s", statusCode, strings.TrimSpace(envelope.Error.Message))
	}

	message := strings.TrimSpace(string(body))
	if message == "" {
		message = http.StatusText(statusCode)
	}
	return fmt.Errorf("openai api error (%d): %s", statusCode, message)
}

func parseDriverTripAIDraft(text string) (*DriverTripAIDraft, error) {
	type driverTripAIDraftAlias struct {
		Prompt             string          `json:"prompt"`
		TripName           string          `json:"tripName"`
		StartCity          string          `json:"startCity"`
		EndCity            string          `json:"endCity"`
		DepartureTimeLocal string          `json:"departureTimeLocal"`
		ArrivalTimeLocal   string          `json:"arrivalTimeLocal"`
		SeatTotal          int             `json:"seatTotal"`
		PriceCent          int             `json:"priceCent"`
		VehicleType        string          `json:"vehicleType"`
		Stops              json.RawMessage `json:"stops"`
		Remark             string          `json:"remark"`
		Suggestions        json.RawMessage `json:"suggestions"`
		Origin             string          `json:"origin"`
		Destination        string          `json:"destination"`
		DepartureTime      string          `json:"departureTime"`
		ArrivalTime        string          `json:"arrivalTime"`
		SeatCount          int             `json:"seatCount"`
		Seats              int             `json:"seats"`
		Response           string          `json:"response"`
	}

	var raw driverTripAIDraftAlias
	if err := json.Unmarshal([]byte(strings.TrimSpace(text)), &raw); err != nil {
		return nil, err
	}

	stops, err := parseJSONStringList(raw.Stops)
	if err != nil {
		return nil, fmt.Errorf("parse stops failed: %w", err)
	}
	suggestions, err := parseJSONStringList(raw.Suggestions)
	if err != nil {
		return nil, fmt.Errorf("parse suggestions failed: %w", err)
	}

	return &DriverTripAIDraft{
		Prompt:             raw.Prompt,
		TripName:           raw.TripName,
		StartCity:          raw.StartCity,
		EndCity:            raw.EndCity,
		DepartureTimeLocal: raw.DepartureTimeLocal,
		ArrivalTimeLocal:   raw.ArrivalTimeLocal,
		SeatTotal:          raw.SeatTotal,
		PriceCent:          raw.PriceCent,
		VehicleType:        raw.VehicleType,
		Stops:              stops,
		Remark:             raw.Remark,
		Suggestions:        suggestions,
		Origin:             raw.Origin,
		Destination:        raw.Destination,
		DepartureTime:      raw.DepartureTime,
		ArrivalTime:        raw.ArrivalTime,
		SeatCount:          raw.SeatCount,
		Seats:              raw.Seats,
		Response:           raw.Response,
	}, nil
}

func (s *AIService) chatPassengerWithRAG(
	ctx context.Context,
	currentUserID uint,
	currentUserRole string,
	messages []AIChatMessage,
	ragResults []*KnowledgeSearchResult,
) (*PassengerAIChatResult, error) {
	knowledgeContext := s.knowledgeService.BuildPromptContext(ragResults)
	contextData := &PassengerAIContext{
		Intent:           PassengerAIIntentGeneral,
		KnowledgeSources: mapKnowledgeSources(ragResults),
	}

	promptMessages := []chatCompletionMessage{
		{
			Role: "system",
			Content: strings.Join([]string{
				"You are a Chinese customer-support assistant for an intercity travel platform.",
				"Answer strictly based on the retrieved knowledge snippets.",
				"If the snippets do not clearly answer the question, say the current knowledge base does not cover the answer and do not invent rules.",
				"Return the final answer by calling reply_directly.",
				knowledgeContext,
			}, "\n\n"),
		},
	}

	for _, item := range messages {
		promptMessages = append(promptMessages, chatCompletionMessage{
			Role:    item.Role,
			Content: item.Content,
		})
	}

	response, err := s.requestChatCompletion(ctx, aiUsageContext{
		UserID:  currentUserID,
		Role:    currentUserRole,
		Feature: model.TokenFeaturePassengerAI,
	}, chatCompletionRequest{
		Model:          s.model,
		Messages:       promptMessages,
		Tools:          passengerReplyTools(),
		ToolChoice:     forceToolChoice("reply_directly"),
		EnableThinking: false,
	})
	if err != nil {
		return nil, err
	}

	message, err := response.firstMessage()
	if err != nil {
		return nil, err
	}

	replyCall, err := firstToolCallByName(message.ToolCalls, "reply_directly")
	if err != nil {
		result, parseErr := parsePassengerAIChatResult(message.Content)
		if parseErr != nil {
			return nil, err
		}
		if result.Intent == "" {
			result.Intent = PassengerAIIntentGeneral
		}
		if len(result.Suggestions) == 0 {
			result.Suggestions = []string{
				"这个规则适用于哪些用户",
				"还有没有例外情况",
				"帮我解释得更简单一点",
			}
		}
		result.Context = contextData
		return result, nil
	}

	result, err := parsePassengerReplyDirectlyArguments(replyCall.Function.Arguments)
	if err != nil {
		return nil, fmt.Errorf("parse rag reply arguments failed: %w", err)
	}
	if result.Intent == "" {
		result.Intent = PassengerAIIntentGeneral
	}
	result.Context = contextData
	return result, nil
}

func shouldUseRAG(query string) bool {
	query = strings.TrimSpace(query)
	if query == "" {
		return false
	}

	liveDataKeywords := []string{
		"我的订单", "我的票", "我这笔订单", "这个订单", "帮我看订单",
		"从", "到", "班次", "车次", "有没有票", "明天", "后天",
	}
	for _, keyword := range liveDataKeywords {
		if strings.Contains(query, keyword) {
			return false
		}
	}

	keywords := []string{
		"规则", "规定", "资格", "实名", "实名认证", "购票条件", "可以买吗", "能买吗",
		"网站", "站内", "平台", "退款", "退票", "禁用", "冻结", "账号状态",
	}
	for _, keyword := range keywords {
		if strings.Contains(query, keyword) {
			return true
		}
	}
	return false
}

func inferKnowledgeCategory(query string) string {
	query = strings.TrimSpace(query)

	switch {
	case strings.Contains(query, "实名"), strings.Contains(query, "资格"), strings.Contains(query, "购票条件"), strings.Contains(query, "未实名"):
		return "ticket-rule"
	case strings.Contains(query, "退款"), strings.Contains(query, "退票"), strings.Contains(query, "改签"):
		return "refund-rule"
	case strings.Contains(query, "账号"), strings.Contains(query, "冻结"), strings.Contains(query, "禁用"), strings.Contains(query, "风控"):
		return "account-rule"
	default:
		return ""
	}
}

func mapKnowledgeSources(results []*KnowledgeSearchResult) []*AIKnowledgeSource {
	sources := make([]*AIKnowledgeSource, 0, len(results))
	for _, item := range results {
		if item == nil {
			continue
		}
		sources = append(sources, &AIKnowledgeSource{
			DocumentID:  item.DocumentID,
			ChunkID:     item.ChunkID,
			Title:       item.Title,
			SectionPath: item.SectionPath,
			Content:     item.Content,
			FinalScore:  fmt.Sprintf("%.4f", item.FinalScore),
		})
	}
	return sources
}

func parsePassengerReplyDirectlyArguments(text string) (*PassengerAIChatResult, error) {
	type alias struct {
		Reply       string          `json:"reply"`
		Response    string          `json:"response"`
		Suggestions json.RawMessage `json:"suggestions"`
	}

	var raw alias
	if err := json.Unmarshal([]byte(strings.TrimSpace(text)), &raw); err != nil {
		return nil, err
	}

	suggestions, err := parseJSONStringList(raw.Suggestions)
	if err != nil {
		return nil, fmt.Errorf("parse suggestions failed: %w", err)
	}

	reply := strings.TrimSpace(raw.Reply)
	response := strings.TrimSpace(raw.Response)
	if reply == "" && response != "" {
		reply = response
	}

	if len(suggestions) == 0 {
		suggestions = []string{
			"帮我查明天杭州到苏州的票",
			"帮我看下我的订单",
			"我这笔订单能退款吗",
		}
	}

	return &PassengerAIChatResult{
		Reply:       reply,
		Response:    response,
		Suggestions: cleanStringList(suggestions),
	}, nil
}

func parsePassengerAIChatResult(text string) (*PassengerAIChatResult, error) {
	text = strings.TrimSpace(text)
	if !looksLikeJSONObject(text) {
		return &PassengerAIChatResult{
			Reply: text,
		}, nil
	}

	var raw struct {
		Reply       string          `json:"reply"`
		Response    string          `json:"response"`
		Suggestions json.RawMessage `json:"suggestions"`
	}
	if err := json.Unmarshal([]byte(text), &raw); err != nil {
		extracted := extractFirstJSONObject(text)
		if extracted == "" {
			return nil, err
		}
		if retryErr := json.Unmarshal([]byte(extracted), &raw); retryErr != nil {
			return nil, err
		}
	}

	suggestions, err := parseJSONStringList(raw.Suggestions)
	if err != nil {
		return nil, fmt.Errorf("parse suggestions failed: %w", err)
	}

	reply := strings.TrimSpace(raw.Reply)
	response := strings.TrimSpace(raw.Response)
	if reply == "" && response != "" {
		reply = response
	}

	if len(suggestions) == 0 {
		suggestions = []string{
			"帮我查明天杭州到苏州的票",
			"帮我看下我的订单",
			"我这笔订单能退款吗",
		}
	}

	return &PassengerAIChatResult{
		Reply:       reply,
		Response:    response,
		Suggestions: cleanStringList(suggestions),
	}, nil
}

func looksLikeJSONObject(text string) bool {
	text = strings.TrimSpace(text)
	return strings.HasPrefix(text, "{") && strings.HasSuffix(text, "}")
}

func extractFirstJSONObject(text string) string {
	start := strings.Index(text, "{")
	if start < 0 {
		return ""
	}

	depth := 0
	for i := start; i < len(text); i++ {
		switch text[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return text[start : i+1]
			}
		}
	}
	return ""
}

func parseJSONStringList(raw json.RawMessage) ([]string, error) {
	if len(raw) == 0 {
		return nil, nil
	}

	var list []string
	if err := json.Unmarshal(raw, &list); err == nil {
		return cleanStringList(list), nil
	}

	var single string
	if err := json.Unmarshal(raw, &single); err == nil {
		single = strings.TrimSpace(single)
		if single == "" {
			return nil, nil
		}

		if strings.Contains(single, "、") || strings.Contains(single, "；") || strings.Contains(single, ",") {
			parts := strings.FieldsFunc(single, func(r rune) bool {
				return r == '、' || r == '；' || r == ','
			})
			return cleanStringList(parts), nil
		}

		return []string{single}, nil
	}

	return nil, errors.New("expected string array or string")
}

func normalizeDriverTripDraft(draft *DriverTripAIDraft) {
	draft.TripName = strings.TrimSpace(draft.TripName)
	draft.StartCity = strings.TrimSpace(draft.StartCity)
	draft.EndCity = strings.TrimSpace(draft.EndCity)
	draft.DepartureTimeLocal = strings.TrimSpace(draft.DepartureTimeLocal)
	draft.ArrivalTimeLocal = strings.TrimSpace(draft.ArrivalTimeLocal)
	draft.VehicleType = normalizeVehicleType(draft.VehicleType)
	draft.Remark = strings.TrimSpace(draft.Remark)
	draft.Stops = cleanStringList(draft.Stops)
	draft.Suggestions = cleanStringList(draft.Suggestions)
}

func reconcileDriverTripDraft(draft *DriverTripAIDraft) {
	draft.Origin = strings.TrimSpace(draft.Origin)
	draft.Destination = strings.TrimSpace(draft.Destination)
	draft.DepartureTime = strings.TrimSpace(draft.DepartureTime)
	draft.ArrivalTime = strings.TrimSpace(draft.ArrivalTime)
	draft.Response = strings.TrimSpace(draft.Response)

	if draft.StartCity == "" && draft.Origin != "" {
		draft.StartCity = draft.Origin
	}
	if draft.EndCity == "" && draft.Destination != "" {
		draft.EndCity = draft.Destination
	}
	if draft.DepartureTimeLocal == "" && draft.DepartureTime != "" {
		draft.DepartureTimeLocal = draft.DepartureTime
	}
	if draft.ArrivalTimeLocal == "" && draft.ArrivalTime != "" {
		draft.ArrivalTimeLocal = draft.ArrivalTime
	}
	if draft.SeatTotal <= 0 && draft.SeatCount > 0 {
		draft.SeatTotal = draft.SeatCount
	}
	if draft.SeatTotal <= 0 && draft.Seats > 0 {
		draft.SeatTotal = draft.Seats
	}
	if draft.TripName == "" && draft.StartCity != "" && draft.EndCity != "" {
		draft.TripName = draft.StartCity + " - " + draft.EndCity + " 班次"
	}
	if draft.Remark == "" && draft.Response != "" {
		draft.Remark = draft.Response
	}
}

func validateDriverTripDraft(draft *DriverTripAIDraft) error {
	if draft.TripName == "" || draft.StartCity == "" || draft.EndCity == "" {
		return errors.New("model output is missing required route fields")
	}
	if draft.StartCity == draft.EndCity {
		return errors.New("model output has the same start and end city")
	}

	departureTime, err := time.ParseInLocation("2006-01-02T15:04", draft.DepartureTimeLocal, time.Local)
	if err != nil {
		return errors.New("model output departureTimeLocal is invalid")
	}
	arrivalTime, err := time.ParseInLocation("2006-01-02T15:04", draft.ArrivalTimeLocal, time.Local)
	if err != nil {
		return errors.New("model output arrivalTimeLocal is invalid")
	}
	if !arrivalTime.After(departureTime) {
		return errors.New("model output arrivalTimeLocal must be later than departureTimeLocal")
	}
	if draft.SeatTotal <= 0 {
		return errors.New("model output seatTotal must be greater than 0")
	}
	if draft.PriceCent <= 0 {
		return errors.New("model output priceCent must be greater than 0")
	}
	if draft.VehicleType == "" {
		return errors.New("model output vehicleType is required")
	}
	if draft.Remark == "" {
		draft.Remark = "请在发布前补充上车点和行李说明。"
	}
	return nil
}

func normalizeVehicleType(raw string) string {
	raw = strings.TrimSpace(raw)
	switch raw {
	case "商务大巴", "城际快线", "拼车专线":
		return raw
	}

	switch {
	case strings.Contains(raw, "拼车"):
		return "拼车专线"
	case strings.Contains(raw, "快线"):
		return "城际快线"
	default:
		return "商务大巴"
	}
}

func normalizeChatMessages(messages []AIChatMessage) []AIChatMessage {
	result := make([]AIChatMessage, 0, len(messages))

	for _, msg := range messages {
		role := strings.ToLower(strings.TrimSpace(msg.Role))
		content := strings.TrimSpace(msg.Content)
		if role != "user" && role != "assistant" {
			continue
		}
		if content == "" {
			continue
		}

		runes := []rune(content)
		if len(runes) > maxPassengerAIMessageLength {
			content = string(runes[:maxPassengerAIMessageLength])
		}

		result = append(result, AIChatMessage{
			Role:    role,
			Content: content,
		})
	}

	if len(result) > maxPassengerAIMessages {
		result = result[len(result)-maxPassengerAIMessages:]
	}

	return result
}

func mapTicketSearchResultsToAITripCards(results []*TicketSearchResult) []*AITripCard {
	items := make([]*AITripCard, 0, len(results))
	for _, result := range results {
		if result == nil {
			continue
		}

		legs := make([]*AITripLeg, 0, len(result.Legs))
		for _, leg := range result.Legs {
			if leg == nil {
				continue
			}
			legs = append(legs, &AITripLeg{
				TripID:        leg.TripID,
				Route:         leg.StartCity + " -> " + leg.EndCity,
				DepartureTime: leg.DepartureTime.Format("2006-01-02 15:04"),
				ArrivalTime:   leg.ArrivalTime.Format("2006-01-02 15:04"),
				SeatAvailable: leg.SeatAvailable,
				PriceCent:     leg.PriceCent,
				VehicleType:   leg.VehicleType,
			})
		}

		suggestions := make([]*AIRouteSuggestion, 0, len(result.Suggestions))
		for _, suggestion := range result.Suggestions {
			if suggestion == nil {
				continue
			}
			suggestions = append(suggestions, &AIRouteSuggestion{
				Route:            suggestion.Route,
				TransferCity:     suggestion.TransferCity,
				FirstLegCount:    suggestion.FirstLegCount,
				SecondLegCount:   suggestion.SecondLegCount,
				TotalOptionCount: suggestion.TotalOptionCount,
				Reason:           suggestion.Reason,
			})
		}

		items = append(items, &AITripCard{
			ID:                 result.ID,
			Kind:               result.Kind,
			Route:              result.StartCity + " -> " + result.EndCity,
			DepartureTime:      result.DepartureTime.Format("2006-01-02 15:04"),
			ArrivalTime:        result.ArrivalTime.Format("2006-01-02 15:04"),
			SeatAvailable:      result.SeatAvailable,
			PriceCent:          result.PriceCent,
			VehicleType:        result.VehicleType,
			TransferCity:       result.TransferCity,
			TransferWaitMinute: result.TransferWaitMinute,
			MatchedCity:        result.MatchedCity,
			MatchRoles:         result.MatchRoles,
			Stops:              result.Stops,
			Legs:               legs,
			Suggestions:        suggestions,
		})
	}
	return items
}

func mapOrdersToAIOrderCards(orders []*model.Order) []*AIOrderCard {
	result := make([]*AIOrderCard, 0, len(orders))
	for _, order := range orders {
		if order == nil {
			continue
		}

		route := ""
		departureTime := ""
		if order.Trip.ID != 0 {
			route = order.Trip.StartCity + " -> " + order.Trip.EndCity
			departureTime = order.Trip.DepartureTime.Format("2006-01-02 15:04")
		}

		result = append(result, &AIOrderCard{
			ID:               order.ID,
			OrderNo:          order.OrderNo,
			Route:            route,
			DepartureTime:    departureTime,
			OrderStatus:      order.OrderStatus,
			PayStatus:        order.PayStatus,
			RefundStatus:     order.RefundStatus,
			RefundReviewNote: strings.TrimSpace(order.RefundReviewNote),
			Amount:           order.Amount,
		})
	}
	return result
}

func buildAIOrderSummary(orders []*model.Order) *AIOrderSummary {
	summary := &AIOrderSummary{
		TotalCount: len(orders),
	}

	for _, order := range orders {
		if order == nil {
			continue
		}

		switch order.OrderStatus {
		case model.OrderStatusPendingPayment:
			summary.PendingPaymentCount++
		case model.OrderStatusPendingVerification:
			summary.PendingVerificationCount++
		case model.OrderStatusCompleted:
			summary.CompletedCount++
		}

		switch order.RefundStatus {
		case model.RefundStatusRequested:
			summary.RefundRequestedCount++
		case model.RefundStatusRejected:
			summary.RefundRejectedCount++
		case model.RefundStatusRefunded:
			summary.RefundedCount++
		}
	}

	return summary
}

func defaultRefundRules() []string {
	return []string{
		"只有已支付订单才能申请退款。",
		"待核销和已完成订单支持申请退款。",
		"已取消订单不支持再次申请退款。",
		"退款申请提交后会进入审核流程。",
		"如果退款已申请或已退款，不能重复提交。",
	}
}

func cleanStringList(items []string) []string {
	result := make([]string, 0, len(items))
	seen := make(map[string]struct{}, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		result = append(result, item)
	}
	return result
}

func lastPassengerMessage(messages []AIChatMessage) string {
	if len(messages) == 0 {
		return ""
	}
	return strings.TrimSpace(messages[len(messages)-1].Content)
}

func buildDirectCurrentDateReply(text string, now time.Time) (*PassengerAIChatResult, bool) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, false
	}

	switch {
	case containsAny(text, "你今年是什么年", "今年是什么年", "今年是哪年", "现在是什么年", "现在是哪年"):
		return &PassengerAIChatResult{
			Reply: fmt.Sprintf("今年是%d年。", now.Year()),
			Suggestions: []string{
				"你也可以继续问今天的日期",
				"如果你要查票，可以直接说起点、终点和日期",
			},
		}, true
	case containsAny(text, "今天几号", "今天是什么日期", "今天多少号", "今天是哪天", "今天是几月几号"):
		return &PassengerAIChatResult{
			Reply: fmt.Sprintf("今天是%s。", now.Format("2006年1月2日")),
			Suggestions: []string{
				"你可以继续问明天或后天的票",
				"你也可以直接说“帮我查杭州到苏州明天的票”",
			},
		}, true
	case containsAny(text, "现在几点", "当前时间", "现在时间"):
		return &PassengerAIChatResult{
			Reply: fmt.Sprintf("现在是%s。", now.Format("2006年1月2日 15:04")),
			Suggestions: []string{
				"你可以继续问具体日期的班次",
			},
		}, true
	default:
		return nil, false
	}
}

func normalizePassengerSearchDate(rawDate string, referenceText string, now time.Time) string {
	rawDate = strings.TrimSpace(rawDate)
	referenceText = strings.TrimSpace(referenceText)

	switch {
	case containsAny(referenceText, "后天"):
		return now.AddDate(0, 0, 2).Format("2006-01-02")
	case containsAny(referenceText, "明天", "明早", "明晚"):
		return now.AddDate(0, 0, 1).Format("2006-01-02")
	case containsAny(referenceText, "今天", "今日"):
		return now.Format("2006-01-02")
	}

	if matches := monthDayPattern.FindStringSubmatch(referenceText); len(matches) == 3 {
		month, monthErr := strconv.Atoi(matches[1])
		day, dayErr := strconv.Atoi(matches[2])
		if monthErr == nil && dayErr == nil && month >= 1 && month <= 12 && day >= 1 && day <= 31 {
			return time.Date(now.Year(), time.Month(month), day, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
		}
	}

	if rawDate == "" {
		return now.Format("2006-01-02")
	}

	if parsed, err := time.ParseInLocation("2006-01-02", rawDate, now.Location()); err == nil {
		return parsed.Format("2006-01-02")
	}

	return now.Format("2006-01-02")
}

func containsAny(text string, keywords ...string) bool {
	for _, keyword := range keywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}
	return false
}
