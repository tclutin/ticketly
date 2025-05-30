package handler

import (
	"context"
	"errors"
	"fmt"
	fsmstate "github.com/tclutin/ticketly/telegram_bot/internal/fsm"
	"github.com/tclutin/ticketly/telegram_bot/internal/keyboard"
	"github.com/tclutin/ticketly/telegram_bot/internal/message"
	"github.com/tclutin/ticketly/telegram_bot/internal/service"
	"github.com/tclutin/ticketly/telegram_bot/pkg/client/ticketly"
	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	"github.com/vitaliy-ukiru/fsm-telebot/v2/fsmopt"
	"gopkg.in/telebot.v4"
	"log/slog"
	"unicode/utf8"
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

	h.bot.Handle(keyboard.BtnCreateTicket, h.CreateMenuTicket())

	h.bot.Handle(keyboard.BtnBack, h.BackToMainMenu())

	h.bot.Handle(keyboard.BtnCancel, h.CancelOperation())
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

		return c.Send("Выберите опцию", keyboard.CreateTicketTypeMenu())
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
	text := c.Text()

	length := utf8.RuneCountInString(text)
	if length > 500 || length < 6 {
		return c.Send("❗ Текст должен быть от 6 до 500 символов. Пожалуйста, введите текст заново.")
	}

	if err := state.Update(context.Background(), "content", text); err != nil {
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
		ticketType, text,
	)

	return c.Send(msg, keyboard.CreateConfirmMenu())
}

func (h *Handler) confirmTicketFSM(c telebot.Context, state fsm.Context) error {
	switch c.Text() {
	case keyboard.BtnConfirmYes:
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

		userId := c.Get("user_id")
		if userId == nil {
			slog.Error("failed to get user_id from context")
			return c.Send("❗ Произошла ошибка при регистрации. Попробуйте позже.")
		}

		id, ok := userId.(uint64)
		if !ok {
			slog.Error("failed to get user_id from context")
			return c.Send("❗ Произошла ошибка при регистрации. Попробуйте позже.")
		}

		//c.Chat().ID
		//нужно было сделать отдельный ендпоинт, чтобы дёргать тикеты человека, и при нажатии на кнопку проверять количество и тд, да уж
		if err := h.srv.CreateTicket(id, c.Sender().ID, ticketType, content); err != nil {
			if errors.Is(err, ticketly.ErrActiveTicketAlreadyExists) {
				if err = h.Cancel(c); err != nil {
					slog.Error("failed to cancel operation", slog.Any("error", err))
					return err
				}

				msg := "❗ У вас уже есть активное обращение\n\n" +
					"Пока оно не будет обработано, создание новых обращений недоступно. " +
					"Наш специалист свяжется с вами в ближайшее время."

				return c.Send(msg, keyboard.CreateMainMenu())
			}

			slog.Error("failed to create ticket", slog.Any("error", err))
			return err
		}

		if err := state.Finish(context.Background(), true); err != nil {
			slog.Error("failed to finish FSM state", slog.Any("error", err))
			return err
		}

		return c.Send(message.SentTicket, keyboard.CreateMainMenu())

	case keyboard.BtnConfirmNo:
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
	case keyboard.BtnTicketTypeSingle:
		if err := state.Update(context.Background(), "type", "only-message"); err != nil {
			slog.Error("failed to update FSM state", slog.Any("error", err))
			return err
		}

		if err := state.SetState(context.Background(), fsmstate.Content); err != nil {
			slog.Error("failed to set FSM state", slog.Any("error", err))
			return err
		}

		return c.Send(message.HelpMessage, keyboard.CreateCancelMenu())

	case keyboard.BtnTicketTypeRealtime:
		if err := state.Update(context.Background(), "type", "realtime-chat"); err != nil {
			slog.Error("failed to update FSM state", slog.Any("error", err))
			return err
		}

		if err := state.SetState(context.Background(), fsmstate.Content); err != nil {
			slog.Error("failed to set FSM state", slog.Any("error", err))
			return err
		}

		return c.Send(message.RealtimePrompt, keyboard.CreateCancelMenu())

	default:
		return nil
	}
}
