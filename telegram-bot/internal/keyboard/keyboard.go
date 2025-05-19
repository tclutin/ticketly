package keyboard

import "gopkg.in/telebot.v4"

func CreateMainMenu() *telebot.ReplyMarkup {
	return &telebot.ReplyMarkup{
		ResizeKeyboard: true,
		ReplyKeyboard: [][]telebot.ReplyButton{
			{
				{Text: "📨 Создать тикет"},
			},
			{
				{Text: "📂 Мои обращения"},
			},
		},
	}
}

func CreateTicketTypeMenu() *telebot.ReplyMarkup {
	return &telebot.ReplyMarkup{
		ResizeKeyboard: true,
		ReplyKeyboard: [][]telebot.ReplyButton{
			{
				{Text: "💬 Realtime-chat"},
				{Text: "📝 Only-one-message"},
			},
			{
				{Text: "🔙 Назад"},
			},
		},
	}
}

func CreateConfirmMenu() *telebot.ReplyMarkup {
	return &telebot.ReplyMarkup{
		ResizeKeyboard: true,
		ReplyKeyboard: [][]telebot.ReplyButton{
			{
				{Text: "✅ Да"},
				{Text: "❌ Нет"},
			},
		},
	}
}

func CreateCancelMenu() *telebot.ReplyMarkup {
	return &telebot.ReplyMarkup{
		ResizeKeyboard: true,
		ReplyKeyboard: [][]telebot.ReplyButton{
			{
				{Text: "❌ Отмена"},
			},
		},
	}
}
