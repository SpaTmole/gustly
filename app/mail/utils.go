package mail

import (
	"fmt"
	"github.com/revel/revel"
)

func MakeActivationLink(key string) string {
	return fmt.Sprintf("%s/activate?activation_key=%s", revel.Config.StringDefault("host", ""), key)
}

func MakeActivationMessage(key string) []byte {
	msg := fmt.Sprintf(
		"Greetings at Gustly!\n\nTo continue using our service, please, proceed with following link: %s \n\n If this wasn't you, please, ignore this message.\n Activation link expires in 12 days. \n Best regards,\nGustly.",
		MakeActivationLink(key),
	)
	return []byte(msg)
}
