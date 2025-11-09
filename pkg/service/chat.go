package service

import (
	"fmt"
	"log"
	"poststone-chat/pkg/db"
	"poststone-chat/pkg/model"
	"poststone-chat/pkg/openai"
	"poststone-chat/pkg/texts"
	"strings"
)

func Answer(message string) (string, error) {
	if strings.TrimSpace(message) == "" {
		return texts.DefaultAnswer, nil
	}

	var monuments []model.Monument
	reqType, value := openai.InterpretQuery(message)
	log.Printf("Interpretated: %s key phrase: %s ", reqType, value)

	if reqType == openai.IntentGeneral {
		return value, nil
	} else if reqType == openai.IntentCheap {
		monuments = db.Cheapest()
	} else if reqType == openai.IntentExpensive {
		monuments = db.MostExpensive()
	} else if reqType == openai.IntentSearch {
		vector, err := openai.EmbedUserPhrase(value)
		if err != nil {
			log.Printf("error while phrase emdedding %v", err)
		}
		monuments = db.SearchByVector(vector)
	}

	if monuments != nil {
		return present(monuments), nil
	}

	return texts.DefaultAnswer, nil
}

func present(monuments []model.Monument) string {
	if len(monuments) == 0 {
		return "К сожалению, мы не нашли подходящих памятников. Пожалуйста, уточните ваш запрос."
	}

	var sb strings.Builder

	sb.WriteString("Вот что я подобрал для Вас:\n\n")

	for _, m := range monuments {
		sb.WriteString(fmt.Sprintf(
			"🔹 <b>%s</b>\n%s\n💰 %.2f BYN\n🔗 <a href=\"%s\">Посмотреть</a>\n\n",
			m.Title,
			m.Description,
			m.Price,
			m.Link,
		))
	}

	sb.WriteString("Если остались вопросы — звоните или приходите к нам лично:\n\n")
	sb.WriteString("📞 <a href=\"tel:+375296911400\">+375 (29) 6-911-400</a>\n")
	sb.WriteString("🏛 г. Минск, ул. Кальварийская 33\n\n")
	sb.WriteString("Могу ли я вам ещё чем-то помочь?")

	return sb.String()
}
