package handler

import (
	"context"
	"fmt"
	fsmstate "github.com/tclutin/ticketly/telegram_bot/internal/fsm"
	"github.com/tclutin/ticketly/telegram_bot/internal/keyboard"
	"github.com/tclutin/ticketly/telegram_bot/internal/message"
	"github.com/tclutin/ticketly/telegram_bot/internal/service"
	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	"github.com/vitaliy-ukiru/fsm-telebot/v2/fsmopt"
	"gopkg.in/telebot.v4"
	"log/slog"
)

type Handler struct {
	dp  fsm.Dispatcher
	mn  *fsm.Manager
	bot *telebot.Bot
	srv service.Ticketly
}

func NewHandler(dp fsm.Dispatcher, mn *fsm.Manager, bot *telebot.Bot, srv service.Ticketly) *Handler {
	return &Handler{
		dp:  dp,
		mn:  mn,
		bot: bot,
		srv: srv,
	}
}

func (h *Handler) Register() {
	// only one message
	h.dp.Dispatch(
		h.mn.New(
			fsmopt.On(telebot.OnText),
			fsmopt.OnStates(fsmstate.Content),
			fsmopt.Do(h.contentTicketFSM),
		),
	)

	h.dp.Dispatch(
		h.mn.New(
			fsmopt.On(telebot.OnText),
			fsmopt.OnStates(fsmstate.Confirm),
			fsmopt.Do(h.confirmTicketFSM),
		),
	)

	h.dp.Dispatch(
		h.mn.New(
			fsmopt.On(telebot.OnText),
			fsmopt.OnStates(fsm.AnyState),
			fsmopt.Do(h.onTextFSM),
		),
	)

	h.bot.Handle("/start", h.Start())

	h.bot.Handle("📨 Создать тикет", h.CreateMenuTicket())

	h.bot.Handle("🔙 Назад", h.BackToMainMenu())

	h.bot.Handle("❌ Отмена", h.CancelOperation())
}

func (h *Handler) Start() telebot.HandlerFunc {
	return func(c telebot.Context) error {
		if err := h.Cancel(c); err != nil {
			slog.Error("failed to cancel operation", slog.Any("error", err))
			return err
		}

		return c.Send(message.WelcomeMessage, keyboard.CreateMainMenu())
	}
}

func (h *Handler) CreateMenuTicket() telebot.HandlerFunc {
	return func(c telebot.Context) error {
		if err := h.Cancel(c); err != nil {
			slog.Error("failed to cancel operation", slog.Any("error", err))
			return err
		}
		return c.Send("Выберите опцию", keyboard.CreateTicketTypeMenu())
	}
}

func (h *Handler) BackToMainMenu() telebot.HandlerFunc {
	return func(c telebot.Context) error {
		if err := h.Cancel(c); err != nil {
			slog.Error("failed to cancel operation", slog.Any("error", err))
			return err
		}

		return c.Send(message.WelcomeMessage, keyboard.CreateMainMenu())
	}
}

func (h *Handler) CancelOperation() telebot.HandlerFunc {
	return func(c telebot.Context) error {
		if err := h.Cancel(c); err != nil {
			slog.Error("failed to cancel operation", slog.Any("error", err))
			return err
		}

		return c.Send("Выберите опцию:", keyboard.CreateTicketTypeMenu())
	}
}

func (h *Handler) Cancel(c telebot.Context) error {
	ctx, _ := h.mn.NewContext(c)

	if err := ctx.Finish(context.Background(), true); err != nil {
		slog.Error("failed to finish FSM state", slog.Any("error", err))
		return err
	}

	return nil
}

func (h *Handler) contentTicketFSM(c telebot.Context, state fsm.Context) error {
	if err := state.Update(context.Background(), "content", c.Text()); err != nil {
		slog.Error("failed to update state", slog.Any("error", err))
		return err
	}

	if err := state.SetState(context.Background(), fsmstate.Confirm); err != nil {
		slog.Error("failed to set content state", slog.Any("error", err))
		return err
	}

	var ticketType string
	if err := state.Data(context.Background(), "type", &ticketType); err != nil {
		slog.Error("failed to update content state", slog.Any("error", err))
		return err
	}

	msg := fmt.Sprintf(
		"📨 Вы заполнили тикет!\n\n"+
			"🔖 *Тип тикета:* %s\n"+
			"📝 *Содержание:* %s\n\n"+
			"✅ Всё верно? Подтвердите отправку или отмените.",
		ticketType, c.Text(),
	)

	return c.Send(msg, keyboard.CreateConfirmMenu())
}

func (h *Handler) confirmTicketFSM(c telebot.Context, state fsm.Context) error {
	switch c.Text() {
	case "✅ Да":
		var ticketType string
		if err := state.Data(context.Background(), "type", &ticketType); err != nil {
			slog.Error("failed to update content state", slog.Any("error", err))
			return err
		}

		var content string
		if err := state.Data(context.Background(), "content", &content); err != nil {
			slog.Error("failed to update content state", slog.Any("error", err))
			return err
		}

		if err := state.Finish(context.Background(), true); err != nil {
			slog.Error("failed to finish FSM state", slog.Any("error", err))
			return err
		}

		return c.Send(message.SentTicket, keyboard.CreateMainMenu())

	case "❌ Нет":
		if err := state.SetState(context.Background(), fsmstate.Content); err != nil {
			slog.Error("failed to set state to Content", slog.Any("error", err))
			return err
		}

		return c.Send("Пожалуйста, введите сообщение заново:", keyboard.CreateCancelMenu())

	default:
		return c.Send("Пожалуйста, выберите одну из кнопок ниже.", keyboard.CreateConfirmMenu())
	}
}

func (h *Handler) onTextFSM(c telebot.Context, state fsm.Context) error {
	switch c.Text() {
	case "📝 Only-one-message":
		if err := state.Update(context.Background(), "type", "only-message"); err != nil {
			slog.Error("failed to update FSM state", slog.Any("error", err))
			return err
		}

		if err := state.SetState(context.Background(), fsmstate.Content); err != nil {
			slog.Error("failed to set FSM state", slog.Any("error", err))
			return err
		}

		return c.Send(message.HelpMessage, keyboard.CreateCancelMenu())

	case "💬 Realtime-chat":
		if err := state.Update(context.Background(), "type", "realtime-chat"); err != nil {
			slog.Error("failed to update FSM state", slog.Any("error", err))
			return err
		}

		if err := state.SetState(context.Background(), fsmstate.Content); err != nil {
			slog.Error("failed to set FSM state", slog.Any("error", err))
			return err
		}

		return c.Send("Опишите вашу проблему для realtime-чата", keyboard.CreateCancelMenu())

	default:
		return nil
	}
}

//
//func (h *Handler) EnsureUserRegisteredMiddleware(next fsm.Handler) fsm.Handler {
//	return func(c telebot.Context, state fsm.Context) error {
//		userId := strconv.FormatInt(c.Sender().ID, 10)
//		username := c.Sender().Username
//
//		if _, err := h.srv.CreateUser(userId, username); err != nil {
//			slog.Error("failed to ensure user registration in FSM",
//				slog.String("user_id", userId),
//				slog.String("username", username),
//				slog.Any("error", err))
//
//			return c.Send("❗ Вы не зарегистрированы в системе. Попробуйте позже.")
//		}
//
//		return next(c, state)
//	}
//}
//
//func (h *Handler) EnsureRegisteredMiddleware(next telebot.HandlerFunc) telebot.HandlerFunc {
//	return func(c telebot.Context) error {
//		userId := strconv.FormatInt(c.Sender().ID, 10)
//		username := c.Sender().Username
//
//		_, err := h.srv.CreateUser(userId, username)
//		if err != nil {
//			slog.Error("failed to ensure user registration",
//				slog.String("user_id", userId),
//				slog.String("username", username),
//				slog.Any("error", err))
//			return c.Send("❗ Произошла ошибка при регистрации. Попробуйте позже.")
//		}
//
//		return next(c)
//	}
//}
