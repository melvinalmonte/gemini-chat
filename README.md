# Gemini Chat

Gemini Chat is a streamlined command-line interface (CLI) tool designed to facilitate real-time conversations with Google's `gemini-pro` model. By leveraging this innovative AI, users can engage in interactive dialogues directly from their terminal. Each chat session is conveniently saved in a markdown file within the `history` directory, ensuring no moment of your interaction is lost.

## Requirements

- **Go**: Gemini Chat is built with Go. Ensure you have Go 1.22.0 or later installed on your system. You can download it from [the official Go website](https://golang.org/dl/).
- **Environment Variable**: To use Gemini Chat, you must have the `GEMINI_API_KEY` environment variable set with your API key for the Gemini model. You can obtain an API key by signing in to [Google AI Studio](https://ai.google.dev/).

## Setup

1. **Clone the Project**: Start by cloning the Gemini Chat repository to your local machine.

    ```sh
    git clone <repository-url>
    cd gemini-chat
    ```

2. **Set Environment Variable**: Before running Gemini Chat, you need to set the `GEMINI_API_KEY` environment variable. This can be done as follows:

   For Unix/Linux/macOS:

    ```sh
    export GEMINI_API_KEY="your_api_key_here"
    ```

   For Windows Command Prompt:

    ```cmd
    set GEMINI_API_KEY=your_api_key_here
    ```

   For Windows PowerShell:

    ```powershell
    $env:GEMINI_API_KEY="your_api_key_here"
    ```

3. **Build the Project**: While in the project directory, compile the application using Go.

    ```sh
    go build
    ```

   This command creates an executable named `gemini-chat` in your current directory.

## Usage

To start a chat session, simply run the executable from your terminal:

```sh
./gemini-chat
```
Follow the prompts to begin chatting with the `gemini-pro` model. Your chat history will be saved in the `.history` directory after each session, with filenames formatted as `chat_history_YYYY_MM_DD_HHMMSS.md`.

To exit the chat, press `CTRL+C` or submit an empty message.

## Example interaction
```
---------------------------------------------------------------------------------------
Terminal chat started. Type something and press enter (CTRL+C to exit).
---------------------------------------------------------------------------------------
Me: Hello, how are you?
Gemini: 
I'm doing well, thank you for asking! How can I assist you today?
Me: What's the weather like?
Gemini: 
I'm not currently connected to real-time data, but I can provide weather forecasting tips if you'd like.
```

## Todos
1. Make it actually useful.
2. Add shift-enter functionality to send messages.
3. Add feature to send files and images.
4. Enhance error handling and user feedback.
