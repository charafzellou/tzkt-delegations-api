package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var (
	// Define a Fiber API config
	config = fiber.Config{
		Prefork:       false,
		CaseSensitive: false,
		StrictRouting: false,
		ServerHeader:  "",
		AppName:       "",
	}

	// Define a Fiber Recoverer config
	recover_config = recover.Config{
		Next:              nil,
		EnableStackTrace:  false,
		StackTraceHandler: recover.ConfigDefault.StackTraceHandler,
	}

	// Define a Fiber Logger config
	logger_config = logger.Config{
		Format:     "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: "01-Jan-2000",
		TimeZone:   "Europe/Paris",
		Output:     setLogFile(),
	}
)
