package main

import {
	"fmt"
	"os"	
}

//устанавливаем переменные окружения

func set_environment() {
	os.Setenv("RABBITMQ_USER", "guest")
	os.Setenv("RABBITMQ_PASS", "guest")
	os.Setenv("RABBITMQ_ADDR", "localhost:5672/")
	os.Setenv("SMTP_PORT", "587")
	os.Setenv("SMTP_HOST", "smtp.mail.ru")
	os.Setenv("EMAIL_LOGIN", "Olegovich99@inbox.ru")
	os.Setenv("EMAIL_PASSWORD", "58PGx1zk5Zxgb4WPwq5i")
}