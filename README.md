# financebot
Telegram-bot to store and visualize expences

## Installation

Following commads will give you an executable binary of bot
```
git clone https://github.com/Masynchin/financebot.git
cd financebot
go build .
```

## Usage

Bot main commands are (this command are also presented by `/start` command):

- `/add <category> <amount>` - add new expence with following category and amount (example: `/add taxi 250`)
- `/get` - show current month expences. This command will show you them via markup table with `category`, `amount` and `delete button` columns. Click on `delete button` will delete selected expence
- `/chart` will send you a "pie chart" of your month expences grouped by category
