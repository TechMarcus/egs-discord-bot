import discord
import json
import os
import time
from apscheduler.schedulers.asyncio import AsyncIOScheduler

async def send_game_info(title, image_url, channel=None, message=None):
    embed = discord.Embed(title=title)
    embed.set_image(url=image_url)
    if message is not None:
        await message.channel.send(embed=embed)
        return
    await channel.send(embed=embed)


async def send_avalible_games(client=None, channel=None, message=None):
    print("sending avalible games")


    print(f"send_avalible_games called at {time.strftime('%X')}")
    filepath = '/app/games_info.json'
    mtime = os.path.getmtime(filepath)
    print(f"File modification time: {mtime}")

    with open(filepath, 'r') as file:
        games_data = json.load(file)
        print(f"Loaded games data: {games_data}")

        if message is not None:
            for game in games_data:
                await send_game_info(game['Name'], game['Picture'], message=message)
            return
        await client.wait_until_ready()
        for game in games_data:
            await send_game_info(game['Name'], game['Picture'], channel=channel)


def daily_check_handler(client, channel):
    scheduler = AsyncIOScheduler()
    scheduler.add_job(send_avalible_games, 'cron', day_of_week='tue', hour=0, minute=1, args=[client, channel])

    print('daily_check_handler loaded')
    scheduler.start()
