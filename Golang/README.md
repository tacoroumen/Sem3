# Discord OS Information Sender

This Go program gathers information about the operating system of a VM and sends it to a Discord webhook.

## Installation

1. Make sure you have Go installed on your system.
2. Clone this repository.
3. Navigate to the directory containing the source code.
4. Run `go build` to build the executable.

## Usage

1. Obtain a Discord webhook URL where you want to send the OS information.
2. Replace `{{webhookURL}}` in the source code with your Discord webhook URL.
3. Run the built executable.

The program will gather OS information using `uname -a` and `lsb_release -a` commands and send it to the specified Discord webhook.

## Dependencies

This program requires Go's standard library only.

## Disclaimer

This program is provided as is without any warranty. Use it at your own risk.

## Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request.

## License

This project is licensed under the [MIT License](../LICENSE).
