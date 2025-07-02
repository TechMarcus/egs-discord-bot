import os
import discord
from bot import Mybot


def main():
    intents = discord.Intents.default()
    bot = Mybot(intents)
    bot.discord_bot(os.environ.get('DISCORD_TOKEN'))


if __name__ == '__main__':
    main()
