package commands

import (
	tele "gopkg.in/telebot.v3"
)


type MarkUp struct{
	Text string
	Unique string
	Data string
}
func SendMarkupSelector(mk []MarkUp, selector *tele.ReplyMarkup) *tele.ReplyMarkup{
	var buttons []tele.Btn
   for _, d := range mk {
	 buttons = append(buttons, selector.Data(d.Text, d.Unique, d.Data))
   }
 // Create rows for each button
 var rows []tele.Row
 for _, btn := range buttons {
	 rows = append(rows, selector.Row(btn))
 }

 // Set all rows at once
 selector.Inline(rows...)
 return selector

}