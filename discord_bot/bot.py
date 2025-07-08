import discord
import asyncio
from bot_tasks import send_avalible_games


class Mybot:
    def __init__(self, intents):
        self.intents = intents
        self.intents.message_content = True
        self.client = discord.Client(intents=self.intents)

    def discord_bot(self, TOKEN):
        client = self.client

        @client.event
        async def on_ready():
            print(f'We have logged in as {client.user}')

            for channel in client.guilds[0].channels:
                if channel.name == "жидівське-лігво":
                    games_channel = client.get_channel(channel.id)
            asyncio.create_task(send_avalible_games(client=client, channel=games_channel))
            await client.close()
            # daily_check_handler(client=client, channel=games_channel)
            # print('on_ready functions loaded')
    

        # @client.event
        # async def on_message(message):
        #     if message.author == client.user:
        #         return

            # if message.author.name == 'maksred_ay':
            #     await message.reply("ХВОЙДІ СЛОВА НЕ ДАВАЛИ")

            # if message.content.startswith('$games'):
            #     asyncio.create_task(send_avalible_games(message=message))
            # if message.content.startswith('$goida'):
            #     await message.channel.send('Hello there!')

        client.run(TOKEN)

