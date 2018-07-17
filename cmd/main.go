package main

import (
    "flag"
    "os"
    "fmt"
)

func main() {
    askCommand := flag.NewFlagSet("Question", flag.ExitOnError)
    questionFlag := askCommand.String("question", "", "Question that you are asking for")

    sendCommand := flag.NewFlagSet("send", flag.ExitOnError)
    recipientFlag := sendCommand.String("recipient", "", "Recipient of your message")
    messageFlag := sendCommand.String("message", "", "Text message")

    switch os.Args[1] {
    case "ask":
        askCommand.Parse(os.Args[2:])
    case "send":
        sendCommand.Parse(os.Args[2:])
    default:
        fmt.Printf("%q is not a valid command.\n", os.Args[1])
        os.Exit(2)
    }

    if askCommand.Parsed() {
        if *questionFlag == "" {
            fmt.Println("Please supply the question using -question.")
            return
        }
        fmt.Printf("You asked: %q\n", *questionFlag)
    }

    if sendCommand.Parsed() {
        if *recipientFlag == "" {
            fmt.Println("Please supply the recipient using -recipient option.")
            return
        }

        if *messageFlag == "" {
            fmt.Println("Please supply the message using the -message option.")
            return
        }

        fmt.Printf("Your message is sent to %q.\n", *recipientFlag)
        fmt.Printf("Message: %q.\n", *messageFlag)
    }
}
