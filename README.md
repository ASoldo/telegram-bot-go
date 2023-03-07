# telegram-bot-go
Telegram Bot in Go. Simple and extensible.

### Telegram Bot Go Lang

This is a Go program that listens for updates from a Telegram Bot API, filters messages based on some rules defined in a semantic cluster, and forwards relevant messages to a specific chat.

The program starts by defining two types, KeywordGroup and SemanticCluster. KeywordGroup is a simple struct that contains a list of keywords. SemanticCluster is a more complex struct that contains a chat ID, a list of KeywordGroup objects, a list of negative words, and minimum and maximum price values. These are used to define the rules for filtering messages.

The program then defines a main function, which starts by defining some command-line flags to specify the Telegram bot API token, the chat ID to forward messages to, and the path to a JSON file containing the semantic clusters. It then parses the command-line flags and checks for missing arguments.

Next, the program reads the contents of the semantic clusters JSON file and unmarshals it into a slice of SemanticCluster objects. If the file is not found or the unmarshalling fails, the program logs an error and exits.

The program then creates a new Telegram bot API client using the provided token and sets the Debug flag to true. It then creates an Update object to retrieve updates from the bot API and starts a loop to receive updates from the bot API.

For each received update, the program checks if the message is not nil and then iterates through each SemanticCluster. For each cluster, the program checks if the message matches the defined rules by calling the isRelevantMessage function. If the message matches the rules, it is forwarded to the specified chat using the forwardMessage function.

The isRelevantMessage function takes a SemanticCluster object and a tgbotapi.Message object and checks if the message matches the rules defined in the cluster. First, it checks if the message comes from a group or channel chat, as defined in the Chat.Type field. If it is not from a group or channel chat, the function returns false.

Next, the function checks if the message contains any of the keywords defined in the KeywordGroup objects. If it does not contain any of the keywords, the function returns false. The function also checks if the message contains any of the negative words defined in the SemanticCluster object. If it contains any negative words, the function returns false.

The function then uses a regular expression to find any price values in the message and checks if they are within the range of minimum and maximum prices defined in the SemanticCluster object. If the price is outside the range, the function returns false.

If the message passes all the rules defined in the SemanticCluster object, the function returns true, indicating that the message is relevant.

The forwardMessage function takes a tgbotapi.BotAPI object, a SemanticCluster object, a tgbotapi.Message object, and a chat ID. It checks if the original message was sent from a channel, and if so, it forwards the message to the chat ID defined in the SemanticCluster object. Otherwise, it forwards the message to the chat ID specified in the command-line arguments.

Overall, this program acts as a filter that forwards only relevant messages based on the defined rules.

```json
[  
    {    
        "chat_id": 6229440871,    
        "keyword_groups": [
            { "keywords": ["Buy", "Purchase", "Get"] },
            { "keywords": ["Dog", "Pet", "Cat"] }
        ],
        "negative_words": ["Miami", "New York", "Boston"],
        "min_price": 200,
        "max_price": 300
    }
]

```
This is a JSON object representing a single SemanticCluster. It has five key-value pairs:

chat_id: an integer value representing the chat ID to forward messages to if they match the defined rules.
keyword_groups: a list of KeywordGroup objects, each containing a list of keywords. The isRelevantMessage function checks if the message contains any of these keywords.
negative_words: a list of strings. The isRelevantMessage function checks if the message contains any of these negative words. If it does, the message is considered not relevant and is not forwarded.
min_price: an integer value representing the minimum price of items in the message for the message to be considered relevant.
max_price: an integer value representing the maximum price of items in the message for the message to be considered relevant.
In this particular example, the SemanticCluster object defines that messages containing any of the words "Buy", "Purchase", or "Get" AND any of the words "Dog", "Pet", or "Cat" AND no negative words of "Miami", "New York", or "Boston" are relevant. Additionally, the message must contain a price between 200 and 300 dollars to be considered relevant. If a message meets all of these criteria, it will be forwarded to the chat with ID xxxxxxxxxx.
