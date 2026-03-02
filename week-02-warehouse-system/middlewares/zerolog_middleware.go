package middlewares

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/config"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

func maskSensitiveData(body []byte) map[string]interface{} {
	var payload map[string]interface{}
	if len(payload) == 0 {
		return payload
	}

	if err := json.Unmarshal(body, &payload); err == nil {
		if _, exists := payload["password"]; exists {
			payload["password"] = "********"
		}
		return payload
	}

	return nil
}

func ActivityLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		reqBody := c.Body()
		safeReqPayload := maskSensitiveData(reqBody)

		queryParams := c.Queries()

		err := c.Next()

		duration := time.Since(start)

		var userID interface{} = "guest"
		var role interface{} = "unauthenticated"

		if uid := c.Locals("user_id"); uid != nil {
			userID = uid
		}
		if r := c.Locals("role"); r != nil {
			role = r
		}

		resBody := c.Response().Body()
		var resPayload interface{}

		contentType := string(c.Response().Header.Peek("Content-Type"))
		if strings.Contains(contentType, "application/json") {
			var parsedRes map[string]interface{}
			if json.Unmarshal(resBody, &parsedRes) == nil {
				resPayload = parsedRes
			}
		} else if strings.Contains(contentType, "text/csv") {
			resPayload = "[FILE_CSV_DOWNLOADED]"
		} else {
			resPayload = "[NON_JSON_RESPONSE]"
		}

		statusCode := c.Response().StatusCode()
		var loggerEvent *zerolog.Event

		if err != nil {
			loggerEvent = config.Log.Error().Err(err)
		} else if statusCode >= 400 && statusCode < 500 {
			loggerEvent = config.Log.Warn()
		} else if statusCode >= 500 {
			loggerEvent = config.Log.Error()
		} else {
			loggerEvent = config.Log.Info()
		}

		// loggerEvent.Str("type", "http_access").Str("method", c.Method()).
		// 	Str("path", c.Path()).Int("status", statusCode).Str("ip_address", c.IP()).Interface("user_id", userID).
		// 	Interface("user_role", role).Str("latency", duration.String()).Str("user_agent", c.Get("User-Agent")).
		// 	Msg("Aktivitas Pengguna")

		loggerEvent.
			Str("type", "advanced_http_access").
			Str("method", c.Method()).
			Str("path", c.Path()).
			Int("status", statusCode).
			Str("ip_address", c.IP()).
			Interface("user_id", userID).
			Interface("user_role", role).
			Str("latency", duration.String()).
			Str("user_agent", c.Get("User-Agent")).
			Interface("query_params", queryParams).       // Log Query (Pencarian/Filter)
			Interface("request_payload", safeReqPayload). // Log apa yang dikirim user
			Interface("response_payload", resPayload).    // Log apa yang dibalas server
			Msg("Aktivitas Pengguna Terekam")

		return err
	}
}
