package keyboard

import "gopkg.in/telebot.v4"

const (
	BtnCreateTicket       = "📨 Создать тикет"
	BtnTicketTypeRealtime = "💬 Чат с поддержкой"
	BtnTicketTypeSingle   = "📝 Отправить запрос"
	BtnBack               = "🔙 Назад"
	BtnConfirmYes         = "✅ Да"
	BtnConfirmNo          = "❌ Нет"
	BtnCancel             = "❌ Отмена"
)

func CreateMainMenu() *telebot.ReplyMarkup {
	return &telebot.ReplyMarkup{
		ResizeKeyboard: true,
		ReplyKeyboard: [][]telebot.ReplyButton{
			{{Text: BtnCreateTicket}},
		},
	}
}

func CreateTicketTypeMenu() *telebot.ReplyMarkup {
	return &telebot.ReplyMarkup{
		ResizeKeyboard: true,
		ReplyKeyboard: [][]telebot.ReplyButton{
			{
				{Text: BtnTicketTypeRealtime},
				{Text: BtnTicketTypeSingle},
			},
			{{Text: BtnBack}},
		},
	}
}

func CreateConfirmMenu() *telebot.ReplyMarkup {
	return &telebot.ReplyMarkup{
		ResizeKeyboard: true,
		ReplyKeyboard: [][]telebot.ReplyButton{
			{
				{Text: BtnConfirmYes},
				{Text: BtnConfirmNo},
			},
		},
	}
}

func CreateCancelMenu() *telebot.ReplyMarkup {
	return &telebot.ReplyMarkup{
		ResizeKeyboard: true,
		ReplyKeyboard: [][]telebot.ReplyButton{
			{{Text: BtnCancel}},
		},
	}
}
