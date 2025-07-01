import discord, datetime, asyncio, json, time
from apscheduler.schedulers.background import BackgroundScheduler

async def send_game_info(title, image_url,channel=None,message=None):
    embed = discord.Embed(title=title)
    embed.set_image(url=image_url)
    if message != None:
         await message.channel.send(embed=embed)
         return
    await channel.send(embed=embed)

async def send_avalible_games(message=None, client=None, channel=None):
    print("sending avalible games")
    with json.loads(open('./games_info.json').read()) as games_data:
        if message != None:
            for game in games_data:
                await send_game_info(game['Name'], game['Picture'], message=message)
            return
        await client.wait_until_ready()
        for game in games_data:
            await send_game_info(game['Name'], game['Picture'], channel=channel)
    

def daily_check_handler(client, channel):
    scheduler = BackgroundScheduler()

    scheduler.add_job(
        lambda: asyncio.create_task(send_avalible_games(client=client, channel=channel)),
        'cron',
        day_of_week='tue',
        hour=0,
        minute=1
    )
    
    print('daily_check_handler loaded')
    scheduler.start()
